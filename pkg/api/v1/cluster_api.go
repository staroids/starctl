package v1

import (
	"encoding/json"
	"fmt"
)

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
