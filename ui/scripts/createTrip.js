window.addEventListener("load", function () {
    if (!checkLoggedIn()) {
        _showError("create-trip", "Трябва да сте влезли в профила си, за да достъпите тази страница.");
        return;
    }
    hasAutomobile();
    _readCities();
});

document.getElementById("create-trip").addEventListener("submit", function (e) {
    e.preventDefault();
});


function hasAutomobile() {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onreadystatechange = function () {
        if (req.readyState === XMLHttpRequest.DONE && req.status == 200) {
            return;
        } else if (req.status == 404) {
            _showError("create-trip", "Нямате регистриран автомобил.");
        }
    };
    req.open('GET', `/users/${getCookie(usernameCookieId)}/automobile`, true);
    req.send();
}

function createTrip() {
    let from = document.getElementById("from").value;
    let to = document.getElementById("to").value;
    let when = new Date(document.getElementById("when").value);
    let price = parseFloat(document.getElementById("price").value);
    let maxPassengers = parseInt(document.getElementById("max-passengers").value);
    let airConditioning = document.getElementById("air-conditioning").checked;
    let smoking = document.getElementById("smoking").checked;
    let pets = document.getElementById("pets").checked;
    let comment = document.getElementById("comment").value;

    var params = {
        from: from,
        to: to,
        when: when,
        maxPassengers: maxPassengers,
        price: price,
        smoking: smoking,
        pets: pets,
        airConditioning: airConditioning,
        comment: comment
    };

    if (!validateParams(params)) {
        return;
    }

    var req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (req.readyState === XMLHttpRequest.DONE && req.status == 200) {
            _showNotification("create-trip","Пътуването беше създадено успешно");
        } else if (req.status == 400) {
            _showError("error-box", "Данните за пътуването не са коректни.");
        }
    };
    req.open('POST', "/trips", true);
    req.send(JSON.stringify(params));

}

function validateParams(params) {
    if (params.from == "") {
        _showError("error-box", "Начална дестинация е задължителна.");
        return false;
    }
    if (params.to == "") {
        _showError("error-box", "Крайна дестинация е задължителна.");
        return false;
    }
    if (isNaN(params.price)) {
        _showError("error-box", "Цената е задължителна");
        return false;
    }
    if (isNaN(params.maxPassengers)) {
        _showError("error-box", "Брой пътници е задължителен.");
        return false;
    }
    if(params.when == undefined){
        _showError("error-box", "Време на тръгване е задължително.");
        return false;
    }
    if (params.when.valueOf() < new Date().valueOf()) {
        _showError("error-box","Времето на тръгване не може да бъде в миналото");
        return false;
    }

    return true;
}

function assertNotBlank(param, msg) {
    if (param == "") {
        _showError("error-box", msg);
        return false;
    }
}

function assertNotNull(param, msg) {
    if (param == null) {
        _showError("error-box", msg);
        return false;
    }
}