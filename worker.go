package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func worker(input Input, tab context.Context) {
	// identify context
	contexts := identifyCtx(input, tab)

	// escape context
	for _, context := range contexts {
		switch {
		case context.Type == "html":
			breakHtml(context, tab)

		case context.Type == "href" || context.Type == "attr":
			if context.Type == "href" {
				breakLink(context, tab)
			}
			// escape attribute
			breakAttr(context, tab)

		case context.Type == "script":
			breakScript(context, tab)
		}
	}
}

func breakScript(context Context, tab context.Context) {
	// first try to close script tag break into html context
	// loop open brackets
	for _, openBracket := range Payloads["openBracket"] {
		u := buildUrl(context, fmt.Sprintf(Canary3, openBracket))
		doc := chromeQuery(u, tab)
		str, err := doc.Html()
		if err != nil && Debug {
			log.Println(err)
		}
		ok := strings.Contains(str, fmt.Sprintf(Canary3, "<")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<"))
		if !ok {
			continue
		}

		// loop close brackets
		for _, closeBracket := range Payloads["closeBracket"] {
			u = buildUrl(context, fmt.Sprintf(Canary3, openBracket+closeBracket))
			doc = chromeQuery(u, tab)
			str, err = doc.Html()
			if err != nil && Debug {
				log.Println(err)
			}
			ok = strings.Contains(str, fmt.Sprintf(Canary3, "<>")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<>"))
			if !ok {
				continue
			}

			// loop backslash
			for _, backslash := range Payloads["backslash"] {
				u = buildUrl(context, fmt.Sprintf(Canary3, backslash))
				doc = chromeQuery(u, tab)
				str, err = doc.Html()
				if err != nil && Debug {
					log.Println(err)
				}
				ok = strings.Contains(str, fmt.Sprintf(Canary3, "/")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "/"))
				if !ok {
					continue
				}

				// loop close tags
				for _, scriptTag := range Payloads["script"] {
					u = buildUrl(context, fmt.Sprintf(Canary3, openBracket+backslash+scriptTag+closeBracket))
					doc = chromeQuery(u, tab)
					str, err = doc.Html()
					if err != nil && Debug {
						log.Println(err)
					}
					ok = strings.Contains(str, fmt.Sprintf(Canary3, openBracket+backslash+scriptTag+closeBracket)) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, openBracket+backslash+scriptTag+closeBracket))
					if ok {
						// now we are in html context, so generate payload for that
						context.Prefix = openBracket + backslash + scriptTag + closeBracket
						breakHtml(context, tab)
					}

				}
			}
		}
	}
}

func breakLink(context Context, tab context.Context) {
	// test javascript href
	for _, payload := range Payloads["href"] {
		u := buildUrl(context, payload)
		//confirmAlert(u, tab)
		confirmAlert(u, tab, context.Attr, context.Selector)
	}
}

func findAttr(context Context, tab context.Context) {
	// loop handlers
	tag := context.Tag
	closeBracket := context.CloseBracket
	openBracket := context.OpenBracket
	attrList := []string{}
	attrData := []Handler{}
	for handler, data := range TagMap[tag] {
		attrList = append(attrList, handler)
		attrData = append(attrData, data)
	}
	for handler, data := range TagMap["*"] {
		attrList = append(attrList, handler)
		attrData = append(attrData, data)
	}
	for i := 0; i < len(attrList); i++ {
		handler := attrList[i]
		data := attrData[i]
		u := buildUrl(context, openBracket+tag+" "+handler+"="+Canary2+closeBracket)
		doc := chromeQuery(u, tab)
		nreflections := doc.Find(fmt.Sprintf("%s[%s=%s]", tag, handler, Canary2)).Length()
		if nreflections == 0 {
			continue
		}

		Results <- Result{Type: "medium", Message: u}

		for _, payload := range data.Payloads {
			u = buildUrl(context, payload)
			// verify payload
			confirmAlert(u, tab, handler, fmt.Sprintf("%s[%s]", tag, handler))
		}

		for _, action := range Payloads["actions"] {
			u = buildUrl(context, openBracket+tag+" "+handler+"="+action+closeBracket)
			doc = chromeQuery(u, tab)
			confirmAlert(u, tab, handler, fmt.Sprintf("%s[%s='%s']", tag, handler, action))

			nreflections = doc.Find(fmt.Sprintf("%s[%s='%s']", tag, handler, action)).Length()
			if nreflections == 0 {
				continue
			}

			Results <- (Result{Type: "high", Message: u})
		}
	}
}

