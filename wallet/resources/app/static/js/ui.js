
$(document).ready(function() {
    registerEvents();
});

function registerEvents() {
    $(".inp").focus(function () {
        //$(this).parent().parent().parent().parent().addClass("active del");
        $(this).closest('.inp_group').addClass("active del");
    });

    $(".inp").blur(function () {
        //$(this).parent().parent().parent().parent().removeClass("active del");
        $(this).closest('.inp_group').removeClass("active del");
    });
}

function togglePopUp(popUpId) {
    let selector = "#" + popUpId;
    if ($(selector).hasClass("view")) {
        $(selector).removeClass("view");
        $(selector).addClass("hide");
    } else {
        $(selector).removeClass("hide");
        $(selector).addClass("view");
    }
}



