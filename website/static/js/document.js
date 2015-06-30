function confirm_doc_expense(dmid, did) {
	$.get('backend/CONFIRM_DOC_EXPENSE?dmid=' + dmid);
	update_expenses(did)
}

function remove_doc_expense(dmid, did) {
	$.get('backend/REMOVE_DOC_EXPENSE?dmid=' + dmid);
	update_expenses(did)
}

function update_document_view(did) {
	new_expense=$.get('document_fragment?did='+did, function(data) {
		document.getElementById('document_'+did).innerHTML=data;
	});
}

function reprocess_document(did) {
	$.get('backend/PROCESS_DOCUMENT?did='+did);
	update_expenses(did)
}

function reclassify_document(did) {
	$.get('backend/RECLASSIFY_DOC?did='+did);
	update_expenses(did)
}

function update_expenses(did) {
	new_expense=$.get('document_all_expense_fragments?did='+did, function(data) {
		document.getElementById('expenses').innerHTML=data;
	});
}
