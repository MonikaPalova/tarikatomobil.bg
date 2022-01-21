// on load logic

window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
    _readCities();
    filterTrips();
});

document.getElementById("search-trips").addEventListener("submit", function (e) {
    e.preventDefault();
});

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

function _removeChildren(obj) {
    while (obj.firstChild) {
        obj.removeChild(obj.firstChild);
    }
}
// TRIPS logic

// GET 
function _getTripsRequest(tripsURL) {
    var xhr = new XMLHttpRequest(),
        method = 'GET',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = tripsURL;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {

            var trips = JSON.parse(xhr.responseText);
            _loadTrips(trips);
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

//TODO filters
// POST
//TODO

//DELETE

// helper methods

function _loadTrips(trips) {
    var tripsSection = document.getElementById("trips");

    _removeChildren(tripsSection);
    trips.forEach(trip => _addTrip(tripsSection, trip));
};

function _addTrip(section, trip) {
    var tripObj = document.createElement("div");
    tripObj.className = "row coloured-field centered-field trip";

    var fromTo = document.createElement("p");
    fromTo.innerHTML = `${trip.from} - ${trip.to}`;
    fromTo.classList.add("from-to");
    tripObj.appendChild(fromTo);

    var tripDetails = document.createElement("div");
    tripDetails.className = "trip-details row";

    if (trip.smoking) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/cig.png" />`);
    }
    if (trip.airConditioning) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/ac.png" />`);
    }
    if (trip.pets) {
        tripDetails.insertAdjacentHTML('afterbegin', `<img class="icon" src="images/paw.png" />`);
    }

    var dateTime = new Date(Date.parse(trip.when));
    let date = ('0' + dateTime.getUTCDate()).slice(-2);
    let month = ('0' + (dateTime.getUTCMonth() + 1)).slice(-2);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-details-string">${date}.${month}.${dateTime.getUTCFullYear()}</div>`);

    let hours = ('0' + (dateTime.getUTCHours() + 2)).slice(-2);
    let minutes = ('0' + dateTime.getUTCMinutes()).slice(-2);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-details-string">${hours}:${minutes}</div>`);

    let price = (trip.price + '0').slice(0, trip.price < 10 ? 4 : 5);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-details-string">${price}лв.</div>`);

    tripObj.appendChild(tripDetails);

    section.appendChild(tripObj);
}

function filterTrips() {
    let _determineBooleanParam = function (yesId, noId) {
        let yesValue = document.getElementById(yesId).checked;
        let noValue = document.getElementById(noId).checked;
        if (yesValue == noValue) {
            return "";
        } else if (yesValue) {
            return "true";
        } else {
            return "false";
        }
    }

    let builder = new TripsURLBuilder();

    // from, to, before, after, maxPrice, airConditioning, smoking, pets
    let from = document.getElementById("start-destination").value;
    let to = document.getElementById("end-destination").value;
    // TODO before, after
    let before = "";
    let after = "";
    let maxPrice = document.getElementById("maxprice").value;
    let airConditioning = _determineBooleanParam("air-conditioning-yes", "air-conditioning-no");
    let smoking = _determineBooleanParam("smoking-yes", "smoking-no");
    let pets = _determineBooleanParam("pets-yes", "pets-no");

    let url = builder //
        .setFrom(from) //
        .setTo(to) //
        .setBefore(before) //
        .setAfter(after) //
        .setMaxPrice(maxPrice) //
        .setAirConditioning(airConditioning) //
        .setSmoking(smoking) //
        .setPets(pets) //
        .fetch();

    console.log(url);
    _getTripsRequest(url);
}

// BUILD TRIPS API URL
let TripsURLBuilder = function () {

    let from = "";
    let to = "";
    let before = "";
    let after = "";
    let maxPrice = "";
    let airConditioning = "";
    let smoking = "";
    let pets = "";

    var url = "/trips";

    let _addQueryParam = function (isFirst, paramName, paramValue) {
        if (paramValue.length == 0) {
            return;
        }
        url += isFirst ? "?" : "&";
        url += `${paramName}=${paramValue}`;
    }

    return {
        setFrom: function (from) {
            this.from = from;
            return this;
        },
        setTo: function (to) {
            this.to = to;
            return this;
        },
        setBefore: function (before) {
            this.before = before;
            return this;
        },
        setAfter: function (after) {
            this.after = after;
            return this;
        },
        setMaxPrice: function (maxPrice) {
            this.maxPrice = maxPrice;
            return this;
        },
        setAirConditioning: function (airConditioning) {
            this.airConditioning = airConditioning;
            return this;
        },
        setSmoking: function (smoking) {
            this.smoking = smoking;
            return this;
        },
        setPets: function (pets) {
            this.pets = pets;
            return this;
        },
        fetch: function () {
            _addQueryParam(true, "maxPrice", this.maxPrice); // maxPrice will always be set
            _addQueryParam(false, "from", this.from);
            _addQueryParam(false, "to", this.to);
            _addQueryParam(false, "before", this.before);
            _addQueryParam(false, "after", this.after);
            _addQueryParam(false, "airConditioning", this.airConditioning);
            _addQueryParam(false, "smoking", this.smoking);
            _addQueryParam(false, "pets", this.pets);

            return url;
        }
    };
};