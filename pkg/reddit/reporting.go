package reddit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (c *Client) ProcessPosts(ctx context.Context, subreddit string) {
	var mu sync.Mutex

	go c.ReportStatistics()

	for {
		fmt.Println("fetching-posts")
		posts, err := c.GetPosts(ctx, subreddit)
		if err != nil {
			fmt.Println(err)
			continue
		}

		mu.Lock()
		c.UpdateStatistics(posts)
		mu.Unlock()

		time.Sleep(1 * time.Minute)
	}
}
