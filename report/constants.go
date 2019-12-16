package report

const (
	pageOpen = `
	<div id="content" class="bg-dark">
	`

	htmlHeader = `<html>
	<head>
		<meta charset="utf-8">
		<title>Dependency Insights</title>
		<link type="text/css" href="./report/display/css/default.css" rel="stylesheet">
		<link type="text/css" href="./report/display/css/jquery.tbltree.css" rel="stylesheet">
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
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
	<script type="text/javascript">	
	$('#tableLicense').tbltree({
		initState: 'expanded',
		treeColumn: 3
	});
	</script>
	<script type="text/javascript">	
	var allHighCells = document.getElementsById("high")
	for(var i = 0, max = allHighCells.length; i < max; i++) {
		var node = allHighCells[i];
		var currentText = node.childNodes[0].nodeValue;

		if (currentText !== "0")
			node.classList.add("table-danger");
	}
	</script>
	</body>
	</html>
	`

	summaryTable = `
	<div class="card">
			<div class="card-body">
				<h5 class="card-title">Summary</h5>
				<table id="table" class="jquery-tbltree table table-bordered table-sm">
					<tr>
						<th colspan="3">Security Risks</th>
						<th colspan="3">Legal Risks</th>
					</tr>
					<tr>
						<th>H</th>
						<th>M</th>
						<th>L</th>
						<th>UL</th>
						<th>RK</th>
						<th>LC</th>
					  </tr>
					<tr>
					  <td id="high">%d</td>
					  <td>%d</td>
					  <td>%d</td>
					  <td>%d</td>
					  <td>%d</td>
					  <td>%d</td>
					</tr>
	`

	tableOpen = `
	<div class="card">
			<div class="card-body">
				<h5 class="card-title">Security Risks</h5>
	<table id="table" class="jquery-tbltree table table-bordered table-sm">
	<tr>
	  <th>Package Name</th>
	  <th>Package Version</th>
	  <th>Declared License(s)</th>
	  <th>Detected License (Confidence)</th>
	</tr>
	<tr row-id="0">
		<td id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td>%s</td>
		<td></td>
		<td></td>
  	</tr>
	`

	tableRow = `<tr row-id="%d" parent-id="%d">
		<td id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		<a href="https://www.npmjs.com/package/%s" target="_blank">
		%s
		</a>
		</td>
		<td class="data">%s</td>
		<td class="data">%s</td>
		<td class="data">%s</td>
	</tr>
	`

	vulnTableOpen = `
	<div class="card">
			<div class="card-body">
				<h5 class="card-title">Security Risks</h5>
	<table id="tableVuln" class="jquery-tbltree table table-bordered table-sm">
	<tr>
		<th colspan="3">Package Totals</th>
		<th></th>
		<th colspan="3">Tree Subtotals</th>
		<th></th>
	</tr>
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
		<td class="data" id="name">
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
		<td class="data" id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		<a href="https://www.npmjs.com/package/%s" target="_blank">
		%s@%s
		</a>
		</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%s</td>
	`

	licenseTableOpen = `
	<div class="card">
			<div class="card-body">
				<h5 class="card-title">License Risks</h5>
	<table id="tableLicense" class="jquery-tbltree table table-bordered table-sm">
	<tr>
		<th colspan="3">Package Totals</th>
		<th></th>
		<th colspan="3">Tree Subtotals</th>
		<th></th>
		<th></th>
		<th></th>
	</tr>
	<tr>
		<th>UL</th>
		<th>RK</th>
		<th>LC</th>
		<th></th>
		<th>UL</th>
		<th>RK</th>
		<th>LC</th>
		<th>Declared</th>
		<th>Detected (Confidence)</th>
		<th>Comment</th>
	</tr>
	<tr row-id="0">
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data" id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
	</tr>
	`

	licenseTableRow = `<tr row-id="%d" parent-id="%d">
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td class="data" id="name">
	<span class="tbltree-indent"></span>
	<span class="tbltree-expander"></span>
	<a href="https://www.npmjs.com/package/%s" target="_blank">
	%s@%s
	</a>
	</td>
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td class="data">%s</td>
	<td class="data">%s</td>
	<td class="data">%s</td>
	`

	tableClose = `
	</table>
	</div>
	</div>
	`

	pageClose = `
	</div>
	`
)
