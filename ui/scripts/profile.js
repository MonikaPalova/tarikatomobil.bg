var username;
var imageData = "n";
window.addEventListener("load", function () {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });

    if (params.name == null) {
        window.location = "index.html";
    }
    username = params.name;
    isThisLoggedInUser();
    getUserInfo(username);
    prepareModal();
    const fileInput = document.getElementById("avatar");

    fileInput.addEventListener("change", (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onloadend = () => {
            imageData = reader.result;
            console.log(imageData);
        };
        reader.readAsDataURL(file);
    });
});

function prepareModal() {
    var reviewModal = document.getElementById("review-modal");
    var editModal = document.getElementById("edit-modal");

    var reviewBtn = document.getElementById("review-btn");
    var editBtn = document.getElementById("edit-btn");
    var closeModalReview = document.getElementById("review-close");
    var closeModalEdit = document.getElementById("edit-close");

    reviewBtn.onclick = function () {
        reviewModal.style.display = "block";
    }

    closeModalReview.onclick = function () {
        reviewModal.style.display = "none";
    }

    editBtn.onclick = function () {
        editModal.style.display = "block";
    }

    closeModalEdit.onclick = function () {
        editModal.style.display = "none";
    }


    window.onclick = function (event) {
        if (event.target == reviewModal) {
            reviewModal.style.display = "none";
        }

        if (event.target == editModal) {
            editModal.style.display = "none";
        }
    }

    var slider = document.getElementById("myRange");
    var output = document.getElementById("value");
    output.innerHTML = slider.value;

    slider.oninput = function () {
        output.innerHTML = this.value;
    }
}

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
        reviewer.innerHTML = review.fromUser;

        textDiv.appendChild(reviewer);

        reviewDiv.appendChild(textDiv);

        reviewDiv.appendChild(document.createElement("hr"));

        document.getElementById("reviews").appendChild(reviewDiv);

    }


}

function writeReview() {
    var comment = document.getElementById("review-text").value;
    var rating = document.getElementById("value").innerText;

    createReview(username, comment, Number(rating));
}

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
        } else {
            location.reload();
        }
    };
    req.open('POST', "/users/" + username + "/reviews", true);
    req.send(JSON.stringify(params));
}

function editProfile() {

    var newEmail = document.getElementById("new-email").value;;
    var newPhone = document.getElementById("new-phone").value;;
    var newPassword = document.getElementById("new-password").value;

    var params = {};
    if (newEmail != "") {
        params.email = newEmail;
    }
    if (newPhone != "") {
        params.phoneNumber = newPhone;
    }
    if (newPassword != "") {
        params.password = newPassword;
    }
    if (newEmail != "") {
        params.email = newEmail;
    }

    var uploadImageReq = uploadPhotoPriorToEdit();
    if (uploadImageReq != null) {
        uploadImageReq.onload = function () {
            if (uploadImageReq.status != 200) {
                loadError("Грешка!");
                console.log(uploadImageReq.responseText);
            } else {
                console.log(uploadImageReq.response);
                params.photoId = uploadImageReq.response.id;
                editRequest(params);
            }
        };
    } else {
        editRequest(params);
    }

}

function editRequest(params) {
    console.log(params);

    var req = new XMLHttpRequest();

    req.onload = function () {
        if (req.status != 200) {
            loadError("Грешка!");
            console.log(req.responseText);
        } else {
            location.reload();
        }
    };
    req.open('PATCH', "/users/" + username, true);
    req.send(JSON.stringify(params));
}

function uploadPhotoPriorToEdit() {
    if (imageData != "n") {
        var req = new XMLHttpRequest();
        req.responseType = 'json';
        var params = {
            "base64Content": imageData
        }
        req.open('POST', "/photos", true);
        req.send(JSON.stringify(params));

        return req;
    } else {
        return null;
    }
}