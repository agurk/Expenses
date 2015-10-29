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
