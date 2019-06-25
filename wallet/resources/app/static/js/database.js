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
    selectMember : function () {
        let message = {"name" : "callDB"}
        message.payload = {
            "api" : "selectMember",
            "args" : ["aa"]
        }
        asticode.loader.show()
        astilectron.sendMessage(message , function (message) {
            asticode.loader.hide()
            console.log()
            // var obj = message.payload
            // var keys = Object.keys(obj);
            // var contents = ''
            // $('#contactData').empty()
            // for ( var i in keys) {
            //     contents += '<tr>'
            //     console.log("key="+keys[i]+ ",  data="+ obj[keys[i]]);
            //     contents += '<td><input type="text" value="'+keys[i]+'"></td>'
            //     contents += '<td><input type="text" value="'+obj[keys[i]]+'"></td>'
            //     contents += '</tr>'
            // }
            // $('#memberData').append(contents)
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