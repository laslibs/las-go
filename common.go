package lasgo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// StructorConfig is used to expose a subset of the configuration options
// provided by the mapstructure package.
//
// See: https://godoc.org/github.com/mitchellh/mapstructure#DecoderConfig
type StructorConfig struct {

	// DecodeHook, if set, will be called before any decoding and any
	// type conversion (if WeaklyTypedInput is on). This lets you modify
	// the values before they're set down onto the resulting struct.
	//
	// If an error is returned, the entire decode will fail with that
	// error.
	DecodeHook mapstructure.DecodeHookFunc

	// If WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//   - negative numbers to overflowed uint values (base 10)
	//   - slice of maps to a merged map
	//   - single values are converted to slices if required. Each
	//     element is weakly decoded. For example: "4" can become []int{4}
	//     if the target type is an int slice.
	//
	WeaklyTypedInput bool
}

// PostUnmarshaler allows you to further modify all results after unmarshaling.
// The ConcreteStruct pointer must implement this interface to make use of this feature.
type PostUnmarshaler interface {

	// PostUnmarshal is called for each row after all results have been fetched.
	// You can use it to further modify the values of each ConcreteStruct.
	PostUnmarshal(ctx context.Context, row, count int) error
}

// DataOptions is used to modify the default behavior.
type DataOptions struct {
	// ConcreteStruct can be set to any concrete struct (not a pointer).
	// When set, the mapstructure package is used to convert the returned
	// results automatically from a map to a struct. The `dbq` struct tag
	// can be used to map column names to the struct's fields.
	//
	// See: https://godoc.org/github.com/mitchellh/mapstructure
	ConcreteStruct interface{}

	// DecoderConfig is used to configure the decoder used by the mapstructure
	// package. If it's not supplied, a default StructorConfig is assumed. This means
	// WeaklyTypedInput is set to true and no DecodeHook is provided.
	//
	// See: https://godoc.org/github.com/mitchellh/mapstructure
	DecoderConfig *StructorConfig

	// ConcurrentPostUnmarshal can be set to true if PostUnmarshal must be called concurrently.
	ConcurrentPostUnmarshal bool
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

func structConvert(ctx context.Context, vals [][]string, header []string, o *DataOptions) ([]interface{}, error) {
	var (
		outStruct = []interface{}{}
	)

	for _, row := range vals {

		// map header to row value
		rowMap := map[string]interface{}{}

		if len(header) != len(row) {
			return nil, fmt.Errorf("length of each row must be same as length of header")
		}

		for idx, field := range row {
			headerI := strings.ToLower(header[idx])
			rowMap[headerI] = field
		}

		res := reflect.New(reflect.TypeOf(o.ConcreteStruct)).Interface()
		if o.DecoderConfig != nil {
			dc := &mapstructure.DecoderConfig{
				DecodeHook:       o.DecoderConfig.DecodeHook,
				ZeroFields:       true,
				TagName:          "las",
				WeaklyTypedInput: o.DecoderConfig.WeaklyTypedInput,
				Result:           res,
			}
			decoder, err := mapstructure.NewDecoder(dc)
			if err != nil {
				return nil, err
			}

			err = decoder.Decode(rowMap)
			if err != nil {
				return nil, err
			}

		} else {
			dc := &mapstructure.DecoderConfig{
				ZeroFields:       true,
				TagName:          "las",
				WeaklyTypedInput: true,
				Result:           res,
			}
			decoder, err := mapstructure.NewDecoder(dc)
			if err != nil {
				return nil, err
			}
			err = decoder.Decode(rowMap)
			if err != nil {
				return nil, err
			}
		}

		if len(outStruct) > 0 {
			csTyp := reflect.TypeOf(reflect.New(reflect.TypeOf(o.ConcreteStruct)).Interface())
			ics := reflect.TypeOf((*PostUnmarshaler)(nil)).Elem()

			if csTyp.Implements(ics) {
				rows := reflect.ValueOf(outStruct)
				count := rows.Len()

				if o.ConcurrentPostUnmarshal && runtime.GOMAXPROCS(0) > 1 {
					g, newCtx := errgroup.WithContext(ctx)

					for i := 0; i < count; i++ {
						i := i
						g.Go(func() error {
							if err := newCtx.Err(); err != nil {
								return err
							}

							row := reflect.ValueOf(rows.Index(i).Interface())
							retVals := row.MethodByName("PostUnmarshal").Call([]reflect.Value{reflect.ValueOf(newCtx), reflect.ValueOf(i), reflect.ValueOf(count)})
							err := retVals[0].Interface()
							if err != nil {
								return xerrors.Errorf("dbq.PostUnmarshal @ row %d: %w", i, err)
							}
							return nil
						})
					}

					if err := g.Wait(); err != nil {
						return nil, err
					}
				} else {
					for i := 0; i < count; i++ {
						if err := ctx.Err(); err != nil {
							return nil, err
						}
						row := reflect.ValueOf(rows.Index(i).Interface())
						retVals := row.MethodByName("lasData.PostUnmarshal").Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(i), reflect.ValueOf(count)})
						err := retVals[0].Interface()
						if err != nil {
							return nil, xerrors.Errorf("lasData.PostUnmarshal @ row %d: %w", i, err)
						}
					}
				}
			}
		}

		outStruct = append(outStruct, res)
	}

	return outStruct, nil
}
