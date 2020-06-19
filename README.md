<!--
SPDX-FileCopyrightText: 2020 Ethel Morgan

SPDX-License-Identifier: CC0-1.0
-->

# Catbus Wake-On-LAN

A simple daemon that emits [Wake-On-LAN "magic packets"](https://en.wikipedia.org/wiki/Wake-on-LAN) when triggered by [Catbus](https://ethulhu.co.uk/catbus).

## Config

```json
{
  "mqttBroker": "tcp://broker.local:1883",
  "devices": {
    "TV": {
      "mac": "aa:bb:cc:dd:ee:ff",
      "topic": "home/living-room/tv/power"
    }
  }
}
```

## Wake-On-LAN

Wake-On-LAN is a protocol to wake or boot devices over LAN using "magic packets".
The magic packet is, in big endian:

- 6 bytes of `0xFF`.
- the MAC address of the device you are waking, 16 times.

For example, for a device with a MAC address `aa:bb:cc:dd:ee:ff`, the packet will be, in hex:

```
ffffffffffffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeffaabbccddeeff
```

This is broadcast over UDP on port 9.
