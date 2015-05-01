function show_all(date) {
	document.getElementById('show_all_toggle').setAttribute('class', 'right list-group-item active');
    document.getElementById('show_all_toggle').setAttribute('onclick', "show_only_expenses('"+date+"')")
	$.get('detailed_expenses_all?date='+date, function(data) {
		document.getElementById('detailed_expenses').innerHTML=data;
	});
}

function show_only_expenses(date) {
	document.getElementById('show_all_toggle').setAttribute('class', 'right list-group-item');
    document.getElementById('show_all_toggle').setAttribute('onclick', "show_all('"+date+"')")
	$.get('detailed_expenses?date='+date, function(data) {
		document.getElementById('detailed_expenses').innerHTML=data;
	});
}
