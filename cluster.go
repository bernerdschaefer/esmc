package esmc

import (
	es "github.com/peterbourgon/elasticsearch"
	"time"
)

type Cluster struct {
	Name    string
	Mode    Mode
	cluster *es.Cluster
}

// Connects to a cluster using the provided Config.
func NewCluster(config Config) *Cluster {
	return &Cluster{
		Name: config.Name(),
		Mode: config.Mode(),
		cluster: es.NewCluster(
			config.Endpoints,
			config.PingInterval(),
			config.PingTimeout(),
		),
	}
}

// Wraps es.Cluster.Shutdown
func (c *Cluster) Shutdown() {
	c.cluster.Shutdown()
}

func (c *Cluster) labels(requestType string, ok bool) map[string]string {
	labels := map[string]string{
		"cluster":      c.Name,
		"cluster_mode": c.Mode.String(),
		"request_type": requestType,
	}

	if ok {
		labels["outcome"] = "success"
	} else {
		labels["outcome"] = "failed"
	}

	return labels
}

// Wraps es.Cluster.Execute with instrumentation.
func (c *Cluster) Execute(f es.Fireable, response interface{}) (err error) {
	defer func(began time.Time) {
		labels := c.labels("execute", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Execute(f, response)
}

// Wraps es.Cluster.Index with instrumentation.
func (c *Cluster) Index(r es.IndexRequest) (_ es.IndexResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("index", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Index(r)
}

// Wraps es.Cluster.Update with instrumentation.
func (c *Cluster) Update(r es.UpdateRequest) (_ es.IndexResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("update", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Update(r)
}

// Wraps es.Cluster.Delete with instrumentation.
func (c *Cluster) Delete(r es.DeleteRequest) (_ es.IndexResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("delete", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Delete(r)
}

// Wraps es.Cluster.Create with instrumentation.
func (c *Cluster) Create(r es.CreateRequest) (_ es.IndexResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("create", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Create(r)
}

// Wraps es.Cluster.Bulk with instrumentation.
func (c *Cluster) Bulk(r es.BulkRequest) (_ es.BulkResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("bulk", err == nil)
		updateRequestMetrics(labels, time.Since(began))
	}(time.Now())

	return c.cluster.Bulk(r)
}

// Wraps es.Cluster.MultiSearch with instrumentation. Implements es.MultiSearcher.
func (c *Cluster) MultiSearch(r es.MultiSearchRequest) (resp es.MultiSearchResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("multi_search", err == nil)

		updateRequestMetrics(labels, time.Since(began))
		if err == nil {
			for _, response := range resp.Responses {
				updateReportedMetrics(labels, response.Took)
			}
		}
	}(time.Now())

	return c.cluster.MultiSearch(r)
}

// Wraps es.Cluster.Search with instrumentation. Implements es.Searcher.
func (c *Cluster) Search(r es.SearchRequest) (response es.SearchResponse, err error) {
	defer func(began time.Time) {
		labels := c.labels("search", err == nil)

		if err == nil && response.TimedOut {
			labels["outcome"] = "timeout"
		}

		updateRequestMetrics(labels, time.Since(began))

		if err == nil {
			updateReportedMetrics(labels, response.Took)
		}
	}(time.Now())

	return c.cluster.Search(r)
}
