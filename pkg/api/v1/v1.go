package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/staroids/starctl/pkg/auth"
)

type V1 struct {
	Auth auth.StaroidAuth
}

func (v *V1) Cluster() *ClusterRequestBuilder {
	return &ClusterRequestBuilder{
		v1: v,
	}
}

func (v *V1) Namespace() *NamespaceRequestBuilder {
	return &NamespaceRequestBuilder{
		v1: v,
	}
}

func (v *V1) Org() *OrgRequestBuilder {
	return &OrgRequestBuilder{
		v1: v,
	}
}

func (v *V1) NewGetRequest(path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", v.Auth.ApiServer(), path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", v.Auth.AccessToken()))
	return req, nil
}

func (v *V1) HttpClient() *http.Client {
	return &http.Client{}
}

type ClusterRequestBuilder struct {
	v1       *V1
	Provider string
	Org      string
}

func (b *ClusterRequestBuilder) WithOrg(provider string, org string) *ClusterRequestBuilder {
	b.Provider = provider
	b.Org = org
	return b
}

func (b *ClusterRequestBuilder) GetAll() (*[]StaroidCluster, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	client := b.v1.HttpClient()
	req, err := b.v1.NewGetRequest(fmt.Sprintf("/orgs/%s/%s/vc", b.Provider, b.Org))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse json response
	clusters := make([]StaroidCluster, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&clusters)
	if err != nil {
		return nil, err
	}

	return &clusters, nil
}

type NamespaceRequestBuilder struct {
	v1          *V1
	Provider    string
	Org         string
	ClusterID   int64
	NamespaceID int64
	Name        string // kubernetes namespace
}

func (b *NamespaceRequestBuilder) WithOrg(provider string, org string) *NamespaceRequestBuilder {
	b.Provider = provider
	b.Org = org
	return b
}

func (b *NamespaceRequestBuilder) WithClusterID(clusterID int64) *NamespaceRequestBuilder {
	b.ClusterID = clusterID
	return b
}

func (b *NamespaceRequestBuilder) WithNamespaceID(namespaceID int64) *NamespaceRequestBuilder {
	b.NamespaceID = namespaceID
	return b
}

func (b *NamespaceRequestBuilder) WithName(name string) *NamespaceRequestBuilder {
	b.Name = name
	return b
}

func (b *NamespaceRequestBuilder) GetAll() (*[]StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()
	req, err := b.v1.NewGetRequest(fmt.Sprintf("/orgs/%s/%s/vc/%d/instance", b.Provider, b.Org, b.ClusterID))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse json response
	namespaces := make([]StaroidNamespace, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&namespaces)
	if err != nil {
		return nil, err
	}

	return &namespaces, nil
}

func (b *NamespaceRequestBuilder) GetAllResources() (*StaroidNamespaceResources, error) {
	if b.Name == "" {
		return nil, fmt.Errorf("Name is not set. call withName()")
	}

	client := b.v1.HttpClient()
	req, err := b.v1.NewGetRequest(fmt.Sprintf("/namespace/%s", b.Name))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse json response
	resources := StaroidNamespaceResources{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&resources)
	if err != nil {
		return nil, err
	}

	return &resources, nil
}

type OrgRequestBuilder struct {
	v1 *V1
}

func (b *OrgRequestBuilder) GetAll() (*[]StaroidOrg, error) {
	client := b.v1.HttpClient()
	req, err := b.v1.NewGetRequest("/orgs/")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse json response
	orgs := make([]StaroidOrg, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&orgs)
	if err != nil {
		return nil, err
	}

	return &orgs, nil
}
