var account;
var loginId;
var loginPwd;
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

const BERITH_UNIT = 18;
const DISPLAY_UNIT = 8;
const VERSION = "v1.5";


let loops = [];
let trListPage = 1;

// meessage.go 와 통신하는 함수
async function sendMessage(methodType, methodName, args) {
    let messagePromise = new Promise(function (resolve) {
        let message = {"name": methodType};
        message.payload = {
            "api": methodName,
            "args": args
        }
        //asticode.loader.show()

        //console.log("Request: ", JSON.stringify(message));syncing
        astilectron.sendMessage(message, function (response) {
            //asticode.loader.hide();
            //console.log("Response: ", JSON.stringify(response));
            resolve(response);
        });
    });
    return messagePromise;
}
async function sendMessage2(methodType, methodName, args) {
    let messagePromise = new Promise(function (resolve) {
        let message = {"name": methodType};
        message.payload = {
            "api": methodName,
            "args": args
        }
        ///asticode.loader.show()

        //console.log("Request: ", JSON.stringify(message));
        astilectron.sendMessage(message, function (response) {
           // asticode.loader.hide();
            //console.log("Response: ", JSON.stringify(response));
            resolve(response);
        });
    });
    return messagePromise;
}

// spa 구조 선언하는 함수
function loadAppContents() {
    $( "#header-content" ).load( "header.html");
    $( "#left-content" ).load( "left.html");
    $( "#main-content" ).load( "main.html");
    $( "#footer-content" ).load( "bottom.html");
}


function clearLoops() {
    loops.forEach(loop => {
        clearInterval(loop);
    })
}

function loadMainContent(htmlName) {
    clearLoops();
    $( "#main-content" ).load( htmlName, function() {
        registerEvents();
    });
}

function loadMainContentWithCallBack(htmlName, callBackFunction) {
    clearLoops();
    $( "#main-content" ).load( htmlName, function() {
        registerEvents();
        callBackFunction();
    });
}

async function hexToDecimal(value) {
    let decimalValue = toDecimal(value);
    // console.log(decimalValue);
    return decimalValue;
}

async function convertAmount(value) {
    return toBerValue(value);
}
function trlistAmount(value) {
    return toBerValue(value);
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

function toBerValue(s) {
    s = s.substr(2);
    var i, j, digits = [0], carry;
    for (i = 0; i < s.length; i += 1) {
        carry = parseInt(s.charAt(i), 16);
        for (j = 0; j < digits.length; j += 1) {
            digits[j] = digits[j] * 16 + carry;
            carry = digits[j] / 10 | 0;
            digits[j] %= 10;
        }
        while (carry > 0) {
            digits.push(carry % 10);
            carry = carry / 10 | 0;
        }
    }
    for (var left = BERITH_UNIT - digits.length ; left >= 0 ; left--) {
        digits.push(0);
    }
    var result = digits.reverse().join('');
    return result.substr(0,result.length-BERITH_UNIT)+"."+result.substr(result.length-BERITH_UNIT, DISPLAY_UNIT)
}

function getRealtimeBalance(method,address,...elems) {
    var send = () => {
        sendMessage2("callApi", method, [address,"latest"]).then(result => {
            if(result && result.payload) {
                var payload = JSON.parse(result.payload);
                
                elems.forEach(elem => {
                    elem.text(toBerValue(payload));
                    elem.trigger("change");
                });
            }
        }).catch(err => {
            console.error(err);
        });
    }
    send();
    loops.push(setInterval(send , 1000))
}

