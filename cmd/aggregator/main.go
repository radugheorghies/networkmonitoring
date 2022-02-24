package main

import "networkmonitoring/pkg/aggregator"

func main() {
	a := aggregator.New()
	a.Run()
}
