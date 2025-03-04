package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/VincNT21/aggreGATOR/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	// Ensure that not more than one argument is passed
	if len(cmd.Args) > 1 {
		return fmt.Errorf("browse command expect no more than one (optional) argument -> <limit>")
	}
	// Deal with no limit given in command args or parse the given number to int32
	var limit int32
	limit = 2
	if len(cmd.Args) == 1 {
		i, err := strconv.ParseInt(cmd.Args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("could'nt parse args to int32: %w", err)
		}
		limit = int32(i)
	}

	// Call get posts query
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts by user: %w", err)
	}

	// Print results
	fmt.Printf("==> Found %d posts from feeds followed by %s\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf(" * %s\n", post.Title)
		fmt.Printf("      published on %v\n", post.PublishedAt.Format("Mon Jan 2"))
		fmt.Printf("      on %s\n", post.FeedName)
		fmt.Printf("      Link: %s\n", post.Url)
		fmt.Println("============")
	}

	return nil
}
