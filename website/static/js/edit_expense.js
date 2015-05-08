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
