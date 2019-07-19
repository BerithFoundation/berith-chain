let qrcode = {
    qrgen :  function(qrtext) {
        console.log("qrtext :: " + qrtext)
        var qrcode = new QRCode(document.getElementById("qrcode"), {
            text : qrtext,
            width: 128,
            height: 128,
            colorDark : "#000000",
            colorLight : "#ffffff",
            correctLevel : QRCode.CorrectLevel.H
        });
        $("#qrcode > img").css({"margin":"auto"});
        $("#qrcode > img").attr("id", "qrImg")
        // $("#qrcode > canvas").attr("id", "qrImg")
    }
}