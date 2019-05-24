package ripe

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

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

func (c *Controller) Search(search string) (payload answer, err error) {
	// Prepare URL
	searchSplitted := strings.Split(search, " ")
	url := *baseURL
	queryValues := url.Query()
	queryValues.Set("q", "("+strings.Join(searchSplitted, " AND ")+")")
	url.RawQuery = queryValues.Encode()
	// Start request
	resp, err := c.client.Get(url.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Decode
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&payload)
	return
}
