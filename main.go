package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JaafarAlMuallim/blog_agg/internal/config"
	"github.com/JaafarAlMuallim/blog_agg/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s := state{cfg, dbQueries}
	commands := commands{handler: make(map[string]func(*state, command) error)}
	commands.register("reset", resetHandler)
	commands.register("login", loginHandler)
	commands.register("register", regHandler)
	commands.register("users", usersHandler)
	commands.register("agg", aggHandler)
	commands.register("addfeed", middlewareLoggedIn(addFeedHandler))
	commands.register("feeds", feedsHandler)
	commands.register("follow", middlewareLoggedIn(followHandler))
	commands.register("unfollow", middlewareLoggedIn(unfollowHandler))
	commands.register("following", middlewareLoggedIn(followingHandler))
	commands.register("browse", browseHandler)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not Enough Args")
		os.Exit(1)
	}
	cmd := command{name: args[1], args: args[1:]}
	if err := commands.run(&s, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
