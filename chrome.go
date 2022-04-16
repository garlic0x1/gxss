package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func chromeQuery(u string, ctx context.Context) *goquery.Document {
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
	return doc
}

func verifyScript(u string, ctx context.Context) bool {
	//log.Println("verifying:", u)
	// perform chrome request
	err := chromedp.Run(ctx,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
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

func identifyCtx(input Input, ctx context.Context) []Context {
	var contexts []Context
	u := buildPayload(input)

	doc := chromeQuery(u, ctx)
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
