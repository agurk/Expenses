function load_expenses_all(date, period) {
	var all = $('#allButton').attr('aria-pressed');
	var ccy = $('#ccy label.active input').val()
	if (all == 'true') {
		do_load_expenses(date, '', ccy, period);
	} else {
		do_load_expenses(date, 'ALL', ccy, period);
	}
}

function load_expenses(date, ccy, period) {
	var all = $('#allButton').attr('aria-pressed');
	if (all == 'true') {
		do_load_expenses(date, 'ALL', ccy, period);
	} else {
		do_load_expenses(date, '', ccy, period);
	}
}

function do_load_expenses(date, all, ccy, period)
{
	$.get('detailed_expenses?date='+date+'&all='+all+'&ccy='+ccy+'&period='+period, function(data) {
		document.getElementById('detailed_expenses').innerHTML=data;
	});
}

function set_specific_ccy(ccy) {
	$('#ccyLabel').html(ccy)
	$('#ccySpecific').val(ccy)
	$('#ccySpecific').click()
}

function ccySetup(date) {
$('#ccyBase').change( function() {
	load_expenses(date)
})
}

function highlight_category_matches(category, date, period) {
	var ccy = $('#ccy label.active input').val();
    do_load_expenses(date, category, ccy, period);
}

function set_ccy(ccy) {
    var newUrl = $.query.set("ccy", ccy).toString();
    window.location.href = newUrl;
}

function goto_previous_period(date) {
    var newUrl = $.query.set("date", date).toString();
    window.location.href = newUrl;
}

function set_period(period) {
    var newUrl = $.query.set("period", period).toString();
    window.location.href = newUrl;
}

 +$(document).ready(function(){$("#overall_expenses").tablesorter({sortList: [[1,0], [0,0]]});}); 
