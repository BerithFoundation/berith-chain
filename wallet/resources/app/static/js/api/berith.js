var account
let berith = {

    blockNumber: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_blockNumber",
            "args" : []
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
           // console.log($('#text'))
            $('#text').val(message.payload)
        })
    },

    coinbase: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_coinbase",
            "args" : []
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            //console.log($('#text'))
            $('#text').val(message.payload)
            account = message.payload
            console.log( "account ::: " + account)
        })
    },

    accounts: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_accounts",
            "args" : []
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            //console.log($('#text'))
            $('#text').val(message.payload)

        })
    },
    getBalance: function () {
        let message = {"name": "callApi"};
        console.log("getBalance account ::: " + account)
        message.payload = {
            "api" : "berith_getBalance",
            "args" : [account,"latest" ]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
            //console.log($('#text'))
            $('#text').val(message.payload)
        })
    },
}