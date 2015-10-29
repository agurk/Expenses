function cursor_date(e, id) {
	e = e || window.event;
	switch (e.keyCode) {
		case 38:
			change_date(1, id);
			break;
		case 40:
			change_date(-1, id);
			break;
	}
}

function change_date(delta, id) {
	var d = new Date(document.getElementById(id).value);
	d.setDate(d.getDate() + delta);
	document.getElementById(id).value = (d.toISOString().slice(0,10));
}
