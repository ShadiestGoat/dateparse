package main

import (
	"os"
	"strings"
	"text/template"
)

const tpl = `package dateparse

// AUTO-GENERATED, DO NOT TOUCH!

// Default date layouts
var DefaultLayoutsDate = []string{
{{ .date }}
}

// Default time layouts
var DefaultLayoutsTime = []string{
{{ .time }}
}
`

// - DDth/DDst/DDnd
// - (DD/)MM(/YY(YY))
//   - DD/MM
//   - DD/MM/YY
//   - DD/MM/YYYY
//   - MM/YY
//   - MM/YYYY
// - YYYY/MM/DD
// - All Formats above, but with `-` instead of `/`
// -
// - (DD(th/st/nd) (of)) Month ((of) YY(YY))
// - (DD(th/st/nd) (of)) Mon ((of) YY(YY))

func dates() []string {
	o := []string{
		"2006-01-02",
	}

	d := "_2"
	m := "1"
	yf := []string{"06", "2006"}

	// short form dates
	sd := [][]string{
		{d, m},
	}

	for _, y := range yf {
		sd = append(sd, [][]string{
			{d, m, y},
			{m, y},
		}...)
	}

	for _, f := range sd {
		o = append(o,
			strings.Join(f, "/"),
			strings.Join(f, "-"),
		)
	}

	// Natural dates
	mo := []string{"Jan", "January"}
	for _, m := range mo {
		o = append(o,
			m,
			strings.Join([]string{d, m}, " "),
		)
		for _, y := range yf {
			o = append(o,
				strings.Join([]string{d, m, y}, " "),
				strings.Join([]string{m, y}, " "),
			)
		}
	}

	return append(o, "2006")
}

const h24, h12 = "15", "3"

func times() []string {
	o := []string{h24, h12 + "PM"}
	p := []string{"4", "5"}

	for i := range p {
		v := p[:i+1]

		o = append(o,
			strings.Join(append([]string{h24}, v...), ":"),
			strings.Join(append([]string{h12}, v...), ":")+"PM",
		)
	}

	return o
}

func list(v []string) string {
	return "\t\"" + strings.Join(v, "\",\n\t\"") + `",`
}

func main() {
	t := template.New(`main`)
	t, err := t.Parse(tpl)
	if err != nil {
		panic("failed to parse tpl: " + err.Error())
	}

	f, err := os.Create("defaults.go")
	if err != nil {
		panic("Failed to touch defaults.go: " + err.Error())
	}

	err = t.Execute(f, map[string]string{
		"date": list(dates()),
		"time": list(times()),
	})

	if err != nil {
		panic("Failed to exec tpl: " + err.Error())
	}

	f.Close()
}
