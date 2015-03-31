function show_raw(data, eid) {
    document.getElementById('expenese_detail_' + eid).innerHTML=data;
    document.getElementById('show_details_'+eid).innerHTML='-'
    var function_string = "$.get('expense_details?eid=" + eid +"', function(data){update_expense_view("+eid+")})"
    document.getElementById('show_details_'+eid).setAttribute('onclick', function_string)
    document.getElementById('classification_'+eid).innerHTML=""
    document.getElementById('classification_'+eid ).setAttribute('onclick', function_string)
    document.getElementById('amount_'+eid).innerHTML=""
}
        
function confirm_expense(eid) {
    $.get('backend/CONFIRM_CLASSIFICATION?eid=' + eid)
	update_expense_view(eid)
}   

function update_expense_view(eid) {
	new_expense=$.get('expense?eid='+eid, function(data) {
		document.getElementById('expense_'+eid).innerHTML=data;
	});
}

function tag_expense(eid, tag) {
	$.get('backend/TAG_EXPENSE?eid='+ eid +'&tag=' + tag)
	update_expense_view(eid)
}
