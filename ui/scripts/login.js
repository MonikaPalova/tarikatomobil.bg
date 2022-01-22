window.addEventListener("load", function () {
    if (checkLoggedIn()) {
        window.location = "index.html";
    } else {
        document.getElementById("error").display = "none";
    }
});

function login() {
    document.getElementById("error").display = "initial";

    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;

    var authHeader = btoa(username + ':' + password);

    var req = new XMLHttpRequest();

    req.onload = function () {
        if (req.status == 200) {
            window.location = "index.html";
        } else {
            loadError(req.responseText);
        }
      };
    req.open('POST', "/login", true);
    req.setRequestHeader('Authorization','Basic ' + authHeader);
    req.send();

}

function register() {
    document.getElementById("error").display = "initial";

    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    var email = document.getElementById("email").value;
    var phone = document.getElementById("phone").value;

    var params = {
        "name": username,
        "password": password,
        "email": email,
        "phoneNumber": phone
    };

    var req = new XMLHttpRequest();

    req.onload = function () {
        if (req.status == 200) {
            window.location = "login.html";
        } else {
            loadError(req.responseText);
        }
      };
    req.open('POST', "/users", true);
    req.send(JSON.stringify(params));

}

function loadError(errorMessage) {
    document.getElementById("error").innerHTML = errorMessage;
    document.getElementById("error").display = "initial";
}