package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Result struct {
	Type    string
	Message string
}

type Input struct {
	URL  string
	Keys []string
}

type Context struct {
	Type  string
	URL   string
	Key   string
	Build string
}

var (
	Canary       = "zzx%dqyj"
	Canary2      = "zzxqyj"
	Canary3      = "zzx%sqyj"
	Queue        = make(chan Input)
	Results      = make(chan Result)
	Confirm      = make(chan string)
	Stop         bool
	ShowType     bool
	Wait         int
	Payloads     map[string]map[string]map[string][]string
	AttrPayloads map[string][]string
	ChromeCtx    context.Context
)

func reader() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var input Input
		err := json.Unmarshal([]byte(s.Text()), &input)
		if err != nil {
			log.Println("Error parsing JSON", err)
			continue
		}
		Queue <- input
	}
	close(Queue)
}

func writer() {
	for res := range Results {
		if ShowType {
			fmt.Println("["+res.Type+"]", res.Message)
		} else {
			fmt.Println(res.Message)
		}
	}
}

func spawnConfirmers(n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)

		// confirmer
		go func() {
			defer wg.Done()
			tab, cancel := chromedp.NewContext(ChromeCtx)

			alert := make(chan bool, 1)
			// attach handler to javascript:alert() for xss confirmation
			chromedp.ListenTarget(tab, func(ev interface{}) {
				if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
					alert <- true
					go func() {
						if err := chromedp.Run(tab,
							page.HandleJavaScriptDialog(true),
						); err != nil {
							log.Println(err)
						}
					}()
				}
			})

			for msg := range Confirm {
				verifyScript(msg, tab, alert)
			}
			cancel()
		}()
	}
	wg.Wait()
	log.Println("Done confirming payloads.")
	close(Results)
}

func spawnWorkers(n int) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)

		// generator
		go func() {
			defer wg.Done()
			tab, cancel := chromedp.NewContext(ChromeCtx)
			defer cancel()

			chromedp.ListenTarget(tab, func(ev interface{}) {
				if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
					go func() {
						if err := chromedp.Run(tab,
							page.HandleJavaScriptDialog(true),
						); err != nil {
							panic(err)
						}
					}()
				}
			})

			for input := range Queue {
				worker(input, tab)
			}
		}()
	}
	wg.Wait()
	log.Println("Done generating payloads.")
	close(Confirm)
}

func main() {
	threads := flag.Int("t", 8, "Number of threads to use.")
	showType := flag.Bool("s", false, "Show result type.")
	stop := flag.Bool("stop", false, "Stop on first confirmed xss.")
	payloads := flag.String("p", "./payloads.yaml", "YAML file of escape patterns and xss payloads.")
	proxy := flag.String(("proxy"), "", "Proxy URL. Example: -proxy http://127.0.0.1:8080")
	swait := flag.Int("wait", 0, "Seconds to wait on page after loading in chrome mode. (Use to wait for AJAX reqs)")
	debugChrome := flag.Bool("debug-chrome", false, "Don't use headless. (slow but fun to watch)")
	flag.Parse()
	Stop = *stop
	Wait = *swait
	ShowType = *showType
	parsePayloads(*payloads)

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
	go spawnConfirmers(*threads)
	writer()
}
