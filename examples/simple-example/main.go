package main

import (
	_ "embed"
	inseki "github.com/ForkBench/Inseki-Core"
)

//go:embed config.json
var configJson string

func main() {

	// ----------------------------- Read the configuration -----------------------------
	err, config := inseki.ReadEmbedConfigFile(configJson)
	if err != nil {
		panic(err)
	}

	// Check if the folder exists
	err = inseki.CheckIfConfigFolderExists(config)
	if err != nil {
		panic(err)
	}

	err, insekiIgnore := inseki.ReadInsekiIgnore(config)
	if err != nil {
		panic(err)
	}

	err, val := inseki.Process("~/Documents", config, insekiIgnore)
	if err != nil {
		return
	}

	// Print the results
	for _, response := range val {
		println(response.String())
	}
}
