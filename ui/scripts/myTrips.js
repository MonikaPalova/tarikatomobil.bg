// on load logic
window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
    if (!checkLoggedIn()) {
        _showError("Трябва да сте влезли в профила си, за да достъпите тази страница.");
        return;
    }

    showCurrentUserTrips();
});

function showCurrentUserTrips() {
    var isDriver=document.getElementById("is-driver-select").value;
    _getMyTripsRequest(isDriver);
}

function _getMyTripsRequest(isDriver) {
    var xhr = new XMLHttpRequest(),
        method = 'GET',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/mytrips?isDriver=${isDriver}`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {

            var trips = JSON.parse(xhr.responseText);
            if (trips.length == 0) {
                let msg = "Не сте планирали участие в пътуване. ";
                msg += isDriver == "true" ? "Създайте ново пътуване и опитайте отново." : "Запишете се в някое пътуване и опитайте отново.";
                _showNotification(msg);
            } else {
                _loadTrips(trips);
            }
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};
