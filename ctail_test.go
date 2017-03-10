package main

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := loadConfig()
	if config == nil {
		t.Error("Configuration must have been loaded")
	}
}

func TestLoadDefaults(t *testing.T) {
	config := defaults()
	if config == nil {
		t.Error("Configuration must have been loaded")
		t.FailNow()
	}

	if config.CheckColoredOutput != false {
		t.Error("default check_colored_output must be true")
	}

	if config.Matches == nil {
		t.Error("defaults matches must exists")
		t.FailNow()
	}

	if len(config.Matches) < 1 {
		t.Error("there must have more that 1 matches configured by default")
	}

}

/*

	{ "color_name" : "default" , "ansi_definition" : "\\e[39m"},
	{ "color_name" : "black", "ansi_definition" : "\\e[30m" },
	{ "color_name" : "red", "ansi_definition" : "\\e[31m" },
	{ "color_name" : "green" , "ansi_definition" : "\\e[32m"},
	{ "color_name" : "yellow" , "ansi_definition" : "\\e[33m"},
	{ "color_name" : "blue" , "ansi_definition" : "\\e[34m" },
	{ "color_name" : "magenta" , "ansi_definition" : "\\e[35m" },
	{ "color_name" : "cyan" , "ansi_definition" : "\\e[36m" },
	{ "color_name" : "light gray" , "ansi_definition" : "\\e[37m" },
	{ "color_name" : "dark gray" , "ansi_definition" : "\\e90[m" },
	{ "color_name" : "light red" , "ansi_definition" : "\\e[91m" },
	{ "color_name" : "light green" , "ansi_definition" : "\\e[92m" },
	{ "color_name" : "light yellow" , "ansi_definition" : "\\e[93m" },
	{ "color_name" : "light blue" , "ansi_definition" : "\\e[94m" },
	{ "color_name" : "light magenta" , "ansi_definition" : "\\e[95m" },
	{ "color_name" : "light cyan" , "ansi_definition" : "\\e[96m" },
	{ "color_name" : "light white" , "ansi_definition" : "\\e[97m" }

*/
func TestGetAnsiColorFromMatch(t *testing.T) {
	config := loadConfig()

	for _, m := range config.Matches {
		colorDefinition := getAnsiColor(m)
		colorName := m.Color
		switch colorName {
		case "default":
			if colorDefinition != "\\e[39m" {
				t.Errorf("color %s must have \\e[39m as definition", colorName)
			}
		case "black":
			if colorDefinition != "\\e[30m" {
				t.Errorf("color %s must have \\e[30m as definition", colorName)
			}
		case "red":
			if colorDefinition != "\\e[31m" {
				t.Errorf("color %s must have \\e[31m as definition", colorName)
			}
		case "green":
			if colorDefinition != "\\e[32m" {
				t.Errorf("color %s must have \\e[32m as definition", colorName)
			}
		case "yellow":
			if colorDefinition != "\\e[33m" {
				t.Errorf("color %s must have \\e[33m as definition", colorName)
			}
		case "blue":
			if colorDefinition != "\\e[34m" {
				t.Errorf("color %s must have \\e[34m as definition", colorName)
			}
		}

	}

}

func TestGetAnsiColorNotEmpty(t *testing.T) {
	config := loadConfig()

	for _, m := range config.Matches {
		color := getAnsiColor(m)
		if color == "" {
			fmt.Printf("color must not be empty")
		}
	}
}
