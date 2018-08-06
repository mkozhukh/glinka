package main

import (
	"fmt"
	"net/url"
	"sort"

	"github.com/fatih/color"
)

func reportStats(data *LinksStore) string {
	text := "\n"

	type kv struct {
		Key   string
		Value int
	}
	var total, linktotal, errors, mixed, external, binary int
	var errorsText string

	domains := make(map[string]int)
	files := []kv{}

	for _, link := range data.Records {
		switch link.Status {
		case StatusBinary:
			binary++
			files = append(files, kv{link.URL, link.Count})
		case StatusExternal:
			external++
			linkURL, _ := url.Parse(link.URL)
			domains[linkURL.Host] += link.Count
		case StatusMixedContent:
			mixed++
		case StatusError:
			errors++
			errorsText += link.Error + "\n'"
		default:
			total++
			if link.Links != nil {
				linktotal += len(link.Links)
			}
		}
	}

	text += fmt.Sprintf("%d errors in %d pages ( %d links )\n", errors, total, linktotal)
	text += fmt.Sprintf("there were %d file links and %d external links\n", binary, external)

	if external != 0 {
		text += "\nExternal domains\n"

		var ss []kv
		for k, v := range domains {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss {
			text += fmt.Sprintf("%4d : %s\n", kv.Value, kv.Key)
		}
	}

	if binary != 0 {
		text += "\nLinked files\n"
		sort.Slice(files, func(i, j int) bool {
			return files[i].Value > files[j].Value
		})

		for _, file := range files {
			text += fmt.Sprintf(" - %4d : %s", file.Value, file.Key)
		}
	}

	if errors != 0 {
		text += color.RedString(errorsText)
	}
	return text
}
