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
			GCG1OldGenerationCount                     uint64 `json:"gc.G1-Old-Generation.count"`
			GGG1OldGenerationTime                      uint64 `json:"gc.G1-Old-Generation.time"`
			GCG1YoungGenerationCount                   uint64 `json:"gc.G1-Young-Generation.count"`
			GCG1YoungGenerationTime                    uint64 `json:"gc.G1-Young-Generation.time"`
			MemoryTotalMax                             uint64 `json:"memory.total.max"`
			MemoryTotalUsed                            uint64 `json:"memory.total.used"`
			MemoryHeapUsed                             uint64 `json:"memory.heap.used"`
			MemoryNonHeapUsed                          uint64 `json:"memory.non-heap.used"`
			MemoryPoolsCodeHeapNonNmethodsUsed         uint64 `json:"memory.pools.CodeHeap-'non-nmethods'.used"`
			MemoryPoolsCodeHeapNonProfiledNmethodsUsed uint64 `json:"memory.pools.CodeHeap-'non-profiled-nmethods'.used"`
			MemoryPoolsCodeHeapProfiledNmethodsUsed    uint64 `json:"memory.pools.CodeHeap-'profiled-nmethods'.used"`
			MemoryPoolsCompressedClassSpaceUsed        uint64 `json:"memory.pools.Compressed-Class-Space.used"`
			MemoryPoolsG1EdenSpaceUsed                 uint64 `json:"memory.pools.G1-Eden-Space.used"`
			MemoryPoolsG1OldGenUsed                    uint64 `json:"memory.pools.G1-Old-Gen.used"`
			MemoryPoolsG1SurvivorSpaceUsed             uint64 `json:"memory.pools.G1-Survivor-Space.used"`
			MemoryPoolsMetaspaceUsed                   uint64 `json:"memory.pools.Metaspace.used"`
			ThreadsBlockedCount                        uint64 `json:"threads.blocked.count"`
			ThreadsCount                               uint64 `json:"threads.count"`
			ThreadsDaemonCount                         uint64 `json:"threads.daemon.count"`
			ThreadsDeadlockCount                       uint64 `json:"threads.deadlock.count"`
			ThreadsNewCount                            uint64 `json:"threads.new.count"`
			ThreadsRunnableCount                       uint64 `json:"threads.runnable.count"`
			ThreadsTerminatedCount                     uint64 `json:"threads.terminated.count"`
			ThreadsTimedWaitingCount                   uint64 `json:"threads.timed_waiting.count"`
			ThreadsWaitingCount                        uint64 `json:"threads.waiting.count"`
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
		"gc_count.old":                         apiResp.Metrics.JVM.GCG1OldGenerationCount,
		"gc_time.old":                          apiResp.Metrics.JVM.GGG1OldGenerationTime,
		"gc_count.young":                       apiResp.Metrics.JVM.GCG1YoungGenerationCount,
		"gc_time.young":                        apiResp.Metrics.JVM.GCG1YoungGenerationTime,
		"memory_used.total_max":                apiResp.Metrics.JVM.MemoryTotalMax,
		"memory_used.total":                    apiResp.Metrics.JVM.MemoryTotalUsed,
		"memory_used.heap":                     apiResp.Metrics.JVM.MemoryHeapUsed,
		"memory_used.non_heap":                 apiResp.Metrics.JVM.MemoryNonHeapUsed,
		"memory_space.code_heap_non_n_methods": apiResp.Metrics.JVM.MemoryPoolsCodeHeapNonNmethodsUsed,
		"memory_space.code_heap_non_profiled_n_methods": apiResp.Metrics.JVM.MemoryPoolsCodeHeapNonProfiledNmethodsUsed,
		"memory_space.code_heap_profiled_n_methods":     apiResp.Metrics.JVM.MemoryPoolsCodeHeapProfiledNmethodsUsed,
		"memory_space.compressed_class":                 apiResp.Metrics.JVM.MemoryPoolsCompressedClassSpaceUsed,
		"memory_space.eden":                             apiResp.Metrics.JVM.MemoryPoolsG1EdenSpaceUsed,
		"memory_space.old":                              apiResp.Metrics.JVM.MemoryPoolsG1OldGenUsed,
		"memory_space.survivor":                         apiResp.Metrics.JVM.MemoryPoolsG1SurvivorSpaceUsed,
		"memory_space.metaspace":                        apiResp.Metrics.JVM.MemoryPoolsMetaspaceUsed,
		"thread_count.blocked":                          apiResp.Metrics.JVM.ThreadsBlockedCount,
		"thread_count.all":                              apiResp.Metrics.JVM.ThreadsCount,
		"thread_count.daemon":                           apiResp.Metrics.JVM.ThreadsDaemonCount,
		"thread_count.deadlock":                         apiResp.Metrics.JVM.ThreadsDeadlockCount,
		"thread_count.new":                              apiResp.Metrics.JVM.ThreadsNewCount,
		"thread_count.runnable":                         apiResp.Metrics.JVM.ThreadsRunnableCount,
		"thread_count.terminated":                       apiResp.Metrics.JVM.ThreadsTerminatedCount,
		"thread_count.timedWaiting":                     apiResp.Metrics.JVM.ThreadsTimedWaitingCount,
		"thread_count.waiting":                          apiResp.Metrics.JVM.ThreadsWaitingCount,
	}, nil
}

