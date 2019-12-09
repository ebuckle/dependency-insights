package insights

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/src-d/go-license-detector.v3/licensedb"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/filer"
)

// ProduceInsights calls the appropriate crawling function for the provided language and then reports on licensing
func ProduceInsights(language string, projectPath string) (*map[string]interface{}, error) {
	_, err := os.Stat(projectPath)

	if err != nil {
		return nil, err
	}

	insightData := make(map[string]interface{})

	// Choose correct walk function depending on the language of the project
	switch language {
	case "nodejs":
		nodeCommand := exec.Command("npm", "ls", "--json")
		nodeCommand.Dir = projectPath
		rawOutput, _ := nodeCommand.Output()
		err = json.Unmarshal(rawOutput, &insightData)
		if err != nil {
			return nil, err
		}
		err = nodeWalk(projectPath, insightData["dependencies"].(map[string]interface{}))
	case "go":
		err = goWalk(projectPath, insightData)
	default:
		err := errors.New("language not recognised")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	performLicenseCheck(insightData["dependencies"].(map[string]interface{}))

	if language == "nodejs" {
		checkVulnerabilities(projectPath, insightData)
	}
	return &insightData, nil
}

/*
func nodeWalk(projectPath string, insightData map[string]interface{}) error {
	if _, err := os.Stat(projectPath + "/node_modules"); err == nil {
		files, err := ioutil.ReadDir(projectPath + "/node_modules")
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				path := projectPath + "/node_modules/" + file.Name()
				if _, err := os.Stat(path + "/package.json"); err == nil {
					jsonFile, err := os.Open(path + "/package.json")

					if err != nil {
						return err
					}

					byteValue, err := ioutil.ReadAll(jsonFile)

					if err != nil {
						return err
					}

					jsonFile.Close()

					var result map[string]interface{}
					err = json.Unmarshal([]byte(byteValue), &result)

					if err != nil {
						return err
					}

					newPackageData := make(map[string]interface{})
					transferNodeData(result, newPackageData, path)
					packageID := newPackageData["name"].(string) + "@" + newPackageData["version"].(string)
					if _, ok := insightData[packageID]; !ok {
						insightData[packageID] = newPackageData
					}

					if _, err := os.Stat(path + "/node_modules"); err == nil {
						dependenciesData := make(map[string]interface{})
						err := nodeWalk(path, dependenciesData)
						insightData[packageID].(map[string]interface{})["dependencies"] = dependenciesData

						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
*/

func nodeWalk(projectPath string, insightData map[string]interface{}) error {
	if _, err := os.Stat(projectPath + "/node_modules"); err == nil {
		files, err := ioutil.ReadDir(projectPath + "/node_modules")
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				path := projectPath + "/node_modules/" + file.Name()
				if _, err := os.Stat(path + "/package.json"); err == nil {
					jsonFile, err := os.Open(path + "/package.json")

					if err != nil {
						return err
					}

					byteValue, err := ioutil.ReadAll(jsonFile)

					if err != nil {
						return err
					}

					jsonFile.Close()

					var result map[string]interface{}
					err = json.Unmarshal([]byte(byteValue), &result)

					if err != nil {
						return err
					}

					newPackageData := make(map[string]interface{})
					transferNodeData(result, newPackageData, path)
					packageID := newPackageData["name"].(string) + "@" + newPackageData["version"].(string)

					mapData(insightData, newPackageData, packageID)

					if _, err := os.Stat(path + "/node_modules"); err == nil {
						err := nodeWalk(path, insightData)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// goWalk walks through installed go packaged to map dependencies
func goWalk(projectPath string, insightData map[string]interface{}) error {
	// Ensure go packages are installed
	if _, err := os.Stat(projectPath + "/vendor"); err == nil {
		sources, err := ioutil.ReadDir(projectPath + "/vendor")
		if err != nil {
			return err
		}
		// Loop through package sources
		for _, source := range sources {
			if source.IsDir() {
				sourcePath := projectPath + "/vendor/" + source.Name()
				authors, err := ioutil.ReadDir(sourcePath)

				if err != nil {
					return err
				}

				// Look through package authors
				for _, author := range authors {
					completeFlag := false
					authorPath := sourcePath + "/" + author.Name()
					dependencies, err := ioutil.ReadDir(authorPath)

					if err != nil {
						return err
					}

					// Check for non-directory in author path to see if we are at the package level
					for _, dependency := range dependencies {
						if !dependency.IsDir() {
							newPackageData := make(map[string]interface{})
							fullPath := authorPath
							newPackageData["path"] = fullPath
							depName := source.Name() + "/" + author.Name()
							insightData[depName] = newPackageData
							completeFlag = true
						}
					}

					if completeFlag {
						continue
					}

					// Loop through packages
					for _, dependency := range dependencies {
						newPackageData := make(map[string]interface{})
						fullPath := authorPath + "/" + dependency.Name()
						newPackageData["path"] = fullPath
						depName := source.Name() + "/" + author.Name() + "/" + dependency.Name()

						insightData[depName] = newPackageData
					}
				}
			}
		}
	}
	return nil
}

// transferNodeData takes existing package data from a package.json and loads it into a packageData struct
func transferNodeData(packageJSON map[string]interface{}, packageData map[string]interface{}, path string) {
	if str, ok := packageJSON["name"].(string); ok {
		packageData["name"] = str
	}
	if str, ok := packageJSON["version"].(string); ok {
		packageData["version"] = str
	}
	if str, ok := packageJSON["description"].(string); ok {
		packageData["description"] = str
	}
	if str, ok := packageJSON["license"].(string); ok {
		packageData["declaredLicenses"] = str
	} else {
		packageData["declaredLicenses"] = "No Declared License"
	}
	packageData["path"] = path
}

// performLicenseCheck takes an existing map of package data and performs a license check on each package
func performLicenseCheck(insightData map[string]interface{}) {
	for _, depI := range insightData {
		if dep, ok := depI.(map[string]interface{}); ok {
			if dep["path"] == nil {
				continue
			}
			filer, err := filer.FromDirectory(dep["path"].(string))

			if err != nil {
				dep["license-analysis"] = err.Error()
				continue
			}

			results, err := licensedb.Detect(filer)

			if err != nil {
				dep["license-analysis"] = err.Error()
			} else {
				dep["license-analysis"] = results
			}

			if dep["dependencies"] != nil {
				performLicenseCheck(dep["dependencies"].(map[string]interface{}))
			}
		}
	}
}

func checkVulnerabilities(projectPath string, insightData map[string]interface{}) error {
	npmAudit := new(map[string]interface{})
	command := exec.Command("npm", "audit", "--json", "--production")
	command.Dir = projectPath
	npmOutput, err := command.Output()
	err = json.Unmarshal(npmOutput, &npmAudit)
	if err != nil {
		return err
	}
	insightData["npmAudit"] = *npmAudit

	for _, advisoryI := range (*npmAudit)["advisories"].(map[string]interface{}) {
		advisory := advisoryI.(map[string]interface{})
		moduleName := advisory["module_name"].(string)

		findings := advisory["findings"].([]interface{})

		for _, findingI := range findings {
			finding := findingI.(map[string]interface{})
			moduleVersion := finding["version"].(string)
			moduleID := moduleName + "@" + moduleVersion

			mapVulnerabilities(insightData["dependencies"].(map[string]interface{}), advisory, moduleID)
		}
	}

	return nil
}

func mapVulnerabilities(insightData map[string]interface{}, advisory map[string]interface{}, moduleID string) {
	for key, depI := range insightData {
		if dep, ok := depI.(map[string]interface{}); ok {
			if dep["version"] != nil && key+"@"+dep["version"].(string) == moduleID {
				dep["npmAudit"] = advisory
			}
			if dep["dependencies"] != nil {
				mapVulnerabilities(dep["dependencies"].(map[string]interface{}), advisory, moduleID)
			}
		}
	}
}

func mapData(insightData map[string]interface{}, moduleData map[string]interface{}, moduleID string) {
	for key, depI := range insightData {
		if dep, ok := depI.(map[string]interface{}); ok {
			if dep["version"] != nil && key+"@"+dep["version"].(string) == moduleID {
				dep["path"] = moduleData["path"]
				dep["declaredLicenses"] = moduleData["declaredLicenses"]
			}
			if dep["dependencies"] != nil {
				mapData(dep["dependencies"].(map[string]interface{}), moduleData, moduleID)
			}
		}
	}
}
