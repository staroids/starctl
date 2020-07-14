package v1

import (
	"fmt"
	"io"
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
	return v.NewRequest("GET", path, nil)
}

func (v *V1) NewDeleteRequest(path string) (*http.Request, error) {
	return v.NewRequest("DELETE", path, nil)
}

func (v *V1) NewPostRequest(path string, body io.Reader) (*http.Request, error) {
	return v.NewRequest("POST", path, body)
}

func (v *V1) NewPutRequest(path string, body io.Reader) (*http.Request, error) {
	return v.NewRequest("PUT", path, body)
}

func (v *V1) NewRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", v.Auth.ApiServer(), path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", v.Auth.AccessToken()))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (v *V1) HttpClient() *http.Client {
	return &http.Client{}
}

func GetApiErrorFromResponse(resp *http.Response, customErrorMessage map[int]string) error {
	message := map[int]string{
		404: "Not found",
		402: "Not authorized",
	}

	for k, v := range customErrorMessage {
		message[k] = v
	}

	if resp.StatusCode == 200 {
		return nil
	}

	if m, ok := message[resp.StatusCode]; ok {
		return fmt.Errorf("%d %s", resp.StatusCode, m)
	}
	return fmt.Errorf("%d", resp.StatusCode)
}
