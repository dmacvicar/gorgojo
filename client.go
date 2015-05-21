package gorgojo

import (
	"github.com/kolo/xmlrpc"
	"net/url"
)

type Client struct {
	bzUrl     string
	rpcClient *xmlrpc.Client
}

func NewClient(siteUrl string) (*Client, error) {

	if _, err := url.Parse(siteUrl); err != nil {
		return nil, err
	}

	client, _ := xmlrpc.NewClient(siteUrl+"/xmlrpc.cgi", nil)
	return &Client{bzUrl: siteUrl, rpcClient: client}, nil
}

func (c *Client) Version() (string, error) {
	ret := struct {
		Version string `xmlrpc:"version"`
	}{}

	if err := c.rpcClient.Call("Bugzilla.version", nil, &ret); err != nil {
		return "", err
	}

	return ret.Version, nil
}

func (c *Client) Search(attrs map[string]interface{}) ([]Bug, error) {
	ret := struct {
		Bugs []Bug `xmlrpc:"bugs"`
	}{}
	if err := c.rpcClient.Call("Bug.search", attrs, &ret); err != nil {
		return []Bug{}, err
	}
	return ret.Bugs, nil
}

func (c *Client) Query() *Query {
	return NewQuery(c)
}
