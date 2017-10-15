package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	ctail := Ctail{Name: "ctail", UnderlyingCmd: "tail", PipedCmd: "awk", Args: []string{"-f"}}
	ctail.LoadConfig("testdata/ctail.config")
	if ctail.Config == nil || reflect.DeepEqual(*ctail.Config, Config{}) {
		t.Error("Configuration must have been loaded")
	}
}

func TestLoadConfigFromBadFile(t *testing.T) {
	ctail := Ctail{Name: "ctail", UnderlyingCmd: "tail", PipedCmd: "awk", Args: []string{"-f"}}
	ctail.LoadConfig("testdata/ctail_bad_syntax.config")
	if ctail.Config != nil && !reflect.DeepEqual(*ctail.Config, Config{}) {
		t.Error("Configuration must be nil or empty")
		t.FailNow()
	}
}

func TestLoadDefaults(t *testing.T) {
	ctail := Ctail{Name: "ctail", UnderlyingCmd: "tail", PipedCmd: "awk", Args: []string{"-f"}}
	ctail.LoadConfig("some_dummy_file_that_doesnt_exist")
	if ctail.Config == nil || reflect.DeepEqual(*ctail.Config, Config{}) {
		t.Error("Default Configuration must have been loaded")
		t.FailNow()
	}

	if len(ctail.Config.Matches) != 4 && len(ctail.Config.ColorDefinitions) != 17 {
		t.Error("Expected 4 matches and 17 color definitions and got matches : " + string(len(ctail.Config.Matches)) + " and color definitions : " + string(len(ctail.Config.ColorDefinitions)))
		t.FailNow()
	}

}

func TestAnsiColor(t *testing.T) {
	m1 := Match{Expression: "expr1", Color: "blue"}
	m2 := Match{Expression: "expr2", Color: "red"}
	config := Config{
		Matches: []Match{
			m1,
			m2,
		},
		ColorDefinitions: []ColorDefinition{
			{Name: "red", ANSIDefinition: "\"\\033[34m\""},
			{Name: "blue", ANSIDefinition: "\"\\033[31m\""},
		},
	}

	c1 := config.AnsiColor(m1)
	if c1 != "\"\\033[31m\"" {
		t.Error("The expected color was:  \"\\033[31m\" and got:  " + c1)
		t.FailNow()
	}
}

func TestDefaultAnsiColor(t *testing.T) {
	config := Config{
		Matches: []Match{},
		ColorDefinitions: []ColorDefinition{
			{Name: "default", ANSIDefinition: "\"\\033[34m\""},
			{Name: "blue", ANSIDefinition: "\"\\033[31m\""},
		},
	}

	dft := config.DefaultAnsiColor()
	if dft != "\"\\033[34m\"" {
		t.Error("The expected color was:  \"\\033[34m\" and got:  " + dft)
		t.FailNow()
	}
}

func TestAnsiColorWhenFallbackToDefaults(t *testing.T) {
	config := Config{}
	m2 := Match{Expression: "expr2", Color: "red"}

	dft := config.AnsiColor(m2)
	if dft != "\"\\033[39m\"" {
		t.Error("The expected color was:  \"\\033[39m\" and got:  " + dft)
		t.FailNow()
	}
}

type MockCtailCmd struct {
	Config *Config
}

func (c *MockCtailCmd) LoadConfig(file string) {
	var cfg Config
	data := []byte(`{
		"colors_definition" : [
			{ "color_name" : "red", "ansi_definition" : "\"\\033[31m\""  }
			],
		"matches" : [
			{ "expression" : "ERROR", "color" : "red"  } 
			]
			}`)
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot  unmarshal defaults. %s \n", err)
	}
	c.Config = &cfg
}

func (c *MockCtailCmd) Run() {
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

	outfile, err := os.Create("/Users/michel/dev/go/src/ctail/testdata/out.txt")
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(outfile)
	defer writer.Flush()
	defer outfile.Close()

	tailCmd := exec.Command("tail", "testdata/example.txt")
	awkCmd := exec.Command("awk", patterns)
	awkCmd.Stdin, _ = tailCmd.StdoutPipe()
	awkCmd.Stdout = outfile

	awkCmd.Start()
	tailCmd.Run()

	//go io.Copy(writer, os.Stdout)
	awkCmd.Wait()
}

func TestRun(t *testing.T) {
	testfile := "/Users/michel/dev/go/src/ctail/testdata/out.txt"
	defer os.Remove(testfile)
	mockCtail := MockCtailCmd{}

	mockCtail.LoadConfig(testfile)
	mockCtail.Run()

	// reading the file and asserting the content
	dat, _ := ioutil.ReadFile(testfile)
	if string(dat) != "\033[31mERROR\033[0m\n" {
		t.Error("Content is not as expected: Expecting \033[31mERROR\033[0m Got : " + string(dat))
		t.FailNow()
	}
}
