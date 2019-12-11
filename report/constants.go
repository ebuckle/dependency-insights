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

	vulnTableOpen = `<table>
	<tr>
		<th>H</th>
		<th>M</th>
		<th>L</th>
		<th></th>
		<th>H</th>
		<th>M</th>
		<th>L</th>
		<th>Info</th>
	</tr>
	`

	vulnTableRow = `<tr>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%s</td>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%s</td>
		`

	tableRow = `<tr>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
		<td>%s</td>
	</tr>
	`

	tableClose = `
	</table>
	`
)
