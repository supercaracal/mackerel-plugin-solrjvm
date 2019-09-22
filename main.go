package main

import (
	"encoding/json"
	"flag"
	"net/http"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

// SolrJVMPlugin for mackerelplugin
type SolrJVMPlugin struct {
	Prefix string
	URL    string
}

// SolrMetricAPIResponse for JSON parsing
type SolrMetricAPIResponse struct {
	Metrics struct {
		JVM struct {
			GCG1OldGenerationCount                     int `json:"gc.G1-Old-Generation.count"`
			GGG1OldGenerationTime                      int `json:"gc.G1-Old-Generation.time"`
			GCG1YoungGenerationCount                   int `json:"gc.G1-Young-Generation.count"`
			GCG1YoungGenerationTime                    int `json:"gc.G1-Young-Generation.time"`
			MemoryTotalMax                             int `json:"memory.total.max"`
			MemoryTotalUsed                            int `json:"memory.total.used"`
			MemoryHeapUsed                             int `json:"memory.heap.used"`
			MemoryNonHeapUsed                          int `json:"memory.non-heap.used"`
			MemoryPoolsCodeHeapNonNmethodsUsed         int `json:"memory.pools.CodeHeap-'non-nmethods'.used"`
			MemoryPoolsCodeHeapNonProfiledNmethodsUsed int `json:"memory.pools.CodeHeap-'non-profiled-nmethods'.used"`
			MemoryPoolsCodeHeapProfiledNmethodsUsed    int `json:"memory.pools.CodeHeap-'profiled-nmethods'.used"`
			MemoryPoolsCompressedClassSpaceUsed        int `json:"memory.pools.Compressed-Class-Space.used"`
			MemoryPoolsG1EdenSpaceUsed                 int `json:"memory.pools.G1-Eden-Space.used"`
			MemoryPoolsG1OldGenUsed                    int `json:"memory.pools.G1-Old-Gen.used"`
			MemoryPoolsG1SurvivorSpaceUsed             int `json:"memory.pools.G1-Survivor-Space.used"`
			MemoryPoolsMetaspaceUsed                   int `json:"memory.pools.Metaspace.used"`
			ThreadsBlockedCount                        int `json:"threads.blocked.count"`
			ThreadsCount                               int `json:"threads.count"`
			ThreadsDaemonCount                         int `json:"threads.daemon.count"`
			ThreadsDeadlockCount                       int `json:"threads.deadlock.count"`
			ThreadsNewCount                            int `json:"threads.new.count"`
			ThreadsRunnableCount                       int `json:"threads.runnable.count"`
			ThreadsTerminatedCount                     int `json:"threads.terminated.count"`
			ThreadsTimedWaitingCount                   int `json:"threads.timed_waiting.count"`
			ThreadsWaitingCount                        int `json:"threads.waiting.count"`
		} `json:"solr.jvm"`
	} `json:"metrics"`
}

var logger = logging.GetLogger("metrics.plugin.solrjvm")

func fetchSolrJVMMetrics(baseURL string) (apiResp SolrMetricAPIResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/solr/admin/metrics?group=jvm", nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "mackerel-plugin-solrjvm")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return
	}

	return
}

