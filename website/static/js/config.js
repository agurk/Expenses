function save_account(aid) {
	$.post('backend/SAVE_ACCOUNT?aid='+aid
		+'&name='+ $('#'+aid+'_account_name').val()
		+'&ccy=' + $('#'+aid+'_account_ccy').val()
		+'&lid=' + $('#'+aid+'_account_lid').val()
		+'&pid=' + $('#'+aid+'_account_pid').val()
	);
}
