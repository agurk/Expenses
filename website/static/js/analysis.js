function set_ccy(ccy) {
    var newUrl = $.query.set("ccy", ccy).toString();
    window.location.href = newUrl;
}
function set_to(date) {
    var newUrl = $.query.set("to", date).toString();
    window.location.href = newUrl;
}
function set_from(date) {
    var newUrl = $.query.set("from", date).toString();
    window.location.href = newUrl;
}
