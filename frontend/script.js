$(document).ready(function() {
    $("#methods").ready(setDisplayBlock())
});

//Display only currently wanted input and select elements based on selected REQUEST TYPE
function setDisplayElementText() {
    const divElem = document.querySelector("div");
    const inputElements = divElem.querySelectorAll("input");
    const ElementTexts = ["Actor ID: ", "Director ID: ", "Movie ID: ", "Revenues ID: "];
    var methodsSelectedValue = document.getElementById("methods").value;
    var requestsSelectedValue = document.getElementById("requests").value;
    var elementText = document.getElementById(methodsSelectedValue.toLowerCase() + "-request-element-span");

    //Set all input values to empty strings
    for (let i = 0; i < inputElements.length; i++) {
        inputElements[i].value = "";
    }

    if (methodsSelectedValue == "GET" || methodsSelectedValue == "PATCH") {
        elementText.innerText = ElementTexts[document.getElementById("requests").selectedIndex];
        return;
    }
    if (methodsSelectedValue == "POST") {
        $(".input-block").hide();
        $("#" + requestsSelectedValue + "-input-block").show();
        return;
    }
}

//Display only currently wanted input and select elements based on selected METHOD TYPE
function setDisplayBlock() {
    var methodsSelectedValue = document.getElementById("methods").value;
    $('#requests').val("actors").change();
    $(".element-block").hide();
    $(".display-block").hide();
    $("#" + methodsSelectedValue.toLowerCase() + "-id-block").show();
}

//Takes parameter from element values, select states, predetermined backend address
//Sends a request based on the method selected
function sendRequest() {
    var domainUrl = "http://localhost:5000/";
    var MethodsObject = document.getElementById("methods");
    var methodsSelectedValue = MethodsObject.options[MethodsObject.selectedIndex].value;
    var requestsSelectedValue = document.getElementById("requests").value;
    $(".element-block").hide();

    //Handle request ID, allow only numbers or empty string
    //In case of null or whitespaces, set to ""
    //In case of a number, build a string with a '/' prefix and trim off whitespaces
    if (requestIDValue == null || requestIDValue.trim() === '') { requestIDValue = "" }
    if (isNaN(requestIDValue)) { alert("Please input a whole number greater than 0 or leave the field empty!"); return; }
    if (requestIDValue != "") { requestIDValue = "/" + requestIDValue.trim() }

    switch (methodsSelectedValue) {
        case "GET":
            var requestIDValue = document.getElementById(methodsSelectedValue.toLowerCase() + "-requestID").value;
            if (requestIDValue != "") {
                if (requestIDValue <= 0) { alert("Please input a whole number greater than 0 or leave the field empty."); return; }
                else { requestIDValue = "/" + requestIDValue.trim() }
            }
            httpGet(domainUrl, requestsSelectedValue, requestIDValue, methodsSelectedValue);
            break;
        case "PATCH":
            var requestIDValue = document.getElementById(methodsSelectedValue.toLowerCase() + "-requestID").value;
            if (requestIDValue != "") { requestIDValue = "/" + requestIDValue.trim() }
            if (requestIDValue == "") { alert("You must input ID to delete."); return; }
            httpDelete(domainUrl, requestsSelectedValue, requestIDValue);
            break;
        case "POST":
            httpPost(domainUrl + requestsSelectedValue, requestsSelectedValue);
            break;
    }
}

//Takes values of a JSON object, descriptive identifier of a <span> element,
//iterates through the objects and filles the <span> inner HTML with its values
function populateSpan(objectValues, requestValue) {
    var htmlElement = document.getElementsByClassName(requestValue + "-span");

    for (let i = 0; i < objectValues.length; i++) {
        if (objectValues[i] == null) { htmlElement[i].innerHTML = "<i>no data</i>"; htmlElement[i].style.color = "grey" }
        else { htmlElement[i].innerHTML = objectValues[i]; htmlElement[i].style.color = "black"}
    }
}

