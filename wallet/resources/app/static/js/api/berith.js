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
            var obj = message.payload
            console.log("msg :: " + message.payload)
            blockNumber = obj
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

    getBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload)
        var val = await hexToDecimal(obj) / 1000000000000000000;
        return val;
    },

    getStakeBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getStakeBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload)
        var val  =await hexToDecimal(obj) / 1000000000000000000;
        // var val  =parseInt(obj, 16);
        return val;
    },

    getRewardBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getRewardBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload);
        var val  =await hexToDecimal(obj) / 1000000000000000000;
        // var val  =parseInt(obj, 16);
        return val;
    },

    sendTransaction: async function (sendAmount , receiverAccount , gasLimit , gasPrice) {
        var valueData = hexConvert.getTxValue(sendAmount).value
        var valueData2 = "0x"+valueData
        var gasLimitValue  = parseInt(gasLimit).toString(16)
        var gasPriceValue = parseInt(gasPrice).toString(16)
        var gasLimitValue2  = "0x"+gasLimitValue
        var gasPriceValue2  = "0x"+gasPriceValue
        console.log( "gasLimitV ::: " +  gasLimitValue2 +"  , gasPriceV ::: " + gasPriceValue2)
        result = await sendMessage("callApi", "berith_sendTransaction", [{from : account , to :receiverAccount , value: valueData2 , gas :gasLimitValue2 , gasPrice: gasPriceValue2 } ]);
        return result;
    },

    stake : async function (stakeAmount , gasLimit, gasPrice) {
        var valueData = hexConvert.getTxValue(stakeAmount).value
        var valueData2 = "0x"+valueData;
        var gasLimitValue  = parseInt(gasLimit).toString(16)
        var gasPriceValue = parseInt(gasPrice).toString(16)
        var gasLimitValue2  = "0x"+gasLimitValue
        var gasPriceValue2  = "0x"+gasPriceValue
        result = await sendMessage("callApi", "berith_stake", [{from : account , value: valueData2 ,gas :gasLimitValue2 , gasPrice: gasPriceValue2  } ]);
        return result;
    },

    stopStaking : async function () {
        result = await sendMessage("callApi", "berith_stopStaking", [{from: account}]);
        return result;
    },
    mining : async function(){
        result = await sendMessage("callApi", "berith_mining", []);
        return result.payload;
    },

    /*rewardToBalance: function (rtbAmount ) {
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
    },*/

    rewardToBalance: async function (rtbAmount) {
        //var valueData = hexConvert.getTxValue(rtbAmount).value;
        var valueData2 = toHex(rtbAmount);
        result = await sendMessage("callApi", "berith_rewardToBalance", [{from : account , value: valueData2 }]);
        return result;
    },



    /*rewardToStake: function (rtsAmount ) {
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
    },*/

    rewardToStake: async function (rtsAmount) {
        // /var valueData = hexConvert.getTxValue(rtsAmount).value;
        var valueData2 = toHex(rtsAmount);
        result = await sendMessage("callApi", "berith_rewardToStake", [{from : account , value: valueData2 }]);
        return result;
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
    startPolling : function () {
        let message = { "name" :  "polling"};
        message.payload = {
            "args" : []
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
        });
    },
    stopPolling : function () {
        let message = { "name" :  "stopPolling"};
        message.payload = {
            "args" : []
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
        });
    },

    exportKeystore: function (pwd) {
        let message = {"name": "exportKeystore"};
        // let password = $('#exportPassword').val()
        //
        // if (!pwd) {
        //     alert("Enter password to export keystore")
        //     return
        // }
        message.payload = {
            "args" : [pwd]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            var bytes = base64ToArrayBuffer(message.payload)
            var blob=new Blob([bytes], {type: "application/zip"});
            var link=document.createElement('a');
            link.href=window.URL.createObjectURL(blob);
            link.download="berith-keystore.zip";
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

