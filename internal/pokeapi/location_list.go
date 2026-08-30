package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (LocationsPage, error) {
	var err error
	var body []byte
	var res *http.Response

	url := baseURL
	if pageURL != nil {
		url = *pageURL
	}

	body, ok := c.cache.Get(url)

	if ok == false {
		fmt.Println("cache miss")
		res, err = c.httpClient.Get(url)
        	if err != nil {
                	return LocationsPage{}, err
        	}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return LocationsPage{}, err
		}

	        if res.StatusCode > 299 {
        	        return LocationsPage{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
        	}

		c.cache.Add(url, body)
	} else {
		fmt.Println("cache hit")
	}

        locations := LocationsPage{}
        err = json.Unmarshal(body, &locations)
        if err != nil {
                return LocationsPage{}, err
        }

	return locations, nil
}
