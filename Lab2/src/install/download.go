package main

import (
	"net/http"
	//"strconv"

	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Time, Title string
}

func readItem(item *html.Node) *Item {
	if chld := getChildren(item); len(chld) == 1 && isText(chld[0]) {
		//log.Info(chld[0].Data)
		return &Item{
			Ref:   getAttr(item, "href"),
			Title: chld[0].Data,
		}
	} else {
		return &Item{

			Ref:   "e",
			Time:  "",
			Title: "error 1",
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	//if isElem(node, "div") && getAttr(node, "class") == "" {
	if isDiv(node, "threadinfo") {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "inner") {
				for h3 := c.FirstChild; h3 != nil; h3 = h3.NextSibling {
					if isElem(h3, "h3") && getAttr(h3, "class") == "threadtitle" {
						for a := h3.FirstChild; a != nil; a = a.NextSibling {
							if isElem(a, "a") {
								if item := readItem(a); item != nil {
									items = append(items, item)
								}

							}
						}
					} else {
						var a Item
						a.Ref = "h"
						a.Time = ""
						a.Title = "hubabuba"
						items = append(items, &a)
					}
				}
			}
			//}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

var items []*Item

func recsearch(node *html.Node/*, items []*Item*/)  {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if isElem(c, "h3") && getAttr(c, "class") == "threadtitle" {
			//log.Info("Node "+c.Data)
			for a := c.FirstChild; a != nil; a = a.NextSibling {
				if isElem(a, "a") {
					if item := readItem(a); item != nil {
						items = append(items, item)
						//log.Info("Items "+strconv.Itoa(len(items)-1)+" "+items[len(items)-1].Title)
					}
				}
			}

		} else {
			recsearch(c)
		}
	}
	//return items
}

func downloadNews() []*Item {
	log.Info("sending request to forums.drom.ru")
	if response, err := http.Get("https://forums.drom.ru/moscow"); err != nil {
		log.Error("request to forum.droms.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from forum.droms.ru", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from forum.droms.ru", "error", err)
			} else {
				log.Info("HTML from forum.droms.ru parsed successfully")
				items = nil
				recsearch(doc)
				return items
			}
		}
	}
	return nil
}

/*func main(){
	downloadNews()
}*/