// FetchMetrics interface for mackerelplugin
func (sj SolrJVMPlugin) FetchMetrics() (map[string]interface{}, error) {
	apiResp, err := fetchSolrJVMMetrics(sj.URL)
	if err != nil {
		logger.Errorf("Failed to fetch Solr JVM metrics. %s", err)
		return nil, err
	}

	return map[string]interface{}{
		"gc_count.old":                                 apiResp.Metrics.JVM.GCG1OldGenerationCount,
		"gc_time.old":                                  apiResp.Metrics.JVM.GGG1OldGenerationTime,
		"gc_count.young":                               apiResp.Metrics.JVM.GCG1YoungGenerationCount,
		"gc_time.young":                                apiResp.Metrics.JVM.GCG1YoungGenerationTime,
		"memory_used.total_max":                        apiResp.Metrics.JVM.MemoryTotalMax,
		"memory_used.total":                            apiResp.Metrics.JVM.MemoryTotalUsed,
		"memory_used.heap":                             apiResp.Metrics.JVM.MemoryHeapUsed,
		"memory_used.non_heap":                         apiResp.Metrics.JVM.MemoryNonHeapUsed,
		"memory_used.code_heap_non_n_methods":          apiResp.Metrics.JVM.MemoryPoolsCodeHeapNonNmethodsUsed,
		"memory_used.code_heap_non_profiled_n_methods": apiResp.Metrics.JVM.MemoryPoolsCodeHeapNonProfiledNmethodsUsed,
		"memory_used.code_heap_profiled_n_methods":     apiResp.Metrics.JVM.MemoryPoolsCodeHeapProfiledNmethodsUsed,
		"memory_used.compressed_class":                 apiResp.Metrics.JVM.MemoryPoolsCompressedClassSpaceUsed,
		"memory_used.eden":                             apiResp.Metrics.JVM.MemoryPoolsG1EdenSpaceUsed,
		"memory_used.old":                              apiResp.Metrics.JVM.MemoryPoolsG1OldGenUsed,
		"memory_used.survivor":                         apiResp.Metrics.JVM.MemoryPoolsG1SurvivorSpaceUsed,
		"memory_used.metaspace":                        apiResp.Metrics.JVM.MemoryPoolsMetaspaceUsed,
		"thread_count.blocked":                         apiResp.Metrics.JVM.ThreadsBlockedCount,
		"thread_count.all":                             apiResp.Metrics.JVM.ThreadsCount,
		"thread_count.daemon":                          apiResp.Metrics.JVM.ThreadsDaemonCount,
		"thread_count.deadlock":                        apiResp.Metrics.JVM.ThreadsDeadlockCount,
		"thread_count.new":                             apiResp.Metrics.JVM.ThreadsNewCount,
		"thread_count.runnable":                        apiResp.Metrics.JVM.ThreadsRunnableCount,
		"thread_count.terminated":                      apiResp.Metrics.JVM.ThreadsTerminatedCount,
		"thread_count.timedWaiting":                    apiResp.Metrics.JVM.ThreadsTimedWaitingCount,
		"thread_count.waiting":                         apiResp.Metrics.JVM.ThreadsWaitingCount,
	}, nil
}

// GraphDefinition interface for mackerelplugin
func (sj SolrJVMPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"gc_count": {
			Label: "JVM GC Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "old", Label: "Old"},
				{Name: "young", Label: "Young"},
			},
		},
		"gc_time": {
			Label: "JVM GC Time",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "old", Label: "Old"},
				{Name: "young", Label: "Young"},
			},
		},
		"memory_used": {
			Label: "JVM Memory Used",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "total_max", Label: "Total Max"},
				{Name: "total", Label: "Total"},
				{Name: "heap", Label: "Heap", Stacked: true},
				{Name: "non_heap", Label: "Non Heap", Stacked: true},
				{Name: "code_heap_non_n_methods", Label: "Code Heap Non N Methods", Stacked: true},
				{Name: "code_heap_non_profiled_n_methods", Label: "Code Heap Non Profiled N Methods", Stacked: true},
				{Name: "code_heap_profiled_n_methods", Label: "Code Heap Profiled N Methods", Stacked: true},
				{Name: "compressed_class", Label: "Compressed Class", Stacked: true},
				{Name: "eden", Label: "Eden", Stacked: true},
				{Name: "old", Label: "Old", Stacked: true},
				{Name: "survivor", Label: "Survivor", Stacked: true},
				{Name: "metaspace", Label: "Metaspace", Stacked: true},
			},
		},
		"thread_count": {
			Label: "JVM Thread Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "blocked", Label: "Blocked"},
				{Name: "all", Label: "All"},
				{Name: "daemon", Label: "Daemon"},
				{Name: "deadlock", Label: "Deadlock"},
				{Name: "new", Label: "New"},
				{Name: "runnable", Label: "Runnable"},
				{Name: "terminated", Label: "Terminated"},
				{Name: "timedWaiting", Label: "TimedWaiting"},
				{Name: "waiting", Label: "Waiting"},
			},
		},
	}
}

// MetricKeyPrefix is implementation of Mackerel PluginWithPrefix interface
func (sj SolrJVMPlugin) MetricKeyPrefix() string {
	return sj.Prefix
}

// @see https://mackerel.io/ja/docs/entry/advanced/go-mackerel-plugin
// @see https://lucene.apache.org/solr/guide/8_1/metrics-reporting.html
func main() {
	optURL := flag.String("url", "http://127.0.0.1:8983", "Solr URL")
	optTempfile := flag.String("tempfile", "mackerel-plugin-solrjvm", "Temp file name")
	flag.Parse()

	sj := SolrJVMPlugin{URL: *optURL, Prefix: "solrjvm"}
	p := mp.NewMackerelPlugin(sj)
	p.Tempfile = *optTempfile
	p.Run()
}
