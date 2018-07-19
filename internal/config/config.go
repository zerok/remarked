package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Sample is the content of the remark.yml file as generated with the
// --init flag.
const Sample = `
# Set the title that should be rendered in the browser's title bar.
title: Slides

# Set the path to a local CSS file for remarked to serve.
# stylesheet: style.css

# Markdown file
markdownFile: slides.md

# RemarkJS allows overriding the default RemarkJS file.
# remarkJS: "https://remarkjs.com/downloads/remark-latest.min.js"

# The folder that should be served under /static. Default: none
# staticFolder: ./static

# If you want to use some of Go's template constructs inside the
# markdown file, enable this:
# markdownAsTemplate: false

# leftActionDelimiter: "{{"
# rightActionDelimiter: "}}"
`

// Config is usually the content of a remarked.yml file. Pretty much
// everything in here can be overloaded with commandline parameters.
type Config struct {
	Title        string `yaml:"title"`
	Stylesheet   string `yaml:"stylesheet"`
	MarkdownFile string `yaml:"markdownFile"`
	RemarkJS     string `yaml:"remarkjs"`

	// StaticFolder specifies the folder which should be served under the
	// /static mountpoint.
	StaticFolder string `yaml:"staticFolder"`

	// The Token should not be read from the config file but should instead be
	// either generated or explicitly set through the command-line flag.
	Token string `yaml:"-"`

	// The FinalStylesheet is the URL of the stylesheet as it is being served
	// by the HTTP server.
	FinalStylesheet string `yaml:"-"`

	MarkdownAsTemplate   bool   `yaml:"markdownAsTemplate"`
	LeftActionDelimiter  string `yaml:"leftActionDelimiter"`
	RightActionDelimiter string `yaml:"rightActionDelimiter"`

	TemplateFile string `yaml:"templateFile"`
}

func (c *Config) String() string {
	return fmt.Sprintf("<Config Title={%v} Stylesheet={%v} MarkdownFile={%v} RemarkJS={%v} Token={%v} FinalStylesheet={%v}>", c.Title, c.Stylesheet, c.MarkdownFile, c.RemarkJS, c.Token, c.FinalStylesheet)
}

// LoadFromPath generates a new Config object from the YAML file available
// through the given path.
func LoadFromPath(path string) (*Config, error) {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if c.LeftActionDelimiter == "" {
		c.LeftActionDelimiter = "{{"
	}
	if c.RightActionDelimiter == "" {
		c.RightActionDelimiter = "}}"
	}
	return &c, yaml.Unmarshal(data, &c)
}
