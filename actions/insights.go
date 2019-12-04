package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

	printResults(*response)
}

// InsightsDockerProject produces insights on the contents of a docker container
func InsightsDockerProject(c *cli.Context) {
	containerID := strings.TrimSpace(strings.ToLower(c.String("conid")))
	projectLanguage := strings.TrimSpace(strings.ToLower(c.String("language")))
	tempFolder := "temp-folder-delete-me"

	os.Mkdir(tempFolder, 0700)

	dockerCommand := exec.Command("docker", "container", "export", containerID, "-o", "output.tar")
	dockerCommand.Dir = tempFolder
	err := dockerCommand.Run()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tarCommand := exec.Command("tar", "-xvf", "output.tar")
	tarCommand.Dir = tempFolder
	tarCommand.Run()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	projectPath := tempFolder + "/app"

	response, err := insights.ProduceInsights(projectLanguage, projectPath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.RemoveAll(tempFolder)

	printResults(*response)
}

func printResults(response map[string]interface{}) {
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
