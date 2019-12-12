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
	// TODO: summary data for licensing
	deepcopy.Copy(vulnerabilityReport, report)
	filterVulnerabilities(&vulnerabilityReport.Dependencies)
}

func printReport(w io.Writer, insightData *insights.NpmReport, vulnerabilityReport *insights.NpmReport) {
	fmt.Fprintf(w, htmlHeader)
	fmt.Fprintf(w, vulnTableOpen, insightData.Name)
	printVulnerabilities(w, &vulnerabilityReport.Dependencies, 0)
	fmt.Fprintf(w, tableClose)
	fmt.Fprintf(w, tableOpen, insightData.Name, insightData.Version)
	printPackages(w, &insightData.Dependencies, 0)
	fmt.Fprintf(w, tableClose)
	fmt.Fprintf(w, htmlFooter)
}

func printPackages(w io.Writer, insightData *map[string]*insights.DependencyData, parentID int) int {
	i := parentID
	for packageName, packageData := range *insightData {
		i++
		var licenseAnalysis string
		if packageData.LicenseAnalysis != nil {
			licenseAnalysis = produceLicenseString(packageData.LicenseAnalysis)
		} else {
			licenseAnalysis = packageData.LicenseAnalysisError
		}
		fmt.Fprintf(w, tableRow, i, parentID, packageName, packageData.Version, packageData.DeclaredLicenses, licenseAnalysis)
		if packageData.Depedencies != nil {
			i = printPackages(w, &packageData.Depedencies, i)
		}
	}
	return i
}

func printVulnerabilities(w io.Writer, vulnerabilityReport *map[string]*insights.DependencyData, parentID int) int {
	i := parentID
	for packageName, packageData := range *vulnerabilityReport {
		i++
		infoString := produceInfoString(packageData.Audit)
		fmt.Fprintf(w, vulnTableRow, i, parentID, packageData.Vulnerabilities.High, packageData.Vulnerabilities.Medium, packageData.Vulnerabilities.Low,
			packageName, packageData.ChildVulnerabilities.High, packageData.ChildVulnerabilities.Medium, packageData.ChildVulnerabilities.Low,
			infoString)
		if packageData.Depedencies != nil {
			i = printVulnerabilities(w, &packageData.Depedencies, i)
		}
	}
	return i
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
