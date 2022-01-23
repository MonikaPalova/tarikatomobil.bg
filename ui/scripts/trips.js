function _removeChildren(obj) {
    while (obj.firstChild) {
        obj.removeChild(obj.firstChild);
    }
}

// helper methods

function _loadTrips(trips) {
    var tripsSection = document.getElementById("trips");

    _removeChildren(tripsSection);
    trips.forEach(trip => _addTrip(tripsSection, trip));
};

function _addTrip(section, trip) {
    var tripObj = document.createElement("div");
    tripObj.className = "row coloured-field centered-field trip-short";

    var fromTo = document.createElement("p");
    fromTo.innerHTML = `${trip.from} - ${trip.to}`;
    fromTo.classList.add("from-to");
    tripObj.appendChild(fromTo);

    var tripDetails = document.createElement("div");
    tripDetails.className = "trip-short-details row";

    if (trip.smoking) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/cig.png" />`);
    }
    if (trip.airConditioning) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/ac.png" />`);
    }
    if (trip.pets) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/paw.png" />`);
    }

    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-short-detail trip-short-detail-string">${_parseDate(trip.when)}</div>`);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-short-detail trip-short-detail-string">${_parseTime(trip.when)}</div>`);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-short-detail trip-short-detail-string">${_parsePrice(trip.price)}</div>`);
    tripDetails.insertAdjacentHTML('beforeend', `<a  class="trip-short-detail" href="trip.html?id=${trip.id}"><img src="images/open.svg"/></a>`);

    tripObj.appendChild(tripDetails);

    section.appendChild(tripObj);
}

function _parsePrice(priceFloat) {
    let price = priceFloat + '';
    if (price.includes(".")) {
        price = (price + '0').slice(0, price < 10 ? 4 : 5);
    } else {
        price += ".00";
    }

    return price + "лв.";
}

function _parseDate(dateTimeStr) {
    let dateTime = new Date(Date.parse(dateTimeStr));
    let date = ('0' + dateTime.getUTCDate()).slice(-2);
    let month = ('0' + (dateTime.getUTCMonth() + 1)).slice(-2);

    return date + "." + month + "." + dateTime.getUTCFullYear();
}

function _parseTime(dateTimeStr) {
    let dateTime = new Date(Date.parse(dateTimeStr));
    let hours = ('0' + (dateTime.getUTCHours() + 2)).slice(-2);
    let minutes = ('0' + dateTime.getUTCMinutes()).slice(-2);

    return hours + ":" + minutes;
}

function _parseDateTime(dateTimeStr) {
    return _parseDate(dateTimeStr) + ", " + _parseTime(dateTimeStr);
}

function _showError(boxId, msg) {
    var tripsSection = document.getElementById(boxId);
    _removeChildren(tripsSection);

    var error = document.createElement("p");
    error.classList.add("trips-error");
    error.innerHTML = msg;

    tripsSection.appendChild(error);
}

function _showNotification(boxId, msg) {
    var tripsSection = document.getElementById(boxId);
    _removeChildren(tripsSection);

    var notification = document.createElement("p");
    notification.classList.add("trips-notification");
    notification.innerHTML = msg;

    tripsSection.appendChild(notification);
}

// Cities list logic

var CITIES;

function _readCities() {
    var xhr = new XMLHttpRequest(),
        method = 'GET',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = '../cities.json';

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {

            CITIES = JSON.parse(xhr.responseText).cities;
            _loadCities();
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

function _loadCities() {
    var citiesList = document.getElementById("cities-list");

    _removeChildren(citiesList);
    CITIES.forEach(city => _addCity(citiesList, city));
};

function _addCity(list, city) {
    var option = document.createElement('option');
    option.value = city.name;
    list.appendChild(option);
}