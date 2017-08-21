package lbc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var (
	key         = os.Getenv("KEY")
	ErrNotFound = errors.New("not found")
)

type Resp struct {
	Results []Result `json:"results,omitempty"`
}

type Result struct {
	Geometry Geometry `json:"geometry,omitempty"`
}

type Geometry struct {
	Location Location `json:"location,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

type LocalizedOffer struct {
	Offer
	Location
}

func Localize(offers []Offer) ([]LocalizedOffer, error) {
	result := make([]LocalizedOffer, len(offers))
	var wg sync.WaitGroup

	for i, _ := range offers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			location, err := localize(offers[i].Addresses)
			if err != nil {
				println(err.Error())

				return
			}

			result[i].Offer = offers[i]
			result[i].Location = *location
		}(i)
	}

	wg.Wait()
	return result, nil
}

func localize(addresses []string) (*Location, error) {
	u, err := url.Parse("https://maps.googleapis.com/maps/api/geocode/json")
	if err != nil {
		return nil, err
	}

	var address string
	for _, a := range addresses {
		address = fmt.Sprintf("%s+%s", address, a)
	}
	q := u.Query()
	q.Add("key", key)
	q.Add("address", address)

	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("wrong status code %d", resp.StatusCode)
	}

	var r Resp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	if len(r.Results) == 0 {
		return nil, ErrNotFound
	}

	return &(r.Results[0].Geometry.Location), nil
}
