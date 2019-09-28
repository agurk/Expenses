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
    var data = {
        "Metadata" : {
            "Confirmed" : true
        }
    }
    $.ajax({
       type : 'PATCH',
       url : 'http://127.0.0.1:8000/expenses/' + eid,
       data : JSON.stringify(data),
       processData : false,
       contentType : 'application/json-patch+json',
    });
	update_expense_view(eid)
}

function update_expense_view(eid) {
	new_expense=$.get('expense_summary?eid='+eid, function(data) {
		document.getElementById('expense_'+eid).innerHTML=data;
	});
}

function tag_expense(eid, tag) {
	//$.get('backend/TAG_EXPENSE?eid='+ eid +'&tag=' + tag)
    var data = {
        "Metadata" : {
            "Tagged" : tag
        }
    }
    $.ajax({
       type : 'PATCH',
       url : 'http://127.0.0.1:8000/expenses/' + eid,
       data : JSON.stringify(data),
       processData : false,
       contentType : 'application/json-patch+json',
    });
	update_expense_view(eid)
}
