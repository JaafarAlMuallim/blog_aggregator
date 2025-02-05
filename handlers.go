package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/JaafarAlMuallim/blog_agg/internal/database"
	"github.com/google/uuid"
)

func loginHandler(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		fmt.Println("Usage: login <username>")
		os.Exit(1)
	}
	user, err := state.dbQueries.GetUser(context.Background(), cmd.args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := state.cfg.SetUser(user.Name); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("User has been set")
	return nil
}

func regHandler(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		fmt.Println("Usage: register <username>")
		os.Exit(1)
	}
	val, err := state.dbQueries.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[1],
		})
	if err != nil {
		os.Exit(1)
		return err
	}
	if err := state.cfg.SetUser(val.Name); err != nil {
		return err
	}
	return nil
}

func resetHandler(state *state, cmd command) error {
	if err := state.dbQueries.DeleteUsers(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

func usersHandler(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		fmt.Println("Usage: follow url")
		os.Exit(1)
	}
	users, err := state.dbQueries.GetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if state.cfg.User == user.Name {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}
	return nil

}

func aggHandler(state *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Println("Add Time Between Req argument")
		os.Exit(1)
	}
	timeBetweenReq, err := time.ParseDuration(cmd.args[1])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting Feeds Every %v\n", timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C {
		err := scrapeFeeds(state)
		if err != nil {
			return err
		}
	}
}

func addFeedHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 3 {
		fmt.Println("Not enough args")
		os.Exit(1)
	}
	ctx := context.Background()
	feed, err := state.dbQueries.CreateFeed(ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[1],
			Url:       cmd.args[2],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return err
	}
	_, err = state.dbQueries.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)
	return nil

}

func feedsHandler(state *state, cmd command) error {
	feeds, err := state.dbQueries.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Feedname)
		fmt.Println(feed.Url)
		fmt.Println(feed.Username)
	}
	return nil
}

func followHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("Usage: follow <url>")
		os.Exit(1)
	}
	ctx := context.Background()
	feed, err := state.dbQueries.GetFeedsByUrl(ctx, cmd.args[1])
	if err != nil {
		return err
	}
	_, err = state.dbQueries.CreateFeedFollow(ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(feed.Name)
	fmt.Println(user.Name)

	return nil
}

func followingHandler(state *state, cmd command, user database.User) error {
	ctx := context.Background()
	follows, err := state.dbQueries.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return err
	}
	for _, follow := range follows {
		fmt.Println(follow.Feedname)
		fmt.Println(follow.Username)
	}
	return nil
}

func unfollowHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("Usage: unfollow <url>")
		os.Exit(1)
	}
	ctx := context.Background()
	err := state.dbQueries.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		Url:    cmd.args[1],
		UserID: user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("Delete Record")
	return nil
}

func browseHandler(state *state, cmd command) error {
	limit := 2
	if len(cmd.args) >= 2 {
		val, err := strconv.Atoi(cmd.args[1])
		if err != nil {
			limit = 2
		} else {
			limit = val
		}
	}
	ctx := context.Background()
	posts, err := state.dbQueries.GetPosts(ctx, int32(limit))
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Println(post.Title, post.PublishedAt)
	}
	return nil
}
