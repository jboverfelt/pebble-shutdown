# pebble-shutdown

Listens for a predefined sequence of pebble commands and
sends a message on a predefined MQTT topic.

The host and the port of the MQTT broker can be specified like so:

```
./pebble-shutdown -host 0.0.0.0:1883
```

The sample above shows the default address if none is specified
