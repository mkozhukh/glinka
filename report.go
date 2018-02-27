package main

import "log"

type spiderReport struct {
	warnings []string
	infos    []string
	errors   []string
}

func newSpiderReport() *spiderReport {
	rep := spiderReport{}
	rep.infos = []string{}
	rep.warnings = []string{}
	rep.errors = []string{}
	return &rep
}

func (r *spiderReport) addError(s string) {
	r.errors = append(r.errors, s)
}

func (r *spiderReport) addWarning(s string) {
	r.warnings = append(r.warnings, s)
}

func (r *spiderReport) addInfo(s string) {
	r.infos = append(r.infos, s)
}

func (r *spiderReport) toString() {
	log.Println("Errors:")
	for i := range r.errors {
		log.Println(r.errors[i])
	}

	log.Println("Warnings:")
	for i := range r.warnings {
		log.Println(r.warnings[i])
	}

	log.Println("Info:")
	for i := range r.infos {
		log.Println(r.infos[i])
	}
}
