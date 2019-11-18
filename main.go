package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ebuckle/dependency-insights/insights"
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

	print, err := json.MarshalIndent(response, "", "\t")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println(string(print))
	os.Exit(0)
}
