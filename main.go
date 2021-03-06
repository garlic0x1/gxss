package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v2"
)

type Result struct {
	Type    string
	Message string
	Kill    chan bool
}

type Input struct {
	URL  string   `yaml:"URL"`
	Keys []string `yaml:"Keys"`
}

type Context struct {
	Type         string
	URL          string
	Key          string
	Selector     string
	Prefix       string
	Tag          string
	Attr         string
	Quote        string
	OpenBracket  string
	CloseBracket string
	BackSlash    string
}

type Handler struct {
	Payloads    []string
	Browsers    []string
	Interaction bool
}

var (
	Canary    = "zzx%dqyj"
	Canary2   = "zzxqyj"
	Canary3   = "zzx%sqyj"
	Queue     = make(chan Input)
	Results   = make(chan Result)
	Interact  bool
	Debug     bool
	ShowType  bool
	Stop      int
	Wait      int
	TagMap    map[string]map[string]Handler
	Payloads  map[string][]string
	ChromeCtx context.Context
)

func reader() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var input Input
		err := yaml.Unmarshal([]byte(s.Text()), &input)
		if err != nil {
			log.Println("Error parsing JSON/YAML", err)
			continue
		}
		Queue <- input
	}
	close(Queue)
}

func writer(sevLimit *int) {
	severity := map[string]int{
		"low":       4,
		"medium":    3,
		"high":      2,
		"confirmed": 1,
	}

	for res := range Results {
		if Stop > 0 && severity[res.Type] <= Stop {
			go func() {
				res.Kill <- true
			}()
		}
		if severity[res.Type] <= *sevLimit && isUniqueOutput(res) {
			if ShowType {
				fmt.Println("["+res.Type+"]", res.Message)
			} else {
				fmt.Println(res.Message)
			}
		}
	}
}

func spawnWorkers(n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			tab, cancel := chromedp.NewContext(ChromeCtx)
			defer cancel()

			kill := make(chan bool, 1)

			// handler that runs when alert pops
			chromedp.ListenTarget(tab, func(ev interface{}) {
				if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
					go func() {
						var u string
						if err := chromedp.Run(tab,
							page.HandleJavaScriptDialog(false),
							chromedp.Location(&u),
						); err != nil && Debug {
							log.Println(err)
						}
						Results <- Result{Type: "confirmed", Message: u, Kill: kill}
					}()
				}
			})

			for input := range Queue {
				live := make(chan bool, 1)
				go func() {
					worker(input, tab, kill)
					live <- true
				}()
				select {
				case _ = <-live:
				case _ = <-kill:
				}
			}
		}()
	}
	wg.Wait()
	close(Results)
}

func main() {
	threads := flag.Int("t", 8, "Number of threads to use.")
	sevLimit := flag.Int("sev", 4, "Filter by severity. (1 is a confirmed alert, 2-4 are high-low.)")
	showType := flag.Bool("s", false, "Show result type.")
	interact := flag.Bool("i", false, "Try to perform handler to trigger payload.")
	showErrors := flag.Bool("debug", false, "Display errors.")
	stop := flag.Int("stop", 0, "Stop on first xss of specified priority. (1 is a confirmed alert, 2-4 are high-low.)")
	payloads := flag.String("p", "./payloads.yaml", "YAML file of escape patterns and xss payloads.")
	proxy := flag.String(("proxy"), "", "Proxy URL. Example: -proxy http://127.0.0.1:8080")
	swait := flag.Int("wait", 0, "Seconds to wait on page after loading in chrome mode. (Use to wait for AJAX reqs)")
	debugChrome := flag.Bool("debug-chrome", false, "Don't use headless. (slow but fun to watch)")
	flag.Parse()
	Debug = *showErrors
	Interact = *interact
	Stop = *stop
	Wait = *swait
	ShowType = *showType
	parsePayloads(*payloads, "tagmap.yaml")

	// check for stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No input detected, use `cat urls.txt | url-miner -w wordlist.txt`")
		os.Exit(1)
	}

	// start browser
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(*proxy),
		chromedp.Flag("headless", !(*debugChrome)))...)

	ChromeCtx = ctx
	defer cancel()

	// these each finish the next when done, finishing the program
	go reader()
	go spawnWorkers(*threads)
	writer(sevLimit)
}
