package nft

import (
	"fmt"
	"testing"
)

func TestPostPost(t *testing.T) {
	createPost("", "testuser", "testurl", "testcontent")
	createPost("", "testuser", "testurl2", "testcontent2")
	createPost("", "testuser", "testurl3", "testcontent3")
	createPost("", "testuser", "testurl4", "testcontent4")
}

func TestGetPost(t *testing.T) {
	posts, err := getPost(114485557671821312)
	if err != nil {
		t.Errorf("Error getting posts: %v", err)
		return
	}
	fmt.Println(posts)
	// for _, post := range posts {
	// 	t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
	// }
}

func TestGetPosts(t *testing.T) {
	posts, err := getPosts([]int64{114485557671821312, 114485557672017920, 114485557672083457})
	if err != nil {
		t.Errorf("Error getting posts: %v", err)
		return
	}
	fmt.Println(posts)
	for _, post := range posts {
		t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
	}
}

func TestGetPostsByTag(t *testing.T) {
	posts, err := getPostsByTagWithJoin("tag", 0, 0)
	if err != nil {
		t.Errorf("Error getting posts: %v", err)
		return
	}
	fmt.Println(posts)
	for _, post := range posts {
		t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
	}
}

func TestGetNewPosts(t *testing.T) {
	{
		posts, err := getNewPosts(0, 1)
		if err != nil {
			t.Errorf("Error getting posts: %v", err)
			return
		}
		fmt.Println(posts)
		for _, post := range posts {
			t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
		}
	}

	{
		posts, err := getNewPosts(114494234766999552, 1)
		if err != nil {
			t.Errorf("Error getting posts: %v", err)
			return
		}
		fmt.Println(posts)
		for _, post := range posts {
			t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
		}
	}
	{
		posts, err := getNewPosts(114494233645154304, 1)
		if err != nil {
			t.Errorf("Error getting posts: %v", err)
			return
		}
		fmt.Println(posts)
		for _, post := range posts {
			t.Logf("Post ID: %s, URL: %s, Content: %s", post.ID, post.URL, post.Content)
		}
	}
}
