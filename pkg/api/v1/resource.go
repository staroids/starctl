package v1

import (
	"fmt"
	"strings"

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
	Type      string `json:"type"`
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

type Commit struct {
	Provider string `json:"provider"`
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	Branch   string `json:"branch"`
	Commit   string `json:"commit"`
}

// commitLoc format is [Provider]/[]
func NewCommitFromCommitLocation(commitLoc string) (*Commit, error) {
	commitHashPos := strings.Index(commitLoc, "#")
	commit := ""
	if commitHashPos > 0 {
		commit = commitLoc[commitHashPos+1:]
	}
	branchPos := strings.Index(commitLoc, ":")
	branch := ""
	if branchPos < 0 {
		return nil, fmt.Errorf("Invalid project flag. No branch info.")
	}
	if commitHashPos > 0 {
		branch = commitLoc[branchPos+1 : commitHashPos]
	} else {
		branch = commitLoc[branchPos+1:]
	}

	if len(branch) == 0 {
		return nil, fmt.Errorf("Invalid project flag. No branch info.")
	}

	proj := strings.Split(commitLoc[:branchPos], "/")
	if len(proj) != 3 {
		return nil, fmt.Errorf("Invalid project flag.")
	}

	return &Commit{
		Provider: proj[0],
		Owner:    proj[1],
		Repo:     proj[2],
		Branch:   branch,
		Commit:   commit,
	}, nil
}
