package report

const (
	pageOpen = `
	<div id="content" class="">
	`

	noneFound = `
	<div>None Found.</div>
	`

	htmlHeader = `<html>
	<head>
		<meta charset="utf-8">
		<title>Dependency Insights</title>
		<link type="text/css" href="./report/display/css/jquery.tbltree.css" rel="stylesheet">
		<link rel="stylesheet" href="./report/display/css/darkly.css" rel="stylesheet">
		<link type="text/css" href="./report/display/css/default.css" rel="stylesheet">
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

	var allMediumCells = document.getElementsByClassName("medium")
	for(var i = 0, max = allMediumCells.length; i < max; i++) {
		var node = allMediumCells[i];
		var currentText = node.childNodes[0].nodeValue;

		if (currentText !== "0" && currentText !== "-") {
			node.classList.add("table-warning");
			var module = node.parentElement.querySelector('#name')
			if (module !== null) {
				module = module.querySelector('a');
				module.classList.add("table-warning");
			}
		}
	}

	var allLowCells = document.getElementsByClassName("low")
	for(var i = 0, max = allLowCells.length; i < max; i++) {
		var node = allLowCells[i];
		var currentText = node.childNodes[0].nodeValue;

		if (currentText !== "0" && currentText !== "-") {
			node.classList.add("table-warning");
			var module = node.parentElement.querySelector('#name')
			if (module !== null) {
				module = module.querySelector('a');
				module.classList.add("table-warning");
			}
		}
	}

	var allHighCells = document.getElementsByClassName("high")
	for(var i = 0, max = allHighCells.length; i < max; i++) {
		var node = allHighCells[i];
		var currentText = node.childNodes[0].nodeValue;

		if (currentText !== "0" && currentText !== "-") {
			node.classList.add("table-danger");
			var module = node.parentElement.querySelector('#name')
			if (module !== null) {
				module = module.querySelector('a');
				module.classList.add("table-danger");
			}
		}
	}
	</script>
	</body>
	</html>
	`

	summaryTable = `
	<div class="card">
			<div class="card-body">
				<h2 class="card-title">Summary</h2>
				<hr>
				<table id="table" class="jquery-tbltree table table-bordered table-sm table-striped table-hover">
				<thead class="thead">
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
				</thead>
					<tr>
					  <td class="high data">%d</td>
					  <td class="medium data">%d</td>
					  <td class="low data">%d</td>
					  <td class="high data">%d</td>
					  <td class="medium data">%d</td>
					  <td class="low data">%d</td>
					</tr>
	`

	tableOpen = `
	<div class="card">
			<div class="card-body">
				<h2 class="card-title">Dependency Tree</h2>
				<hr>
	<table id="table" class="jquery-tbltree table table-bordered table-sm table-striped table-hover">
	<thead class="thead">
	<tr>
	  <th>Package Name</th>
	  <th>Package Version</th>
	  <th>Declared License(s)</th>
	  <th>Detected License (Confidence)</th>
	</tr>
	</thead>
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
		<td class="">%s</td>
		<td class="">%s</td>
		<td><div class="overflow-auto">%s</div></td>
	</tr>
	`

	vulnTableOpen = `
	<div class="card">
			<div class="card-body">
				<h2 class="card-title">Security Risks</h2>
				<hr>
	<table id="tableVuln" class="jquery-tbltree table table-bordered table-sm table-striped table-hover">
	<thead class="thead">
	<tr>
		<th colspan="3">Package Totals</th>
		<th></th>
		<th colspan="3">Tree Subtotals</th>
		<th></th>
	</tr>
	<tr>
		<th class="data">H</th>
		<th class="data">M</th>
		<th class="data">L</th>
		<th></th>
		<th class="data">H</th>
		<th class="data">M</th>
		<th class="data">L</th>
		<th>Info</th>
	</tr>
	</thead>
	<tr row-id="0">
		<td class="data high">-</td>
		<td class="data medium">-</td>
		<td class="data low">-</td>
		<td class="name" id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="">-</td>
	</tr>
	`

	vulnTableRow = `<tr row-id="%d" parent-id="%d">
		<td class="data high">%d</td>
		<td class="data medium">%d</td>
		<td class="data low">%d</td>
		<td class="name" id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		<a href="https://www.npmjs.com/package/%s" target="_blank">
		%s@%s
		</a>
		</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="data">%d</td>
		<td class="">%s</td>
	`

	licenseTableOpen = `
	<div class="card">
			<div class="card-body">
				<h2 class="card-title">License Risks</h2>
				<hr>
	<table id="tableLicense" class="jquery-tbltree table table-bordered table-sm table-striped table-hover">
	<thead class="thead">
	<tr>
		<th colspan="3">Package Totals</th>
		<th></th>
		<th colspan="3">Tree Subtotals</th>
		<th></th>
		<th></th>
		<th></th>
	</tr>
	<tr>
		<th class="data">UL</th>
		<th class="data">RK</th>
		<th class="data">LC</th>
		<th class="name"></th>
		<th class="data">UL</th>
		<th class="data">RK</th>
		<th class="data">LC</th>
		<th>Declared</th>
		<th>Detected (Confidence)</th>
		<th>Keyword Hits</th>
		<th>Comment</th>
	</tr>
	</thead>
	<tr row-id="0">
		<td class="data high">-</td>
		<td class="data medium">-</td>
		<td class="data low">-</td>
		<td class="name" id="name">
		<span class="tbltree-indent"></span>
		<span class="tbltree-expander"></span>
		%s
		</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td class="data">-</td>
		<td><div class="overflow-auto">-</div></td>
		<td><div class="overflow-auto">-</div></td>
		<td><div class="overflow-auto">-</div></td>
		<td class="">-</td>
	</tr>
	`

	licenseTableRow = `<tr row-id="%d" parent-id="%d">
	<td class="data high">%d</td>
	<td class="data medium">%d</td>
	<td class="data low">%d</td>
	<td class="name" id="name">
	<span class="tbltree-indent"></span>
	<span class="tbltree-expander"></span>
	<a href="https://www.npmjs.com/package/%s" target="_blank">
	%s@%s
	</a>
	</td>
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td class="data">%d</td>
	<td><div class="overflow-auto">%s</div></td>
	<td><div class="overflow-auto">%s</div></td>
	<td><div class="overflow-auto">%s</div></td>
	<td class="">%s</td>
	`

	tableClose = `
	</table>
	</div>
	</div>
	`

	pageClose = `
	</div>
	`

	cardOpen = `
	<div class="card">
			<div class="card-body">
				<h2 class="card-title">%s</h2>
				<hr>
	`
	cardClose = `
	</div>
	</div>`
)
