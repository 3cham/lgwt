package concurrency

import (
	"reflect"
	"testing"
	"time"
)

func mockWebsiteChecker(url string) bool {
	if url == "wjf://kjlkjaf.jslkfj" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {

	websites := []string{
		"http://google.com",
		"http://wikipedia.org",
		"http://otto.de",
		"wjf://kjlkjaf.jslkfj",
	}

	expected := map[string]bool{
		"http://google.com":    true,
		"http://wikipedia.org": true,
		"http://otto.de":       true,
		"wjf://kjlkjaf.jslkfj": false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Check websites returns incorrect! got %v, expected %v", got, expected)
	}
}

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}
