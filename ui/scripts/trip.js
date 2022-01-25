var IS_OWNER;
var PASSENGER_NAMES;
var tripID = _getTripID();

// on load logic
window.addEventListener("load", function () {
    _getTripRequest();
    _getTripPassengersRequest();
});

function _getTripID() {
    return new URLSearchParams(window.location.search).get("id");
}

// GET 
function _getTripRequest() {
    var xhr = new XMLHttpRequest(),
        method = 'GET',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/trips/${tripID}`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {

            let trip = JSON.parse(xhr.responseText);
            _loadTrip(trip);
        } else if (xhr.status == 404) {
            _removeChildren(document.getElementById("trip"));
            _showNotification("trip", "Това пътуване не съществува. Вече е започнало или е изтрито.");
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

function _loadTrip(trip) {
    document.getElementById("from").innerHTML = trip.from;
    document.getElementById("to").innerHTML = trip.to;
    document.getElementById("when").innerHTML = _parseDateTime(trip.when);
    IS_TRIP_FINISHED = new Date(trip.when) < new Date();
    document.getElementById("driver").innerHTML = trip.driverName;
    document.getElementById("driver-link").href = "/profile.html?name=" + trip.driverName;
    document.getElementById("price").innerHTML = _parsePrice(trip.price);
    document.getElementById("airConditioning").innerHTML = trip.airConditioning ? "Да" : "Не";
    document.getElementById("smoking").innerHTML = trip.smoking ? "Да" : "Не";
    document.getElementById("pets").innerHTML = trip.pets ? "Да" : "Не";
    document.getElementById("comment").innerHTML = trip.comment;
}

function _getTripPassengersRequest() {
    if (PASSENGER_NAMES != undefined) {
        console.log("passengers already taken");
        return;

    }
    var xhr = new XMLHttpRequest(),
        method = 'GET',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/trips/${tripID}/passengers`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {

            let passengers = JSON.parse(xhr.responseText);
            PASSENGER_NAMES = passengers.map(passenger => passenger.name)
            if (isDriver()) {
                _addPassengers();
            }
            _addTripButtons();
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

function _addTripButtons() {
    _removeChildren("buttons");
    if (!checkLoggedIn()) {
        return; // can do nothing
    }
    if (isDriver()) {
        _addDriverButtons();
    } else if (isPassenger()) {
        _addUnsubcribeButton();
    } else { // user that is not a driver or passenger
        _addSubcribeButton();
    }

    if (isTripFinished() && isPassenger()) {
        _addAddReviewButton();
    }
}

function isDriver() {
    if (IS_OWNER != undefined) {
        return IS_OWNER;
    }
    let tripDriver = document.getElementById("driver").innerHTML;
    let currentUser = getCookie(usernameCookieId);

    IS_OWNER = tripDriver == currentUser;
    return IS_OWNER;
}

function isPassenger() {
    let currentUser = getCookie(usernameCookieId);
    return PASSENGER_NAMES.includes(currentUser);
}


function _addSubcribeButton() {
    if (isTripFinished()) {
        return;
    }
    let buttons = document.getElementById("buttons");

    buttons.insertAdjacentHTML('beforeend', `<button onclick="subscribeToTrip()" class="trip-button">Запиши се</button>`);
}

function _addUnsubcribeButton() {
    if (isTripFinished()) {
        return;
    }
    let buttons = document.getElementById("buttons");

    buttons.insertAdjacentHTML('beforeend', `<button onclick="unsubscribeToTrip()" class="trip-button">Отпиши се</button>`);
}

function _addDriverButtons() {
    if (isTripFinished()) {
        return;
    }
    let buttons = document.getElementById("buttons");

    buttons.insertAdjacentHTML('beforeend', `<button onclick="deleteTrip()" class="trip-button">Изтрий пътуване</button>`);
}

function _addAddReviewButton() {
    let buttons = document.getElementById("buttons");

    let driver = document.getElementById("driver").innerHTML;
    buttons.insertAdjacentHTML('beforeend', `<button onclick="location.href='/profile.html?name=${driver}'" class="trip-button">Напиши ревю</button>`);
}

function _addPassengers() {
    let passengersSection = document.getElementById("passengers");
    _removeChildren(passengersSection);

    if (PASSENGER_NAMES.length == 0) {
        _showNotification("passengers", "Няма записани участници все още.");
    } else {
        passengersSection.insertAdjacentHTML('beforeend', `<header>Записани пътници</header>`);
        PASSENGER_NAMES.forEach(passengerName => _addPassenger(passengersSection, passengerName));
    }
}

function _addPassenger(section, passenger) {
    let passengerObj = document.createElement("div");
    passengerObj.className = "passenger row";

    passengerObj.insertAdjacentHTML('beforeend', `<p>${passenger}</p>`);
    passengerObj.insertAdjacentHTML('beforeend', `<a href="/profile.html?name=${passenger}"><img src="images/person.svg"/></a>`);
    passengerObj.insertAdjacentHTML('beforeend', `<button onclick='removePassenger("${passenger}")'><img src="images/person_Remove.svg"/></button>`);

    section.appendChild(passengerObj);
}

function removePassenger(name) {
    if (isDriver()) {
        _removePassengerRequestByDriver(name);
    }
}

function isTripFinished() {
    return IS_TRIP_FINISHED;
}

function subscribeToTrip() {
    if (isPassenger() || isDriver()) {
        return;
    }
    _addPassengerRequest();
}

function _addPassengerRequest() {
    var xhr = new XMLHttpRequest(),
        method = 'POST',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/trips/${tripID}/passengers`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            _reloadPage();
        } else if (xhr.status === 400) {
            _showError("error-box", "Няма свободни места.");
        }
    };
    xhr.open(method, url, true);
    xhr.send();
}

function unsubscribeToTrip() {
    if (isPassenger()) {
        _removePassengerRequestByPassenger();
    }
}

function _removePassengerRequestByPassenger() {
    // TODO are you sure?

    var xhr = new XMLHttpRequest(),
        method = 'DELETE',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/trips/${tripID}/passengers`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            _reloadPage();
        } else if (xhr.status === 404) {
            _showError("error-box", "Пътуването не съществува или сте отписан от него.");
        }
    };
    xhr.open(method, url, true);
    xhr.send();
}

