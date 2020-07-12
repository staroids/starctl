package v1

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type StaroidSke struct {
	ID     string `json:"name"`
	Cloud  string `json:"cloud"`
	Region string `json:"region"`
}

type StaroidCluster struct {
	ID    int64      `json:"id"`
	Name  string     `json:"name"`
	Ske   StaroidSke `json:"ske"`
	OrgID int64      `json:"orgId"`
	Type  string     `json:"type"`
}

type StaroidNamespace struct {
	ID        int64  `json:"id"`
	Namespace string `json:"name"`
	Alias     string `json:"instanceName"`
	Phase     string `json:"phase"`
	Status    string `json:"status"`
	Access    string `json:"access"`
	URL       string `json:"url"`
}

func (n *StaroidNamespace) ServiceURL(serviceName string, port int) string {
	return fmt.Sprintf("https://p%d-%s--%s", port, serviceName, n.URL[len("https://"):])
}

type StaroidOrg struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ID       int64  `json:"id"`
}

type StaroidNamespaceResources struct {
	Services v1.ServiceList `json:"services"`
}
