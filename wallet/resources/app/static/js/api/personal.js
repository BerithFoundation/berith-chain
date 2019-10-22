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
    /*hasAddress: function (address) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_hasAddress",
            "args" : [address]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            var obj = message.payload
            return obj;
        })
    },*/
    hasAddress : async function (address) {
        result = await sendMessage("callApi", "personal_hasAddress", [address]);
        return result.payload;
    },
    /*unlockAccount: function () {
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
    },*/


    unlockAccount : async function (account, password, time) {
        result = await sendMessage("callApi", "personal_unlockAccount", [account, password, time]);
        return result.payload;
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
            $('#lockAccount').val(message.payload)
        })
    },
    // 개인키 조회 api
    getPrivateKey:async function (getPrivateAdd , getPrivatePwd ) {
        result = await sendMessage("callApi", "personal_privateKey" , [getPrivateAdd,getPrivatePwd])
        return result
    },
    importRawKey : function (importRawKeyAdd , importRawKeyPwd) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_importRawKey",
            "args" : [importRawKeyAdd,importRawKeyPwd ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            var obj =message.payload

            if ( obj == "account already exists") {
                $('#privateKeyGroup').addClass('error')
                $('#err1').html("키스토어 파일이 이미 존재합니다.")
                return
            }
            var obj2 = JSON.parse(obj)
            console.log( "obj ::: " + obj)
            console.log( "obj2 ::: " + obj2)
            location.href="keystoreRestore2.html?importRawKey="+importRawKeyAdd+"&pwd="+importRawKeyPwd+"&add="+obj2;
        })
    },
}