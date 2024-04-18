package reddit

import (
	"fmt"
	"sort"
	"sync"
)

type Statistics struct {
	PostsByUser    map[string]int
	PostsByUpvotes map[string]int
	mutex          sync.Mutex
}

func NewStatistics() *Statistics {
	return &Statistics{
		PostsByUser:    make(map[string]int),
		PostsByUpvotes: make(map[string]int),
	}
}

func (s *Statistics) Update(post Post) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.PostsByUser[post.Author]++
	s.PostsByUpvotes[post.ID] = post.UpVotes
}

func (s *Statistics) TopUsers(n int) []string {
	type userPostCount struct {
		User  string
		Count int
	}

	var userCounts []userPostCount
	for user, count := range s.PostsByUser {
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

func (s *Statistics) TopPosts(n int) []string {
	type postUpvotes struct {
		PostID  string
		UpVotes int
	}

	var postUpvotesList []postUpvotes
	for postID, upvotes := range s.PostsByUpvotes {
		postUpvotesList = append(postUpvotesList, postUpvotes{PostID: postID, UpVotes: upvotes})
	}

	sort.Slice(postUpvotesList, func(i, j int) bool {
		return postUpvotesList[i].UpVotes > postUpvotesList[j].UpVotes
	})

	var topPosts []string
	for i := 0; i < n && i < len(postUpvotesList); i++ {
		topPosts = append(topPosts, postUpvotesList[i].PostID)
	}

	return topPosts
}

func (s *Statistics) Report() {
	fmt.Println("Statistics report:")
	fmt.Println("Top users:")
	topUsers := s.TopUsers(5) // Report top 5 users
	for i, user := range topUsers {
		fmt.Printf("%d. %s\n", i+1, user)
	}
	fmt.Println("Top posts:")
	topPosts := s.TopPosts(5) // Report top 5 posts by upvotes
	for i, postID := range topPosts {
		fmt.Printf("%d. %s (Upvotes: %d)\n", i+1, postID, s.PostsByUpvotes[postID])
	}
	fmt.Println("--------------------")
}
