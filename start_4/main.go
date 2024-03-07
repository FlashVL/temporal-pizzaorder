package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	b := prometheus.LinearBuckets(310, 10, 80)

	fmt.Println(b)
}
