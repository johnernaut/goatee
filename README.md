goatee
======

A Redis-backed notification server written in Go.

[![Build Status](https://travis-ci.org/johnernaut/goatee.png?branch=master)](https://travis-ci.org/johnernaut/goatee)

**Note:** This project is *alphaware*.  For the time being, it's API is bound to change as features are continually added and enhancements are made.

##Installation
`go get github.com/johnernaut/goatee`

`import "github.com/johnernaut/goatee"`

##Usage
**goatee** works by listening on a channel via [Redis Pub/Sub](http://redis.io/topics/pubsub) and then sending the received message to connected clients via [WebSockets](http://en.wikipedia.org/wiki/WebSocket) and has fallback support for **long polling**.

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
    // subscribe to one or many redis channels
    err := goatee.CreateServer([]string{"achannel", "anotherchannel"})

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
========
#### Redis
With **goatee** running and your web browser connected to the socket, you should now be able to test message sending from Redis to your client (browser).  Run `redis-cli` and publish a message to the channel you subscribed to in your Go server.  E.x. `publish 'mychannel' 'mymessage'`

## Tests
`go test github.com/johnernaut/goatee`

## Authors
- [johnernaut](https://github.com/johnernaut)
- you
