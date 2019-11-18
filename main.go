package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	projectPath := args[1]
	projectLanguage := args[2]

	response, err := insights.ProduceInsights(projectLanguage, projectPath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	print, err := json.MarshalIndent(response, "", "		")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println(print)
	os.Exit(0)
}
