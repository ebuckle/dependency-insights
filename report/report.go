package report

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/ebuckle/dependency-insights/insights"
	"github.com/getlantern/deepcopy"
	"github.com/pkg/browser"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/api"
)

type pair struct {
	Key   string
	Value float32
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ProduceReport takes the raw json data from a dependency analysis and produces an HTML report
func ProduceReport(insightData *insights.NpmReport) {
	vulnerabilityReport := new(insights.NpmReport)
	licenseReport := new(insights.NpmReport)
	buildReportData(insightData, vulnerabilityReport, licenseReport)
	report, err := os.Create("report.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	printReport(report, insightData, vulnerabilityReport, licenseReport)
	browser.OpenURL("report.html")
}

func buildReportData(report *insights.NpmReport, vulnerabilityReport *insights.NpmReport, licenseReport *insights.NpmReport) {
	deepcopy.Copy(vulnerabilityReport, report)
	deepcopy.Copy(licenseReport, report)
	filterVulnerabilities(&vulnerabilityReport.Dependencies)
	filterProblemLicenses(&licenseReport.Dependencies)
}

func printReport(w io.Writer, insightData *insights.NpmReport, vulnerabilityReport *insights.NpmReport, licenseReport *insights.NpmReport) {
	fmt.Fprintf(w, htmlHeader)
	fmt.Fprintf(w, pageOpen)

	fmt.Fprintf(w, summaryTable, insightData.ChildVulnerabilities.High, insightData.ChildVulnerabilities.Medium, insightData.ChildVulnerabilities.Low,
		insightData.ChildLicenseData.Unknown, insightData.ChildLicenseData.RiskyKeywords, insightData.ChildLicenseData.LicenseCompatability)
	fmt.Fprintf(w, tableClose)

	fmt.Fprintf(w, vulnTableOpen, insightData.Name)
	printVulnerabilities(w, &vulnerabilityReport.Dependencies, 0)
	fmt.Fprintf(w, tableClose)

	fmt.Fprintf(w, licenseTableOpen, insightData.Name)
	printLicenseData(w, &licenseReport.Dependencies, 0)
	fmt.Fprintf(w, tableClose)

	fmt.Fprintf(w, tableOpen, insightData.Name, insightData.Version)
	printPackages(w, &insightData.Dependencies, 0)
	fmt.Fprintf(w, tableClose)

	fmt.Fprintf(w, pageClose)
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
		fmt.Fprintf(w, tableRow, i, parentID, packageName, packageName, packageData.Version, packageData.DeclaredLicenses, licenseAnalysis)
		if packageData.Dependencies != nil {
			i = printPackages(w, &packageData.Dependencies, i)
		}
	}
	return i
}

func printVulnerabilities(w io.Writer, vulnerabilityReport *map[string]*insights.DependencyData, parentID int) int {
	i := parentID
	for packageName, packageData := range *vulnerabilityReport {
		i++
		infoString := produceInfoString(packageData.Audit)
		fmt.Fprintf(w, vulnTableRow, i, parentID, packageData.Vulnerabilities.High, packageData.Vulnerabilities.Medium, packageData.Vulnerabilities.Low, packageName,
			packageName, packageData.Version, packageData.ChildVulnerabilities.High, packageData.ChildVulnerabilities.Medium, packageData.ChildVulnerabilities.Low,
			infoString)
		if packageData.Dependencies != nil {
			i = printVulnerabilities(w, &packageData.Dependencies, i)
		}
	}
	return i
}

func printLicenseData(w io.Writer, licenseReport *map[string]*insights.DependencyData, parentID int) int {
	i := parentID
	for packageName, packageData := range *licenseReport {
		i++
		fmt.Fprintf(w, licenseTableRow, i, parentID, packageData.LicenseData.Unknown, packageData.LicenseData.RiskyKeywords, packageData.LicenseData.LicenseCompatability, packageName,
			packageName, packageData.Version, packageData.ChildLicenseData.Unknown, packageData.ChildLicenseData.RiskyKeywords, packageData.ChildLicenseData.LicenseCompatability,
			packageData.DeclaredLicenses, produceLicenseString(packageData.LicenseAnalysis), packageData.LicenseData.Comment)
		if packageData.Dependencies != nil {
			i = printLicenseData(w, &packageData.Dependencies, i)
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
	returnString := "<ul>"
	sortedList := make(pairList, len(licenseAnalysis))
	i := 0
	for licenseName, licenseData := range licenseAnalysis {
		sortedList[i] = pair{licenseName, licenseData.Confidence}
		i++
	}
	sort.Sort(sort.Reverse(sortedList))

	for _, pair := range sortedList {
		licenseName := pair.Key

		returnString += "<li>" + licenseName + "(" + fmt.Sprintf("%.2f%%", (licenseAnalysis[licenseName].Confidence*100)) + ")</li>"
		if licenseAnalysis[licenseName].Files != nil {
			returnString += "<ul>"
			for file, confidence := range licenseAnalysis[licenseName].Files {
				returnString += "<li>" + file + "(" + fmt.Sprintf("%.2f%%", (confidence*100)) + ")</li>"
			}
			returnString += "</ul>"
		}
	}
	returnString += "</ul>"
	return returnString
}

func filterVulnerabilities(vulnerabilityReport *map[string]*insights.DependencyData) {
	for packageName, packageData := range *vulnerabilityReport {
		if compareVulnValues(packageData.Vulnerabilities) && compareVulnValues(packageData.ChildVulnerabilities) {
			delete(*vulnerabilityReport, packageName)
		} else if packageData.Dependencies != nil {
			filterVulnerabilities(&packageData.Dependencies)
		}
	}
}

func filterProblemLicenses(licenseReport *map[string]*insights.DependencyData) {
	for packageName, packageData := range *licenseReport {
		if compareLicenseValues(packageData.LicenseData) && compareLicenseValues(packageData.ChildLicenseData) {
			delete(*licenseReport, packageName)
		} else if packageData.Dependencies != nil {
			filterProblemLicenses(&packageData.Dependencies)
		}
	}
}

func compareVulnValues(vulnStruct *insights.Vulnerabilities) bool {
	if vulnStruct.High == vulnStruct.Medium && vulnStruct.Medium == vulnStruct.Low && vulnStruct.Low == 0 {
		return true
	}
	return false
}

func compareLicenseValues(licenseData *insights.LicenseData) bool {
	if licenseData.Unknown == licenseData.LicenseCompatability && licenseData.LicenseCompatability == licenseData.RiskyKeywords && licenseData.RiskyKeywords == 0 {
		return true
	}
	return false
}
