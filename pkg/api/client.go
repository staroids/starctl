package api

import (
	v1 "github.com/staroids/starctl/pkg/api/v1"
	"github.com/staroids/starctl/pkg/auth"
)

// StaroidClient is rest api client for staroid.com
type StaroidClient struct {
	Auth auth.StaroidAuth
}

func (c *StaroidClient) V1() *v1.V1 {
	return &v1.V1{
		Auth: c.Auth,
	}
}
