package main

import (
	"reflect"
	"testing"
	"time"
)

func TestCheckWebsite(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://news.ycombinator.com",
		"http://github.com",
	}
	want := map[string]bool{
		"http://google.com":           true,
		"http://news.ycombinator.com": true,
		"http://github.com":           false,
	}
	got := checkWebsites(mockWebsiteChecker, websites)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < 100; i++ {
		urls[i] = "an url"
	}
	b.ResetTimer() // reset the counter before actually run
	for i := 0; i < b.N; i++ {
		checkWebsites(slowStubWebsitechecker, urls)
	}
}

func mockWebsiteChecker(site string) bool {
	if site == "http://github.com" {
		return false
	}
	return true
}

func slowStubWebsitechecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}