// GraphDefinition interface for mackerelplugin
func (sj SolrJVMPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"gc_count": {
			Label: "JVM GC Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "old", Label: "Old", AbsoluteName: true},
				{Name: "young", Label: "Young", AbsoluteName: true},
			},
		},
		"gc_time": {
			Label: "JVM GC Time",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "old", Label: "Old", AbsoluteName: true},
				{Name: "young", Label: "Young", AbsoluteName: true},
			},
		},
		"memory_used": {
			Label: "JVM Memory Used",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "total_max", Label: "Total Max", AbsoluteName: true},
				{Name: "total", Label: "Total", AbsoluteName: true},
				{Name: "heap", Label: "Heap", AbsoluteName: true, Stacked: true},
				{Name: "non_heap", Label: "Non Heap", AbsoluteName: true, Stacked: true},
			},
		},
		"memory_space": {
			Label: "JVM Memory Space",
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: "code_heap_non_n_methods", Label: "Code Heap Non N Methods", AbsoluteName: true, Stacked: true},
				{Name: "code_heap_non_profiled_n_methods", Label: "Code Heap Non Profiled N Methods", AbsoluteName: true, Stacked: true},
				{Name: "code_heap_profiled_n_methods", Label: "Code Heap Profiled N Methods", AbsoluteName: true, Stacked: true},
				{Name: "compressed_class", Label: "Compressed Class", AbsoluteName: true, Stacked: true},
				{Name: "eden", Label: "Eden", AbsoluteName: true, Stacked: true},
				{Name: "old", Label: "Old", AbsoluteName: true, Stacked: true},
				{Name: "survivor", Label: "Survivor", AbsoluteName: true, Stacked: true},
				{Name: "metaspace", Label: "Metaspace", AbsoluteName: true, Stacked: true},
			},
		},
		"thread_count": {
			Label: "JVM Thread Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "blocked", Label: "Blocked", AbsoluteName: true, Stacked: true},
				{Name: "all", Label: "All", AbsoluteName: true},
				{Name: "daemon", Label: "Daemon", AbsoluteName: true},
				{Name: "deadlock", Label: "Deadlock", AbsoluteName: true, Stacked: true},
				{Name: "new", Label: "New", AbsoluteName: true, Stacked: true},
				{Name: "runnable", Label: "Runnable", AbsoluteName: true, Stacked: true},
				{Name: "terminated", Label: "Terminated", AbsoluteName: true, Stacked: true},
				{Name: "timedWaiting", Label: "TimedWaiting", AbsoluteName: true, Stacked: true},
				{Name: "waiting", Label: "Waiting", AbsoluteName: true, Stacked: true},
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
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	sj := SolrJVMPlugin{URL: *optURL, Prefix: "solrjvm"}
	p := mp.NewMackerelPlugin(sj)
	p.Tempfile = *optTempfile
	p.Run()
}
