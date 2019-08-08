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

    /*
    getBalance: function (address) {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name": "callApi"};
                message.payload = {
                    "api": "berith_getBalance",
                    "args": [address, "latest"]
                }
                asticode.loader.show()
                astilectron.sendMessage(message, function (message) {
                    asticode.loader.hide();
                    var obj = JSON.parse(message.payload)
                    var val  =parseInt(obj, 16);
                    mainBalance = val
                    resolve(mainBalance)
                }) // astilectron
            }) // settimeout
        }) // promise
    },
    */


    getBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload)
        var val = await hexToDecimal(obj);
        return val;
    },

    /*
    getStakeBalance: function (address) {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name": "callApi"};
                message.payload = {
                    "api" : "berith_getStakeBalance",
                    "args" : [address,"latest" ]
                }
                asticode.loader.show()
                astilectron.sendMessage(message, function(message) {
                    asticode.loader.hide();
                    var obj =JSON.parse(message.payload)
                    var val = parseInt(obj,16)
                    stakeBalance = val
                    resolve(stakeBalance)
                }) // astilectron
            }) // settimeout
        }) // promise
    },*/

    getStakeBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getStakeBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload)
        var val  =parseInt(obj, 16);
        return val;
    },

    /*
    getRewardBalance: function (address) {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name": "callApi"};
                message.payload = {
                    "api" : "berith_getRewardBalance",
                    "args" : [address,"latest" ]
                }
                asticode.loader.show()
                astilectron.sendMessage(message, function(message) {
                    asticode.loader.hide();
                    var obj =JSON.parse(message.payload)
                    var val = parseInt(obj,16)
                    rewardBalance = val
                    // totalBalance = mainBalance + stakeBalance + rewardBalance
                    // console.log( "totBalance ::: " + totalBalance)
                    resolve(rewardBalance)
                })
            }) // settimeout
        }) // promise
    },
*/
    getRewardBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getRewardBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload);
        var val  =parseInt(obj, 16);
        return val;
    },

    /*sendTransaction: function (sendAmount , sendAccount) {
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
    },*/


    sendTransaction: async function (sendAmount , receiverAccount) {
        let valueData2 = toHex(sendAmount);
        result = await sendMessage("callApi", "berith_sendTransaction", [{from : account , to :receiverAccount , value: valueData2 } ]);
        return result;
    },

    /*stakeTransaction: function (stakeAmount ) {
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
    },*/

    stake : async function (stakeAmount) {
        var valueData = toHex(stakeAmount);
        var valueData2 = "0x"+valueData;
        result = await sendMessage("callApi", "berith_stake", [{from : account , value: valueData } ]);
        return result;
    },

    stopStaking : async function () {
        result = await sendMessage("callApi", "berith_stopStaking", [{from: account}]);
        return result;
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
        var valueData = hexConvert.getTxValue(rtbAmount).value;
        var valueData2 = "0x"+valueData;
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
        var valueData = hexConvert.getTxValue(rtsAmount).value;
        var valueData2 = "0x"+valueData;
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