//Takes a domain URL, endpoint string, variable ID and a method type
//Builds a URL string and sends a request
//In case variable ID is an empty string, get and display all of the data
function httpGet(domainUrl, endpointValue, requestID, method) {
    $(document).ready(function () {
        var urlString = domainUrl + endpointValue + requestID;
        console.log(urlString);
        $.ajax({
            url: urlString,
            type: method,
            dataType: 'json',
            success: function(response) {
                $("#" + endpointValue + "-element").show();
                var objectValues = Object.values(response);
                console.log(requestID);
                if (requestID != "") { populateSpan(objectValues, endpointValue); return; }
                else {
                    var currentObject = Object.values(objectValues[objectValues.length - 1])
                    populateSpan(currentObject, endpointValue)
                    for (let i = objectValues.length - 2; i >= 0; i--) {
                        currentObject = Object.values(objectValues[i]);
                        var clonedDiv = $('#' + endpointValue + "-element").clone();
                        clonedDiv.attr("id", "#" + endpointValue + "-element" + i);
                        $('#' + endpointValue + "-element").after(clonedDiv);
                        populateSpan(currentObject, endpointValue);
                    }
                }
            },
            error: function (response, textStatus, errorThrown) {
                console.log(textStatus + " - " + errorThrown)
                alert(response.responseText);
            }
        });
    });
}

//Takes a domain URL, endpoint string and variable ID
//Builds a URL string and sends a POST request to soft-delete an item
function httpDelete(domainUrl, endpointValue, requestID) {
    $(document).ready(function () {
        var urlString = domainUrl + endpointValue + requestID;
        $.ajax({
            url: urlString,
            type: "POST",
            crossDomain: true,
            //dataType: 'json',
            success: function(response) {
                alert(response);
            },
            error: function (response, textStatus, errorThrown) { 
                console.log(textStatus + " - " + errorThrown)
                alert(response.responseText);
            }
        });
    });
}

//Takes a domain URL, endpoint string and values from input fields,
//handles their values and builds a valid JSON-able object
//Valid object is parsed into JSON string and sent as POST request body
function httpPost(urlString, endpointValue) {
    var bodyObject = {};
    var selectedItems = document.getElementsByClassName(endpointValue + "-input-field");
    var returnKey = "";

    for (let i = 0; i < selectedItems.length; i++) {
        itemName = selectedItems[i].name;

        //Handle all of the input data and prepare the object for request body
        if (itemName.includes("ID") || itemName.includes("Takings")) {
            if (itemName.includes("Takings") && selectedItems[i].value.length == 0) {
                itemValue = null;
            } else {
                returnKey = itemValue;
                var number = parseInt(selectedItems[i].value);
                if (isNaN(number) || number <= 0) {
                    alert(selectedItems[i].name + " must be a whole number greater than 0!");
                    return;
                }
            }
            itemValue = number;
        }
        if (itemName.includes("Name") || itemName.includes("nationality") || itemName.includes("Lang") || itemName.includes("Certificate")) {
            if (selectedItems[i].value.length == 0 || selectedItems[i].value.trim() === '') {
                itemValue = null
                emptyCounter++;
            } else {
                itemValue = selectedItems[i].value.replace(/[^a-zA-Z ]/g, "").trim();
            }
        }
        if (itemName.includes("gender")) {
            itemValue = selectedItems[i].value;
        }
        if (itemName.includes("release") || itemName.includes("Birth")) {
            if (selectedItems[i].value.length == 0 || selectedItems[i].value == null) {
                alert("Date field must not be empty.");
                return;
            }
            itemValue = selectedItems[i].value;
        }
        //Store itemValue as value in the object under itemName key
        bodyObject[itemName] = itemValue;
    }

    //Send the request
    $.post(urlString, JSON.stringify(bodyObject), function(data) {
        alert("Created object '" + endpointValue + "' with ID '" + Object.values(data)[0] + "'");
    }, "json")
        .fail(function (response) {
            alert("Error: " + response.responseText);
        });
}