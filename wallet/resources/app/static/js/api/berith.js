let berith = {
    blockNumber: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_blockNumber",
            "args" : ["asdasd", "11111111"]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            console.log($('#text'))
            $('#text').val(message.payload)
        })
    },
    getBalance: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_getBalance",
            "args" : ["0x78c2b0dfde452677ccd0cd00465e7cca0e3c5353", "latest"]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            console.log($('#text'))
            $('#text').val(message.payload)
        })
    },
}