package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (LocationsPage, error) {
	url := baseURL
	if pageURL != nil {
		url = *pageURL
	}

	res, err := c.httpClient.Get(url)
        if err != nil {
                return LocationsPage{}, err
        }
	defer res.Body.Close()

        body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationsPage{}, err
	}

        if res.StatusCode > 299 {
                return LocationsPage{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
        }

        locations := LocationsPage{}
        err = json.Unmarshal(body, &locations)
        if err != nil {
                return LocationsPage{}, err
        }

	return locations, nil
}
