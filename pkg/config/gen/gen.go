package main

import (
	cfg "github.com/conductorone/baton-trello/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("trello", cfg.Config)
}
