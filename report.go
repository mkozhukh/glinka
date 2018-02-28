package main

import (
	"fmt"
	"github.com/fatih/color"
)

type spiderReport struct {
	externals map[string]int
	warnings  []string
	infos     []string
	errors    []string
}

func newSpiderReport() *spiderReport {
	rep := spiderReport{}
	rep.infos = []string{}
	rep.warnings = []string{}
	rep.errors = []string{}
	rep.externals = make(map[string]int)

	return &rep
}

func (r *spiderReport) addError(s string) {
	r.errors = append(r.errors, s)
}

func (r *spiderReport) addWarning(s string) {
	r.warnings = append(r.warnings, s)
}

func (r *spiderReport) addExternal(s string) {
	r.externals[s] = r.externals[s] + 1
}

func (r *spiderReport) toString() {
	if len(r.errors) > 0 {
		for i := range r.errors {
			color.Red(r.errors[i] + "\n")
		}
	}
	if len(r.warnings) > 0 {
		for i := range r.warnings {
			color.Yellow(r.warnings[i] + "\n")
		}
	}

	if len(r.externals) > 0 {
		color.Cyan("External links:\n")
		for link, count := range r.externals {
			fmt.Printf("[%d] - %s\n", count, link)
		}
	}
}
