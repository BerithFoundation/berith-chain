let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', async function() {
            index.listen();

            let responseValue = await sendMessage("init", "", [])
            console.log(responseValue)
            onPageload(); // all pages should include this method to guarantee a successful page load.
            includeHTML(); // for templating
        });
    },
    listen: function() {
        //폴링 리시브등록?
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "notify_show":
                    asticode.loader.show();
                    break;
                case "notify_hide":
                    asticode.loader.hide();
                    break;
                case "polling":
                    // $('#polling').val(message.payload)
                    // syncingData(message.payload)
                    // syncingData(message.payload)
                    // console.log("footer ::: " + message.payload)
                    break;
                case "coinbase":
                    berith.coinbase()
                    break;
            }
        });
    },
    nextPage : function () {
        location.href="http://google.com"
    },
}
