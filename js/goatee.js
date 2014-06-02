'use strict';

(function() {
    // goatee
    function goatee(host) {
        this.host = host;

        if (WebSocket != 'undefined') {
             var socket = new goatee.Socket(this.host);
             return socket;
        } else {
            console.log("Your browser doesn't support websockets.");
        }
    };

    // goatee.util
    (function() {
        var util = goatee.util = {};

        util.merge = function(target, additional, deep, lastseen) {
            var seen = lastseen || [],
                depth = typeof deep == 'undefined' ? 2 : deep,
                prop;

            for (prop in additional) {
                if (additional.hasOwnProperty(prop) && util.indexOf(seen, prop) < 0) {
                    if (typeof target[prop] !== 'object' || !depth) {
                        target[prop] = additional[prop];
                        seen.push(additional[prop]);
                    } else {
                        util.merge(target[prop], additional[prop], depth - 1, seen);
                    }
                }
            }

            return target;
        }

        util.indexOf = function(arr, o, i) {
            for (var j = arr.length, i = i < 0 ? i + j < 0 ? 0 : i + j : i || 0;
                i < j && arr[i] !== o; i++) {}

            return j <= i ? -1 : i;
        };
    })('undefined' != typeof goatee ? goatee : (this.goatee = {}), this);

    // goatee.Socket
    (function() {
        goatee.Socket = Socket;

        function Socket(host) {
            this.socket = {};
            this.socket.date = new Date();
            this.socket.action = "bind";
            this.connect(host);
        };

        // Connect to a websocket on the given host
        Socket.prototype.connect = function(host) {
            goatee.socket = new WebSocket(host);
            goatee.socket.onopen = function() {
                console.log('open');
            };
        };

        // Bind to the websocket channel and listen for incoming messages.
        // Make sure the socket has finished connecting first.
        Socket.prototype.bind = function(channel, fn) {
            this.socket.channel = channel;
            var that = this;

            var t = setTimeout(function() {
                if (goatee.socket.readyState === 1) {
                    goatee.socket.send(JSON.stringify(that.socket));
                    clearTimeout(t);
                } else {
                    Socket.bind(channel, fn);
                }
            }, 50);

            goatee.socket.onmessage = function(data) {
                fn(data);
            };
        };

        // Send data to the given websocket.
        Socket.prototype.send = function(data) {
            this.socket.action = "message";
            goatee.util.merge(this.socket, data);

            if (goatee.socket.readyState !== 1) {
                console.log('You must be bound to a channel to send data to the server.');
                return;
            }

            goatee.socket.send(JSON.stringify(this.socket));
        };
    })('undefined' != typeof goatee ? goatee : (this.goatee = {}), this);

    this.goatee = goatee;
}).call(this);

(function() {
    // Current version
    goatee.VERSION = '0.1.0';

    // Established socket instances
    goatee.sockets = [];
}).call(this);