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
    selectTxInfo : async function(txAccount){
        result = await sendMessage("callDB" , "selectTxInfo" , [txAccount])
        return result
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
            // // console.log("member :: " + message.payload)
            var obj = message.payload
            $('#memberData').empty()
            var contents = ''
            // // console.log("ADD :: " + obj.Address)
            // // console.log("id :: " + obj.ID)
            // // console.log("pwd :: " + obj.Password)
            contents += '<tr>'
            contents += '<td><input type="text" value="'+obj.Address+'"></td>'
            contents += '<td><input type="text" value="'+obj.ID+'"></td>'
            contents += '<td><input type="text" value="'+obj.Password+'"></td>'
            contents += '</tr>'
            $('#memberData').append(contents)
        })
    },
    insertTxInfo : function(blockNumber , txAdd , txType , txAmount  , hash , gasLimit, gasPrice ,gasUsed) {
        let message = {"name" : "callDB"}

        message.payload = {
            "api" : "insertTxInfo",
            "args" : [blockNumber ,txAdd , txType ,txAmount , hash , gasLimit , gasPrice ,gasUsed]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            if( message == null || message == ""){
                return
            }
            var obj = message.payload
            // // console.log( "insertTxInfo ::: " + obj)
            // resolve(obj)
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
                    // // console.log( "insertContact ::: " + obj)
                    resolve(obj)
                })
            }) // settimeout
        }) // promise
    },
    updateContact : function (contactAdd ) {
        return new Promise(resolve => {
            setTimeout(()=> {
                let message = {"name" : "callDB"}
                message.payload = {
                    "api" : "updateContact",
                    "args" : [contactAdd ]
                }
                asticode.loader.show()
                astilectron.sendMessage(message , function (message) {
                    asticode.loader.hide()
                    var obj = message.payload
                    resolve(obj)
                })
            }) // settimeout
        }) // promise
    },
    updateMember : function (memberName , memberPwd) {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "updateMember",
            "args" : [memberName , memberPwd]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            var obj = message.payload
            // // console.log("obj :::  " + obj)
            if (obj != "" || obj != undefined){
                succesChange()
            }
        });
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
                // console.log("ADD :: " + obj.Address)
                // console.log("id :: " + obj.ID)
                // console.log("pwd :: " + obj.Password)
                // console.log("private :: " + obj.PrivateKey)
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
                $('#err1').html("키스토어 복원에 실패하였습니다.")
                return
            }else{
                // console.log("ADD :: " + obj.Address)
                // console.log("id :: " + obj.ID)
                // console.log("pwd :: " + obj.Password)
                // console.log("private :: " + obj.PrivateKey)
                location.href="keystoreRestoreConfirm.html";
            }
        });
    }


}