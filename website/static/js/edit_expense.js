function duplicate_expense(eid) {  
    $.post('backend/DUPLICATE_EXPENSE?eid='+eid);
}

function reprocess_expense(eid) {
	$.post('backend/REPROCESS_EXPENSE?eid='+eid);
}

function merge_expense(eid, merge_eid) {
	$.post('backend/MERGE_EXPENSE?eid='+eid);
	discard_pin();
}

function merge_expense_commission(eid, merge_eid) {
	$.post('backend/MERGE_EXPENSE_COMMISSION?eid='+eid);
	discard_pin();
}

function save_expense(eid) {
	if (verify_expense_data()) {
	    $.post('backend/SAVE_EXPENSE?eid='+eid
				+'&amount='+		$('#exAmount').val()
				+'&date='+			$('#exDate').val()
				+'&description='+	$('#exDesc').val()
				+'&classification='+$('#exClass').val()
				+'&fxAmount='+		$('#exFXAmount').val()
				+'&fxCCY='+			$('#exFXCCY').val()
				+'&fxRate='+		$('#exFXRate').val()
				+'&commission='+	$('#exCommission').val()
				+'&ccy='+			$('#exCCY').text()
				+'&aid='+			$('#exAccount').val()
				+'&documents='+	get_dids()
		);
	} else {
		alert ('Invalid date');
	}
}

function get_dids() {
	dids=""
	$('.document-thumbnail').each( function(index){
		if ( dids != "") {
			dids = dids + ';'
		}
		dids = dids + this.id
	})
	return dids
}

function verify_expense_data() {
	var d = new Date(document.getElementById('exDate').value);
	if(d instanceof Date && isFinite(d))
	{
		return true;
	} else {
		return false;
	}
}

function cursor_date(e) {
	e = e || window.event;
	switch (e.keyCode) {
		case 38:
			change_date(1);
			break;
		case 40:
			change_date(-1);
			break;
	}
}

function change_date(delta) {
	var d = new Date(document.getElementById('exDate').value);
	d.setDate(d.getDate() + delta);
	document.getElementById('exDate').value = (d.toISOString().slice(0,10));
}

