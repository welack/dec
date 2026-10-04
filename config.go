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

type SubmenuConfig struct {
	Name string
	Url  string
	Children []SubmenuConfig
}

type Config struct {
	Main    MainConfig
	Submenu []SubmenuConfig
}
