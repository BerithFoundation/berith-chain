let personal = {
    newAccount: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_newAccount",
            "args" : ["1234"]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#newAccount').val(message.payload)
        })
    },
    unlockAccount: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_unlockAccount",
            "args" : [account , "1234" , 0]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#unlockAccount').val(message.payload)
        })
    },
    lockAccount: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_lockAccount",
            "args" : [account ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#unlockAccount').val(message.payload)
        })
    },
    lockAccount: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_lockAccount",
            "args" : [account ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#unlockAccount').val(message.payload)
        })
    }

    //personal_importRawKey
}