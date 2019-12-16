package insights

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/src-d/go-license-detector.v3/licensedb"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/api"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/filer"
)

// DependencyData is the data structure representing a single dependency and its sub dependencies
type DependencyData struct {
	Version              string                     `json:"version"`
	From                 string                     `json:"from"`
	Resolved             string                     `json:"resolved"`
	Dependencies         map[string]*DependencyData `json:"dependencies"`
	Path                 string                     `json:"path"`
	Audit                map[string]interface{}     `json:"audit"`
	LicenseAnalysis      map[string]api.Match       `json:"licenseAnalysis"`
	LicenseAnalysisError string                     `json:"licenseAnalysisError"`
	DeclaredLicenses     string                     `json:"declaredLicenses"`
	Vulnerabilities      *Vulnerabilities           `json:"Vulnerabilities"`
	ChildVulnerabilities *Vulnerabilities           `json:"childVulnerabilities"`
	LicenseData          *LicenseData               `json:"licenseData"`
	ChildLicenseData     *LicenseData               `json:"childLicenseData"`
}

// NpmReport contains information about the parent project, dependencies, npm audit and project issues
type NpmReport struct {
	Dependencies         map[string]*DependencyData
	Version              string
	Name                 string
	Problems             []string
	Audit                map[string]interface{}
	ChildVulnerabilities *Vulnerabilities
	ChildLicenseData     *LicenseData
}

type packageJSONData struct {
	Name             string
	Version          string
	Description      string
	DeclaredLicenses string
	Path             string
}

// Vulnerabilities stores the number of vulnerabilities a package has
type Vulnerabilities struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
}

// LicenseData stores data about licensing issues for a package
type LicenseData struct {
	Unknown              int    `json:"unknown"`
	RiskyKeywords        int    `json:"riskyKeywords"`
	LicenseCompatability int    `json:"licenseCompatability"`
	Comment              string `json:"comment"`
}

