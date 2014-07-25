(function () {
    "use strict";
    
    function goatee(key) {
        validateKey(key);
        
        this.event_emitter = new goatee.CommandDispatcher();
    }
    
    var prototype = goatee.prototype;
    
    prototype.bind = function(name, callback) {
        this.event_emitter.bind(name, callback);
        return this;
    };
    
    prototype.emit = function(name, data) {
        this.event_emitter.emit(name, data);
        return this;
    };
    
    function validateKey(key) {
        if (key === null || key === undefined)
            console.warn("An API key must be supplied.");
    }
    
    this.goatee = goatee;
}).call(this);

(function () {
    "use strict";

    function CommandDispatcher() {
        this.callbacks = {};
    }

    var prototype = CommandDispatcher.prototype;

    prototype.bind = function (name, callback, ctx) {
        this.callbacks[name] = this.callbacks[name] || [];
        this.callbacks[name].push({
            fn: callback,
            ctx: ctx
        });

        return this;
    };

    prototype.emit = function (eventName, data) {
        var i;
        var callbacks = this.callbacks[eventName];
        
        if (callbacks && callbacks.length > 0) {
            for (i = 0; i < callbacks.length; i++) {
                callbacks[i].fn.call(callbacks[i].ctx || window, data);
            }
        }
        
        return this;
    };
    
    goatee.CommandDispatcher = CommandDispatcher;
}).call(this);

var gt = new goatee("ABC123LOLFU");

gt.bind('msg', function(data) {
    console.log(data);
});

gt.emit('msg', 'some rly kewl txt lmao');
