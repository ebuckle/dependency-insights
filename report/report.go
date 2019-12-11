package report

import (
	"fmt"
	"io"
	"os"

	"github.com/ebuckle/dependency-insights/insights"
	"github.com/getlantern/deepcopy"
	"github.com/pkg/browser"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/api"
)

// ProduceReport takes the raw json data from a dependency analysis and produces an HTML report
func ProduceReport(insightData *insights.NpmReport) {
	vulnerabilityReport := new(insights.NpmReport)
	buildReportData(insightData, vulnerabilityReport)
	report, err := os.Create("report.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	printReport(report, insightData, vulnerabilityReport)
	browser.OpenURL("report.html")
}

func buildReportData(report *insights.NpmReport, vulnerabilityReport *insights.NpmReport) {
	// TODO: summary data for deps/sub deps
	deepcopy.Copy(vulnerabilityReport, report)
	filterVulnerabilities(&vulnerabilityReport.Dependencies)
	// licenseReport := report
	//spew.Dump(vulnerabilityReport.Dependencies)
}

func printReport(w io.Writer, insightData *insights.NpmReport, vulnerabilityReport *insights.NpmReport) {
	fmt.Fprintf(w, htmlHeader, "")
	fmt.Fprintf(w, vulnTableOpen)
	printVulnerabilities(w, &vulnerabilityReport.Dependencies, "")
	fmt.Fprintf(w, tableClose)
	fmt.Fprintf(w, tableOpen)
	printPackages(w, &insightData.Dependencies, "")
	fmt.Fprintf(w, tableClose)
}

func printPackages(w io.Writer, insightData *map[string]*insights.DependencyData, spacing string) {
	i := 1
	for packageName, packageData := range *insightData {
		formatString := ""
		if i == len(*insightData) {
			formatString = spacing + "└ " + packageName
		} else {
			formatString = spacing + "├ " + packageName
		}
		licenseAnalysis := "No License Data Found"
		if packageData.LicenseAnalysis != nil {
			licenseAnalysis = produceLicenseString(packageData.LicenseAnalysis)
		} else {
			licenseAnalysis = packageData.LicenseAnalysisError
		}
		fmt.Fprintf(w, tableRow, formatString, packageData.Version, packageData.DeclaredLicenses, licenseAnalysis)
		if packageData.Depedencies != nil {
			newSpacing := spacing + "│  "
			if i == len(*insightData) {
				newSpacing = spacing + "   "
			}
			printPackages(w, &packageData.Depedencies, newSpacing)
		}
		i++
	}
}

func printVulnerabilities(w io.Writer, vulnerabilityReport *map[string]*insights.DependencyData, spacing string) {
	i := 1
	for packageName, packageData := range *vulnerabilityReport {
		formatString := ""
		if i == len(*vulnerabilityReport) {
			formatString = spacing + "└ " + packageName
		} else {
			formatString = spacing + "├ " + packageName
		}
		infoString := produceInfoString(packageData.Audit)
		fmt.Fprintf(w, vulnTableRow, packageData.Vulnerabilities.High, packageData.Vulnerabilities.Medium, packageData.Vulnerabilities.Low,
			formatString, packageData.ChildVulnerabilities.High, packageData.ChildVulnerabilities.Medium, packageData.ChildVulnerabilities.Low,
			infoString)
		if packageData.Depedencies != nil {
			newSpacing := spacing + "│  "
			if i == len(*vulnerabilityReport) {
				newSpacing = spacing + "----"
			}
			printVulnerabilities(w, &packageData.Depedencies, newSpacing)
		}
		i++
	}
}

func produceInfoString(auditData map[string]interface{}) string {
	returnString := ""
	for _, vulnDataI := range auditData {
		vulnData := vulnDataI.(map[string]interface{})
		returnString += vulnData["url"].(string) + "\n"
	}
	return returnString
}

func produceLicenseString(licenseAnalysis map[string]api.Match) string {
	returnString := ""
	for licenseName, licenseData := range licenseAnalysis {
		returnString += licenseName + "(" + fmt.Sprintf("%f", licenseData.Confidence) + ")\t"
	}
	return returnString
}

func filterVulnerabilities(vulnerabilityReport *map[string]*insights.DependencyData) {
	for packageName, packageData := range *vulnerabilityReport {
		if compareVulnValues(packageData.Vulnerabilities) && compareVulnValues(packageData.ChildVulnerabilities) {
			delete(*vulnerabilityReport, packageName)
		} else if packageData.Depedencies != nil {
			filterVulnerabilities(&packageData.Depedencies)
		}
	}
}

func compareVulnValues(vulnStruct *insights.Vulnerabilities) bool {
	if vulnStruct.High == vulnStruct.Medium && vulnStruct.Medium == vulnStruct.Low && vulnStruct.Low == 0 {
		return true
	}
	return false
}
