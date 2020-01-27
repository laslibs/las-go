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
	if len(l.content) < 1 {
		return []string{}, fmt.Errorf("file has no header content, make sure to create an instance of LasType with Las function")
	}
	hP := regexp.MustCompile("~C(?:\\w*\\s*)*\n\\s*")
	spl := strings.Split(hP.Split(l.content, 2)[1], "~")[0]
	headers := removeComment(spl)
	if len(headers) < 1 {
		return []string{}, fmt.Errorf("file has no header content, make sure to create an instance of LasType with Las function")
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

// func metadata(str string) {
// 	// const str = await this.blobString;
// 	// const sB = (str as string)
// 	// 	.trim()
// 	// 	.split(/~V(?:\w*\s*)*\n\s*/)[1]
// 	// 	.split(/~/)[0];
// 	sB := strings.Split(regexp.MustCompile("~V(?:\\w*\\s*)*\n\\s*").Split(str, 2)[1], "~")[0]
// 	// const sw = Las.removeComment(sB);
// 	sw := removeComment(sB)
// 	// const refined = sw
// 	// 	.split('\n')
// 	// 	.map(m => m.split(/\s{2,}|\s*:/).slice(0, 2))
// 	// 	.filter(f => Boolean(f));
// 	for _, val := range sw {

// 	}
// 	// const res = refined.map(r => r[1]);
// 	// const wrap = res[1].toLowerCase() === 'yes' ? true : false;
// 	// if ([+res[0], wrap].length < 0) {
// 	// 	throw new LasError("Couldn't get metadata");
// 	// }
// 	// return [+res[0], wrap];
// }

func main() {
	las, err := Las("sample/example1.las")
	if err != nil {
		panic(err)
	}
	fmt.Println(las.Data())
}
