// on load logic
window.addEventListener("load", function () {
    _readCities();
    filterTrips();
});

document.getElementById("search-trips").addEventListener("submit", function (e) {
    e.preventDefault();
});

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
            if (trips.length == 0) {
                _showNotification("trips","Няма пътувания, които съответсват на търсенето.");
            } else {
                _loadTrips(trips);
            }
        } else if(xhr.status === 400){
            _showError("trips","Има проблем със взимането на пътувания. Моля, свържете се с администратор.")
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

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

    let from = document.getElementById("start-destination").value;
    let to = document.getElementById("end-destination").value;
    var dateStr = document.getElementById("trip-date").value;
    if (dateStr != "") {
        let date = new Date(dateStr);
        if (date.valueOf() < new Date().setHours(0, 0, 0, 0).valueOf()) {
            _showError("trips","Датата за пътуване не може да бъде в миналото");
            return;
        }
        let before = _formatDate(new Date(date.getTime() + 28 * 60 * 60 * 1000));
        let after = _formatDate(new Date(date.getTime() - 20 * 60 * 60 * 1000));
        builder = builder.setBefore(before).setAfter(after);
    } else {
        let after = _formatDate(new Date());
        builder = builder.setAfter(after);
    }
    let maxPrice = document.getElementById("maxprice").value;
    let airConditioning = _determineBooleanParam("air-conditioning-yes", "air-conditioning-no");
    let smoking = _determineBooleanParam("smoking-yes", "smoking-no");
    let pets = _determineBooleanParam("pets-yes", "pets-no");

    let url = builder //
        .setFrom(from) //
        .setTo(to) //
        .setMaxPrice(maxPrice) //
        .setAirConditioning(airConditioning) //
        .setSmoking(smoking) //
        .setPets(pets) //
        .fetch();

    _getTripsRequest(url);
}

// BUILD TRIPS API URL
let TripsURLBuilder = function () {

    var url = "/trips";

    let _addQueryParam = function (isFirst, paramName, paramValue) {
        if (paramValue.length == 0) {
            return;
        }
        url += isFirst ? "?" : "&";
        url += `${paramName}=${paramValue}`;
    }

    return {
        from: "",
        to: "",
        before: "",
        after: "",
        maxPrice: "",
        airConditioning: "",
        smoking: "",
        pets: "",

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

// helper methods
function _formatDate(date) {
    let day = date.toLocaleString('default', { day: '2-digit' });
    let month = date.toLocaleString('default', { month: 'short' });
    let year = date.toLocaleString('default', { year: 'numeric' });
    return year + '-' + month + '-' + day;
}