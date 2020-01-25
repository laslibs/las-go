package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
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
	trimmedStr := strings.TrimSpace(str)
	strVec := strings.Split(trimmedStr, "\n")
	var result strings.Builder
	for _, line := range strVec {
		strTrimLeft := strings.TrimSpace(line)
		if !strings.HasPrefix(strTrimLeft, "#") {
			result.WriteString(fmt.Sprintf("%s\n", strTrimLeft))
		}
	}
	return strings.TrimSpace(result.String())
}

//Header - returns the header in the las file
func (l *LasType) Header() ([]string, error) {
	if len(l.content) < 1 {
		return []string{}, fmt.Errorf("file has no header content, make sure to create an instance of LasType with Las function")
	}
	hP := regexp.MustCompile("~C(?:\\w*\\s*)*\n\\s*")
	spl := hP.Split(l.content, 2)[1]
	hStr := strings.Split(spl, "~")[0]
	hStr = removeComment(hStr)
	if len(hStr) < 1 {
		return []string{}, fmt.Errorf("file has no header content, make sure to create an instance of LasType with Las function")
	}
	headers := strings.Split(hStr, "\n")
	for i, val := range headers {
		headers[i] = regexp.MustCompile("\\s+[.]").Split(strings.TrimSpace(val), 2)[0]
	}
	return headers, nil
}

func main() {
	test, err := Las("sample/example1.las")
	if err != nil {
		panic(err)
	}
	header, hErr := test.Header()
	if hErr != nil {
		panic(hErr)
	}
	fmt.Println(header)
}
