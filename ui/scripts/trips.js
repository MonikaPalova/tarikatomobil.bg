// on load logic

window.addEventListener("load", function () {
    _readCities();
    _loadTrips();
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

            CITIES = JSON.parse(JSON.parse(JSON.stringify(xhr.responseText))).cities;
            _loadCities();
        }
    };
    xhr.open(method, url, true);
    xhr.send();
};

function _loadCities() {
    var citiesList = document.getElementById("cities-list");

    CITIES.forEach(city => _addCity(citiesList, city));
};

function _addCity(list, city) {
    var option = document.createElement('option');
    option.value = city.name;
    list.appendChild(option);
}


// Load trips logic

var SMOKING_ICON_IMG_NAME = "cig.png";
var AIR_CONDITIONING_ICON_IMG_NAME = "ac.png";
var PETS_ALLOWED_ICON_IMG_NAME = "paw.png";

function _loadTrips() {

}

function _getIcon(imageFileName) {
    var img = document.createElement("img");
    img.classList.add("icon");
    img.src = "images/" + imageFileName;

    return img;
}

function _getTripDateTime(dataStr) {
    var obj = document.createElement("div");
    obj.classList.add("trip-date-time");
    obj.innerHTML = dataStr;

    return obj;
}