let miner = {
    setBerithbase: function (setAccount) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "miner_setBerithbase",
            "args" : [setAccount]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            var obj  = message.payload
            console.log("obj :: " +obj)
            $('#setBerithbase').val(message.payload)
            location.href="main.html?Address="+setAccount;
        })
    },
    miningStart: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "miner_start",
            "args" : []
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#miningStart').val(message.payload)
        })
    },
    miningStop: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "miner_stop",
            "args" : []
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#miningStop').val(message.payload)
        })
    }
}