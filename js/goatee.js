/*!
 * goatee.js
 *
 * Inspired heavily by the Pusher
 * client library.
 */

(function () {
    "use strict";

    function goatee(key, url) {
        validateKey(key);
        var self = this;

        this.channel = null;

        this.event_emitter = new goatee.CommandDispatcher();
        this.connection = new goatee.ConnectionManager(url);

        this.connection.bind('connected', function() {
            self.connected = true;
            console.log("Connected to goatee server.");
        });

        this.connection.bind('message', function(data) {
            console.log(self.channel);
            self.event_emitter.emit(self.channel, data);
        });

        this.connection.bind('error', function(err) {
            console.warn("Error connecting to goatee server: " + err);
        });

        this.connection.bind('closed', function() {
            console.warn("Connection to the goatee server was closed.")
        });

        this.connect();
    }

    var prototype = goatee.prototype;

    goatee.connected = false;

    prototype.bind = function(name, callback) {
        this.channel = name;
        this.connection.subscribe(name);
        this.event_emitter.bind(name, callback);
        return this;
    };

    prototype.emit = function(name, data) {
        this.connection.send_event(name, data);
        return this;
    };

    prototype.connect = function() {
        this.connection.connect();
    };

    function validateKey(key) {
        if (key === null || key === undefined)
            console.warn("An API key must be supplied.");
    }

    this.goatee = goatee;
}).call(this);

/**** Command Dispatcher ****/
(function () {
    "use strict";

    function CommandDispatcher() {
        this.callbacks = new CallStore();
    }

    var prototype = CommandDispatcher.prototype;

    prototype.bind = function (name, callback, ctx) {
        this.callbacks.add(name, callback, ctx);

        return this;
    };

    prototype.emit = function (eventName, data) {
        var i;
        var callbacks = this.callbacks.get(eventName);

        if (callbacks && callbacks.length > 0) {
            for (i = 0; i < callbacks.length; i++) {
                callbacks[i].fn.call(callbacks[i].ctx || window, data);
            }
        }

        return this;
    };

    // fix readonly error for callbacks...
    function CallStore() {
      this._callbacks = {};
    }

    CallStore.prototype.get = function(name) {
      return this._callbacks[name];
    };

    CallStore.prototype.add = function(name, callback, ctx) {
        this._callbacks[name] = this._callbacks[name] || [];
        this._callbacks[name].push({
          fn: callback,
          context: ctx
        });
    };

    goatee.CommandDispatcher = CommandDispatcher;
}).call(this);

/**** Util ****/
(function() {
  "use strict";

  goatee.Utils = {
    extend: function(proto) {
      for (var i = 1; i < arguments.length; i++) {
        var extensions = arguments[i];
        for (var property in extensions) {
          if (extensions[property] && extensions[property].constructor &&
              extensions[property].constructor === Object) {
            proto[property] = goatee.Utils.extend(
              proto[property] || {}, extensions[property]
            );
          } else {
            proto[property] = extensions[property];
          }
        }
      }

      return proto;
    }
  };
}).call(this);

/**** Connection Manager ****/
(function() {
  "use strict";

  function ConnectionManager(url) {
      goatee.CommandDispatcher.call(this);

      this.url = url;
      this.connection = null;
      this.state = "initialized";
  }

  var prototype = ConnectionManager.prototype;
  goatee.Utils.extend(prototype, goatee.CommandDispatcher.prototype);

  prototype.connect = function() {
    var self = this;
    var compatible = this.checkCompatibility();

    if (compatible) {
        this.connection = new WebSocket(this.url);
        this.connection.onopen = function() { self.onOpen(); };
        this.connection.onmessage = function(data) { self.onMessage(data); };
        this.connection.onerror = function(err) { self.onError(err); };
        this.connection.onclose = function(err) { self.onClose(); };
    }
  };

  prototype.subscribe = function(name) {
      var self = this;

      this.waitForConnection(function() {
          var data = {
            channel: name,
            action: 'bind',
            token: 'ABC123'
          };

          self.connection.send(JSON.stringify(data));
      });
  };

  prototype.onOpen = function() {
      this.state = 'connected';
      this.emit('connected');
  };

  prototype.onError = function(err) {
      this.state = 'error';
      this.emit('error', err);
  };

  prototype.onClose = function() {
      this.state = 'closed';
      this.emit('closed');
  };

  prototype.onMessage = function(data) {
      this.emit('message', data);
  };

  prototype.send_event = function(name, data) {
      if (this.state !== 'connected') {
        this.emit('error', 'Not connected to the goatee server.');
        return;
      }

      var data = {
          channel: name,
          payload: data,
          action: 'message',
          token: 'ABC123'
      };

      this.connection.send(JSON.stringify(data));
  }

  // used to make sure there's a connection before a send is made
  prototype.waitForConnection = function(callback) {
      var self = this;

      setTimeout(function() {
          if (self.state === 'connected') {
              if (callback !== null) {
                callback();
              }
              return;
          } else {
            self.waitForConnection(callback);
          }
      }, 5);
  };

  prototype.checkCompatibility = function() {
      if (window.WebSocket === undefined) {
        this.emit('error', "WebSockets aren't supported on this browser.")
        return false;
      }

      return true;
  };

  goatee.ConnectionManager = ConnectionManager;
}).call(this);
