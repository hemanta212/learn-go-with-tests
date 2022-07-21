package main

import (
	"fmt"
	"os"
	"strings"

	blogposts "github.com/hemanta212/blogposts"
)

func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
	if err != nil {
		fmt.Println("Error reading...")
	}
	// log.Println(posts)
	for _, post := range posts {
		fmt.Printf("------ %s ------\n", post.Title)
		fmt.Printf("%q\n", post.Description+"...")
		fmt.Printf(":: %s\n", strings.Join(post.Tags, ", "))
		fmt.Printf("\t%s\n\t\t---------------\n\n", post.Body)
	}
}
