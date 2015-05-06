function get_id() {
    return document.getElementById('item_id').innerHTML;
}

function get_type()   {
    return document.getElementById('item_type').innerHTML;
}

function set_pinned() {
    var url = 'pinned?pinned_type='+ get_type() + '&pinned_id='+ get_id();
    var data = $.ajax({method: "POST", url: url, async: false });
    document.getElementById('pinned').innerHTML=data.responseText;
}

function discard_pin() {
    var data = $.ajax({method: "DELETE", url: "pinned", async: false });
    document.getElementById('pinned').innerHTML=data.responseText;
}

function pin(type, id) {
    var url = 'pinned?pinned_type='+ get_type() + '&pinned_id='+ get_id();
    var data = $.ajax({method: "PUT", url: url, async: false });
    document.getElementById('pinned').innerHTML=data.responseText;
}

function load_pinned()  {
	var data=$.ajax({method: "GET", url: "pinned", async: false } );
	document.getElementById('pinned').innerHTML=data.responseText;
}

window.onload = load_pinned
