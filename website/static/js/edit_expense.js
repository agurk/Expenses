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
    $.post('backend/SAVE_EXPENSE?eid='+eid
			+'&amount='+		$('#exAmount').val()
			+'&date='+			$('#exDate').val()
			+'&description='+	$('#exDesc').val()
			+'&classification='+$('#exClass').val()
			+'&fxAmount='+		$('#exFXAmount').val()
			+'&fxCCY='+			$('#exFXCCY').val()
			+'&fxRate='+		$('#exFXRate').val()
			+'&commission='+	$('#exCommission').val()
	);

}

