package main

import (
	"encoding/json"
	"log"

	"github.com/christophe-dufour/lbc"
)

const url = "https://www.leboncoin.fr/annonces/offres/midi_pyrenees"

func main() {
	o, err := lbc.Parse(url)
	h(err)

	offers, err := lbc.Localize(o)
	h(err)

	b, err := json.Marshal(offers)
	h(err)

	println(string(b))
}

func h(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
