package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

//LasType is the main type definition
type LasType struct {
	path    string
	content string
}

//Las creates an instance of LasType
func Las(path string) (*LasType, error) {
	bs, err := ioutil.ReadFile(path)
	l := LasType{path: path}
	l.content = string(bs)
	return &l, err
}

func removeComment(str string) string {
	trimmedStr := strings.Trim(str, " ")
	strVec := strings.Split(trimmedStr, "\n")
	var result strings.Builder
	for _, line := range strVec {
		strTrimLeft := strings.TrimLeft(line, " ")
		if strTrimLeft[0] != '#' {
			result.WriteString(fmt.Sprintf("%s\n", strTrimLeft))
		}
	}
	return result.String()
}

func main() {
}
