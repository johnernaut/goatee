goatee
======

A Redis-backed notification server written in Go.

##Installation
`go get github.com/johnernaut/goatee`

`import "github.com/johnernaut/goatee"`

##Usage
**goatee** works by listening on a channel via [Redis Pub/Sub](http://redis.io/topics/pubsub) and then sending the received message to connected clients via [WebSockets](http://en.wikipedia.org/wiki/WebSocket).

### Configuration
**goatee** will look for a JSON configuration file in a `config` folder at the root of your project with the following names based on your environment: `development.json`, `production.json`, `etc`.  By default `config/development.json` will be used but you can also specify a `GO_ENV` environment variable and the name of that will be used instead.

```javascript
// example json configuration
// specify redis and websocket hosts
{
  "redis": {
    "host": "localhost:6379"
  },
  "web": {
    "host": "localhost:1235"
  }
}
```

======
#### Server
```go
package main

import (
    "github.com/johnernaut/goatee"
    "log"
)

func main() {
    err := goatee.CreateServer("achannel") // pass in the redis channel you'd like to subscribe to

    if err != nil {
        log.Fatal("Error: ", err.Error())
    }
}
```
========
#### Client
```html
<!doctype html>
<html>
  <head><title>goatee</title></head>
  <body>
    <h1>Websocket Messages:</h1>
    <ul id="ws">
    </ul>
  </body>
  <script src="http://code.jquery.com/jquery-1.10.1.min.js"></script>
  <script>
    var ws = new WebSocket("ws://localhost:1235/"); // based on the websocket host set in your config file
    ws.onopen = function() {
      console.log('opened!');
    }
    ws.onclose = function(e) {
      console.log('closed!' + e.code);
    }
    ws.onerror = function(error) {
      console.log("Error: " + error);
    }
    ws.onmessage = function(e) {
      $('#ws').append('<li>' + e.data + '</li>');
    }
  </script>
</html>
```

## Authors
- [johnernaut](https://github.com/johnernaut)
- you
