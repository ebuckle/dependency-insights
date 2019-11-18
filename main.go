package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ebuckle/dependency-insights/insights"
	"github.com/pkg/browser"
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

	_ = ioutil.WriteFile("output.json", print, 0644)
	root, _ := os.Getwd()
	browser.OpenURL("file:///" + root + "/output.json")
	os.Exit(0)
}
