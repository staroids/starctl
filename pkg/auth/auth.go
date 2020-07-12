package auth

import (
	"fmt"
	"os"

	"github.com/staroids/starctl/pkg/constants"
)

// StaroidAuth handles authentication
type StaroidAuth struct {
	accessToken string
	apiServer   string
}

// CheckAuth checks authentication
func (s *StaroidAuth) CheckAuth() error {
	token := s.AccessToken()
	if token == "" {
		return fmt.Errorf("Please set %s environment variable.", constants.EnvStaroidAccessToken)
	}
	return nil
}

// AccessToken returns Access Token
func (s *StaroidAuth) AccessToken() string {
	s.accessToken = os.Getenv(constants.EnvStaroidAccessToken)
	return s.accessToken
}

func (s *StaroidAuth) ApiServer() string {
	s.apiServer = os.Getenv(constants.EnvStaroidApiServer)
	if s.apiServer == "" {
		s.apiServer = constants.ApiServer
	}
	return s.apiServer
}
