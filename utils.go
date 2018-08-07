package main

func getParents(data *LinksStore, url string, count int) []string {
	p := make([]string, 0)
	for _, link := range data.Records {
		for i := range link.Links {
			if link.Links[i] == url {
				p = append(p, link.URL)
				if len(p) == count {
					return p
				}
			}
		}
	}

	return p
}
