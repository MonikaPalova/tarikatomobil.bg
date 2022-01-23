function createReview(forUsername, comment, rating) {
    username = getCookie(usernameCookieId);
    var params = {
        "forUser": forUsername,
        "rating": rating,
        "comment": comment
    };

    var req = new XMLHttpRequest();

    req.onload = function () {
        if (req.status != 200) {
            loadError("Грешка!");
            console.log(req.responseText);
        }
      };
    req.open('POST', "/users/" + username + "/reviews", true);
    req.send(JSON.stringify(params));
}