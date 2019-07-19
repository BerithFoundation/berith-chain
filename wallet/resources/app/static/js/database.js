let database = {
    selectContact : function () {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "selectContact",
            "args" : ["soni"]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            var obj = message.payload
            console.log("contact :: " +message.payload)
            console.log("sdfsdfsdf :: " + Array.isArray(obj))
            var keys = Object.keys(obj);
            var contents = ''
            $('#contactData').empty()
            for ( var i in keys) {
                contents += '<tr>'
                console.log("key="+keys[i]+ ",  data="+ obj[keys[i]]);
                contents += '<td><input type="text" value="'+keys[i]+'"></td>'
                contents += '<td><input type="text" value="'+obj[keys[i]]+'"></td>'
                contents += '</tr>'
            }
            $('#contactData').append(contents)
        })
    },
    checkLogin : function(memberName, memberPwd) {
        let message = {"name" : "callDB"}
        var result
        message.payload = {
            "api" : "checkLogin",
            "args" : [memberName,memberPwd ]
        }
        astilectron.sendMessage(message , function (message) {
            if ( message == undefined){
                console.log("no exists !!! ")
                $('#idGroup').addClass('error')
                $('.error_txt').html("일치하는 아이디가 존재하지 않습니다.")
                result = false
                return
            }
            var obj = message.payload
            if(memberPwd != obj.Password) {
                $('#pwdGroup').addClass('error')
                $('.error_txt').html("비밀번호가 일치하지 않습니다.")
                result = false
                return
            }else {
                $('#idGroup').removeClass('error')
                $('#pwdGroup').removeClass('error')
                console.log( "objAddress :: "  +obj.Address)
                personal.hasAddress(obj.Address)
                result = true
                return
            }
        })
        return result
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
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "insertContact",
            "args" : [contactAdd , contactName]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            // 성공 메세지 처리
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