package portfolio

import (
	//"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func EmptyStrArr() []string {
	return []string{}
}

func GetCellStr(f *excelize.File, sheet, axis string, defVal string) string {
	cellVal, err := f.GetCellValue(sheet, axis)
	if err != nil {
		return defVal
	} else {
		return strings.TrimSpace(cellVal)
	}
}

func GetCellInt(f *excelize.File, sheet, axis string, defVal int) int {
	cellVal, err := f.GetCellValue(sheet, axis)
	//fmt.Println("getCellInt cellVal=", cellVal, "err=", err)
	if err != nil {
		return defVal
	} else {
		tval, err2 := strconv.Atoi(cellVal)
		fmt.Println("getCellInt cellVal=", cellVal, "tval=", tval, "err2=", err2)
		if err2 == nil {
			return tval
		} else {
			return defVal
		}
	}
}

func GetCellFloat64(f *excelize.File, sheet, axis string, defVal float64) float64 {
	cellVal, err := f.GetCellValue(sheet, axis)
	if err != nil {
		return defVal
	} else {
		tval, err := strconv.ParseFloat(cellVal, 64)
		if err == nil {
			return tval
		} else {
			return defVal
		}
	}
}

func GetCellFloat(f *excelize.File, sheet, axis string, defVal float64) float64 {
	return GetCellFloat64(f, sheet, axis, defVal)
}

func GetCellFloat32(f *excelize.File, sheet, axis string, defVal float32) float32 {
	cellVal, err := f.GetCellValue(sheet, axis)
	if err != nil {
		return defVal
	} else {
		tval, err := strconv.ParseFloat(cellVal, 32)
		if err == nil {
			return float32(tval)
		} else {
			return defVal
		}
	}
}

func GetCellBool(f *excelize.File, sheet, axis string, defVal bool) bool {
	cellVal, err := f.GetCellValue(sheet, axis)
	if err != nil {
		return defVal
	} else {
		tval := strings.ToLower(strings.TrimSpace(cellVal))
		if tval <= "" {
			return defVal
		} else if tval == "y" || tval == "yes" || tval == "t" || tval == "true" || tval == "1" {
			return true
		} else {
			return false
		}
	}
}

// Return a numeric index for a named column that allows that
// column to be retrieved from a row by numeric index instead
// of axis.  For Example AA would return 26.  Note row indexing
// is zero based while excel indexing is 1 based so the index
// would be one lower.   Helpful when looping across columns
// in a row.
func GetColNdx(colName string, defVal int) int {
	if colName <= " " {
		return defVal
	}
	andx, err := excelize.ColumnNameToNumber(colName)
	if err != nil {
		fmt.Println("Error getting colNdx for ", colName)
		return defVal
	} else {
		return andx - 1 // adjust to zero based to allow indexing directly into row
	}
}

// Returns the the string value from the cell at ndx parsed as
// a 32 bit float.  Similar to GetCellFloat32 except it uses
// a row as a input and numeric indexing rather than the axis
func GetRowValFloat32(arow []string, ndx int, defVal float32) float32 {
	if ndx < 0 || arow == nil || ndx >= len(arow) {
		return defVal
	}
	sval := arow[ndx]
	tval, err := strconv.ParseFloat(sval, 32)
	if err == nil {
		return float32(tval)
	} else {
		return defVal
	}
}

// Returns the the string value from the cell at ndx parsed as
// a 64 bit float.  Similar to GetCellFloat64 except it uses
// a row as a input and numeric indexing rather than the axis
func GetRowValFloat64(arow []string, ndx int, defVal float64) float64 {
	if ndx < 0 || arow == nil || ndx >= len(arow) {
		return defVal
	}
	sval := arow[ndx]
	tval, err := strconv.ParseFloat(sval, 64)
	if err == nil {
		return tval
	} else {
		return defVal
	}
}

// Returns the the string value from the cell at ndx parsed as
// a int.  Similar to GetCellInt except it uses
// a row as a input and numeric indexing rather than the axis
func GetRowValInt(arow []string, ndx int, defVal int) int {
	if ndx < 0 || arow == nil || ndx >= len(arow) {
		return defVal
	}
	sval := arow[ndx]
	tval, err := strconv.Atoi(sval)
	if err == nil {
		return tval
	} else {
		return defVal
	}
}

// Returns the the string array
// from field in row
func GetRowValStrArr(arow []string, ndx int, defVal []string) []string {

	if ndx < 0 || arow == nil || ndx >= len(arow) {
		return defVal
	}

	sval := strings.Replace(strings.TrimSpace(arow[ndx]), ",", " ", -1)
	return strings.Fields(sval)
}

// Returns the the string from field in row
// lookup by column name rather than ndx.
// normally it is better to lookup the ndx for the colname
// and resuse especially in loops but this method provides
// the shorthand with autolookup.
func GetRowValStrCN(arow []string, colName string, defVal string) string {

	if colName <= " " {
		return defVal
	}
	ndx, err := excelize.ColumnNameToNumber(colName)
	if err != nil {
		return defVal
	}
	return GetRowValStr(arow, ndx-1, defVal)
}

// Returns the the string from field in row
// or default value if the ndx is out of range
func GetRowValStr(arow []string, ndx int, defVal string) string {

	if ndx < 0 || arow == nil || ndx >= len(arow) {
		return defVal
	}
	return strings.TrimSpace(arow[ndx])
}

func GetRowValBool(arow []string, ndx int, defVal bool) bool {
	if ndx == -1 || arow == nil || ndx >= len(arow) {
		return defVal
	}
	if ndx >= len(arow) {
		return defVal
	}
	tval := strings.ToLower(strings.TrimSpace(arow[ndx]))

	if tval <= "" {
		return defVal
	} else if tval == "y" || tval == "yes" || tval == "t" || tval == "true" || tval == "1" {
		return true
	} else {
		return false
	}
}

// Make axis for column number and row.
func MakeAxisCNR(colNum int, row int) string {
	colName, _ := excelize.ColumnNumberToName(colNum)
	return colName + strconv.Itoa(row)
}

// Make axis for column and row.
func MakeAxis(colName string, row int) string {
	return colName + strconv.Itoa(row)
}
