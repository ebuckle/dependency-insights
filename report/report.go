package report

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/browser"
	"gopkg.in/src-d/go-license-detector.v3/licensedb/api"
)

// ProduceReport takes the raw json data from a dependency analysis and produces an HTML report
func ProduceReport(insightData map[string]interface{}) {
	buildReportData(insightData["dependencies"].(map[string]interface{}))
	report, err := os.Create("report.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	printReport(report, insightData)
	browser.OpenURL("report.html")
}

func buildReportData(insightData map[string]interface{}) int {
	// TODO: summary data for deps/sub deps
	return 0
}

func printReport(w io.Writer, insightData map[string]interface{}) {
	fmt.Fprintf(w, htmlHeader, "")
	fmt.Fprintf(w, tableOpen)
	printPackages(w, insightData["dependencies"].(map[string]interface{}), "")
	fmt.Fprintf(w, tableClose)
}

func printPackages(w io.Writer, insightData map[string]interface{}, spacing string) {
	i := 1
	for packageName, packageDataI := range insightData {
		formatString := ""
		packageData := packageDataI.(map[string]interface{})
		if i == len(insightData) {
			formatString = spacing + "└ " + packageName
		} else {
			formatString = spacing + "├ " + packageName
		}
		licenseAnalysis := "No License Data Found"
		if packageData["license-analysis"] != nil {
			licenseAnalysis = produceLicenseString(packageData["license-analysis"].(map[string]api.Match))
		}
		fmt.Fprintf(w, tableRow, formatString, packageData["version"], packageData["declaredLicenses"], licenseAnalysis)
		if packageData["dependencies"] != nil {
			newSpacing := spacing + "│  "
			if i == len(insightData) {
				newSpacing = spacing + "   "
			}
			printPackages(w, packageData["dependencies"].(map[string]interface{}), newSpacing)
		}
		i++
	}
}

func produceLicenseString(licenseAnalysis map[string]api.Match) string {
	returnString := ""
	for licenseName, licenseData := range licenseAnalysis {
		returnString += licenseName + "(" + fmt.Sprintf("%f", licenseData.Confidence) + ")\t"
	}
	return returnString
}
