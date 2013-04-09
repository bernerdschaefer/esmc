package esmc_test

import (
	"encoding/json"
	"github.com/bernerdschaefer/esmc"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func newCluster(spec string) *esmc.Cluster {
	config, err := esmc.NewConfig(spec)
	if err != nil {
		panic(err)
	}
	return esmc.NewCluster(config)
}

func clusterSpec(endpoint string, name string, mode string) string {
	uri, _ := url.Parse(endpoint)
	uri.Scheme = "es"
	uri.RawQuery = url.Values{
		"name":          {name},
		"mode":          {mode},
		"ping_interval": {"-1s"},
	}.Encode()

	return uri.String()
}

func newServer(f func(*json.Encoder)) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f(json.NewEncoder(w))
		}),
	)
}
