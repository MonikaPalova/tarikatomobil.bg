var username;
var imageData = "n";
var imageDataEdit = "n";

window.addEventListener("load", function () {
    onLoadHideNotLoggedIn();
    if (!checkLoggedIn()) {
        window.location = "index.html";
    }

    username = getCookie(usernameCookieId);
    retrieveCar();
    prepareModal();

    const fileInput = document.getElementById("car-img");

    fileInput.addEventListener("change", (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onloadend = () => {
            imageData = reader.result;
        };
        reader.readAsDataURL(file);
    });

    document.getElementById("car-img-edit").addEventListener("change", (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onloadend = () => {
            imageDataEdit = reader.result;
        };
        reader.readAsDataURL(file);
    });
});

function retrieveCar() {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
            populateCarDetails(req.response)
        } else if (req.status == 404) {
            showNotificationNoCar()
        } else {
            loadError(req.responseText);
        }
    };
    req.open('GET', "/users/" + username + "/automobile", true);
    req.send();
}

function populateCarDetails(details) {
    var createSection = document.getElementById("create-car");
    createSection.style.display = "none";

    var detailsSection = document.getElementById("car-details");
    detailsSection.style.display = "initial";

    document.getElementById("reg-num").innerHTML = details.regNumber;
    document.getElementById("comment").innerHTML = details.comment;
    setPhoto(details.photoId);

}

function showNotificationNoCar() {
    var createSection = document.getElementById("create-car");
    createSection.style.display = "initial";

    var detailsSection = document.getElementById("car-details");
    detailsSection.style.display = "none";
  
}

function deleteCar() {
    var req = new XMLHttpRequest();
    req.onload = function () {
        if (req.status == 204) {
            location.reload();
        } else {
            loadError(req.responseText);
        }
    };
    req.open('DELETE', "/users/" + username + "/automobile", true);
    req.send();
}

function createCar() {
    document.getElementById("error").innerHTML = "";
    var reg = document.getElementById("reg").value;
    var comment = document.getElementById("comment-input").value;

    params = {

    };
    
    if (reg == undefined || reg == "") {
        document.getElementById("error").innerHTML = "Трябва да посочите регистрационен номер!";
        return;
    }

    params.regNumber = reg;
    if (comment != undefined && comment != "") {
        params.comment = comment;
    }

    var uploadImageReq = uploadPhotoPrior(imageData);
    if (uploadImageReq != null) {
        uploadImageReq.onload = function () {
            if (uploadImageReq.status != 200) {
                loadError("Грешка!");
            } else {
                params.photoId = uploadImageReq.response.id;
                createCarRequest(params);
            }
        };
    } else {
        createCarRequest(params);
    }
}

function uploadPhotoPrior(imageData) {
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

function createCarRequest(params) {
    var req = new XMLHttpRequest();
        req.responseType = 'json';
        req.open('POST', "/users/" + username + "/automobile", true);
        req.send(JSON.stringify(params));

        req.onload = function () {
            if (req.status == 200) {
                location.reload();
            } else {
                document.getElementById("error").innerHTML = "Грешка! Вече имате въведен автомобил в системата или автомобил с въведения регистрационен номер съществува!";
            }
        };
}

function setPhoto(id) {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.onload = function () {
        if (req.status == 200) {
            var photoBase64 = req.response.base64Content;
            var photoAttribute = document.getElementById('car-photo');
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

function prepareModal() {
    var editModal = document.getElementById("edit-modal");
    var editBtn = document.getElementById("edit-btn");
    var closeModalEdit = document.getElementById("edit-close");

    editBtn.onclick = function () {
        editModal.style.display = "block";
    }

    closeModalEdit.onclick = function () {
        editModal.style.display = "none";
    }


    window.onclick = function (event) {
        if (event.target == editModal) {
            editModal.style.display = "none";
        }
    }
}

function editCar() {
    params = {};
    var comment = document.getElementById("comment-edit").value;

    if (comment != undefined && comment != "") {
        params.comment = comment;
    }

    var uploadImageReq = uploadPhotoPrior(imageDataEdit);
    if (uploadImageReq != null) {
        uploadImageReq.onload = function () {
            if (uploadImageReq.status != 200) {
                loadError("Грешка!");
            } else {
                params.photoId = uploadImageReq.response.id;
                editCarRequest(params);
            }
        };
    } else {
        editCarRequest(params);
    }
}

function editCarRequest(params) {
    var req = new XMLHttpRequest();
    req.responseType = 'json';
    req.open('PATCH', "/users/" + username + "/automobile", true);
    req.send(JSON.stringify(params));

    req.onload = function () {
        if (req.status == 200) {
            location.reload();
        } else {
            loadError(req.responseText);
        }
    };
}