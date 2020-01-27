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

func removeComment(str string) []string {
	trimmedStr := strings.TrimSpace(str)
	strVec := strings.Split(trimmedStr, "\n")
	var result []string
	for _, line := range strVec {
		fStr := strings.TrimSpace(line)
		if !strings.HasPrefix(fStr, "#") && len(fStr) > 0 {
			result = append(result, fStr)
		}
	}
	return result
}

func chunk(s []string, n int) (store [][]string) {
	for i := 0; i < len(s); i += n {
		if i+n >= len(s) {
			store = append(store, s[i:])
		} else {
			store = append(store, s[i:i+n])
		}
	}
	return
}

//Header - returns the header in the las file
func (l *LasType) Header() ([]string, error) {
	err := "file has no header content, make sure to create an instance of LasType with Las function"
	if len(l.content) < 1 {
		return []string{}, fmt.Errorf(err)
	}
	hP := regexp.MustCompile("~C(?:\\w*\\s*)*\n\\s*")
	spl := strings.Split(hP.Split(l.content, 2)[1], "~")[0]
	headers := removeComment(spl)
	if len(headers) < 1 {
		return []string{}, fmt.Errorf(err)
	}
	for i, val := range headers {
		headers[i] = regexp.MustCompile("\\s+[.]").Split(strings.TrimSpace(val), 2)[0]
	}
	return headers, nil
}

// Data returns the data section in the file
func (l *LasType) Data() [][]string {
	hds, err := l.Header()
	if err != nil {
		panic("No data in file")
	}
	sB := regexp.MustCompile("~A(?:\\w*\\s*)*\n").Split(l.content, 2)[1]
	sBs := regexp.MustCompile("\\s+").Split(strings.TrimSpace(sB), -1)
	return chunk(sBs, len(hds))
}

// Version - returns the version of the las file
func (l *LasType) Version() (version string) {
	version, _ = metadata(l.content)
	return
}

// Wrap - returns the version of the las file
func (l *LasType) Wrap() (wrap bool) {
	_, wrap = metadata(l.content)
	return
}

// ColumnCount - Returns the number of columns in a .las file
func (l *LasType) ColumnCount() (count int) {
	count = len(l.Data())
	return
}

// RowCount - Returns the number of rowa in a .las file
func (l *LasType) RowCount() (count int) {
	header, _ := l.Header()
	// TODO: handle error
	count = len(header)
	return
}

// metadata - picks out version and wrap state of the file
func metadata(str string) (version string, wrap bool) {
	sB := strings.Split(regexp.MustCompile("~V(?:\\w*\\s*)*\n\\s*").Split(str, 2)[1], "~")[0]
	sw := removeComment(sB)
	accum := [][]string{}
	for _, val := range sw {
		current := regexp.MustCompile("\\s{2,}|\\s*:").Split(val, -1)[0:2]
		accum = append(accum, current)
	}
	version = accum[0][1]
	if strings.ToLower(accum[1][1]) == "yes" {
		wrap = true
	} else {
		wrap = false
	}
	return
}

func main() {
	las, err := Las("sample/example1.las")
	if err != nil {
		panic(err)
	}
	fmt.Println(metadata(las.content))
}
