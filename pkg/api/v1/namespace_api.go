package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

	err = GetApiErrorFromResponse(resp, map[int]string{
		409: "Already exists",
	})
	if err != nil {
		return nil, err
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

func (b *NamespaceRequestBuilder) GetById(namespaceId int64) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()
	req, err := b.v1.NewGetRequest(
		fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d",
			b.Provider, b.Org, b.ClusterID, namespaceId))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse json response
	namespaces := StaroidNamespace{}
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

func (b *NamespaceRequestBuilder) Delete(alias string) (*StaroidNamespace, error) {
	ns, err := b.Get(alias)
	if err != nil {
		return nil, err
	}

	return b.DeleteById(ns.ID)
}

func (b *NamespaceRequestBuilder) DeleteById(namespaceID int64) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	req, err := b.v1.NewDeleteRequest(fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d", b.Provider, b.Org, b.ClusterID, namespaceID))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = GetApiErrorFromResponse(resp, map[int]string{})
	if err != nil {
		return nil, err
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

func (b *NamespaceRequestBuilder) Start(alias string) (*StaroidNamespace, error) {
	ns, err := b.Get(alias)
	if err != nil {
		return nil, err
	}

	return b.StartById(ns.ID)
}

func (b *NamespaceRequestBuilder) StartById(namespaceID int64) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	req, err := b.v1.NewPutRequest(fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d/resume", b.Provider, b.Org, b.ClusterID, namespaceID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = GetApiErrorFromResponse(resp, map[int]string{})
	if err != nil {
		return nil, err
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

func (b *NamespaceRequestBuilder) Stop(alias string) (*StaroidNamespace, error) {
	ns, err := b.Get(alias)
	if err != nil {
		return nil, err
	}

	return b.StartById(ns.ID)
}

func (b *NamespaceRequestBuilder) StopById(namespaceID int64) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	req, err := b.v1.NewPutRequest(fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d/pause", b.Provider, b.Org, b.ClusterID, namespaceID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = GetApiErrorFromResponse(resp, map[int]string{})
	if err != nil {
		return nil, err
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

func (b *NamespaceRequestBuilder) Get(alias string) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	// find namespace by alias
	namespaces, err := b.GetAll()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	var found *StaroidNamespace = nil
	for _, ns := range *namespaces {
		if ns.Alias == alias {
			found = &ns
			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("Alias %s not found", alias)
	}

	return found, nil
}
