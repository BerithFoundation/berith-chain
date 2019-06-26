let berith = {

    blockNumber: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_blockNumber",
            "args" : []
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            console.log("msg :: " + message.payload)
            $('#blockNumber').val(message.payload)
        })
    },

    coinbase: function () {
        let message = {"name": "callApi"};

        message.payload = {
            "api" : "berith_coinbase",
            "args" : []
        }

        //asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
          //  asticode.loader.hide();
            $('#coinbase').val(message.payload)
            account = JSON.parse(message.payload)
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
            asticode.loader.hide();
            $('#accounts').val(message.payload)
        })
    },

    getBalance: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_getBalance",
            "args" : [account,"latest" ]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#getBalance').val(message.payload)
        })
    },

    getStakeBalance: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_getStakeBalance",
            "args" : [account,"latest" ]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#getStakeBalance').val(message.payload)
        })
    },

    getRewardBalance: function () {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_getRewardBalance",
            "args" : [account,"latest" ]
        }

        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#getRewardBalance').val(message.payload)
        })
    },

    sendTransaction: function (sendAmount , sendAccount) {
        var valueData = hexConvert.getTxValue(sendAmount).value
        var valueData2 = "0x"+valueData
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_sendTransaction",
            "args" : [{from : account , to :sendAccount , value: valueData2 } ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#sendResult').val(message.payload)
        })
    },

    stakeTransaction: function (stakeAmount ) {
        var valueData = hexConvert.getTxValue(stakeAmount).value
        var valueData2 = "0x"+valueData
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_stake",
            "args" : [{from : account , value: valueData2 } ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#stakeResult').val(message.payload)
        })
    },

    rewardToBalance: function (rtbAmount ) {
        var valueData = hexConvert.getTxValue(rtbAmount).value
        var valueData2 = "0x"+valueData
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_rewardToBalance",
            "args" : [{from : account , value: valueData2 } ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#rtbResult').val(message.payload)
        })
    },
    rewardToStake: function (rtsAmount ) {
        var valueData = hexConvert.getTxValue(rtsAmount).value
        var valueData2 = "0x"+valueData
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_rewardToStake",
            "args" : [{from : account , value: valueData2 } ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#rtsResult').val(message.payload)
        })
    },

}