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
}