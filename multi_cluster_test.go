package esmc_test

import (
	"encoding/json"
	"github.com/bernerdschaefer/esmc"
	es "github.com/peterbourgon/elasticsearch"
	"reflect"
	"testing"
	"time"
)

func TestMultiClusterSearch(t *testing.T) {
	expected := es.SearchResponse{Took: 134}

	cluster1 := newServer(func(e *json.Encoder) {
		e.Encode(expected)
	})
	defer cluster1.Close()
	cluster2 := newServer(func(e *json.Encoder) {
		time.Sleep(10 * time.Millisecond)
		e.Encode(es.SearchResponse{Took: 1340})
	})
	defer cluster2.Close()

	mc := esmc.MultiCluster{
		newCluster(clusterSpec(cluster1.URL, "cluster1", "on")),
		newCluster(clusterSpec(cluster2.URL, "cluster2", "on")),
	}
	defer func() {
		for _, c := range mc {
			c.Shutdown()
		}
	}()

	got, err := mc.Search(es.SearchRequest{
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

func TestMultiClusterSearchDarkMode(t *testing.T) {
	expected := es.SearchResponse{Took: 134}

	cluster1 := newServer(func(e *json.Encoder) {
		time.Sleep(10 * time.Millisecond)
		e.Encode(expected)
	})
	defer cluster1.Close()
	cluster2 := newServer(func(e *json.Encoder) {
		e.Encode(es.SearchResponse{Took: 1340})
	})
	defer cluster2.Close()

	mc := esmc.MultiCluster{
		newCluster(clusterSpec(cluster1.URL, "cluster1", "on")),
		newCluster(clusterSpec(cluster2.URL, "cluster2", "dark")),
	}
	defer func() {
		for _, c := range mc {
			c.Shutdown()
		}
	}()

	got, err := mc.Search(es.SearchRequest{
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
