# merisbot-detect
Detecting meris-bot IPs in realtime. Useful for keeping HTTP RPS down.
You should implement some rate-limiting anyways (e.g. https://www.haproxy.com/blog/four-examples-of-haproxy-rate-limiting/)

This code reads IPs from stdin and outputs only IPs that are meris bots (it checks if the IP has tcp ports 2000 and 5678 open).

Warning: this sends several SYN-packets per IP and tries establishing a full TCP handshake.
It could be optimized for production usage by just waiting for the SYN-ACK. But it's good enough for now. :) 

```
┌─────────────┐              ┌─────────────┐
│Meris Bot#1  ├────────────► │             │
└─────────────┘              │             │
                             │             │
┌─────────────┐              │             │                          ┌────────────┐
│Meris Bot#2  ├────────────► │detectbot.go ├────────────┬──────────►  │Meris Bot#1 │
└─────────────┘              │             │            │             └────────────┘
                             │             │            │
┌─────────────┐              │             │            │             ┌────────────┐
│Normal Client├────────────► │             │            └──────────►  │Meris Bot#2 │
└─────────────┘              └─────────────┘                          └────────────┘
```

## Required tools
- golang (tested with 1.17.1)
- ipset
- iptables

## Example usage:

### Check existing file
```
cat test-ips | go run detectbot.go
```

### Continuously check a logfile, drop all bot traffic via iptables/ipset
```
ipset -N meris-ip iphash
iptables -A INPUT -m set --match-set meris-ip src -j DROP
tail -f /var/log/mylogfile | awk '${print $1} | go run detectbot.go | while read line; do sudo ipset add meris-ip $line; done
```
