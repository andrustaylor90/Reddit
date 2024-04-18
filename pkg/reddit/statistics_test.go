package reddit

import (
	"testing"
)

func TestStatistics_Update(t *testing.T) {
	stats := NewStatistics()

	stats.Update(Post{Author: "user1"})
	stats.Update(Post{Author: "user2"})
	stats.Update(Post{Author: "user1"})

	if stats.PostsByUser["user1"] != 2 {
		t.Errorf("Expected 2 posts for user1, got %d", stats.PostsByUser["user1"])
	}
	if stats.PostsByUser["user2"] != 1 {
		t.Errorf("Expected 1 post for user2, got %d", stats.PostsByUser["user2"])
	}
}

func TestStatistics_TopUsers(t *testing.T) {
	stats := NewStatistics()

	stats.Update(Post{Author: "user1"})
	stats.Update(Post{Author: "user2"})
	stats.Update(Post{Author: "user1"})

	topUsers := stats.TopUsers(2)

	if len(topUsers) != 2 {
		t.Errorf("Expected 2 top users, got %d", len(topUsers))
	}

	if topUsers[0] != "user1" {
		t.Errorf("Expected top user to be user1, got %s", topUsers[0])
	}

	if topUsers[1] != "user2" {
		t.Errorf("Expected second top user to be user2, got %s", topUsers[1])
	}
}
