# ctail

Naive implementation in Go of a colorized `tail` command. 

# Install

The utility has to be built. A Go compiler is needed.

```sh
$ go build ctail
```

a `ctail` executable file will be generated.

# Usage

The utility accepts the same arguments as `tail` do.

```sh
$ ctail -100f some_file.log
```

# Config

`ctail` provides the ability to specify the patterns to find and colors to print those patterns. The configure is stored on a `json` file `ctail.config` that must be location on `$HOME` folder. Patterns can be regular expressions.

Here is an example 

```json
{
    "colors_definition": [ 
        { "color_name" : "default" , "ansi_definition" : "\"\\033[39m\""},
        { "color_name" : "black", "ansi_definition" : "\"\\033[30m\"" },
        { "color_name" : "red", "ansi_definition" : "\"\\033[31m\"" },
        { "color_name" : "green" , "ansi_definition" : "\"\\033[32m\""},
        { "color_name" : "yellow" , "ansi_definition" : "\"\\033[33m\"" },
        { "color_name" : "blue" , "ansi_definition" : "\"\\033[34m\"" },
        { "color_name" : "magenta" , "ansi_definition" : "\"\\033[35m\"" },
        { "color_name" : "cyan" , "ansi_definition" : "\"\\033[36m\"" },
        { "color_name" : "light gray" , "ansi_definition" : "\"\\033[37m\"" },
        { "color_name" : "dark gray" , "ansi_definition" : "\"\\033[90m\"" },
        { "color_name" : "light red" , "ansi_definition" : "\"\\033[91m\"" },
        { "color_name" : "light green" , "ansi_definition" : "\"\\033[92m\"" },
        { "color_name" : "light yellow" , "ansi_definition" : "\"\\033[93m\"" },
        { "color_name" : "light blue" , "ansi_definition" : "\"\\033[94m\"" },
        { "color_name" : "light magenta" , "ansi_definition" : "\"\\033[95m\"" },
        { "color_name" : "light cyan" , "ansi_definition" : "\"\\033[96m\"" },
        { "color_name" : "light white" , "ansi_definition" : "\"\\033[97m\"" }
    ],
    "matches" : [
    {"expression": "Error|ERROR", "color" : "red" },
    {"expression": "DEBUG", "color" : "green" },
    {"expression": "Warning", "color" : "yellow" },
    {"expression": "Exception", "color" : "red" }
    ]
}
```

# Defaults

If no `ctail.config` file is find the `ctail` will default to some predefined values.

Assumed defaults:
```json
{
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
}
```

# Tests

Some performance tests still need to be done.
The utility is only tested on Mac OS X.