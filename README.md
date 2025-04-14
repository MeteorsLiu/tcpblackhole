# A TCP test tool
Simple tool for debugging with TCP.

## Howto

By default, `-addr` default to `127.0.0.1`, `port` default to `9999`

### Echo 

Setup a echo server.

Echo server always reply with your sent data.

```shell
tcpblackhole -mode "echo"
```

### Blackhole

Setup a blackhole server

Blackhole server blackholes all the incoming traffic, and send data filled with zero to client if there's someone reading.


```shell
tcpblackhole -mode "relay"
```

### Relay

Setup a relay server

Relay server forward all income traffic to `endpoint` and reply with its response.


```shell
tcpblackhole -mode "blackhole" -endpoint 127.0.0.1:1234
```
