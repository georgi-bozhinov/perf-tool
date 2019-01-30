package main

import (
	"flag"
	"perf-tool/pkg/perf"
)

var amount = flag.Int("t", 100, "Request amount")
var url = flag.String("u", "", "Url to call")

func main() {
	flag.Parse()

	perf.PerfTest(*amount, *url)
}
