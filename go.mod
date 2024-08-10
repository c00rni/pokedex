module github.com/c00rni/pokedex

go 1.22.5

replace example.com/corni/pokeAPI v0.0.0 => ../pokeAPI

require (
	example.com/corni/pokeAPI v0.0.0
	github.com/mtslzr/pokeapi-go v1.4.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
)
