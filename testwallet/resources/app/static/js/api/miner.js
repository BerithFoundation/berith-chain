let miner = {
    /*setBerithbase: function (setAccount) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "miner_setBerithbase",
            "args" : [setAccount]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            var obj  = message.payload
            console.log("obj1 :: " +setAccount)
            console.log("obj2 :: " +obj)
            return obj;
        });
    },*/

    setBerithbase : async function (address) {
        result = await sendMessage("callApi", "miner_setBerithbase", [address]);
        return result.payload;
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