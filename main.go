package main

import (
	"fmt"
	"sort"
)

func main() {
	config := readConfig()
	client := NewPlexClient(config)

	client.GetLibraries()
	media := client.GetMovies()
	videos := media.Videos

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].AddedAt < videos[j].AddedAt
	})

	v := videos[len(videos)-1]
	fmt.Printf("%d %s\n", v.AddedAt, v.Title)
	// for i, m := range videos {
	// 	fmt.Printf("[%3d] %s\n", i, m.Title)
	// }
}
