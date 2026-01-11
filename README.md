# PokéDex CLI

A CLI mini-pokemon game I built during [Boot.dev](https://boot.dev)'s Golang course. This is my first project in Go, so it was really about getting quicker with the language and familiar with http requests in Go. 
Highlight for me was learning caching with automatic reaping using goroutines and channels.

Now I'll hand the mic to Gemini....

Some Features:
* **Real-time Data:** Integration with the [PokéAPI](https://pokeapi.co/) to fetch locations and Pokémon stats.
* **Custom In-Memory Cache:** A high-performance cache implementation to reduce API latency and minimize network requests.
* **Automatic Cache Reaping:** Background goroutines automatically clean up stale cache entries after a configurable duration.
* **Interactive REPL:** A smooth command-line interface with a dedicated command registry.

---

## Installation

1. **Clone the repository:**
```bash
git clone https://github.com/OriElbaz/pokedex.git
cd pokedex

```


2. **Build the project:**
```bash
go build -o pokedex

```


3. **Run the game:**
```bash
./pokedex

```



---

## How to Play

Once the program is running, you can use the following commands:

### Exploration

* `map`: Displays the next 20 location areas in the Pokémon world.
* `mapb`: Displays the previous 20 location areas.
* `explore <location-name>`: See a list of all Pokémon species found in a specific area.

### Catching Pokemon & Pokedex

* `catch <pokemon-name>`: Attempt to catch a Pokémon. Be careful—higher-level Pokémon are harder to catch!
* `inspect <pokemon-name>`: View the stats (Height, Weight, HP, Attack, etc.) of a Pokémon you have successfully caught.
* `pokedex`: List all the Pokémon currently in your collection.

### System

* `help`: Displays the manual and all available commands.
* `exit`: Safely closes the Pokedex program.
