package kubectl

import (
	"encoding/json"
	"fmt"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/executor"
)

type Pod struct {
	Name      string
	Namespace string
	Phase     string
	Ready     bool
	Restarts  int
	Message   string
}

type Ingress struct {
	Name      string
	Namespace string
	Host      string
	Paths     []string
}

type Runner struct {
	exec    executor.Executor
	context string
}

func NewRunner(exec executor.Executor, context string) *Runner {
	return &Runner{exec: exec, context: context}
}

func (r *Runner) contextArgs() []string {
	if r.context != "" {
		return []string{"--context", r.context}
	}
	return nil
}

func (r *Runner) GetPods(namespace, labelSelector string) ([]Pod, error) {
	args := append(r.contextArgs(), "get", "pods", "-n", namespace, "-l", labelSelector, "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return nil, fmt.Errorf("kubectl get pods: %w", err)
	}

	var raw struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Status struct {
				Phase             string `json:"phase"`
				ContainerStatuses []struct {
					Ready        bool `json:"ready"`
					RestartCount int  `json:"restartCount"`
				} `json:"containerStatuses"`
				Message string `json:"message"`
			} `json:"status"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("parsing pod list: %w", err)
	}

	pods := make([]Pod, 0, len(raw.Items))
	for _, item := range raw.Items {
		ready := false
		restarts := 0
		if len(item.Status.ContainerStatuses) > 0 {
			ready = item.Status.ContainerStatuses[0].Ready
			restarts = item.Status.ContainerStatuses[0].RestartCount
		}
		pods = append(pods, Pod{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Phase:     item.Status.Phase,
			Ready:     ready,
			Restarts:  restarts,
			Message:   item.Status.Message,
		})
	}
	return pods, nil
}

func (r *Runner) GetLogs(podName, namespace string, tail int) (string, error) {
	args := append(r.contextArgs(), "logs", podName, "-n", namespace, fmt.Sprintf("--tail=%d", tail))
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return "", fmt.Errorf("kubectl logs %s: %w", podName, err)
	}
	return out, nil
}

func (r *Runner) GetIngresses(namespace string) ([]Ingress, error) {
	args := append(r.contextArgs(), "get", "ingress", "-n", namespace, "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return nil, fmt.Errorf("kubectl get ingress: %w", err)
	}

	var raw struct {
		Items []struct {
			Metadata struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
			} `json:"metadata"`
			Spec struct {
				Rules []struct {
					Host string `json:"host"`
					HTTP struct {
						Paths []struct {
							Path string `json:"path"`
						} `json:"paths"`
					} `json:"http"`
				} `json:"rules"`
			} `json:"spec"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return nil, fmt.Errorf("parsing ingress list: %w", err)
	}

	ingresses := make([]Ingress, 0, len(raw.Items))
	for _, item := range raw.Items {
		var host string
		var paths []string
		if len(item.Spec.Rules) > 0 {
			host = item.Spec.Rules[0].Host
			for _, p := range item.Spec.Rules[0].HTTP.Paths {
				paths = append(paths, p.Path)
			}
		}
		ingresses = append(ingresses, Ingress{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Host:      host,
			Paths:     paths,
		})
	}
	return ingresses, nil
}

func (r *Runner) Version() (string, error) {
	args := append(r.contextArgs(), "version", "--client", "-o", "json")
	out, err := r.exec.Run("kubectl", args...)
	if err != nil {
		return "", err
	}
	var v struct {
		ClientVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"clientVersion"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		return out, nil
	}
	return v.ClientVersion.GitVersion, nil
}
