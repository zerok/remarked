package main

import (
	"html/template"
	"io/ioutil"
	"strconv"
	"strings"
)

type templateFuncs struct{}

func (f *templateFuncs) FuncMap() template.FuncMap {
	return template.FuncMap{
		"loadCode":  f.LoadCode,
		"markLines": f.MarkLines,
		"counter":   f.Counter,
	}
}

type count struct {
	Current int
	First   bool
	Last    bool
}

func (f *templateFuncs) Counter(from int, to int, step int) ([]count, error) {
	result := make([]count, 0, 0)
	cur := from
	for {
		result = append(result, count{
			Current: cur,
			First:   step == from,
			Last:    step == to,
		})
		if cur == to {
			break
		}
		cur += step
	}
	return result, nil
}

func (f *templateFuncs) LoadCode(path string) (template.HTML, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return template.HTML(data), nil
}

func (f *templateFuncs) MarkLines(lineNumbers, data template.HTML) (template.HTML, error) {
	var result []string
	parsedLineNumbers := parseLineRanges(string(lineNumbers))
	for idx, line := range strings.Split(string(data), "\n") {
		_, highlight := parsedLineNumbers[idx+1]
		result = append(result, highlightLine(line, highlight))
	}
	return template.HTML(strings.Join(result, "\n")), nil
}

func parseLineRanges(ranges string) map[int]struct{} {
	result := make(map[int]struct{})
	rangeSegments := strings.Split(ranges, ",")
	for _, r := range rangeSegments {
		segments := strings.SplitN(r, "-", 2)
		var err error
		var start, end int64
		switch len(segments) {
		case 1:
			line, err := strconv.ParseInt(r, 10, 32)
			if err != nil {
				continue
			}
			start, end = line, line
		case 2:
			start, err = strconv.ParseInt(segments[0], 10, 32)
			if err != nil {
				continue
			}
			end, err = strconv.ParseInt(segments[1], 10, 32)
			if err != nil {
				continue
			}
		}

		if start != 0 {
			for i := start; i <= end; i++ {
				result[int(i)] = struct{}{}
			}
		}
	}
	return result
}

func highlightLine(line string, doHighlight bool) string {
	if !doHighlight {
		return line
	}
	return "*" + strings.TrimPrefix(line, " ")
}
