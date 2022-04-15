package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func worker(input Input, ctx context.Context) {
	// identify context
	contexts := identifyCtx(input, ctx)

	// escape context
	for _, context := range contexts {
		switch {
		case context.Type == "html":
			// see if we can use angle brackets
			_, ok := testBrackets(context, ctx)
			if !ok {
				log.Println("Cannot reflect brackets")
				return
			}
			// search for payloads
			for tag, _ := range Payloads {
				// find a valid tag
				if tryTag(tag, context, ctx) {
					for handler, pocs := range Payloads[tag] {
						// find a valid handler
						if tryHandler(handler, tag, context, ctx) {
							for _, poc := range pocs["payloads"] {
								message := context.URL + "?" + context.Key + "=" + url.QueryEscape(poc)
								// verify payload
								if verifyScript(message, ctx) {
									Results <- Result{Type: "html", Message: message}
									if Stop {
										return
									}
								}
							}
						}
					}
				}
			}
		case context.Type == "attr":
			// determine key
			// break context
		case context.Type == "href":
		case context.Type == "script":
		case context.Type == "style":
		}
	}
}

func verifyScript(u string, ctx context.Context) bool {
	// perform chrome request
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
	)
	if err != nil {
		log.Println(u, err)
	}

	select {
	case ret := <-Alert:
		return ret
	default:
		return false
	}
}

func testBrackets(context Context, ctx context.Context) (string, bool) {
	str := url.QueryEscape(fmt.Sprintf(Canary3, "<>"))
	u := fmt.Sprintf("%s?%s=%s", context.URL, context.Key, str)
	var document string

	// perform chrome request
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
		chromedp.OuterHTML(`html`, &document),
	)
	if err != nil {
		log.Println(u, err)
	}

	// analyze response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		log.Println(err)
	}

	if strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<>")) {
		return "<>", true
	}
	if strings.Contains(document, str) {
		return "<>", true
	}
	return "", false
}

func tryHandler(handler string, tag string, context Context, ctx context.Context) bool {
	str := fmt.Sprintf("<%s %s=%s>", tag, handler, Canary2)
	u := fmt.Sprintf("%s?%s=%s", context.URL, context.Key, str)
	var document string

	// perform chrome request
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
		chromedp.OuterHTML(`html`, &document),
	)
	if err != nil {
		log.Println(u, err)
	}

	// analyze response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		log.Println(err)
	}

	nreflections := doc.Find(fmt.Sprintf("%s[%s=%s]", tag, handler, Canary2)).Length()
	if nreflections > 0 {
		return true
	}
	return false
}

func tryTag(tag string, context Context, ctx context.Context) bool {
	c1 := make(chan bool, 1)

	go func() {
		str := fmt.Sprintf("<%s %s=1>", tag, Canary2)
		u := fmt.Sprintf("%s?%s=%s", context.URL, context.Key, str)
		var document string

		// perform chrome request
		err := chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.Sleep(time.Duration(Wait)),
			chromedp.OuterHTML(`html`, &document),
		)
		if err != nil {
			log.Println(u, err)
		}

		// analyze response
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
		if err != nil {
			log.Println(err)
		}

		nreflections := doc.Find(fmt.Sprintf("%s[%s]", tag, Canary2)).Length()

		c1 <- (nreflections > 0)
	}()

	select {
	case result := <-c1:
		return result
	case <-time.After(time.Duration(5) * time.Second):
		return false
	}
}

func identifyCtx(input Input, ctx context.Context) []Context {
	var contexts []Context
	var document string

	u := buildPayload(input)

	// perform chrome request
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
		chromedp.OuterHTML(`html`, &document),
	)
	if err != nil {
		log.Println(err)
	}

	// analyze response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		log.Println(err)
	}

	doc.Find("*").Each(func(_ int, node *goquery.Selection) {
		n := node.Get(0)
		for _, attr := range n.Attr {
			for i, key := range input.Keys {
				if strings.Contains(attr.Val, fmt.Sprintf(Canary, i)) {
					contexts = append(contexts, Context{Type: "attr", URL: input.URL, Key: key})
				}
			}
		}
	})

	for i, key := range input.Keys {
		if strings.Contains(doc.Text(), fmt.Sprintf(Canary, i)) {
			contexts = append(contexts, Context{Type: "html", URL: input.URL, Key: key})
		}
	}

	return contexts
}
