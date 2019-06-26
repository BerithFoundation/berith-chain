let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
        })

        let message = {"name": "init"};
        asticode.loader.show()
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();
        })
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
                    $('#polling').val(message.payload)
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
};