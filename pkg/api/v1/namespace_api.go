package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type NamespaceStartRequestMessage struct {
	Commit
	InstanceName string `json:"instanceName"`
}

type NamespaceRequestBuilder struct {
	v1          *V1
	Provider    string
	Org         string
	ClusterID   int64
	NamespaceID int64
	Name        string // kubernetes namespace
	Commit      *Commit
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

func (b *NamespaceRequestBuilder) WithCommit(commit *Commit) *NamespaceRequestBuilder {
	b.Commit = commit
	return b
}

func (b *NamespaceRequestBuilder) Create(alias string) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	if b.Commit == nil {
		return nil, fmt.Errorf("Commit is not set. call withCommit()")
	}

	client := b.v1.HttpClient()

	requestBody := NamespaceStartRequestMessage{
		Commit:       *b.Commit,
		InstanceName: alias,
	}
	jsonValue, _ := json.Marshal(&requestBody)
	jsonData := bytes.NewBuffer(jsonValue)

	req, err := b.v1.NewPostRequest(
		fmt.Sprintf("/orgs/%s/%s/vc/%d/instance", b.Provider, b.Org, b.ClusterID),
		jsonData,
	)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("Project, branch or commit not found")
	} else if resp.StatusCode == 402 {
		return nil, fmt.Errorf("Not authorized")
	} else if resp.StatusCode == 402 {
		return nil, fmt.Errorf("Payment required")
	} else if resp.StatusCode == 428 {
		return nil, fmt.Errorf("Precondition required")
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d Error", resp.StatusCode)
	}

	// parse json response
	namespaces := StaroidNamespace{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&namespaces)
	if err != nil {
		return nil, err
	}

	return &namespaces, nil
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
