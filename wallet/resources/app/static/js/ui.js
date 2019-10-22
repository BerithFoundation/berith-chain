
$(document).ready(function() {
    registerEvents();
});

function registerEvents() {
    $(".inp").focus(function () {
        //$(this).parent().parent().parent().parent().addClass("active del");
        $(this).closest('.inp_group').addClass("active");
    });

    $(".inp").blur(function () {
        //$(this).parent().parent().parent().parent().removeClass("active del");
        $(this).closest('.inp_group').removeClass("active");
    });
    $(".inp").on("propertychange change keyup paste input", function() {
        var currentVal = $(this).val();
        if(currentVal == ""){
            $(this).closest('.inp_group').removeClass("del");
        }else{
            $(this).closest('.inp_group').addClass("del");
        }
    });
    $('.del').click(function () {
        $(this).siblings('input').val("")
        $(this).closest('.inp_group').removeClass("del");
    });
    $('.icon').click(function () {
        // console.log( "icon click !! ")
        if( $(this).hasClass("hide_word")){
            $(this).removeClass("hide_word")
            $(this).addClass("view_word")
            $(this).siblings('input').prop("type", "text");
        }else {
            $(this).removeClass("view_word")
            $(this).addClass("hide_word")
            $(this).siblings('input').prop("type", "password");
        }
    })
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



