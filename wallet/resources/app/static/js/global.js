
var account
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

        // console.log("Request: ", JSON.stringify(message))
        astilectron.sendMessage(message, function (response) {
            asticode.loader.hide();
            // console.log("Response: ", JSON.stringify(response))
            resolve(response);
        });
    });
    return messagePromise;
}

