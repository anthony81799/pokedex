package pokeapi

import (
	"net/http"
	"time"

	"github.com/anthony81799/pokedex/internal/pokecache"
)

type Client struct {
	cache      pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cachInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cachInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
