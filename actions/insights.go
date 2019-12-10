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
	"github.com/ebuckle/dependency-insights/report"
	"github.com/urfave/cli"
	"gopkg.in/src-d/go-git.v4"
)

// StartInsights calls the appropriate setup function for the project type, then starts the insights process
func StartInsights(c *cli.Context) {
	tempFolder := setupTemp()
	var projectPath string
	projectLanguage := strings.TrimSpace(strings.ToLower(c.String("language")))
	switch c.Command.FullName() {
	case "local":
		projectPath = strings.TrimSpace(strings.ToLower(c.String("path")))
	case "docker":
		projectPath = setupDockerProject(c, tempFolder)
	case "git":
		projectPath = setupGitProject(c, tempFolder)
	}

	response, err := insights.ProduceInsights(projectLanguage, projectPath)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printResults(*response)

	teardownTemp(tempFolder)
}

func setupDockerProject(c *cli.Context, tempFolder string) string {
	containerID := strings.TrimSpace(strings.ToLower(c.String("conid")))

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

	return projectPath
}

func setupGitProject(c *cli.Context, tempFolder string) string {
	gitURL := strings.TrimSpace(strings.ToLower(c.String("url")))
	projectLanguage := strings.TrimSpace(strings.ToLower(c.String("language")))

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

	return projectPath
}

func printResults(response map[string]interface{}) {
	file, err := json.MarshalIndent(response, "", "\t")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile("output.json", file, 0644)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	report.ProduceReport(response)

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
