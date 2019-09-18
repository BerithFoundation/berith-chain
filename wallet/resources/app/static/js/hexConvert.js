let hexConvert = {
    getTxValue: function(value){

    var rs = {
        result : false,
        value : "",
        error : "",
    }

    try{
        var u = Math.ceil(value)
        var f = Math.floor(value);
        var d = value % f;
        if( f == 0 && u > 0){
            d = 1
        }

        var txValue = "";

        if(d > 0){
            var dot = value.toString().replace(f.toString() + ".", "");

            var len = dot.length;

            var temp = value.toString().replace(".", "");
            for(var i = 0; i< 18 - len; i++){
                temp += "0";
            }
            txValue = BigInt(temp).toString(16);

        } else {

            txValue = BigInt(f.toString() + "000000000000000000").toString(16);

        }
        // console.log("txValue ::: " + txValue)
        rs.result = true;
        rs.value = txValue;

        return rs;
    }
    catch(err){
        rs.result = false;
        rs.error = err;
        return rs
    }
    },
    pad :function (n, width) {
        n = n + '';
        return n.length >= width ? n : new Array(width - n.length + 1).join('0') + n;
    },

    HexToValueString: function (hex){
        var t = hexToDec(hex);

        var big = BigInt(t.toString());

        var t2 = big.dividedToIntegerBy("1e18"); //0

        var m = big.mod("1e18"); //127

        return t2 + "." + this.pad(m.toString(), 18);
    }
}
