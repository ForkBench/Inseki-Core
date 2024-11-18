package main

import (
	_ "embed"
	inseki "github.com/ForkBench/Inseki-Core"
)

func main() {
	structures := inseki.ImportStructures("structures")
	inseki.ResearchAndAnalyse("~/Documents/", structures)
}
