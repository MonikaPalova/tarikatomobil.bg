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
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-detail trip-detail-string">${date}.${month}.${dateTime.getUTCFullYear()}</div>`);

    let hours = ('0' + (dateTime.getUTCHours() + 2)).slice(-2);
    let minutes = ('0' + dateTime.getUTCMinutes()).slice(-2);
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-detail trip-detail-string">${hours}:${minutes}</div>`);

    let price = trip.price + '';
    if(price.includes(".")){
    price = (trip.price + '0').slice(0, trip.price < 10 ? 4 : 5);
    } else {
        price += ".00";
    }
    tripDetails.insertAdjacentHTML('beforeend', `<div class="trip-detail trip-detail-string">${price}лв.</div>`);

    // TODO
    tripDetails.insertAdjacentHTML('beforeend', `<a  class="trip-detail" href="ADD_URL_TO_THIS_TRIP"><img src="images/open.svg"/></a>`);

    tripObj.appendChild(tripDetails);

    section.appendChild(tripObj);
}


function _showError(msg) {
    var tripsSection = document.getElementById("trips");
    _removeChildren(tripsSection);

    var error = document.createElement("p");
    error.classList.add("trips-error");
    error.innerHTML = msg;

    tripsSection.appendChild(error);
}

function _showNotification(msg) {
    var tripsSection = document.getElementById("trips");
    _removeChildren(tripsSection);

    var notification = document.createElement("p");
    notification.classList.add("trips-notification");
    notification.innerHTML = msg;

    tripsSection.appendChild(notification);
}