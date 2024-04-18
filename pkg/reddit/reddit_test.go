package reddit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetPosts(t *testing.T) {
	// Create a new HTTP server to mock the Reddit APId
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with sample JSON data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
            "data": {
                "children": [
                    {
                        "data": {
                            "id": "post1",
                            "author": "user1",
                            "ups": 100
                        }
                    },
                    {
                        "data": {
                            "id": "post2",
                            "author": "user2",
                            "ups": 50
                        }
                    }
                ]
            }
        }`))
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, "", "", "")

	posts, err := client.GetPosts(nil, "testsubreddit")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	if posts[0].Author != "user1" {
		t.Errorf("Expected author of the first post to be user1, got %s", posts[0].Author)
	}

	if posts[1].UpVotes != 50 {
		t.Errorf("Expected upvotes of the second post to be 50, got %d", posts[1].UpVotes)
	}
}
