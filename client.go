package gorgojo

import (
	"encoding/xml"
	"fmt"
	"github.com/kolo/xmlrpc"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	bzUrl     string
	apiUrl    string
	rpcClient *xmlrpc.Client
}

func NewClient(siteUrl string) (*Client, error) {
	plugins := AllPlugins()

	// convert shortcuts to real bugzilla site urls
	var err error
	for _, plugin := range plugins {
		siteUrl, err = plugin.TransformSiteUrlHook(siteUrl)
		if err != nil {
			return nil, err
		}
	}

	if _, err := url.Parse(siteUrl); err != nil {
		return nil, err
	}

	// convert the site urls to APIs
	var apiUrl string
	for _, plugin := range plugins {
		apiUrl, err = plugin.TransformApiUrlHook(siteUrl)
		if err != nil {
			return nil, err
		}
	}

	client, _ := xmlrpc.NewClient(apiUrl+"/xmlrpc.cgi", nil)
	return &Client{bzUrl: siteUrl, apiUrl: apiUrl, rpcClient: client}, nil
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

func (c *Client) Search(attrs map[string][]interface{}) ([]Bug, error) {
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

func (c *Client) namedQueryToUrl(name string) string {
	// it is validated before
	u, _ := url.Parse(c.apiUrl)

	u.Path = "/buglist.cgi"

	q := u.Query()
	q.Set("cmdtype", "runnamed")
	q.Set("namedcmd", name)
	q.Set("ctype", "atom")
	u.RawQuery = q.Encode()

	return u.String()
}

type atomFeed struct {
	Entries []struct {
		Link struct {
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
	} `xml:"entry"`
}

// return the list of bug ids in the atom feed
func parseAtomFeed(r io.Reader) ([]int, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return []int{}, err
	}

	ids := make([]int, 0)
	// parse the atom feed, look for bug ids
	atomFeed := atomFeed{}
	xml.Unmarshal(content, &atomFeed)

	for _, entry := range atomFeed.Entries {
		url, err := url.Parse(entry.Link.Href)
		if err != nil {
			return []int{}, err
		}

		idStr := url.Query()["id"][0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (c *Client) fetchNamedQueryFromUrl(urlStr string) ([]Bug, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return []Bug{}, err
	}

	// keep the credentials across redirects
	// I have to admit I am not 100% sure this is necessary
	redirectPolicyFunc := func(req *http.Request, via []*http.Request) error {
		// keep the credentials
		req.URL.User = url.User
		if len(via) > 10 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	}

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	// get the atom feed
	resp, err := client.Get(urlStr)
	if err != nil {
		return []Bug{}, err
	}

	ids, err := parseAtomFeed(resp.Body)
	if err != nil {
		return []Bug{}, err
	}

	// collect all the bug ids here
	query := c.Query()

	for _, id := range ids {
		query.Field("id", id)
	}
	return query.Result()
}

// retrieves the bug list for the named server-side query
func (c *Client) FetchNamedQuery(name string) ([]Bug, error) {
	url := c.namedQueryToUrl(name)
	return c.fetchNamedQueryFromUrl(url)
}
