package pokeapi

import (
	"net/http"
	"time"

	"pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache *pokecache.Cache
}

func NewClient(timeout time.Duration, c *pokecache.Cache) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,

		},
		cache: c,
	}
}
