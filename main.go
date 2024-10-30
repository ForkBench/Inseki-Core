package main

import (
	_ "embed"
	"github.com/ForkBench/Inseki-Core/tools"
)

//go:embed config.json
var configJson string

func main() {

	// ----------------------------- Read the configuration -----------------------------
	err, config := tools.ReadEmbedConfigFile(configJson)
	if err != nil {
		panic(err)
	}

	// Check if the folder exists
	err = tools.CheckIfConfigFolderExists(config)
	if err != nil {
		panic(err)
	}

	err, insekiIgnore := tools.ReadInsekiIgnore(config)
	if err != nil {
		panic(err)
	}

	err, val := tools.Process("~/Documents", config, insekiIgnore)
	if err != nil {
		return
	}

	// Print the results
	for _, response := range val {
		println(response.String())
	}
}
