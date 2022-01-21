var sessionCookieId = "TARIKATOMOBIL-SESSION-ID";

function checkLoggedIn() {
    var session = getCookie(sessionCookieId);
    
    return session!=null;
}

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
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

function logout() {
    console.log("logout");
    document.cookie = sessionCookieId +'=; expires=Thu, 01 Jan 1970 00:00:00 UTC;';  
    window.location = "index.html";
}

function onLoadHideNotLoggedIn() {
    if (checkLoggedIn()) {
        document.querySelectorAll(".logged-in").forEach(a=>a.style.display = "initial");
        document.querySelectorAll(".not-logged-in").forEach(a=>a.style.display = "none");
    } else {
        document.querySelectorAll(".logged-in").forEach(a=>a.style.display = "none");
        document.querySelectorAll(".not-logged-in").forEach(a=>a.style.display = "initial");
    }
}