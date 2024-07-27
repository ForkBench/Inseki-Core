package main

import (
	_ "embed"
	"os"

	"inseki-core/tools"
)

//go:embed config.json
var configJson string

func main() {
	// Read the config file
	config := tools.ReadEmbedConfigFile(configJson)

	// Check if the folder exists
	tools.CheckIfConfigFolderExists(config)

	insekiignore := tools.ReadInsekiIgnore(config)

	structures := tools.ImportStructure(config, insekiignore)

	if len(structures) == 0 {
		println("No structures found")
		os.Exit(0)
	}

	patterns := tools.ExtractNames(structures, false)
	associations := tools.StringNodeToAssociation(patterns)

	tools.ExploreFolder("~/Documents/", insekiignore, tools.FilterWithPatternMap(&associations, func(filepath string, association tools.Association) {
		println(filepath)
	}))

}
