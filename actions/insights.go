package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ebuckle/dependency-insights/insights"
	"github.com/pkg/browser"
	"github.com/urfave/cli"
)

// InsightsLocalProject produces dependency insights for locally saved projects
func InsightsLocalProject(c *cli.Context) {
	projectPath := strings.TrimSpace(strings.ToLower(c.String("path")))
	projectLanguage := strings.TrimSpace(strings.ToLower(c.String("language")))

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
