package main

import (
	"embed"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type Language string

const (
	LanguageEnglish  Language = "english"
	LanguageFrench            = "french"
	LanguageJapanese          = "japanese"
	LanguageChinese           = "chinese"
)

// Pokemon struct represents the data structure for a Pokémon
type Pokemon struct {
	Path  string              `msgpack:"path"`
	Names map[Language]string `msgpack:"names"`
	Forms []string            `msgpack:"forms"`
}

// Embed assets directory
//
//go:embed assets/*
var assets embed.FS

const (
	rootDir         = "assets"
	shinyRate       = 1.0 / 128.0
	colorscriptsDir = "colorscripts"
	regularSubdir   = "regular"
	shinySubdir     = "shiny"
)

// Generation ranges for Pokémon
var generations = map[string][2]int{
	"1": {1, 151},
	"2": {152, 251},
	"3": {252, 386},
	"4": {387, 493},
	"5": {494, 649},
	"6": {650, 721},
	"7": {722, 809},
	"8": {810, 898},
}

// readPokemon reads the pokemon.json file from the embedded assets
func readPokemon() []Pokemon {
	file, err := assets.ReadFile(filepath.Join(rootDir, "pokemons.msgpack"))
	if err != nil {
		panic(err)
	}

	var pokemon []Pokemon
	if err := msgpack.Unmarshal(file, &pokemon); err != nil {
		panic(err)
	}

	return pokemon
}

// showPokemonByName displays Pokémon information based on its name
func showPokemon(pokemon Pokemon, showTitle, shiny bool, lang Language) {
	colorSubdir := regularSubdir
	if shiny {
		colorSubdir = shinySubdir
	}

	name := pokemon.Names[lang]
	if showTitle {
		if shiny {
			fmt.Printf("%s (shiny)\n", name)
		} else {
			fmt.Println(name)
		}
	}

	fp := filepath.Join(rootDir, colorscriptsDir, colorSubdir, strings.ToLower(pokemon.Path))
	content, err := assets.ReadFile(fp)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	os.Stdout.Write(content)
}

// showRandomPokemon displays a random Pokémon based on specified generations
func showRandomPokemon(generationsStr string, showTitle, shiny bool, lang Language) {
	var startGen, endGen string
	genList := strings.Split(generationsStr, ",")

	if len(genList) > 1 {
		startGen = genList[rand.Intn(len(genList))]
		endGen = startGen
	} else if strings.Contains(generationsStr, "-") {
		parts := strings.Split(generationsStr, "-")
		startGen, endGen = parts[0], parts[1]
	} else {
		startGen = generationsStr
		endGen = startGen
	}

	pokemon := readPokemon()
	startIdx, ok := generations[startGen]
	if !ok {
		fmt.Printf("invalid generation '%s'\n", generationsStr)
		os.Exit(1)
	}

	endIdx, ok := generations[endGen]
	if !ok {
		fmt.Printf("invalid generation '%s'\n", generationsStr)
		os.Exit(1)
	}

	randomIdx := rand.Intn(endIdx[1]-startIdx[0]+1) + startIdx[0]

	if !shiny && rand.Float64() <= shinyRate {
		shiny = true
	}

	showPokemon(pokemon[randomIdx-1], showTitle, shiny, lang)
}

func main() {
	fs := flag.NewFlagSet("pokego", flag.ExitOnError)

	noTitleFlag := fs.Bool("no-title", false, "Do not display pokemon name")
	shinyFlag := fs.Bool("shiny", false, "Show the shiny version of a pokemon instead")
	fs.BoolVar(shinyFlag, "s", *shinyFlag, "--shiny")

	generationFlag := fs.String("generation", "1-8", "Generation number or range filter")
	fs.StringVar(generationFlag, "g", *generationFlag, "--generation")

	languageFlag := fs.String("language", string(LanguageEnglish), fmt.Sprintf("Language to print the pokemon name in, one of (%v)", []Language{LanguageEnglish, LanguageFrench, LanguageJapanese, LanguageChinese}))
	fs.StringVar(languageFlag, "l", string(LanguageEnglish), "--language")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		return
	}

	showRandomPokemon(*generationFlag, !*noTitleFlag, *shinyFlag, *(*Language)(languageFlag))
}
