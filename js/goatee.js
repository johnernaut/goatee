'use strict';

var goatee = ('undefined' === typeof goatee ? {} : goatee);
goatee.VERSION = '0.1.0';

(function() {

    // socket instance
    goatee.socket = {};

    // goatee
    (function() {
        goatee.connect = function(host) {
            this.host = host;

            if (WebSocket != 'undefined') {
                 var socket = new goatee.Socket(this.host);
                 return socket;
            } else {
                console.log("Your browser doesn't support websockets.");
            }
        };
    })('object' === typeof goatee ? goatee : (this.goatee = {}), this);

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

        Socket.prototype.connect = function(host) {
            goatee.socket = new WebSocket(host);
            goatee.socket.onopen = function() {
                console.log('open');
            };
        };

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

})();