function _removePassengerRequestByDriver(passengerName) {
    // TODO are you sure?
    let deleteUrl = `/trips/${tripID}/passengers`;
    if (passengerName != "") {
        deleteUrl += "?kickUser=" + passengerName;
    }

    var xhr = new XMLHttpRequest(),
        method = 'DELETE',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = deleteUrl;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            _reloadPage();
        } else if (xhr.status === 400) {
            _showError("error-box", "Нямате право да отписвате пътници от това пътуване.");
        } else if (xhr.status === 404) {
            _showError("error-box", "Пътуването не съществува или потребителят вече е отписан.");
        }
    };
    xhr.open(method, url, true);
    xhr.send();
}

function deleteTrip() {
    if (isDriver()) {
        _deleteTripRequest();
    }
}

function _deleteTripRequest() {
    // TODO are you sure?

    var xhr = new XMLHttpRequest(),
        method = 'DELETE',
        overrideMimeType = 'application/json',
        scheme = 'HTTP',
        url = `/trips/${tripID}`;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 204) {
            _reloadPage();
        } else if (xhr.status === 404) {
            _showError("error-box", "Вие не притежавате това пътуване или то вече е изтрито.");
        } else if (xhr.status === 409) {
            _showError("error-box", "Има записани пътници в това пътуване. Моля, изтрийте ги, за да продължите.");
        }
    };
    xhr.open(method, url, true);
    xhr.send();
}

function _reloadPage() {
    location.reload();
}
