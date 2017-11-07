package server

import (
	"context"
	api "github.com/SHyx0rmZ/go-bitbucket/server"
	"github.com/concourse/atc/auth/bitbucket"
	"net/http"
)

type client struct {
	endpoint string
}

func NewClient(endpoint string) bitbucket.Client {
	return &client{
		endpoint: endpoint,
	}
}

func (c *client) CurrentUser(httpClient *http.Client) (string, error) {
	bc, err := api.NewClient(context.TODO(), httpClient, c.endpoint)
	if err != nil {
		return "", err
	}

	return bc.CurrentUser()
}

func (c *client) Teams(httpClient *http.Client, role bitbucket.Role) ([]string, error) {
	return nil, nil
}

func (c *client) Projects(httpClient *http.Client) ([]string, error) {
	bc, err := api.NewClient(context.TODO(), httpClient, c.endpoint)
	if err != nil {
		return nil, err
	}

	ps, err := bc.Projects()
	if err != nil {
		return nil, err
	}

	s := make([]string, len(ps))
	for i, p := range ps {
		s[i] = p.GetKey()
	}

	return s, nil
}
