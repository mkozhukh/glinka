package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func reportErrors(data *LinksStore) string {
	var errorsText string

	for _, link := range data.Records {
		if link.Status == StatusError {
			from := strings.Join(getParents(data, link.URL, 1), "; ")
			errorsText += fmt.Sprintf("%3d %s\nfrom %s\n%s\n\n", link.Status, link.URL, from, link.Error)
		}
	}

	return errorsText
}

func reportStats(data *LinksStore) string {
	text := "\n"

	type kv struct {
		Key   string
		Value int
	}
	var total, linktotal, errors, mixed, external, binary int

	domains := make(map[string]int)
	files := []kv{}
	mixedPages := []string{}

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
			mixedPages = append(mixedPages, link.URL)
		case StatusError:
			errors++
		default:
			total++
			if link.Links != nil {
				linktotal += len(link.Links)
			}
		}
	}

	text += fmt.Sprintf("%d errors in %d pages ( %d links )\n", errors, total, linktotal)
	text += fmt.Sprintf("there were %d file links and %d external links\n", binary, external)

	if mixed != 0 {
		text += color.RedString("\nMixed content detected\n")
		for i := range mixedPages {
			text += " - " + mixedPages[i] + "\n"
		}

	}
	if external != 0 {
		text += "\nExternal domains\n"

		var ss []kv
		for k, v := range domains {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		for _, kv := range ss[0:5] {
			text += fmt.Sprintf("%4d : %s\n", kv.Value, kv.Key)
		}
	}

	if binary != 0 {
		text += "\nLinked files\n"
		sort.Slice(files, func(i, j int) bool {
			return files[i].Value > files[j].Value
		})

		for _, file := range files[0:5] {
			text += fmt.Sprintf(" - %4d : %s\n", file.Value, file.Key)
		}
	}

	return text
}
