package reddit

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) ProcessPosts(ctx context.Context, subreddit string) {
	go c.ReportStatistics(ctx)

	for {
		fmt.Println("fetching-posts")
		_, err := c.GetPosts(ctx, subreddit)
		if err != nil {
			fmt.Println(err)
			continue
		}

		time.Sleep(1 * time.Minute)
	}
}

func (c *Client) ReportStatistics(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Statistics.Report()
		case <-ctx.Done():
			return
		}
	}
}
