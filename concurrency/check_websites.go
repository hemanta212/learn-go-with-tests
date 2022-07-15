package main

type WebsiteChecker func(string) bool

type result struct {
	string
	bool
}

func checkWebsites(wc WebsiteChecker, sites []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, site := range sites {
		go func(site string) {
			resultChannel <- result{site, wc(site)}
		}(site)
	}

	for i := 0; i < len(sites); i++ {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}
