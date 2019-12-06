package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ebuckle/dependency-insights/insights"
	"github.com/pkg/browser"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
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
	tempFolder := setupTemp()

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

	teardownTemp(tempFolder)

	printResults(*response)
}

// InsightsGitProject produces insights on the contents of the given git repository
func InsightsGitProject(c *cli.Context) {
	gitURL := strings.TrimSpace(strings.ToLower(c.String("url")))
	projectLanguage := strings.TrimSpace(strings.ToLower(c.String("language")))
	tempFolder := setupTemp()

	_, err := git.PlainClone(tempFolder, false, &git.CloneOptions{
		URL:      gitURL,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	projectPath := tempFolder

	err = installDependencies(projectPath, projectLanguage)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	response, err := insights.ProduceInsights(projectLanguage, projectPath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	teardownTemp(tempFolder)

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

func setupTemp() string {
	tempFolder := "temp-folder-delete-me"
	os.Mkdir(tempFolder, 0700)
	return tempFolder
}

func teardownTemp(tempFolder string) {
	os.RemoveAll(tempFolder)
}

func installDependencies(projectPath string, projectLanguage string) error {
	var err error
	switch projectLanguage {
	case "nodejs":
		npmCommand := exec.Command("npm", "install")
		npmCommand.Dir = projectPath
		err = npmCommand.Run()
	/*
		case "go":
			goCommand := exec.Command("dep", "ensure", "-v")
			goCommand.Dir = projectPath
			err = goCommand.Run()
	*/
	default:
		err = errors.New("language not currently supported")
	}
	return err
}
