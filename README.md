# merisbot-detect
Detecting meris-bot IPs

This code reads IPs from stdin and outputs all IPs that are meris bots.
It checks if they have tcp ports 2000 and 5678 open.

## required tools
apt-get install ipset iptables

## Example usage:

### Check existing file
cat test-ips | go run detectbot.go

### Continuously check a logfile, drop all bot traffic via iptables/ipset

ipset -N meris-ip iphash
iptables -A INPUT -m set --match-set meris-ip src -j DROP

tail -f /var/log/mylogfile | awk '${print $1} | go run detectbot.go | while read line; do sudo ipset add meris-ip $line; done

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
