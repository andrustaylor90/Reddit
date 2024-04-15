package reddit

import (
	"fmt"
	"sort"
	"time"
)

func (c *Client) UpdateStatistics(posts []Post) {
	for _, post := range posts {
		c.PostsByUser[post.Author]++
	}
}

func (c *Client) TopUsers(n int) []string {
	type userPostCount struct {
		User  string
		Count int
	}

	var userCounts []userPostCount
	for user, count := range c.PostsByUser {
		userCounts = append(userCounts, userPostCount{User: user, Count: count})
	}

	sort.Slice(userCounts, func(i, j int) bool {
		return userCounts[i].Count > userCounts[j].Count
	})

	var topUsers []string
	for i := 0; i < n && i < len(userCounts); i++ {
		topUsers = append(topUsers, userCounts[i].User)
	}

	return topUsers
}

func (c *Client) ReportStatistics() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Statistics report:")
			fmt.Println("Top users:")
			topUsers := c.TopUsers(5) // Report top 5 users
			for i, user := range topUsers {
				fmt.Printf("%d. %s\n", i+1, user)
			}
			fmt.Println("--------------------")
		}
	}
}
