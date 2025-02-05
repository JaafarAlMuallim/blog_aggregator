package main

import (
	"context"

	"github.com/JaafarAlMuallim/blog_agg/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.dbQueries.GetUser(context.Background(), s.cfg.User)
		if err != nil {
			return err
		}
		if err := handler(s, cmd, user); err != nil {
			return err
		}
		return nil
	}
}
