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

func chunk(s []int, n int) (store [][]int) {
	for i := 0; i < len(s); {
		if i+n >= len(s) {
			store = append(store, s[i:])
		} else {
			store = append(store, s[i:i+n])
		}
		i += n
	}
	return
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

// func convertToValue(str string)

// func (l *LasType) Data() {
// 	// const s = await this.blobString;
// 	//   const hds = await this.header();
// 	hds, err := l.Header()
// 	if err != nil {

// 	}
// 	sB := regexp.MustCompile("~A(?:\\w*\\s*)*\n").Split(l.content, 2)[1]
// 	sBs := regexp.MustCompile("\\s+").Split(strings.TrimSpace(sB), -1)
// 	//   const totalheadersLength = hds.length;
// 	//   const sB = (s as string)
// 	//     .split(/~A(?:\w*\s*)*\n/)[1]
// 	//     .trim()
// 	//     .split(/\s+/)
// 	//    // .map(m => Las.convertToValue(m.trim()));
// 	//   if (sB.length < 0) {
// 	//     throw new LasError('No data/~A section in the file');
// 	//   }
// 	//   const con = Las.chunk(sB, totalheadersLength);
// 	//   return con;
// }

func main() {
	test := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(chunk(test, 3))
}
