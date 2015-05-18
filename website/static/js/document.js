function confirm_doc_expense(dmid, did) {
	$.get('backend/CONFIRM_DOC_EXPENSE?dmid=' + dmid);
}

function remove_doc_expense(dmid, did) {
	$.get('backend/REMOVE_DOC_EXPENSE?dmid=' + dmid);
	update_document_view(did);
}

function update_document_view(did) {
	new_expense=$.get('document_fragment?did='+did, function(data) {
		document.getElementById('document_'+did).innerHTML=data;
	});
}
