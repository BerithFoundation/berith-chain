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
        message.payload = {
            "api" : "selectMember",
            "args" : ["aa"]
        }
        astilectron.sendMessage(message , function (message) {
            var obj = message.payload
            console.log("array :: " + Array.isArray(obj))
            obj.forEach((v) => {
                console.log("Address :: " +v.Address)
                console.log("ID :: " +v.ID)
                console.log("Password :: " +v.Password)
            })
        })
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
            // 성공 메세지 처리
        })
    },

}