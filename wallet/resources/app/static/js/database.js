let database = {
    selectContact : function () {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name" : "callDB"}
                message.payload = {
                    "api" : "selectContact",
                    "args" : ["soni"]
                }
                asticode.loader.show()
                astilectron.sendMessage(message , function (message) {
                    asticode.loader.hide()
                    if ( message == undefined){
                        var obj = ""
                        resolve(obj)
                        return
                    }else{
                        var obj = message.payload
                        resolve(obj)
                    }
                }) // astilectron
            }) // settimeout
        }) // promise
    },
    checkLogin : async function (memberName, memberPwd) {
        result = await sendMessage("callDB", "checkLogin", [memberName, memberPwd]);
        return result;
    },

    selectMember : function () {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "selectMember",
            "args" : ["aa"]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            console.log("member :: " + message.payload)
            var obj = message.payload
            $('#memberData').empty()
            var contents = ''
            console.log("ADD :: " + obj.Address)
            console.log("id :: " + obj.ID)
            console.log("pwd :: " + obj.Password)
            contents += '<tr>'
            contents += '<td><input type="text" value="'+obj.Address+'"></td>'
            contents += '<td><input type="text" value="'+obj.ID+'"></td>'
            contents += '<td><input type="text" value="'+obj.Password+'"></td>'
            contents += '</tr>'
            $('#memberData').append(contents)
        })
    },
    insertContact : function (contactAdd , contactName) {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name" : "callDB"}
                message.payload = {
                    "api" : "insertContact",
                    "args" : [contactAdd , contactName]
                }
                asticode.loader.show()
                astilectron.sendMessage(message , function (message) {
                    asticode.loader.hide()
                    var obj = message.payload
                    console.log( "insertContact ::: " + obj)
                    resolve(obj)
                })
            }) // settimeout
        }) // promise
    },
    updateContact : function (contactAdd , contactName) {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "updateContact",
            "args" : [contactAdd , contactName]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            var obj = message.payload
            var keys = Object.keys(obj)
            for ( var i in keys) {
                console.log("add : " +keys[i]+ " , name : "  + obj[keys[i]])
            }
        })
    },
    insertMember : function (memberName , memberPwd) {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "insertMember",
            "args" : [memberName , memberPwd]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            var obj = message.payload
            if(obj == "err"){
                $('#idGroup').addClass('error')
                $('#err1').html("이미 존재하는 아이디 입니다.")
                return
            }else{
                console.log("ADD :: " + obj.Address)
                console.log("id :: " + obj.ID)
                console.log("pwd :: " + obj.Password)
                console.log("private :: " + obj.PrivateKey)
                location.href="createAccountConfirm.html?Address="+obj.Address+"&ID="+obj.ID+"&PrivateKey="+obj.PrivateKey;
            }
        });
    },
    restoreMember : function (add , id , pwd , privateKey) {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "restoreMember",
            "args" : [add , id , pwd , privateKey]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            var obj = message.payload
            if(obj == "err"){
                $('#idGroup').addClass('error')
                $('#err1').html("이미 존재하는 아이디 입니다.")
                return
            }else{
                console.log("ADD :: " + obj.Address)
                console.log("id :: " + obj.ID)
                console.log("pwd :: " + obj.Password)
                console.log("private :: " + obj.PrivateKey)
                location.href="keystoreRestoreConfirm.html";
            }
        });
    }

}