package portfolio

import (
	"os"
	//"encoding/json"
	//"fmt"
	//"math"
	"strconv"
	"strings"
	//"github.com/360EntSecGroup-Skylar/excelize"
	"regexp"
)

func IsDir(fname string) bool {
	file, err := os.Open(fname)
	if err != nil {
		return false
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func IsFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Process string array and remove any duplicates
// return a new array.
func UniqueStrArr(pslice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range pslice {
		found, _ := keys[entry]
		if found == false {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Parse a string containing a key=val key=val key=val
// where the items are separated by space or comma
// as a map containing one item for each value in the
// string.  Keys can not contain whitespace or comma
func ParseStrFloatMap(strIn string) map[string]float32 {
	strIn = strings.Replace(strIn, ",", " ", -1)
	strIn = strings.Replace(strIn, "\n", " ", -1)
	tout := make(map[string]float32)
	items := strings.Split(strIn, " ")
	for _, sstr := range items {
		rowarr := strings.Split(sstr, "=")
		if len(rowarr) > 1 {
			key := strings.TrimSpace(rowarr[0])
			valstr := strings.TrimSpace(rowarr[1])
			val, err := strconv.ParseFloat(valstr, 32)
			if err == nil {
				tout[key] = float32(val)
			}
		}
	}
	return tout
}

type StrStrMap map[string]string

func (tin StrStrMap) GetFloat(name string, defVal float32) float32 {
	sval, found := tin[name]
	if found == false {
		return defVal
	}
	val, err := strconv.ParseFloat(sval, 32)
	if err != nil {
		return defVal
	}
	return float32(val)
}

var SpecialFieldsParmMatch, _ = regexp.Compile(`\.=\w+=[\.\w-]+`)

// parse a string for embedded fields following the
// form of .=varname=value  such as .=minDuration=10
// that may be embedded in description or notes fields
// in a tool like version1 or jira.
// known limit is that value can only be one word it will
// terminate at the next space.
func ParseSpecialFields(str string) map[string]string {
	ms := SpecialFieldsParmMatch.FindAllIndex([]byte(str), -1)
	tout := make(map[string]string)
	for _, m := range ms {
		subStr := str[m[0]+2 : m[1]]
		tarr := strings.SplitN(subStr, "=", 2)
		key := tarr[0]
		lckey := strings.TrimSpace(strings.ToLower(key))
		if len(tarr) > 1 {
			tout[key] = tarr[1]
			tout[lckey] = tarr[1]
		} else {
			tout[subStr] = ""
		}
	}
	return tout
}

// Return a string with all characters other than those from a..z, A..Z, 0..9
// converted to spaces.  Strip any leading or trailing spaces
var RegPatStripNonWordChar *regexp.Regexp

func StripNonWordChar(str string) string {
	if RegPatStripNonWordChar == nil {
		RegPatStripNonWordChar, _ = regexp.Compile("[^a-zA-Z0-9]+")
	}
	return strings.TrimSpace(RegPatStripNonWordChar.ReplaceAllString(str, " "))
}
