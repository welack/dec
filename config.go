package main

import (
	_ "github.com/BurntSushi/toml"
)

type MainConfig struct {
	Title    string
	Subtitle string
	Url      string
	Output   string
	Amount   int
}

type ColorschemeConfig struct {
	BG             string `toml:"bg"`
	Text           string `toml:"text"`
	Primary        string `toml:"primary"`
	PrimaryHover   string `toml:"primary_hover"`
	PrimaryText    string `toml:"primary_text"`
	Muted          string `toml:"muted"`
	MutedOnPrimary string `toml:"muted_on_primary"`
	Border         string `toml:"border"`
}

type SubmenuConfig struct {
	Name     string
	Url      string
	Children []SubmenuConfig
}

type Config struct {
	Main        MainConfig
	Colorscheme ColorschemeConfig
	Submenu     []SubmenuConfig
}
