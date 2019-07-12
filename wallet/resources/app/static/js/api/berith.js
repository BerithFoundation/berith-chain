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
            $('#loginID').val(message.payload)
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

    pendingTransactions: function ( ) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_pendingTransactions",
            "args" : []
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#pendingResult').val(message.payload)
        })
    },

    updateAccount : function (updateAccountAdd , updateAccountPwd , updateAccountNewPwd) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "berith_updateAccount",
            "args" : [updateAccountAdd,updateAccountPwd,updateAccountNewPwd]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#updateAccountResult').val(message.payload)
            console.log("updateAccount ::: " + message.payload)
        })
    },

    exportKeystore: function () {
        let message = {"name": "exportKeystore"};
        let password = $('#exportPassword').val()

        if (!password) {
            alert("Enter password to export keystore")
            return
        }
        message.payload = {
            "args" : [password]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            var bytes = base64ToArrayBuffer(message.payload)
            var blob=new Blob([bytes], {type: "application/zip"});
            var link=document.createElement('a');
            link.href=window.URL.createObjectURL(blob);
            link.download="berith-keystore.data";
            link.click();
            asticode.loader.hide();
        })
    },
    importKeystore: function (e) {
        asticode.loader.show()
        let file = document.getElementById("keystoreFile").files[0];
        let message = {"name": "importKeystore"};
        let password = $('#importPassword').val()

        if (!password) {
            alert("Enter keystore backup file password")
            return
        }
        message.payload = {
            "args" : [file.path, password]
        }
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
        })
    },
}

function base64ToArrayBuffer(base64) {
    var binaryString = window.atob(base64);
    var binaryLen = binaryString.length;
    var bytes = new Uint8Array(binaryLen);
    for (var i = 0; i < binaryLen; i++) {
        var ascii = binaryString.charCodeAt(i);
        bytes[i] = ascii;
    }
    return bytes;
}

