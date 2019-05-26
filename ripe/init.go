package ripe

import (
	"net/http"
	"net/url"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

const (
	exampleURL = "https://apps.db.ripe.net/db-web-ui/api/rest/fulltextsearch/select.json?facet=true&format=xml&hl=true&q=(example+AND+example)&start=0&wt=json"
)

var (
	baseURL *url.URL
)

func init() {
	var err error
	if baseURL, err = url.Parse(exampleURL); err != nil {
		panic(err)
	}
}

// New will return an initialized and ready to use controller
func New(timeout time.Duration) (c *Client) {
	c = &Client{
		client: cleanhttp.DefaultPooledClient(),
	}
	c.client.Timeout = timeout
	return
}

// Client holds an http client used by methods to perform searchs
type Client struct {
	client *http.Client
}