func breakHtml(context Context, tab context.Context) {
	// loop open brackets
	for _, openBracket := range Payloads["openBracket"] {
		u := buildUrl(context, fmt.Sprintf(Canary3, openBracket))
		doc := chromeQuery(u, tab)
		str, err := doc.Html()
		if err != nil && Debug {
			log.Println(err)
		}
		ok := strings.Contains(str, fmt.Sprintf(Canary3, "<")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<"))
		if !ok {
			continue
		}

		// loop close brackets
		for _, closeBracket := range Payloads["closeBracket"] {
			u = buildUrl(context, fmt.Sprintf(Canary3, openBracket+closeBracket))
			doc = chromeQuery(u, tab)
			str, err = doc.Html()
			if err != nil && Debug {
				log.Println(err)
			}
			ok = strings.Contains(str, fmt.Sprintf(Canary3, "<>")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<>"))
			if !ok {
				continue
			}

			// loop tags
			for tag, _ := range TagMap {
				u = buildUrl(context, openBracket+tag+" "+Canary2+"=x"+closeBracket)
				doc = chromeQuery(u, tab)
				nreflections := doc.Find(fmt.Sprintf("%s[%s]", tag, Canary2)).Length()
				if nreflections == 0 {
					continue
				}

				Results <- Result{Type: "low", Message: u}

				//context.Prefix = openBracket + tag + " "
				context.Tag = tag
				context.OpenBracket = openBracket
				context.CloseBracket = closeBracket
				findAttr(context, tab)
			}
		}
	}
}

func breakAttr(context Context, tab context.Context) {

	// loop escapes
	for _, quote := range Payloads["quotes"] {
		u := buildUrl(context, quote+Canary2+"="+quote)
		doc := chromeQuery(u, tab)
		nreflections := doc.Find(fmt.Sprintf("*[%s='']", Canary2)).Length()
		tag := goquery.NodeName(doc.Find(fmt.Sprintf("*[%s='']", Canary2)))
		if nreflections == 0 {
			continue
		}
		Results <- Result{Type: "low", Message: u}

		// loop handlers
		for handler, _ := range TagMap[tag] {
			u := buildUrl(context, quote+handler+"="+quote+Canary2)
			doc := chromeQuery(u, tab)
			nreflections = doc.Find(fmt.Sprintf("*[%s='%s']", handler, Canary2)).Length()
			if nreflections == 0 {
				continue
			}

			Results <- Result{Type: "medium", Message: u}

			// loop actions
			for _, action := range Payloads["actions"] {
				u := buildUrl(context, quote+handler+"="+quote+action)
				doc := chromeQuery(u, tab)
				nreflections = doc.Find(fmt.Sprintf("*[%s='%s']", handler, action)).Length()
				if nreflections == 0 {
					continue
				}
				confirmAlert(u, tab, handler, fmt.Sprintf("%s[%s]", context.Selector, handler))

				Results <- Result{Type: "high", Message: u}
			}
		}

		// break into html and then try that
		for _, bracket := range Payloads["closeBracket"] {
			u := buildUrl(context, quote+bracket+Canary2)
			doc := chromeQuery(u, tab)

			item := doc.Find(context.Selector)
			htmlText, err := item.Html()
			if err != nil && Debug {
				log.Println(err)
			}
			if strings.Contains(htmlText, Canary2) {
				context.Prefix = quote + bracket
				breakHtml(context, tab)
			}
		}
	}
}
