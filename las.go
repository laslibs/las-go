package main

import (
	"errors"
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

// index returns the position of an element in a slice of strings
func index(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
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

func pattern(str string) *regexp.Regexp {
	return regexp.MustCompile(str)
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

// WellProps contains basic definition of a single well measurement
type WellProps struct {
	unit        string
	description string
	value       string
}

//Header - returns the header in the las file
func (l *LasType) Header() ([]string, error) {
	err := "file has no header content, make sure to create an instance of LasType with Las function"
	if len(l.content) < 1 {
		return []string{}, fmt.Errorf(err)
	}
	hP := pattern("~C(?:\\w*\\s*)*\n\\s*")
	spl := strings.Split(hP.Split(l.content, 2)[1], "~")[0]
	headers := removeComment(spl)
	if len(headers) < 1 {
		return []string{}, fmt.Errorf(err)
	}
	for i, val := range headers {
		headers[i] = pattern("\\s+[.]").Split(strings.TrimSpace(val), 2)[0]
	}
	return headers, nil
}

// Data returns the data section in the file
func (l *LasType) Data() [][]string {
	hds, err := l.Header()
	if err != nil {
		panic("No data in file")
	}
	sB := pattern("~A(?:\\w*\\s*)*\n").Split(l.content, 2)[1]
	sBs := pattern("\\s+").Split(strings.TrimSpace(sB), -1)
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

// Column returns entry of an individual log, say gamma ray
func (l *LasType) Column(key string) []string {
	header, _ := l.Header()
	// TODO: handle error
	data := l.Data()
	var res []string
	keyIndex := index(header, key)
	// TODO: handle when keyIndex is -1
	for _, val := range data {
		res = append(res, val[keyIndex])
	}
	return res
}

// func (l *LasType) HeaderAndDesc() {
// 	// const cur = (await this.property('curve')) as object;
// 	curve, _ := property(l.content, "curve")
// 	// TODO: handle error
// 	//   const hd = Object.keys(cur);
// 	//   const descr = Object.values(cur).map((c, i) => (c.description === 'none' ? hd[i] : c.description));
// 	//   const obj: { [key: string]: string } = {};
// 	//   hd.map((_, i) => (obj[hd[i]] = descr[i]));
// 	//   if (Object.keys(obj).length < 0) {
// 	//     throw new LasError('Poorly formatted ~curve section in the file');
// 	//   }
// 	//   return obj;
// }

// CurveParams - Returns Curve Parameters
func (l *LasType) CurveParams() map[string]WellProps {
	curve, _ := property(l.content, "curve")
	// TODO: handle error
	return curve
}

// WellParams - Returns Overrall Well Parameters
func (l *LasType) WellParams() map[string]WellProps {
	well, _ := property(l.content, "well")
	// TODO: handle error
	return well
}

// LogParams - Returns Log Parameters
func (l *LasType) LogParams() map[string]WellProps {
	param, _ := property(l.content, "param")
	// TODO: handle error
	return param
}

// metadata - picks out version and wrap state of the file
func metadata(str string) (version string, wrap bool) {
	sB := strings.Split(pattern("~V(?:\\w*\\s*)*\n\\s*").Split(str, 2)[1], "~")[0]
	sw := removeComment(sB)
	accum := [][]string{}
	for _, val := range sw {
		current := pattern("\\s{2,}|\\s*:").Split(val, -1)[0:2]
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

func property(str string, key string) (property map[string]WellProps, err error) {
	err = errors.New("property cannot be found")
	property = make(map[string]WellProps)

	regDict := map[string]string{
		"curve": "~C(?:\\w*\\s*)*\\n\\s*",
		"param": "~P(?:\\w*\\s*)*\\n\\s*",
		"well":  "~W(?:\\w*\\s*)*\\n\\s*",
	}
	prop, ok := regDict[key]
	if !ok {
		return
	}
	substr := pattern(prop).Split(str, 2)
	var sw []string
	if len(substr) > 1 {
		sw = removeComment(strings.Split(substr[1], "~")[0])
	}
	if len(sw) > 0 {
		for _, val := range sw {
			root := pattern("\\s*[.]\\s+").ReplaceAllString(val, "   none   ")
			title := pattern("[.]|\\s+").Split(root, 2)[0]
			unit := pattern("\\s+").Split(pattern("^\\w+\\s*[.]*s*").Split(root, 2)[1], 2)[0]
			desc := strings.TrimSpace(strings.Split(root, ":")[1])
			desc = pattern("\\d+\\s*").ReplaceAllString(desc, "")
			if len(desc) < 1 {
				desc = "none"
			}
			vD := pattern("\\s{2,}\\w*\\s{2,}").Split(strings.Split(root, ":")[0], -1)
			var value string
			if len(vD) > 2 && len(vD[len(vD)-1]) > 0 {
				value = strings.TrimSpace(vD[len(vD)-2])
			} else {
				value = strings.TrimSpace(vD[len(vD)-1])
			}
			property[title] = WellProps{unit, desc, value}
		}
		return property, nil
	}
	return
}

func main() {
	las, err := Las("sample/example1.las")
	if err != nil {
		panic(err)
	}
	fmt.Println(las.CurveParams())
}