// ProduceInsights calls the appropriate crawling function for the provided language and then reports on licensing
func ProduceInsights(language string, projectPath string) (*NpmReport, error) {
	_, err := os.Stat(projectPath)

	if err != nil {
		return nil, err
	}

	insightData := new(NpmReport)

	// Choose correct walk function depending on the language of the project
	switch language {
	case "nodejs":
		nodeCommand := exec.Command("npm", "ls", "--json", "--production")
		nodeCommand.Dir = projectPath
		rawOutput, _ := nodeCommand.Output()
		err = json.Unmarshal(rawOutput, insightData)
		if err != nil {
			return nil, err
		}
		err = nodeWalk(projectPath, &insightData.Dependencies)
	// case "go":
	// err = goWalk(projectPath, insightData)
	default:
		err := errors.New("language not recognized")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	insightData.ChildLicenseData = calculateLicenseTotals(&insightData.Dependencies)
	performLicenseCheck(&insightData.Dependencies)

	if language == "nodejs" {
		err := checkVulnerabilities(projectPath, insightData)

		if err != nil {
			return nil, err
		}
	}
	return insightData, nil
}

func nodeWalk(projectPath string, dependencies *map[string]*DependencyData) error {
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

					newPackageData := new(packageJSONData)
					transferNodeData(result, newPackageData, path)
					packageID := newPackageData.Name + "@" + newPackageData.Version

					mapData(dependencies, newPackageData, packageID)

					if _, err := os.Stat(path + "/node_modules"); err == nil {
						err := nodeWalk(path, dependencies)
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
func transferNodeData(packageJSON map[string]interface{}, packageData *packageJSONData, path string) {
	if str, ok := packageJSON["name"].(string); ok {
		packageData.Name = str
	}
	if str, ok := packageJSON["version"].(string); ok {
		packageData.Version = str
	}
	if str, ok := packageJSON["description"].(string); ok {
		packageData.Description = str
	}

	if packageJSON["license"] != nil {
		packageData.DeclaredLicenses = parseLicenseDeclaration(packageJSON["license"])
	} else if packageJSON["licenses"] != nil {
		packageData.DeclaredLicenses = parseLicenseDeclaration(packageJSON["licenses"])
	} else {
		packageData.DeclaredLicenses = "No Declared License Found"
	}

	packageData.Path = path
}

func parseLicenseDeclaration(licenseDeclaration interface{}) string {
	if str, ok := licenseDeclaration.(string); ok {
		return str
	} else if byteOut, err := json.Marshal(licenseDeclaration); err != nil {
		return string(byteOut)
	} else {
		return fmt.Sprintf("%#+v", licenseDeclaration)
	}
}

// performLicenseCheck takes an existing map of package data and performs a license check on each package
func performLicenseCheck(dependencies *map[string]*DependencyData) {
	for _, dep := range *dependencies {
		if dep.Path == "" {
			continue
		}

		filer, err := filer.FromDirectory(dep.Path)

		if err != nil {
			dep.LicenseAnalysis = nil
			dep.LicenseAnalysisError = err.Error()
			continue
		}

		results, err := licensedb.Detect(filer)

		if err != nil {
			dep.LicenseAnalysisError = err.Error()
		} else {
			dep.LicenseAnalysis = results
		}

		if dep.Dependencies != nil {
			performLicenseCheck(&dep.Dependencies)
		}
	}
}

func checkVulnerabilities(projectPath string, insightData *NpmReport) error {
	// npmAudit := new(map[string]interface{})
	var npmAudit map[string]interface{}
	command := exec.Command("npm", "audit", "--json", "--production")
	command.Dir = projectPath
	npmOutput, err := command.Output()
	err = json.Unmarshal(npmOutput, &npmAudit)
	if err != nil {
		return err
	}
	insightData.Audit = npmAudit

	if insightData.Audit["error"] != nil {
		errorMessage := insightData.Audit["error"].(map[string]interface{})["detail"].(string)
		return errors.New(errorMessage)
	}

	for id, advisoryI := range insightData.Audit["advisories"].(map[string]interface{}) {
		advisory := advisoryI.(map[string]interface{})
		moduleName := advisory["module_name"].(string)

		findings := advisory["findings"].([]interface{})

		for _, findingI := range findings {
			finding := findingI.(map[string]interface{})
			moduleVersion := finding["version"].(string)
			moduleID := moduleName + "@" + moduleVersion

			mapVulnerabilities(&insightData.Dependencies, id, advisory, moduleID)
		}
	}

	insightData.ChildVulnerabilities = calculateVulnerabilityTotals(&insightData.Dependencies)

	return nil
}

func mapVulnerabilities(dependencies *map[string]*DependencyData, id string, advisory map[string]interface{}, moduleID string) {
	for key, dep := range *dependencies {
		if dep.Version != "" && key+"@"+dep.Version == moduleID {
			if dep.Audit == nil {
				dep.Audit = map[string]interface{}{}
			}
			dep.Audit[id] = advisory
			switch advisory["severity"].(string) {
			case "high":
				dep.Vulnerabilities.High++
			case "medium":
				dep.Vulnerabilities.Medium++
			case "low":
				dep.Vulnerabilities.Low++
			}
		}
		if dep.Dependencies != nil {
			mapVulnerabilities(&dep.Dependencies, id, advisory, moduleID)
		}
	}
}

func sumVulnerabilities(parentTally *Vulnerabilities, childTally *Vulnerabilities) {
	parentTally.High += childTally.High
	parentTally.Medium += childTally.Medium
	parentTally.Low += childTally.Low
}

func calculateVulnerabilityTotals(dependencies *map[string]*DependencyData) *Vulnerabilities {
	vulnerabilityTally := new(Vulnerabilities)
	for _, dep := range *dependencies {
		if dep.Dependencies != nil {
			sumVulnerabilities(dep.ChildVulnerabilities, calculateVulnerabilityTotals(&dep.Dependencies))
		}
		sumVulnerabilities(dep.ChildVulnerabilities, dep.Vulnerabilities)
		sumVulnerabilities(vulnerabilityTally, dep.ChildVulnerabilities)
	}
	return vulnerabilityTally
}

func mapData(dependencies *map[string]*DependencyData, packageData *packageJSONData, moduleID string) {
	for key, dep := range *dependencies {
		if dep.Vulnerabilities == nil {
			dep.Vulnerabilities = new(Vulnerabilities)
			dep.ChildVulnerabilities = new(Vulnerabilities)
			dep.LicenseData = new(LicenseData)
			dep.ChildLicenseData = new(LicenseData)
		}
		if dep.Version != "" && key+"@"+dep.Version == moduleID {
			dep.Path = packageData.Path
			dep.DeclaredLicenses = packageData.DeclaredLicenses
			checkLicensing(dep)
		}
		if dep.Dependencies != nil {
			mapData(&dep.Dependencies, packageData, moduleID)
		}
	}
}

func checkLicensing(dependency *DependencyData) {
	if dependency.DeclaredLicenses == "No Declared License Found" || dependency.DeclaredLicenses == "UNLICENSED" {
		dependency.LicenseData.Unknown++
		dependency.LicenseData.Comment += "Declared license unclear.\n"
	}
}

func calculateLicenseTotals(dependencies *map[string]*DependencyData) *LicenseData {
	licenseTally := new(LicenseData)
	for _, dep := range *dependencies {
		if dep.Dependencies != nil {
			sumLicensing(dep.ChildLicenseData, calculateLicenseTotals(&dep.Dependencies))
		}
		sumLicensing(dep.ChildLicenseData, dep.LicenseData)
		sumLicensing(licenseTally, dep.ChildLicenseData)
	}
	return licenseTally
}

func sumLicensing(parentTally *LicenseData, childTally *LicenseData) {
	parentTally.Unknown += childTally.Unknown
	parentTally.RiskyKeywords += childTally.RiskyKeywords
	parentTally.LicenseCompatability += childTally.LicenseCompatability
}
