# merisbot-detect
Detecting meris-bot IPs in realtime. Useful for keeping HTTP RPS down.
You should nevertheless implement some rate-limiting (e.g. https://www.haproxy.com/blog/four-examples-of-haproxy-rate-limiting/).

This code reads IPs from stdin and outputs only IPs that are meris bots (it checks if the IP has tcp ports 2000 and 5678 open).

Warning: this sends several SYN-packets per IP and tries establishing a full TCP connection.
It could be optimized for production usage by just waiting for the SYN-ACK. But it's good enough for now. :)
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


## Q & A

Some very important people from radware (https://twitter.com/Radware) from Cyber Threat Intelligence (https://twitter.com/geenensp) & their Head of Research (https://twitter.com/hypoweb) have been very interested in this code, so I'm doing a Q&A!

**Q:** Is this production-ready code?

**A:**: That might depend on your production... (It should be pretty obvious it's not, it's my first Go program I wrote for a little fun!)


**Q:** From Pascal Geenens (https://twitter.com/geenensp/status/1437032163057311745)
"I would advise against this tool. It's blocking all MikroTek routers in the world. Not every MikroTek router is an active bot. This tool will generate to many false positives. Consider taking this down, people might get in trouble because of this!"

**A:** It's blocking all systems that have port 2000 and port 5678 open. That's the intended purpose. Only use when under attack. Don't run random bullshit you find on the internet you don't understand. Mikrotek has a marketshare of 0.1% - it's better to be unavailable for some people than being down for everyone (not every Mikrotek user in the world is accessing your site during the DDoS, so please factor that in when calculating things).


**Q:** From Daniel Smith (https://twitter.com/hypoweb/status/1437411192088055808)
"Slippery slope to blindly block every device just because a threat exists. But then again no one in this industry wants to do any real work. Why invest time in being a professional when you can just block and censor everyone because a few bad apples exist?"

**A:** I'm not sure what Daniel means, as I don't work in the Infosec industry. I would surely appreciate if someone took down meris.


**Q:** From Daniel Smith (https://twitter.com/hypoweb/status/1437419028142690312)
In my opinion, what you did is a total joke. Real victims will need a preemptive solution.  Or wait, DDoS'er never comes back, right?
 
**A:** This specifically aims to protects against repeated attacks. If Daniel wants a proactive solution he can scan the whole internet for Mikrotik devices and then block those premptively when a DDoS starts.


**Q:** From Pascal Geenens (https://twitter.com/geenensp/status/1437403030568251403)
"Start by caching those IPs that you tested. In current state I can use your DDoS mitigation to attack any random IP by sending you a flood of spoofed TCP SYN packets. For each SYN packet you will create 2 TCP sessions that timeout after 1sec."

**A:** As the README.md says, it "tries establishing a full TCP connection". It cannot be used for reflection (and would be a lot less ineffective than using any random UDP resolver for reflection). Caching non-Mikrotik IPs is a good idea in order not to retest them. One should make sure to do some bash "uniq" magic or could add that to the code.


**Q:** From Pascal Geenens (https://twitter.com/geenensp/status/1437404468270731273)
"I don’t think your script will mitigate when under attack as it will fail running out of memory caused by too many goroutines that are blocked on the 1sec socket timeout."

**A:** Pascal has not tested the script. Pascal doesn't know I have proper hardware with 64 Cores/128GB RAM available and it's really hard to exhaust that with a list of IP addresses and very tiny go code.


**Q:** From Pascal Geenens (https://twitter.com/geenensp/status/1437420528210481153)
So you will active this one you detected your server is suffering. You did notice the attacks last no longer than 2mins, yes?
 
**A:** Yes, that's why I recommended continuously looking at logfiles. The more IPs you block via ipset, the faster the server should become again (unless the attacker uses something else than the Layer 7 HTTP-Pipeling attack...)


**Q:** From Pascal Greens (https://twitter.com/geenensp/status/1437405203389628424)
"And lastly, you are not annoying the people that manage or buy Mikrotik routers, they will say their router forwarded the packets but it is the other side that refused the connection."

**A:** Well it is indeed true and they are right saying that. They won't be able to access the website, though - so my educated guess is, they will probably be a bit annoyed? 


**Q:** From Pascal Geenens (again!!!) (https://twitter.com/geenensp/status/1437434501479735297)
"You really want to filter MikroTik, then change your program to check for response  \x01\x00\x00\x00 from port 2000 and take out port 5678 check.  Ref: https://github.com/samm-git/btest-opensource"

**A:** You said "lastly" before (and some other unfriendly stuff https://twitter.com/geenensp/status/1437422071059976201). But this is not a bad idea actually! Please send a PR. :)


**Q:** Why is the Q&A longer than the actual code?

**A:** I guess I got carried aways with this while watching a TV show. Also, some people with clever titles and clever remarks on twitter can't read code, which was fun.
