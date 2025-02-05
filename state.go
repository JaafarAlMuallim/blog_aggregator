package main

import (
	"github.com/JaafarAlMuallim/blog_agg/internal/config"
	"github.com/JaafarAlMuallim/blog_agg/internal/database"
)

type state struct {
	cfg       *config.Config
	dbQueries *database.Queries
}

type command struct {
	name string
	args []string
}
