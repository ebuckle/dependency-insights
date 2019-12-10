package report

import (
	"fmt"
	"io"
	"os"
)

// ProduceReport takes the raw json data from a dependency analysis and produces an HTML report
func ProduceReport(insightData map[string]interface{}) {
	totalRisks := buildReportData(insightData["dependencies"].(map[string]interface{}))
	println(totalRisks)
	report, err := os.Create("report.html")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	printReport(report, insightData)
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
		fmt.Fprintf(w, tableRow, formatString, packageData["version"], packageData["declaredLicenses"], "Pending")
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
