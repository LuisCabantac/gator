package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/LuisCabantac/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	for _, ff := range feedFollows {
		if ff.Url == feed.Url {
			fmt.Printf("%s\n%s\n", ff.FeedName, ff.UserName)
			return nil
		}
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("%s\n%s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	err = s.db.DeleteFeedFollowByUserAndFeedID(ctx, database.DeleteFeedFollowByUserAndFeedIDParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	for _, ff := range feedFollows {
		fmt.Println(ff.FeedName)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	for _, f := range feeds {
		fmt.Printf("%s\n%s\n%s\n", f.Name, f.Url, f.UserName)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	length := 2

	if len(cmd.Args) > 1 {
		rawLength := cmd.Args[0]
		inputLength, err := strconv.Atoi(rawLength)
		if err == nil || inputLength < 1 {
			length = inputLength
		}
	}

	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(length),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%+v\n", post)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, f := range feed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, f.PubDate)
		if err != nil {
			continue
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			PublishedAt: publishedAt,
			Title:       f.Title,
			Url:         f.Link,
			Description: f.Description,
			FeedID: uuid.NullUUID{
				UUID:  nextFeed.ID,
				Valid: true,
			},
		})
		if err != nil {
			continue
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("couldn't add feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't add feed follow: %w", err)
	}

	return nil
}
