package data

import "testing"

// TestDraftPostFiltering seeds the in-memory catalog directly (no DB needed)
// and verifies test-slug drafts stay out of every public surface.
func TestDraftPostFiltering(t *testing.T) {
	mu.Lock()
	oldPosts, oldBySlug := load.Posts, load.postBySlug
	load.Posts = []Post{
		{Slug: "test-scorpion-season", Title: "Test Post"},
		{Slug: "ants-in-spring", Title: "Real Post"},
	}
	load.postBySlug = mapPosts(load.Posts)
	mu.Unlock()
	defer func() {
		mu.Lock()
		load.Posts, load.postBySlug = oldPosts, oldBySlug
		mu.Unlock()
	}()

	if !IsDraftPost(Post{Slug: "test-x"}) {
		t.Error("IsDraftPost: test- prefix should be a draft")
	}
	if IsDraftPost(Post{Slug: "ants-in-spring"}) {
		t.Error("IsDraftPost: normal slug should not be a draft")
	}

	pubs := PublishedPosts()
	if len(pubs) != 1 || pubs[0].Slug != "ants-in-spring" {
		t.Errorf("PublishedPosts leaked a draft: %+v", pubs)
	}
	if _, ok := PublishedPostBySlug("test-scorpion-season"); ok {
		t.Error("PublishedPostBySlug: draft slug must not resolve")
	}
	if _, ok := PublishedPostBySlug("ants-in-spring"); !ok {
		t.Error("PublishedPostBySlug: real slug must resolve")
	}
	if _, ok := PublishedPostBySlug("nope"); ok {
		t.Error("PublishedPostBySlug: unknown slug must not resolve")
	}
}
