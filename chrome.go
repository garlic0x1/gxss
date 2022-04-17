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

func chromeQuery(u string, tab context.Context) *goquery.Document {
	var document string

	// perform chrome request
	err := chromedp.Run(tab,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
		chromedp.OuterHTML(`html`, &document),
	)
	if err != nil && Debug {
		log.Println("Error from chromeQuery()", u, err)
	}

	// analyze response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		log.Println(err)
	}
	return doc
}

func confirmAlert(u string, tab context.Context, handler string, selector string) {
	// perform chrome request
	err := chromedp.Run(tab,
		chromedp.Navigate(u),
		chromedp.Sleep(time.Duration(Wait)),
		//chromedp.ResetViewport(),
	)
	if err != nil && Debug {
		log.Println("Error from chromeQuery()", u, err)
	}

	if Interact {
		// pop alerts that require interaction
		cmd := fmt.Sprintf("document.querySelector(\"%s\").%s()", selector, handler)
		err = chromedp.Run(tab,
			chromedp.Evaluate(cmd, nil),
		)
		if err != nil && Debug {
			log.Println("Error executing handler", selector, handler, err)
		}
	}
}

func identifyCtx(input Input, tab context.Context) []Context {
	var contexts []Context
	u := buildPayload(input)

	doc := chromeQuery(u, tab)
	doc.Find("*").Each(func(_ int, node *goquery.Selection) {
		n := node.Get(0)
		for _, attr := range n.Attr {
			for i, key := range input.Keys {
				if strings.Contains(attr.Val, fmt.Sprintf(Canary, i)) {
					contexts = append(contexts, Context{
						Type:   "attr",
						URL:    input.URL,
						Prefix: "",
						Key:    key,
					})
				}
			}
		}
	})

	for i, key := range input.Keys {
		if strings.Contains(doc.Text(), fmt.Sprintf(Canary, i)) {
			contexts = append(contexts, Context{
				Type:   "html",
				URL:    input.URL,
				Prefix: "",
				Key:    key,
			})
		}
	}

	return contexts
}
