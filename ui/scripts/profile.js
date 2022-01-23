var username;
window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
      });

    if (params.name == null) {
        window.location = "index.html";
    }
    username = params.name;
    isThisLoggedInUser();
    getUserInfo(username);
});

function isThisLoggedInUser() {
    if (checkLoggedIn()) {
        if (getCookie(usernameCookieId) == username) {
            document.getElementById("logged-in-user-profile").style.display = "initial";
            document.getElementById("not-logged-in-user-profile").style.display = "none";
        } else {
            document.getElementById("logged-in-user-profile").style.display = "none";
            document.getElementById("not-logged-in-user-profile").style.display = "initial";
        }
    } 
}

function populateUserInfo(userInfo) {
    setPhoto(userInfo['photoId']);
    setContactInfo(userInfo);
} 

function getUserInfo(username) {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
           populateUserInfo(req.response);
           loadReviews(req.response.name)
        } else {
            loadError(req.responseText);
        }
      };
    req.open('GET', "/users/" + username, true);
    req.send();
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

function setContactInfo(userInfo) {
    document.getElementById('username').innerHTML = userInfo.name;
    document.getElementById('email').innerHTML = userInfo.email;
    document.getElementById('phone').innerHTML = userInfo.phoneNumber;
    document.getElementById('driver').innerHTML = "Шофирал: " + userInfo.timesDriver + " пъти";
    document.getElementById('passenger').innerHTML = "Пътувал: " + userInfo.timesPassenger + " пъти";
}

function loadReviews(username) {
    // users/newUser/reviews?relation=for
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
            populateReviews(req.response)
        } else {
            loadError(req.responseText);
        }
      };
    req.open('GET', "/users/" + username + "/reviews?relation=for", true);
    req.send();
}

function populateReviews(reviews) {
    for (i = 0; i < reviews.length; i++) {
        let review = reviews[i];
        
        let reviewDiv = document.createElement("div");
        reviewDiv.className = "review";
    
        let rating = review.rating;
        let ratingDiv = document.createElement("div");
        ratingDiv.className = "rating";
        let j = 0;
        for (j ; j < rating; j++) {
            let span = document.createElement("span");
            span.className = "fa fa-star checked";
            ratingDiv.appendChild(span);
        } 
    
        for (j ; j <5; j++) {
            let span = document.createElement("span");
            span.className = "fa fa-star";
            ratingDiv.appendChild(span);
        }
    
        reviewDiv.appendChild(ratingDiv);

        let textDiv = document.createElement("div");
        let reviewText = document.createElement("h2");
        reviewText.innerHTML = review.comment;
        textDiv.appendChild(reviewText);

        let reviewer = document.createElement("p");
        reviewer.className = "reviewer";
        reviewer.innerHTML = review.fromUser;

        textDiv.appendChild(reviewer);

        reviewDiv.appendChild(textDiv);

        reviewDiv.appendChild(document.createElement("hr"));

        document.getElementById("reviews").appendChild(reviewDiv);

    }


}

