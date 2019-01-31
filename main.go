package main

import (
	"flag"
	"perf-tool/pkg/perf"
)

var amount = flag.Int("t", 100, "Request amount")
var url = flag.String("u", "", "Url to call")
var concurrency = flag.Int("c", 1, "Number of requests to run at the same time")

func main() {
	flag.Parse()

	perf.RunPerfTest(*amount, *url, *concurrency)
}
