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
								// build payload
								u, err := url.Parse(context.URL)
								if err != nil {
									log.Fatal(err)
								}
								q := u.Query()
								q.Add(context.Key, poc)
								u.RawQuery = q.Encode()
								message := u.String()

								// verify payload
								if verifyScript(message, ctx) {
									Results <- (Result{Type: "confirmed", Message: message})
									if Stop {
										return
									}
								}
							}
						}
					}
				}
			}

		case context.Type == "href" || context.Type == "attr":
			if context.Type == "href" {
				// test javascript href
				for _, payload := range AttrPayloads["href"] {
					message := context.URL + "?" + context.Key + "=" + url.QueryEscape(payload)
					// verify payload
					if verifyScript(message, ctx) {
						Results <- (Result{Type: "confirmed", Message: message})
						if Stop {
							return
						}
					}
				}
			}
			// escape attribute
			breakAttr(context, ctx)

		case context.Type == "script":
		case context.Type == "style":
		}
	}
}

func breakAttr(context Context, ctx context.Context) {

	// loop escapes
	for _, escape := range AttrPayloads["escapeAttr"] {
		u := buildUrl(context, escape+Canary2+"="+escape)
		doc := chromeQuery(u, ctx)
		nreflections := doc.Find(fmt.Sprintf("*[%s]", Canary2)).Length()
		if nreflections == 0 {
			continue
		}

		// loop handlers
		for _, handler := range AttrPayloads["handlers"] {
			u := buildUrl(context, escape+handler+"="+escape+Canary2)
			doc := chromeQuery(u, ctx)
			nreflections = doc.Find(fmt.Sprintf("*[%s='%s']", handler, Canary2)).Length()
			if nreflections == 0 {
				continue
			}

			// loop actions
			for _, action := range AttrPayloads["actions"] {
				u := buildUrl(context, escape+handler+"="+escape+action)
				doc := chromeQuery(u, ctx)
				nreflections = doc.Find(fmt.Sprintf("*[%s='%s']", handler, action)).Length()
				if nreflections == 0 {
					continue
				}

				Results <- Result{
					Type:    "high",
					Message: u,
				}
			}
		}
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

	doc := chromeQuery(u, ctx)

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

		doc := chromeQuery(u, ctx)
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
