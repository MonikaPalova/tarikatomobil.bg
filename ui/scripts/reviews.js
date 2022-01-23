var username;
window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
    if (!checkLoggedIn()) {
       window.location = "index.html"
    }

    username = getCookie(usernameCookieId);
    showCurrentUserReviews();
});

function showCurrentUserReviews() {
    var username = getCookie(usernameCookieId);

    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
            populateReviews(req.response)
        } else {
            loadError(req.responseText);
        }
    };
    req.open('GET', "/users/" + username + "/reviews", true);
    req.send();
}

function populateReviews(reviews) {
    for (i = 0; i < reviews.length; i++) {
        let review = reviews[i];

        let deleteBtn = document.createElement("button");
        deleteBtn.className = "delete-btn";
        deleteBtn.innerHTML = "Изтрий";

        let reviewDiv = document.createElement("div");
        reviewDiv.className = "review";

        reviewDiv.appendChild(deleteBtn);

        let rating = review.rating;
        let ratingDiv = document.createElement("div");
        ratingDiv.className = "rating";
        let j = 0;
        for (j; j < rating; j++) {
            let span = document.createElement("span");
            span.className = "fa fa-star checked";
            ratingDiv.appendChild(span);
        }

        for (j; j < 5; j++) {
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
        reviewer.innerHTML ="За " + review.forUser;

        textDiv.appendChild(reviewer);

        deleteBtn.onclick = function() {
            deleteReview(review.id);
        }

        reviewDiv.appendChild(textDiv);

        reviewDiv.appendChild(document.createElement("hr"));

        document.getElementById("reviews").appendChild(reviewDiv);

    }


}

function deleteReview(id) {
    var req = new XMLHttpRequest();
    req.onload = function () {
        if (req.status == 204) {
            location.reload();
        } else {
            loadError(req.responseText);
        }
    };
    req.open('DELETE', "/users/" + username + "/reviews/" + id, true);
    req.send();
}
