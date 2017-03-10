package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var config Config

func main() {
	//prog := os.Args[0]

	args := os.Args[1:]

	config := loadConfig()
	fmt.Printf("config : %v", config)
	tailCmd := exec.Command("tail", args...)

	//awkCmd := exec.Command("awk", "// {print \"\\033[37m\" $0 \"\\033[39m\"} /Exception/ {print \"\\033[32m\" $0 \"\\033[39m\" }")
	awkCmd := exec.Command("awk", "// {print \"\\033[37m\" $0 \"\\033[39m\"} /Exception/ {print \"\\033[32m\" $0 \"\\033[39m\" }")
	awkCmd.Stdin, _ = tailCmd.StdoutPipe()
	awkCmd.Stdout = os.Stdout

	awkCmd.Start()
	tailCmd.Run()
	defer awkCmd.Wait()

}

func loadConfig() *Config {
	home := os.Getenv("HOME")
	file := home + "/.ctailrc"
	if _, err := os.Stat(file); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot stat file %s : %s. Assuming defaults \n", file, err)
		return defaults()
	} else {
		return loadConfigFromFile(file)
	}
}

func loadConfigFromFile(file string) *Config {
	if data, err := ioutil.ReadFile(file); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read from file %s : %s. Assuming defaults \n", file, err)
		return defaults()
	} else {
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot  unmarshal : %s. Assuming defaults \n", err)
			return defaults()
		}
		return &config
	}
}

func defaults() *Config {
	matches := loadDefaultMatches()
	return &Config{
		CheckColoredOutput: false,
		Matches:            matches,
	}
}

func loadDefaultMatches() []Match {
	var m []Match
	data := []byte(`[{"INFO":"orange","DEBUG":"blue","WARN":"yellow"}]`)
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot  unmarshal defaults. %s \n", err)
	}
	return m
}

func getAnsiColor(m Match) string {
	c := m.Color
	for _, cd := range config.ColorDefinition {
		if cd.Name == c {
			return cd.ANSIDefinition
		}
	}
	return getDefaultAnsiColor()
}

func getDefaultAnsiColor() string {
	for _, cd := range config.ColorDefinition {
		if cd.Name == "default" {
			return cd.ANSIDefinition
		}
	}
	// no default found. assume white
	return "\\e[39m"
}

type Match struct {
	Expression string `json:"expression"`
	Color      string `json:"color"`
}

type ColorDefinition struct {
	Name           string `json:"color_name"`
	ANSIDefinition string `json:"ansi_definition"`
}

type Config struct {
	CheckColoredOutput bool              `json:"check_colored_output"`
	Matches            []Match           `json:"matches"`
	ColorDefinition    []ColorDefinition `json:"colors_definition"`
}
