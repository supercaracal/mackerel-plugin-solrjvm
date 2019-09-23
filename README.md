mackerel-plugin-solrjvm
=====================

Apache Solr JVM metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-solrjvm [-url=<url>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.solrjvm]
command = "/path/to/mackerel-plugin-solrjvm"
```

## See also
* [Apache Solr Metrics Reporting](https://lucene.apache.org/solr/guide/8_1/metrics-reporting.html)
* [go-mackerel-pluginを利用してカスタムメトリックプラグインを作成する](https://mackerel.io/ja/docs/entry/advanced/go-mackerel-plugin)
* [mkr plugin installに対応したプラグインを作成する](https://mackerel.io/ja/docs/entry/advanced/make-plugin-corresponding-to-installer)

## About `-XX:+PerfDisableSharedMem` option

https://cwiki.apache.org/confluence/display/solr/ShawnHeisey#ShawnHeisey-Currentexperiments

> The PerfDisableSharedMem option is there because of something that is called [the four month bug](https://www.evanjones.ca/jvm-mmap-pause.html).
