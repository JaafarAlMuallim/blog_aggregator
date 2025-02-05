package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	// "time"
	//
	// "github.com/JaafarAlMuallim/blog_agg/internal/database"
	"github.com/JaafarAlMuallim/blog_agg/internal/database"
	"github.com/JaafarAlMuallim/blog_agg/internal/rss"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	val, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rssFeed rss.RSSFeed
	if err := xml.Unmarshal(val, &rssFeed); err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.dbQueries.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	err = s.dbQueries.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return err
	}
	fmt.Println("Fetching " + feed.Url)
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	for _, item := range rssFeed.Channel.Item {
		pub, err := time.Parse(layout, item.PubDate)
		if err != nil {
			return err
		}
		_, err = s.dbQueries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: pub,
			FeedID:      feed.ID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
