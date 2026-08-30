package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetLocationArea(name string) (LocationArea, error) {
	var err error
	var body []byte
	var res *http.Response

	url := baseURL + "/" + name

	body, ok := c.cache.Get(url)

	if ok == false {
		res, err = c.httpClient.Get(url)
		if err != nil {
			return LocationArea{}, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return LocationArea{}, err
		}

		if res.StatusCode > 299 {
			return LocationArea{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}

		c.cache.Add(url, body)
	}

	locationArea := LocationArea{}
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, err
	}

	return locationArea, nil
}
