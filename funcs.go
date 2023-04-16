package contour

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"time"
)

// TODO: Add These
// - https://github.com/moul/funcmap/blob/master/funcmap.go#L11
// - https://github.com/gpmd/gotemplate/blob/master/template.go
// - https://api.rubyonrails.org/classes/Array.html#method-i-to_sentence
// - https://pkg.go.dev/html/template#pkg-functions
// - From filtru: formatPrice (int + float/decimal ones) + relative date (maybe lets just timediff outside of this lib) + inArray + uniq

// AllFuncs returns a full list of all the functions available in this package
func AllFuncs() template.FuncMap {
	return template.FuncMap{
		"split":     strings.Split,
		"join":      strings.Join,
		"title":     strings.Title,
		"trimspace": strings.TrimSpace,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"plural":    Plural,
		"now":       Now,
		"in_array":  inArray,
	}
}

func inArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}

// Plural returns plural or singular form of a word depending on the count
func Plural(num int, sin string, plu string, zero string) string {
	word := sin

	if num == 0 && zero != "" {
		word = zero
	} else if num != 1 && plu != "" {
		word = plu
	}

	return fmt.Sprintf("%d %s", num, word)
}

// Now returns the current local time.
func Now() time.Time {
	return time.Now()
}
