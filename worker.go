package main

import (
	"context"
	"fmt"
	"log"
	"strings"
)

func worker(input Input, ctx context.Context, workerId int) {
	// identify context
	contexts := identifyCtx(input, ctx)

	// escape context
	for _, context := range contexts {
		switch {
		case context.Type == "html":
			breakHtml(context, ctx, workerId)

		case context.Type == "href" || context.Type == "attr":
			if context.Type == "href" {
				breakLink(context, ctx, workerId)
			}
			// escape attribute
			breakAttr(context, ctx, workerId)

		case context.Type == "script":
		}
	}
}

func breakLink(context Context, ctx context.Context, workerId int) {
	// test javascript href
	for _, payload := range AttrPayloads["href"] {
		u := buildUrl(context, payload)
		_ = u
		Confirm <- ConfType{
			URL:    u,
			Worker: workerId,
		}
	}
}

func breakHtml(context Context, ctx context.Context, workerId int) {

	// loop open brackets
	for _, openBracket := range AttrPayloads["openBracket"] {
		u := buildUrl(context, fmt.Sprintf(Canary3, openBracket))
		doc := chromeQuery(u, ctx)
		str, err := doc.Html()
		if err != nil {
			log.Println(err)
		}
		ok := strings.Contains(str, fmt.Sprintf(Canary3, "<")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<"))
		if !ok {
			continue
		}

		// loop close brackets
		for _, closeBracket := range AttrPayloads["closeBracket"] {
			u = buildUrl(context, fmt.Sprintf(Canary3, openBracket+closeBracket))
			doc = chromeQuery(u, ctx)
			str, err = doc.Html()
			if err != nil {
				log.Println(err)
			}
			ok = strings.Contains(str, fmt.Sprintf(Canary3, "<>")) || strings.Contains(doc.Text(), fmt.Sprintf(Canary3, "<>"))
			if !ok {
				continue
			}

			// loop tags
			for tag, _ := range Payloads {
				u = buildUrl(context, openBracket+tag+" "+Canary2+"=x"+closeBracket)
				doc = chromeQuery(u, ctx)
				nreflections := doc.Find(fmt.Sprintf("%s[%s]", tag, Canary2)).Length()
				if nreflections == 0 {
					continue
				}

				// loop handlers
				for handler, data := range Payloads[tag] {
					u = buildUrl(context, openBracket+tag+" "+handler+"="+Canary2+closeBracket)
					doc = chromeQuery(u, ctx)
					nreflections = doc.Find(fmt.Sprintf("%s[%s=%s]", tag, handler, Canary2)).Length()
					if nreflections == 0 {
						continue
					}

					// loop payloads
					for _, payload := range data["payloads"] {
						u = buildUrl(context, payload)
						// verify payload
						Confirm <- ConfType{
							URL:    u,
							Worker: workerId,
						}
					}

					// loop requireds
					c := 0
					for _, required := range data["requires"] {
						c++

						u = buildUrl(context, openBracket+tag+" "+handler+"="+Canary2+" "+required+"="+closeBracket)
						doc = chromeQuery(u, ctx)

						nreflections = doc.Find(fmt.Sprintf("%s[%s='%s'][%s]", tag, handler, Canary2, required)).Length()
						if nreflections == 0 {
							continue
						}

						for _, action := range AttrPayloads["actions"] {
							u = buildUrl(context, openBracket+tag+" "+handler+"="+action+" "+required+"="+closeBracket)
							doc = chromeQuery(u, ctx)

							nreflections = doc.Find(fmt.Sprintf("%s[%s='%s']", tag, handler, action)).Length()
							if nreflections == 0 {
								continue
							}

							Confirm <- ConfType{
								URL:    u,
								Worker: workerId,
							}
							Results <- (Result{Type: "high", Message: u})
						}
					}

					if c > 0 {
						continue
					}

					for _, action := range AttrPayloads["actions"] {
						u = buildUrl(context, openBracket+tag+" "+handler+"="+action+closeBracket)
						doc = chromeQuery(u, ctx)

						nreflections = doc.Find(fmt.Sprintf("%s[%s='%s']", tag, handler, action)).Length()
						if nreflections == 0 {
							continue
						}

						Confirm <- ConfType{
							URL:    u,
							Worker: workerId,
						}
						//Results <- (Result{Type: "medium", Message: u})
					}
				}
			}
		}
	}
}

func breakAttr(context Context, ctx context.Context, workerId int) {

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

				Confirm <- ConfType{
					URL:    u,
					Worker: workerId,
				}
				Results <- Result{
					Type:    "high",
					Message: u,
				}
			}
		}
	}
}
