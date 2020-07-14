package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/staroids/starctl/pkg/constants"
	corev1 "k8s.io/api/core/v1"
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

func (b *NamespaceRequestBuilder) GetById(namespaceID int64) (*StaroidNamespace, error) {
	return b.namespaceOP(namespaceID, "GET", "")
}

func (b *NamespaceRequestBuilder) GetShellService() (*corev1.Service, error) {
	resources, err := b.GetAllResources()
	if err != nil {
		return nil, err
	}

	var shellService *corev1.Service = nil
	for _, service := range resources.Services.Items {
		if service.ObjectMeta.Labels[constants.K8S_LABEL_KEY_RESOURCE_SYSTEM] == constants.K8S_LABEL_VALUE_RESOURCE_SYSTEM_SHELL {
			shellService = &service
			break
		}
	}
	return shellService, nil
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
	return b.namespaceOP(namespaceID, "DELETE", "")
}

func (b *NamespaceRequestBuilder) Start(alias string) (*StaroidNamespace, error) {
	ns, err := b.Get(alias)
	if err != nil {
		return nil, err
	}

	return b.StartById(ns.ID)
}

func (b *NamespaceRequestBuilder) StartById(namespaceID int64) (*StaroidNamespace, error) {
	return b.namespaceOP(namespaceID, "PUT", "resume")
}

func (b *NamespaceRequestBuilder) Stop(alias string) (*StaroidNamespace, error) {
	ns, err := b.Get(alias)
	if err != nil {
		return nil, err
	}

	return b.StartById(ns.ID)
}

func (b *NamespaceRequestBuilder) StopById(namespaceID int64) (*StaroidNamespace, error) {
	return b.namespaceOP(namespaceID, "PUT", "pause")
}

func (b *NamespaceRequestBuilder) namespaceOP(namespaceID int64, method string, op string) (*StaroidNamespace, error) {
	if b.Provider == "" || b.Org == "" {
		return nil, fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return nil, fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	opPath := ""
	if op != "" {
		opPath = fmt.Sprintf("/%s", op)
	}

	path := fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d%s", b.Provider, b.Org, b.ClusterID, namespaceID, opPath)
	req, err := b.v1.NewRequest(method, path, nil)
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

func (b *NamespaceRequestBuilder) ShellStartById(namespaceID int64) error {
	if b.Provider == "" || b.Org == "" {
		return fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	path := fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d/shell", b.Provider, b.Org, b.ClusterID, namespaceID)
	req, err := b.v1.NewPostRequest(path, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = GetApiErrorFromResponse(resp, map[int]string{})
	if err != nil {
		return err
	}
	return nil
}

func (b *NamespaceRequestBuilder) ShellStopById(namespaceID int64) error {
	if b.Provider == "" || b.Org == "" {
		return fmt.Errorf("Org information is not set. call withOrg()")
	}

	if b.ClusterID == int64(0) {
		return fmt.Errorf("Cluster ID is not set. call withClusterID()")
	}

	client := b.v1.HttpClient()

	path := fmt.Sprintf("/orgs/%s/%s/vc/%d/instance/%d/shell", b.Provider, b.Org, b.ClusterID, namespaceID)
	req, err := b.v1.NewDeleteRequest(path)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = GetApiErrorFromResponse(resp, map[int]string{})
	if err != nil {
		return err
	}
	return nil
}
