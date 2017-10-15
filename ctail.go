package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// The configuration file name. The file must be placed on user home folder
const ConfigFile = "ctail.config"

// ctail will do the same as tail command and add also the capability
// to colorize the output following some regexps defined in a
// json configuration file ($HOME/.ctail.config).
//
// ctail accepts the same arguments as tail command. To run it just:
//
// ctail [args] filename
func main() {
	args := os.Args[1:]
	ctail := Ctail{Name: "ctail", UnderlyingCmd: "tail", PipedCmd: "awk", Args: args}

	ctail.Run()
}

type Cmd interface {
	Run()
	LoadConfig(filepath string) *Config
}

// Run will create the commands, loads the configration and then start executing
func (c *Ctail) Run() {
	home := os.Getenv("HOME")
	file := home + "/" + ConfigFile
	c.LoadConfig(file)

	// by default print all lines as they are.
	var patterns string = ""
	matches := c.Config.Matches
	for _, m := range matches {
		patterns = patterns + "/" + m.Expression + "/ { print " + c.Config.AnsiColor(m) + " $0 \"\\033[0m\";next} "
	}
	patterns = patterns + "{print $0}"

	tailCmd := exec.Command(c.UnderlyingCmd, c.Args...)
	awkCmd := exec.Command(c.PipedCmd, patterns)
	awkCmd.Stdin, _ = tailCmd.StdoutPipe()
	awkCmd.Stdout = os.Stdout

	awkCmd.Start()
	tailCmd.Run()
	awkCmd.Wait()
}

// Will load the configration from the file specified. If the file
// doesn't exist then will load some defaults.
func (c *Ctail) LoadConfig(file string) {
	if _, err := os.Stat(file); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot stat file %s : %s. Assuming defaults \n", file, err)
		c.Config = defaults()
	} else {
		c.loadConfigFromFile(file)
	}
}

func (c *Ctail) loadConfigFromFile(file string) {
	var conf Config
	if data, err := ioutil.ReadFile(file); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read from file %s : %s. Assuming defaults \n", file, err)
		c.Config = defaults()
	} else {
		if err := json.Unmarshal(data, &conf); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot  unmarshal : %s. Assuming defaults \n", err)
			c.Config = defaults()
		}
		c.Config = &conf
	}
}

func defaults() *Config {
	var c Config
	data := []byte(`{
						"colors_definition": [ 
							{ "color_name" : "default" , "ansi_definition" : "\"\\033[39m\""},
							{ "color_name" : "black", "ansi_definition" : "\"\\033[30m\"" },
							{ "color_name" : "red", "ansi_definition" : "\"\\033[31m\"" },
							{ "color_name" : "green" , "ansi_definition" : "\"\\033[32m\""},
							{ "color_name" : "yellow" , "ansi_definition" : "\"\\033[33m\"" },
							{ "color_name" : "blue" , "ansi_definition" : "\"\\033[34m\"" }
						],
						"matches" : [
							{"expression": "Error|ERROR", "color" : "red" },
							{"expression": "Debug|DEBUG", "color" : "green" },
							{"expression": "Warning|Warn|WARN|WARNING", "color" : "yellow" },
							{"expression": "Info|INFO", "color" : "blue" }
					]

				}`)
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot  unmarshal defaults. %s \n", err)
	}
	return &c
}

// Return the ansi color loaded by the configration.
// If no ansi color is found then return the default.
func (config Config) AnsiColor(m Match) string {
	c := m.Color
	for _, cd := range config.ColorDefinitions {
		if cd.Name == c {
			return cd.ANSIDefinition
		}
	}
	return config.DefaultAnsiColor()
}

// Will return the default ansi color defined in the configration
// if no "default" is specified then White (\033[39m') is assumed
func (config Config) DefaultAnsiColor() string {
	for _, cd := range config.ColorDefinitions {
		if cd.Name == "default" {
			return cd.ANSIDefinition
		}
	}
	// no default found. assume white
	return "\"\\033[39m\""
}

type Ctail struct {
	Name          string
	UnderlyingCmd string
	PipedCmd      string
	Args          []string
	Config        *Config
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
	Matches          []Match           `json:"matches"`
	ColorDefinitions []ColorDefinition `json:"colors_definition"`
}
