package report

const (
	htmlHeader = `<html>
	<head>
		<meta charset="utf-8">
		<title>Dependency Insights</title>
		<link type="text/css" href="./report/display/css/jquery.tbltree.css" rel="stylesheet">
		<link type="text/css" href="./report/display/css/styles.css" rel="stylesheet">
		<link type="text/css" href="./report/display/css/pygment_trac.css" rel="stylesheet">
		<link type="text/css" href="./report/display/css/css.css" rel="stylesheet">
    </head>
    <body>
		<div id="doctitle">Dependency Insights</div>
	`
	htmlFooter = `
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<script type="text/javascript" src="./report/display/js/jquery.tbltree.js"></script>
	<script type="text/javascript">
		$('#table').tbltree({
			initState: 'expanded',
		});
	</script>
	<script type="text/javascript">	
	$('#tableVuln').tbltree({
		initState: 'expanded',
		treeColumn: 3
	});
	</script>
	</body>
	</html>
	`

	tableOpen = `
	<table id="table" class="jquery-tbltree">
	<tr>
	  <th>Package Name</th>
	  <th>Package Version</th>
	  <th>Declared License(s)</th>
	  <th>Predicted License(s)</th>
	</tr>
	<tr row-id="0">
		<td>
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td>%s</td>
		<td></td>
		<td></td>
  	</tr>
	`

	vulnTableOpen = `
	<table id="tableVuln" class="jquery-tbltree">
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
	<tr row-id="0">
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
	</tr>
	`

	vulnTableRow = `<tr row-id="%d" parent-id="%d">
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%s</td>
		`

	tableRow = `<tr row-id="%d" parent-id="%d">
		<td>
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">%s</td>
		<td class="data">%s</td>
		<td class="data">%s</td>
	</tr>
	`

	tableClose = `
	</table>
	`
)
