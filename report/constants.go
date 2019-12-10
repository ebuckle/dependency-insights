package report

const (
	htmlHeader = `<html>
    <head>
        <title>Dependency Insights</title>
        %s
    </head>
    <body>
        <div id="doctitle">Dependency Insights</div>
	`
	htmlFooter = `</body>
	</html>
	`

	tableOpen = `<table>
	<tr>
	  <th>Package Name</th>
	  <th>Package Version</th>
	  <th>Declared License(s)</th>
	  <th>Predicted License(s)</th>
	</tr>
	`

	tableRow = `<tr>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
	</tr>`

	tableClose = `
	</table>`
)
