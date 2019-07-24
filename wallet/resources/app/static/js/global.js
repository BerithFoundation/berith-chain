var account
var mainBalance
var blockNumber
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
        astilectron.sendMessage(message, function (message) {
            asticode.loader.hide();
            resolve(message);
        });
    });
    return messagePromise;
}

