
function toDecimal(value) {
    return toBigNumber(value).toString(10)
};

var toBigNumber = function(number) {
    /*jshint maxcomplexity:5 */
    number = number || 0;
    if (isBigNumber(number))
        return number;

    if (isString(number) && (number.indexOf('0x') === 0 || number.indexOf('-0x') === 0)) {
        return new BigNumber(number.replace('0x',''), 16);
    }

    return new BigNumber(number.toString(10), 10);
};

var isBigNumber = function (object) {
    return object instanceof BigNumber ||
        (object && object.constructor && object.constructor.name === 'BigNumber');
};

var isString = function (object) {
    return typeof object === 'string' ||
        (object && object.constructor && object.constructor.name === 'String');
};

var toHex = function (val) {
    /*jshint maxcomplexity: 8 */

    if (isBoolean(val))
        return fromDecimal(+val);

    if (isBigNumber(val))
        return fromDecimal(val);

    if (typeof val === 'object')
        return fromUtf8(JSON.stringify(val));

    // if its a negative number, pass it through fromDecimal
    if (isString(val)) {
        if (val.indexOf('-0x') === 0)
            return fromDecimal(val);
        else if(val.indexOf('0x') === 0)
            return val;
        else if (!isFinite(val))
            return fromAscii(val);
    }

    return fromDecimal(val);
};

var isBoolean = function (object) {
    return typeof object === 'boolean';
};

var fromDecimal = function (value) {
    var number = toBigNumber(value);
    var result = number.toString(16);

    return number.lessThan(0) ? '-0x' + result.substr(1) : '0x' + result;
};

var fromUtf8 = function(str) {
    str = utf8.encode(str);
    var hex = "";
    for(var i = 0; i < str.length; i++) {
        var code = str.charCodeAt(i);
        if (code === 0)
            break;
        var n = code.toString(16);
        hex += n.length < 2 ? '0' + n : n;
    }

    return "0x" + hex;
};

var utf8 = {
    'version': '2.1.2',
    'encode': utf8encode,
    'decode': utf8decode
};

var fromUtf8 = function(str) {
    str = utf8.encode(str);
    var hex = "";
    for(var i = 0; i < str.length; i++) {
        var code = str.charCodeAt(i);
        if (code === 0)
            break;
        var n = code.toString(16);
        hex += n.length < 2 ? '0' + n : n;
    }

    return "0x" + hex;
};
