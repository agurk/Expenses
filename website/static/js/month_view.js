function load_expenses_all(date) {
	var all = $('#allButton').attr('aria-pressed');
	var ccy = $('#ccy label.active input').val()
	if (all == 'true') {
		do_load_expenses(date, 'false', ccy);
	} else {
		do_load_expenses(date, 'true', ccy);
	}
}

function load_expenses(date, ccy) {
	var all = $('#allButton').attr('aria-pressed');
	do_load_expenses(date, all, ccy);
}

function do_load_expenses(date, all, ccy)
{
	$.get('detailed_expenses?date='+date+'&all='+all+'&ccy='+ccy, function(data) {
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
