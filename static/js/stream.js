(function() {
    'use strict';

    document.addEventListener('DOMContentLoaded', function () {
        var elements = document.getElementsByClassName('stream');
        if (elements !== null) {
            for (var i = 0; i < elements.length; i++) {
                var element = elements[i];
                if (element.hasAttribute('data-socket')) {
                    createWebSocket(element.getAttribute('data-socket'), element);
                }
            }
        }
    });

    //////////

    function createWebSocket(uri, element) {

        console.debug('coucou');
        var webSocket = new WebSocket(uri, 'test');
        webSocket.onopen = function () {
            console.debug('[Stream] OnOpen');
        };
        webSocket.onmessage = function (event) {
            console.debug('[Stream] OnMessage');
            console.debug('[Stream] ', event.data);
            element.innerHTML = element.innerHTML + event.data;
        };
        webSocket.onclose = function () {
            console.debug('[Stream] OnClose');
             window.scrollTo(0, 0);
            location.reload();
        };
    }
})();
