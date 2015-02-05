function update_expense(eid) {  
    $.post('backend/CHANGE_CLASSIFICATION?eid='+eid + '&cid='+ document.getElementById ('current_expense_'+eid).value);
    $.post('backend/CHANGE_AMOUNT?eid='+eid + '&amount='+ document.getElementById ('current_amount_'+eid).value);
	update_expense_view(eid)
}
