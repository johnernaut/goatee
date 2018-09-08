goatee
======

A Redis-backed notification server written in Go.

[![Build Status](https://travis-ci.org/johnernaut/goatee.png?branch=master)](https://travis-ci.org/johnernaut/goatee)

**Client library:** [goatee.js](https://github.com/johnernaut/goatee.js)

## Installation
`go get github.com/johnernaut/goatee`

`import "github.com/johnernaut/goatee"`

## Usage
**goatee** works by listening on a channel via [Redis Pub/Sub](http://redis.io/topics/pubsub) and then sending the received message to connected clients via [WebSockets](http://en.wikipedia.org/wiki/WebSocket).  Clients may create channels to listen on by using the [goatee client library](https://github.com/johnernaut/goatee.js).

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

### Server
```go
package main

import (
    "github.com/johnernaut/goatee"
    "log"
)

func main() {
    // subscribe to one or many redis channels
    err := goatee.CreateServer()

    if err != nil {
        log.Fatal("Error: ", err.Error())
    }
}
```

### Client
An example of how to use the [goatee client library](https://github.com/johnernaut/goatee.js) can be found in the `examples` folder.


### Redis
With **goatee** running and your web browser connected to the socket, you should now be able to test message sending from Redis to your client (browser).  Run `redis-cli` and publish a message to the channel you subscribed to in your Go server.  By default, **goatee** expects your Redis messages to have a specified JSON format to send to the client with the following details:
* `payload`
* `created_at (optional)`

E.x. `publish 'mychannel' '{"payload": "mymessage which is a string, etc."}'`

## Tests
`go test github.com/johnernaut/goatee`

## Authors
- [johnernaut](https://github.com/johnernaut)
- you
