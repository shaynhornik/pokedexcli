package pokeapi

type LocationArea struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonRef `json:"pokemon"`
}

type PokemonRef struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
