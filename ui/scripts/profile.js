window.addEventListener("load", function () {
    if (checkLoggedIn()) {
        getUserInfo();
    } else {
        document.getElementById("error").display = "none";
    }
});

function populateUserInfo(userInfo) {
    console.log(userInfo);
    setPhoto(userInfo['photoId']);

} 

function getUserInfo() {
    var session = getCookie(sessionCookieId);
    var username = getCookie(usernameCookieId);

    if (session == null) {
        return null;
    }

    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
           populateUserInfo(req.response);
        } else {
            loadError(req.responseText);
        }
      };
    req.open('GET', "/users/" + username, true);
    req.send();

    return response;
}

function setPhoto(id) {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
            var photoBase64 = req.response.base64Content;
            var photoAttribute = document.getElementById('profile-photo');
            photoAttribute.setAttribute(
                'src', photoBase64
            );
        } else {
            loadError(req.responseText);
        }
      };
    req.open('GET', "/photos/" + id, true);
    req.send();
}