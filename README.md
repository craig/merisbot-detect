# merisbot-detect
Detecting meris-bot IPs in realtime. Useful for keeping HTTP RPS down. This is a pretty stupid idea, but I was feeling like coding something in go. :)
You should instead implement some rate-limiting (e.g. https://www.haproxy.com/blog/four-examples-of-haproxy-rate-limiting/).

This code reads IPs from stdin and outputs only IPs that are meris bots (it checks if the IP has tcp ports 2000 and 5678 open).

Warning: this sends several SYN-packets per IP and tries establishing a full TCP connection; it could be optimized by just waiting for the SYN-ACK. 
Make sure to do something like piping your logfile through uniq so you don't retry known good IPs repeatedly.

It works like a filtering pipe, extracting only the Meris Bot IPs. Input format is one line per IP + linebreak.

```
INPUT (stdin)                                                          OUTPUT (stdout)

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

## Example usage

### Check existing file
```
cat test-ips | go run detectbot.go
```

### Check a logfile, drop all bot traffic via iptables/ipset
```
ipset -N meris-ip iphash
iptables -A INPUT -m set --match-set meris-ip src -j DROP
tail -n 1000 /var/log/mylogfile | awk '${print $1} | sort | uniq | go run detectbot.go | while read line; do sudo ipset add meris-ip $line; done
```
