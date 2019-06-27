let admin = {
    nodeInfo: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "admin_nodeInfo",
            "args" : []
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            console.log($('#text'))
            $('#nodeInfo').val(message.payload)
        })
    },
}