package esmc_test

import (
	es "github.com/peterbourgon/elasticsearch"
	"encoding/json"
	"reflect"
	"testing"
)

func TestClusterSearch(t *testing.T) {
	expected := es.SearchResponse{Took: 134}
	server := newServer(func(e *json.Encoder) {
		e.Encode(expected)
	})
	defer server.Close()

	cluster := newCluster(clusterSpec(server.URL, "cluster1", "on"))
	defer cluster.Shutdown()

	got, err := cluster.Search(es.SearchRequest{
		es.SearchParams{},
		es.MatchAllQuery(),
	})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatal("got != epxected")
	}
}
