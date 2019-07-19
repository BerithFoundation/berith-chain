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
    hasAddress: function (address) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_hasAddress",
            "args" : [address]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            var obj = message.payload
            console.log("hasAddress :: " + obj)
            if (obj == "true"){
                miner.setBerithbase(address);
                console.log( "trtrtrttr")
                return
            }else{
                console.log( "ffffffff")
                $('#idGroup').addClass('error')
                $('.error_txt').html("PC에 저장된 Keystore File이 없어 로그인 할 수 없습니다.\n" +
                    "Keystore File 복원을 진행해 주세요.\n")
                return
            }

            //
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
            $('#lockAccount').val(message.payload)
        })
    },
    getPrivateKey: function (getPrivateAdd , getPrivatePwd ) {
        let message = {"name": "callApi"};
        message.payload = {
            "api" : "personal_privateKey",
            "args" : [getPrivateAdd,getPrivatePwd ]
        }
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            asticode.loader.hide();
            $('#getPrivateResult').val(message.payload)
        })
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
            var obj2 = JSON.parse(obj)
            console.log( "obj ::: " + obj)
            console.log( "obj2 ::: " + obj2)
            location.href="keystoreRestore2.html?importRawKey="+importRawKeyAdd+"&pwd="+importRawKeyPwd+"&add="+obj2;
        })
    },
}