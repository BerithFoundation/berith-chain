let berith = {

    // blockNumber: function () {
    //     let message = {"name": "callApi"};
    //     message.payload = {
    //         "api" : "berith_blockNumber",
    //         "args" : []
    //     }
    //     asticode.loader.show()
    //     astilectron.sendMessage(message, function(message) {
    //         asticode.loader.hide();
    //         var obj = message.payload
    //         console.log("msg :: " + message.payload)
    //         blockNumber = obj
    //         $('#loginID').val(message.payload)
    //     })
    // },
    blockNumber : async function (address) {
        result = await sendMessage2("callApi", "berith_blockNumber", []);
        var obj = JSON.parse(result.payload)
        //var val = await convertAmount(obj)
        return obj;
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
        var val = await convertAmount(obj)
        return val;
    },

    getStakeBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getStakeBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload)
        var val  =await convertAmount(obj)
        return val;
    },
    getRewardBalance : async function (address) {
        result = await sendMessage("callApi", "berith_getRewardBalance", [address,"latest"]);
        var obj = JSON.parse(result.payload);
        var val  =await convertAmount(obj) / 1000000000000000000;
        // var val  =parseInt(obj, 16);
        return val;
    },
    getRealGasUsed: async function(hash){
        result = await sendMessage("callApi" , "berith_getTransactionReceipt" , [hash])
        var val = JSON.parse(result.payload)
        return val;
    },
    sendTransaction: async function (sendAmount , receiverAccount , gasLimit , gasPrice, nonce) {
        var valueData = hexConvert.getTxValue(sendAmount).value
        var valueData2 = "0x"+valueData
        var gasLimitValue  = parseInt(gasLimit).toString(16)
        var gasPriceValue = parseInt(gasPrice).toString(16)
        var gasLimitValue2  = "0x"+gasLimitValue
        var gasPriceValue2  = "0x"+gasPriceValue
        var param = {from : account , to :receiverAccount , value: valueData2 , gas :gasLimitValue2 , gasPrice: gasPriceValue2} 
        if(nonce){
            param.nonce = "0x"+parseInt(nonce).toString(16);
        }
        result = await sendMessage("callApi", "berith_sendTransaction", [param]);
        return result;
    },

    stake : async function (stakeAmount , gasLimit, gasPrice, nonce) {
        var valueData = hexConvert.getTxValue(stakeAmount).value
        var valueData2 = "0x"+valueData;
        var gasLimitValue  = parseInt(gasLimit).toString(16)
        var gasPriceValue = parseInt(gasPrice).toString(16)
        var gasLimitValue2  = "0x"+gasLimitValue
        var gasPriceValue2  = "0x"+gasPriceValue
        var param = {from : account , value: valueData2 ,gas :gasLimitValue2 , gasPrice: gasPriceValue2} 
        if(nonce){
            param.nonce = "0x"+parseInt(nonce).toString(16);
        }
        result = await sendMessage("callApi", "berith_stake", [param]);
        return result;
    },

    stopStaking : async function (gasLimit, gasPrice, nonce) {
        var gasLimitValue  = parseInt(gasLimit).toString(16)
        var gasPriceValue = parseInt(gasPrice).toString(16)
        var gasLimitValue2  = "0x"+gasLimitValue
        var gasPriceValue2  = "0x"+gasPriceValue
        var param = {from: account , gas: gasLimitValue2 , gasPrice: gasPriceValue2}
        if(nonce){
            param.nonce = "0x"+parseInt(nonce).toString(16);
        }
        result = await sendMessage("callApi", "berith_stopStaking", [param]);
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
            // console.log("updateAccount ::: " + message.payload)
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

    exportKeystore: function (add, pwd , id ) {
        let message = {"name": "exportKeystore"};
        message.payload = {
            "args" : [add ,pwd, id ]
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
    
    getTransactionByHash: async function(hash) {
        return await sendMessage2("callApi","berith_getTransactionByHash",[hash]);
    },

    resendTransaction: function(hash, gasPrice) {
        asticode.loader.show()
        let gasValue = '0x'+(gasPrice*1000000000).toString(16);
        return new Promise((resolve,reject) => {
            let txdata;
            sendMessage2("callApi","berith_getTransactionByHash",[hash]).then(result => {
                if(result && result.payload) {
                    txdata = JSON.parse(result.payload);
                    txdata.gasPrice = gasValue;
                    if(txdata.base == 2)
                        txdata.base = "stake";
                    else
                        txdata.base = "main";

                    if(txdata.target == 2)
                        txdata.target = "stake";
                    else
                        txdata.target = "main";
                    return sendMessage2("callApi","berith_sendTransaction",[{from:txdata.from, to:txdata.to, value:txdata.value, nonce:txdata.nonce, gas:txdata.gas, gasPrice:gasValue, data:txdata.input, base:txdata.base, target:txdata.target}]);
                }
                return Promise.reject(new error("invalid transaction hash"));
            }).then(result => {
                if(result && result.payload) {
                    resolve(txdata)
                    asticode.loader.hide();
                }
                return Promise.reject(new error("failed to send transaction"));
            }).catch(err => {
                reject(err)
                asticode.loader.hide();
            });
        });
    },
    cancelTransaction: function(hash) {
        asticode.loader.show();
        return new Promise((resolve,reject) => {
            let txdata;
            sendMessage2("callApi","berith_getTransactionByHash",[hash]).then(result => {
                if(result && result.payload) {
                    txdata = JSON.parse(result.payload);
                    var gasPrice = (Math.floor(Number(txdata.gasPrice)/1000000000*11)/10 + 0.1).toFixed(1)
                    if(Number(gasPrice) <= 1)
                        gasPrice = "1";
                    gasPrice = (Number(gasPrice) * 1000000000).toString(16);
                    txdata.gasPrice = "0x"+gasPrice;
                    txdata.to = txdata.from;
                    txdata.value = "0x0";
                    txdata.input = "0x";
                    txdata.base = 1;
                    txdata.target = 1;
                    return sendMessage2("callApi","berith_sendTransaction",[{from:txdata.from, to:txdata.to, value:txdata.value, nonce:txdata.nonce, gas:txdata.gas, gasPrice:txdata.gasPrice, data:txdata.input}]);
                }
                return Promise.reject(new error("invalid transaction hash"));
            }).then(result => {
                if(result && result.payload) {
                    resolve(txdata)
                    asticode.loader.hide();
                }
                return Promise.reject(new error("failed to send transaction"));
            }).catch(err => {
                reject(err)
                asticode.loader.hide();
            });
        });
    } 
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

