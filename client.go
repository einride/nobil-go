package nobil

import "net/http"

const (
	defaultSearchURL = "https://nobil.no/api/server/search.php"
	defaultDumpURL   = "https://nobil.no/api/server/datadump.php"
)

type Client struct {
	apiKey     string
	searchURL  string
	dumpURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
		searchURL:  defaultSearchURL,
		dumpURL:    defaultDumpURL,
	}
}

func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func (c *Client) SetSearchURL(searchURL string) {
	c.searchURL = searchURL
}

func (c *Client) SetDumpURL(dumpURL string) {
	c.dumpURL = dumpURL
}
