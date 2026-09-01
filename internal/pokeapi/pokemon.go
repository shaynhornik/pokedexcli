package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	var err error
	var body []byte
	var res *http.Response

	url := pokemonURL + "/" + name

	body, ok := c.cache.Get(url)

	if ok == false {
		res, err = c.httpClient.Get(url)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)

		if err != nil {
			return Pokemon{}, err
		}

		if res.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}

		c.cache.Add(url, body)
	}

	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
