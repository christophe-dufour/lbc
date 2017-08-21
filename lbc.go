package lbc

import (
	"net/http"

	"github.com/christophe-dufour/lbc/token"
	"golang.org/x/net/html"
)

type Offer struct {
	Title     string
	Link      string
	ImgURL    string
	Price     string
	Addresses []string
}

func Parse(URL string) ([]Offer, error) {
	var offers []Offer
	var current Offer

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return offers, nil
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t, err := token.Parse(z)
			if err == token.ErrUnknownToken {
				continue
			}
			if err != nil {
				return nil, err
			}

			switch t.Kind {
			case token.Title:
				current.Title = t.Content
			case token.Price:
				current.Price = t.Content
			case token.Link:
				current.Link = t.Content
			case token.Img:
				current.ImgURL = t.Content
			case token.Address:
				current.Addresses = append(current.Addresses, t.Content)
			}
		case tt == html.EndTagToken:
			if z.Token().Data == "h2" {
				offers = append(offers, current)
				current = Offer{}
			}
		}
	}

}
