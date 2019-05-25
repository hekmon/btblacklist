package ripe

import (
	"net/http"
	"net/url"

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

func New() *Controller {
	return &Controller{
		client: cleanhttp.DefaultPooledClient(),
	}
}

type Controller struct {
	client *http.Client
}
