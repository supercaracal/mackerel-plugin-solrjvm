package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/solr/admin/metrics":
		fmt.Fprintf(w, fetchJSON("test/stats.json"))
	default:
		fmt.Fprintf(w, "{}")
	}
})

func fetchJSON(path string) string {
	json, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(json)
}

func TestGraphDefinition(t *testing.T) {
	solrJVM := SolrJVMPlugin{}
	actual := solrJVM.GraphDefinition()

	if actual["gc_count"].Label != "JVM GC Count" {
		t.Fatalf("Expected: %s, Actual: %s", "JVM GC Count", actual["gc_count"].Label)
	}
}

func TestFetchMetrics(t *testing.T) {
	testServer := httptest.NewServer(testServerHandler)
	defer testServer.Close()

	solrJVM := SolrJVMPlugin{URL: testServer.URL}
	actual, err := solrJVM.FetchMetrics()

	if err != nil {
		t.Fatal(err)
	}

	var expected uint64
	expected = 12
	if actual["gc_count.young"] != expected {
		t.Fatalf("Expected: %d, Actual: %d", expected, actual["gc_count.young"])
	}
}

func TestMetricKeyPrefix(t *testing.T) {
	solrJVM := SolrJVMPlugin{Prefix: "solrjvm"}
	actual := solrJVM.MetricKeyPrefix()

	if actual != "solrjvm" {
		t.Fatalf("Expected: %s, Actual: %s", "solrjvm", actual)
	}
}
