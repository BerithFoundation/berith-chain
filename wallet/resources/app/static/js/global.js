var account;
var loginId;
var mainBalance
var blockNumber
var stakeBalance
var rewardBalance
var totalBalance
//var privateAccount
// var transAmount
// var transAccount
// var stakeAmount
// var stakeAccount

async function sendMessage(methodType, methodName, args) {
    let messagePromise = new Promise(function (resolve) {
        let message = {"name": methodType};
        message.payload = {
            "api": methodName,
            "args": args
        }
        asticode.loader.show()

        //console.log("Request: ", JSON.stringify(message));
        astilectron.sendMessage(message, function (response) {
            asticode.loader.hide();
            //console.log("Response: ", JSON.stringify(response));
            resolve(response);
        });
    });
    return messagePromise;
}

function loadAppContents() {
    $( "#header-content" ).load( "header.html");
    $( "#left-content" ).load( "left.html");
    $( "#main-content" ).load( "main.html");
    $( "#footer-content" ).load( "bottom.html");
}

function loadMainContent(htmlName) {
    $( "#main-content" ).load( htmlName);
}

function loadMainContentWithCallBack(htmlName, callBackFunction) {
    $( "#main-content" ).load( htmlName, callBackFunction);
}


async function hexToDecimal(value) {
    let decimalValue = toDecimal(value);
    console.log(decimalValue);
    return decimalValue;
}

function getWholePart(value) {
    return (value + "").split(".")[0];
}

function getDecimalPart(value) {
    let decimalPart = (value + "").split(".")[1];
    if (!decimalPart) {
        decimalPart = 0;
    }
    return decimalPart;
}

function togglePopUp(popUpId) {
    let selector = "#" + popUpId;
    if ($(selector).hasClass("view")) {
        $(selector).removeClass("view");
        $(selector).addClass("hide");
    } else {
        $(selector).removeClass("hide");
        $(selector).addClass("view");
    }
}

