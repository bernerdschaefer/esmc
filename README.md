# esmc

`esmc` is a wrapper for
[elasticsearch](https://github.com/peterbourgon/elasticsearch).

The primary function of `esmc` is to provide support for interacting with
multiple elasticsearch clusters, and secondarily to provide per-request
instrumentation (via
[prometheus](https://github.com/prometheus/client_golang)).

```go
mc := esmc.MustNewMultiCluster("es://:9200?name=one&mode=on es://:9201?name=two&mode=dark")

// Can be used as an es.Searcher or es.MultiSearcher
mc.Search(&es.SearchRequest{
	Params: es.SearchParams{
		Indices: []string{"twitter"},
		Types:   []string{"tweet"},
	},
	Query:   es.MatchAllQuery(),
})

// Can be iterated to access invidivual clusters, which implement the full
// es.Cluster API.
for _, c := range mc {
	c.Index(&es.IndexRequest{
		es.IndexParams{Index: "twitter", Type: "tweet", Id: "1"},
		map[string]interface{}{},
	})
}
```

## Multiple Clusters

`esmc` provides support for multiple elasticsearch clusters, with the goal of
simplifying failover and testing.

### Search Behavior

The default behavior for `Search` and `MultiSearch` is the following:

* the request is sent to *all* clusters in "on" or "dark" mode.
* the first response from an "on" cluster is returned to the user.
* all other responses are discarded (after reporting metrics).

## Instrumentation

The following metrics are exposed through the prometheus client's registry:

* `elasticsearch_client_requests`: A counter of the total number of requests to
  an ES cluster
* `elasticsearch_client_request_total_duration_nanoseconds`: The total amount
  of time spent executing requests (nanoseconds)
* `elasticsearch_client_request_durations_nanoseconds`: The amounts of time
  spent executing requests (nanoseconds)
* `elasticsearch_client_reported_request_total_duration_nanoseconds`: The total
  amount of time spent executing requests as reported by elasticsearch
(nanoseconds)
* `elasticsearch_client_reported_request_durations_nanoseconds`: The amounts of
  time spent executing requests as reported by elasticsearch (nanoseconds)

Each metric is labelled with:

* `cluster`: the cluster name (from the cluster's config)
* `cluster_mode`: the cluster's mode (on, off, dark)
* `request_type`: search, multi_search, index, create, update, delete, bulk, execute
* `outcome`: success, failure
