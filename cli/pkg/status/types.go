package status

type Health string

const (
	Healthy  Health = "healthy"
	Degraded Health = "degraded"
	Warning  Health = "warning"
	Unknown  Health = "unknown"
)

type ReleaseStatus struct {
	Name      string `json:"name"`
	Health    Health `json:"status"`
	ReadyPods int    `json:"ready_pods"`
	TotalPods int    `json:"total_pods"`
	Error     string `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
	Logs      string `json:"logs,omitempty"`
}

type Group struct {
	Name     string          `json:"group"`
	Releases []ReleaseStatus `json:"releases"`
}

type Report struct {
	Context  string  `json:"context"`
	Groups   []Group `json:"groups"`
	Healthy  int     `json:"healthy"`
	Degraded int     `json:"degraded"`
	Warning  int     `json:"warning"`
}

var releaseGroups = map[string]string{
	"cert-manager":           "INFRASTRUCTURE",
	"kube-prometheus-stack":  "INFRASTRUCTURE",
	"nginx-ingress":          "INFRASTRUCTURE",
	"zookeeper":              "KAFKA",
	"kafka":                  "KAFKA",
	"schema-registry":        "KAFKA",
	"ksql-server":            "KAFKA",
	"mongodb":                "STORAGE",
	"postgresql":             "STORAGE",
	"redis":                  "STORAGE",
	"minio":                  "STORAGE",
	"radar-appserver":        "RADAR SERVICES",
	"management-portal":      "RADAR SERVICES",
	"radar-fitbit-connector": "RADAR SERVICES",
	"radar-s3-connector":     "RADAR SERVICES",
	"kratos":                 "IDENTITY",
	"hydra":                  "IDENTITY",
}

func GroupFor(releaseName string) string {
	if g, ok := releaseGroups[releaseName]; ok {
		return g
	}
	return "OTHER"
}
