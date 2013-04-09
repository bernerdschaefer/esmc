package esmc

import (
	es "github.com/peterbourgon/elasticsearch"
	"log"
	"strings"
)

type MultiCluster []*Cluster

// Accepts a space-separated list of cluster URIs (see NewConfig), logging a
// fatal error if the MultiCluster cannot be configured.
func MustNewMultiCluster(specs string) MultiCluster {
	mc := MultiCluster{}
	on := 0

	for _, spec := range strings.Split(specs, " ") {
		config, err := NewConfig(spec)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if config.Mode() == OnMode {
			on++
		}

		mc = append(mc, NewCluster(config))
	}

	if len(mc) == 0 {
		log.Fatalf("no elasticsearch clusters defined")
	}

	if on == 0 {
		log.Fatalf("no elasticsearch clusters in on mode")
	}

	return mc
}

// Implements es.Searcher
func (mc MultiCluster) Search(r es.SearchRequest) (es.SearchResponse, error) {
	type pair struct {
		es.SearchResponse
		error
	}

	// This channel will potentially be slightly over-allocated, since only
	// clusters in OnMode will send back their responses.
	ch := make(chan pair, len(mc))

	for _, cluster := range mc {
		switch cluster.Mode {
		case OnMode:
			go func(s es.Searcher) {
				resp, err := s.Search(r)
				ch <- pair{resp, err}
			}(cluster)
		case DarkMode:
			go cluster.Search(r)
		}
	}

	resp := <-ch
	return resp.SearchResponse, resp.error
}

// Implements es.MultiSearcher
func (mc MultiCluster) MultiSearch(r es.MultiSearchRequest) (es.MultiSearchResponse, error) {
	type pair struct {
		es.MultiSearchResponse
		error
	}

	// This channel will potentially be slightly over-allocated, since only
	// clusters in OnMode will send back their responses.
	ch := make(chan pair, len(mc))

	for _, cluster := range mc {
		switch cluster.Mode {
		case OnMode:
			go func(s es.MultiSearcher) {
				resp, err := s.MultiSearch(r)
				ch <- pair{resp, err}
			}(cluster)
		case DarkMode:
			go cluster.MultiSearch(r)
		}
	}

	resp := <-ch
	return resp.MultiSearchResponse, resp.error
}
