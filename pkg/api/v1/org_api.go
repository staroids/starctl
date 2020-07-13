package v1

import "encoding/json"

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
