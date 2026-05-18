package status

import (
	"fmt"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/kubectl"
)

func Collect(kr *kubectl.Runner, namespace string, releases []string) (*Report, error) {
	report := &Report{}
	groupMap := map[string]*Group{}

	for _, release := range releases {
		label := fmt.Sprintf("app.kubernetes.io/name=%s", release)
		pods, err := kr.GetPods(namespace, label)
		if err != nil {
			pods = nil
		}

		rs := classifyRelease(release, pods)

		if rs.Health == Degraded && len(pods) > 0 {
			logs, _ := kr.GetLogs(pods[0].Name, namespace, 20)
			rs.Logs = logs
		}

		groupName := GroupFor(release)
		if _, ok := groupMap[groupName]; !ok {
			groupMap[groupName] = &Group{Name: groupName}
		}
		groupMap[groupName].Releases = append(groupMap[groupName].Releases, rs)

		switch rs.Health {
		case Healthy:
			report.Healthy++
		case Degraded:
			report.Degraded++
		case Warning:
			report.Warning++
		}
	}

	orderedGroups := []string{"INFRASTRUCTURE", "KAFKA", "STORAGE", "RADAR SERVICES", "IDENTITY", "OTHER"}
	for _, g := range orderedGroups {
		if grp, ok := groupMap[g]; ok {
			report.Groups = append(report.Groups, *grp)
		}
	}

	return report, nil
}

func classifyRelease(name string, pods []kubectl.Pod) ReleaseStatus {
	rs := ReleaseStatus{Name: name}
	if len(pods) == 0 {
		rs.Health = Unknown
		rs.Message = "no pods found"
		return rs
	}

	rs.TotalPods = len(pods)
	for _, p := range pods {
		if p.Ready {
			rs.ReadyPods++
		}
		if p.Message != "" {
			rs.Error = p.Message
		}
	}

	switch {
	case rs.ReadyPods == rs.TotalPods:
		rs.Health = Healthy
	case rs.ReadyPods == 0:
		rs.Health = Degraded
		if rs.Error == "" {
			rs.Error = pods[0].Phase
		}
	default:
		rs.Health = Warning
		rs.Message = fmt.Sprintf("%d/%d pods ready", rs.ReadyPods, rs.TotalPods)
	}
	return rs
}
