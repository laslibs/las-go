package lasgo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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

// DataStruct just like Data but returns output in specified struct format
func (l *LasType) DataStruct(opt *DataOptions) []interface{} {

	var (
		o     *DataOptions
		store [][]string
	)

	if opt != nil {
		o = opt
	}

	hds, err := l.Header()
	if err != nil {
		panic("No data in file")
	}

	store = l.Data()

	if o != nil {
		ctx := context.Background()

		output, err := structConvert(ctx, &store, hds, o)
		if err != nil {
			panic(err)
		}

		return output
	}

	return nil
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
	header, _ := l.Header()
	// TODO: handle error
	count = len(header)
	return

}

// RowCount - Returns the number of rowa in a .las file
func (l *LasType) RowCount() (count int) {
	count = len(l.Data())
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

// HeaderAndDesc return the name and description of each log entry
func (l *LasType) HeaderAndDesc() map[string]string {
	// const cur = (await this.property('curve')) as object;
	curve, _ := property(l.content, "curve")
	res := make(map[string]string)
	// TODO: handle error
	for key, val := range curve {
		res[key] = val.description
	}
	// TODO: handle when res is empty
	return res
}

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

// Other returns extra information stored in ~other section
func (l *LasType) Other() string {
	// TODO: make case insensitive
	som := pattern("~O(?:\\w*\\s*)*\n\\s*").Split(l.content, 2)
	if len(som) > 1 {
		res := pattern("\n\\s*").ReplaceAllString(strings.Split(som[1], "~")[0], " ")
		return strings.Join(removeComment(res), "\n")
	}
	return ""
}

// ToCSV creates a csv file using data and header
func (l *LasType) ToCSV(filename string) {
	file, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		panic(err)
	}
	// close file when function call ends
	defer file.Close()
	header, _ := l.Header()
	// TODO: handle error
	file.WriteString(strings.Join(header, ",") + "\n")
	for _, val := range l.Data() {
		// TODO: don't include \n at the last line
		file.WriteString(strings.Join(val, ",") + "\n")
	}
}
