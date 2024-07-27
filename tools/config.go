package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	InsekiPath    string `json:"insekiPath"`
	StructurePath string `json:"structurePath"`
}

func ReadEmbedConfigFile(configJson string) Config {
	var config Config
	err := json.Unmarshal([]byte(configJson), &config)
	if err != nil {
		panic(err)
	}

	// TODO: Refactor
	config.InsekiPath = TranslateDir(config.InsekiPath)
	config.StructurePath = TranslateDir(config.StructurePath)

	return config
}

func CheckIfConfigFolderExists(config Config) {
	// Check if the folder InsekiPath exists
	if _, err := os.Stat(config.InsekiPath); os.IsNotExist(err) {
		fmt.Println("The folder does not exist")

		// Create the folder
		err := os.Mkdir(config.InsekiPath, 0755)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Folder %s created\n", config.InsekiPath)
	}

	// Check if the folder StructurePath exists
	if _, err := os.Stat(config.StructurePath); os.IsNotExist(err) {
		fmt.Println("The folder does not exist")

		// Create the folder
		err := os.Mkdir(config.StructurePath, 0755)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Folder %s created\n", config.StructurePath)
	}
}
