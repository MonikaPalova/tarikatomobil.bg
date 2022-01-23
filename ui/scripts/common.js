var sessionCookieId = "TARIKATOMOBIL-SESSION-ID";
var usernameCookieId = "TARIKATOMOBIL-USERNAME";

window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
});

function checkLoggedIn() {
    var session = getCookie(sessionCookieId);
    return session != null;
}

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return null;
}

function setCookie(cname, value) {
    document.cookie = cname + '=' + value + '; expires=Thu, 1 Jan 2100 12:00:00 UTC; path=/';
}

function logout() {
    document.cookie = sessionCookieId + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC;';
    document.cookie = usernameCookieId + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC;';
    location.reload();
}

function onLoadHideNotLoggedIn() {
    if (checkLoggedIn()) {
        document.querySelectorAll(".logged-in").forEach(a => a.style.display = "initial");
        document.querySelectorAll(".not-logged-in").forEach(a => a.style.display = "none");

        var profileLink = document.getElementById("my-profile-link");
        if (profileLink) {
            profileLink.href += "?name=" + getCookie(usernameCookieId);
        }
        console.log(profileLink.href);
    } else {
        document.querySelectorAll(".logged-in").forEach(a => a.style.display = "none");
        document.querySelectorAll(".not-logged-in").forEach(a => a.style.display = "initial");
    }
}

function loadError(err) {
    alert(err);
}

function removeChildren(obj) {
    while (obj.firstChild) {
        obj.removeChild(obj.firstChild);
    }
}