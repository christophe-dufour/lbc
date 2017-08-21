package token

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

var ErrUnknownToken = errors.New("unknown token")

type Kind int

const (
	Title Kind = iota
	Price
	Link
	Img
	Address
)

type Token struct {
	html.Token
	Kind    Kind
	Content string
}

func Parse(z *html.Tokenizer) (*Token, error) {
	t := z.Token()

	if isAddress(t) {
		return &Token{
			Token:   t,
			Kind:    Address,
			Content: strings.TrimSpace(get("content", t.Attr)),
		}, nil
	}

	if isTitle(t) {
		z.Next()
		return &Token{
			Token:   t,
			Kind:    Title,
			Content: strings.TrimSpace(z.Token().Data),
		}, nil
	}

	if isPrice(t) {
		return &Token{
			Token:   t,
			Kind:    Price,
			Content: get("content", t.Attr),
		}, nil
	}
	if isLink(t) {
		return &Token{
			Token:   t,
			Kind:    Link,
			Content: get("href", t.Attr),
		}, nil
	}
	if isIMG(t) {
		return &Token{
			Token:   t,
			Kind:    Img,
			Content: get("data-imgsrc", t.Attr),
		}, nil
	}

	return nil, ErrUnknownToken
}

func isIMG(t html.Token) bool {
	return t.Data == "span" && contains(t.Attr, html.Attribute{Key: "class", Val: "lazyload"})
}

func isLink(t html.Token) bool {
	return t.Data == "a"
}

func isTitle(t html.Token) bool {
	return t.Data == "h2" && contains(t.Attr, html.Attribute{Key: "class", Val: "item_title"}) &&
		contains(t.Attr, html.Attribute{Key: "itemprop", Val: "name"})
}
func isPrice(t html.Token) bool {
	return t.Data == "h3" && contains(t.Attr, html.Attribute{Key: "class", Val: "item_price"}) &&
		contains(t.Attr, html.Attribute{Key: "itemprop", Val: "price"})
}

func isOffer(t html.Token) bool {
	return t.Data == "li" && contains(t.Attr, html.Attribute{Key: "itemtype", Val: "http://schema.org/Offer"})
}

func isAddress(t html.Token) bool {
	return t.Data == "meta" && contains(t.Attr, html.Attribute{Key: "itemprop", Val: "address"})
}

func contains(atts []html.Attribute, att html.Attribute) bool {
	for _, a := range atts {
		if a == att {
			return true
		}
	}

	return false
}

func get(key string, atts []html.Attribute) string {
	for _, att := range atts {
		if att.Key == key {
			return att.Val
		}
	}
	return ""
}
