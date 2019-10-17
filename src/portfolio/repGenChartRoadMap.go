package portfolio

import (
	//"encoding/json"
	"fmt"

	"strconv"
	//"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// TODO: Generalize This so we pass in the Divisor and Divisor Label so we can
//  produce easily for week,  month, quarter, year
func (pln *Planner) GenChartRoadMapFlexUnit(sheetName string, divisor int, divLabel string) {
	pln.Proj.SortByOutputPriority()

	f := pln.files.GetWithLastSeg("roadmap", false)
	format := `{"fill":{"type":"pattern","color":["#94d3a2"],"pattern":1}}`
	styleID, styleErr := f.NewStyle(format)
	if styleErr != nil {
		fmt.Println("L223: Error creating style for format ", format)
	}

	progFormat := `{"fill":{"type":"pattern","color":["#0403a2"],"pattern":1}}`
	LateFormat := `{"fill":{"type":"pattern","color":["#ff6961"],"pattern":1}}`
	progStyle, progStyleErr := f.NewStyle(progFormat)
	lateStyle, _ := f.NewStyle(LateFormat)
	if progStyleErr != nil {
		fmt.Println("L223: Error creating style for format ", progFormat)
	}

	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "Roadmap by "+divLabel)
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "B", "B", 50)
	chartOffset := 3
	numWeek := pln.Proj.LastDayWorked / divisor
	lastCol := numWeek + chartOffset

	// Set Col Widths for enough cells to take the Weeks
	lastColStr, cerr := excelize.ColumnNumberToName(lastCol)
	if cerr != nil {
		f.SetColWidth(sheetName, "C", lastColStr, 3)
	} else {
		f.SetColWidth(sheetName, "C", "ZZ", 3)
	}

	// Add Week # Headers
	for weekNdx := 0; weekNdx <= numWeek; weekNdx += 1 {
		colName, cerr := excelize.ColumnNumberToName(weekNdx + chartOffset)
		if cerr != nil {
			fmt.Println("L235 Error convColNumToName col=", weekNdx, " err=", cerr)
		} else {
			axis := colName + "2"
			if weekNdx < 11 || weekNdx%2 == 0 {
				f.SetCellValue(sheetName, axis, weekNdx+1)
			}
		}
	}

	outRowNdx := 3
	for _, proj := range pln.Proj.Items {
		isProg := false
		if proj.TotDirectResourceUnit == 0 && len(proj.Children) > 0 {
			isProg = true
		}

		if isProg && outRowNdx != 2 {
			outRowNdx += 1
		}
		axRow := strconv.Itoa(outRowNdx)
		f.SetCellValue(sheetName, "A"+axRow, proj.Id)
		f.SetCellValue(sheetName, "B"+axRow, leftSlice(proj.Name, 45))
		startWeek := proj.StartWorkDay / divisor
		stopWeek := proj.EndWorkDay / divisor
		fmt.Println("L77: proj.Id=", proj.Id, "startDay=", proj.StartWorkDay, "endDay=", proj.EndWorkDay, "startWeek=", startWeek, "stopWeek=", stopWeek)
		for weekNdx := startWeek; weekNdx <= stopWeek; weekNdx++ {
			firstDayNum := float32(weekNdx) * float32(divisor)
			lastDayNum := firstDayNum + float32(divisor) - 1
			colNum := weekNdx + chartOffset
			colName, cerr := excelize.ColumnNumberToName(colNum)
			fmt.Println("L82: proj.id=", proj.Id, "weekNdx=", weekNdx, "startWeek=", startWeek, "stopWeek=", stopWeek, "colNum=", colNum, " colName=", colName)
			if cerr != nil {
				fmt.Println("L235: Error convColNumToName col=", weekNdx, " err=", cerr)
			} else {
				axis := colName + axRow
				wrkStyle := styleID
				if proj.MaxCompletionDays > 0 && float32(lastDayNum) > proj.MaxCompletionDays {
					wrkStyle = lateStyle
				} else if isProg {
					wrkStyle = progStyle
				}
				fmt.Println("L94: proj.Id=", proj.Id, "weekNdx=", weekNdx, " colName=", colName, "colNum=", colNum, "axis=", axis)
				f.SetCellStyle(sheetName, axis, axis, wrkStyle)
				//f.SetCellValue(sheetName, axis, "X")
			}
		}
		outRowNdx++
	}
	f.SetActiveSheet(sheetNdx)

}

/*
func (pln *Planner) GenChartRoadMapWeek(sheetName string) {
	pln.Proj.SortByOutputPriority()
	f := pln.files.GetWithLastSeg("roadmap", false)
	format := `{"fill":{"type":"pattern","color":["#94d3a2"],"pattern":1}}`
	styleID, styleErr := f.NewStyle(format)
	if styleErr != nil {
		fmt.Println("L223: Error creating style for format ", format)
	}

	progFormat := `{"fill":{"type":"pattern","color":["#0403a2"],"pattern":1}}`
	progStyle, progStyleErr := f.NewStyle(progFormat)
	if progStyleErr != nil {
		fmt.Println("L223: Error creating style for format ", progFormat)
	}

	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", sheetName)
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	f.SetCellValue(sheetName, "C1", "Weeks")
	f.SetColWidth(sheetName, "B", "B", 50)

	chartOffset := 3
	numWeek := pln.Proj.LastDayWorked / 7
	lastCol := numWeek + chartOffset

	// Set Col Widths for enough cells to take the Weeks
	lastColStr, cerr := excelize.ColumnNumberToName(lastCol)
	if cerr != nil {
		f.SetColWidth(sheetName, "C", lastColStr, 3)
	} else {
		f.SetColWidth(sheetName, "C", "ZZ", 3)
	}

	// Add Week # Headers
	for weekNdx := 0; weekNdx <= numWeek; weekNdx += 2 {
		colName, cerr := excelize.ColumnNumberToName(weekNdx + chartOffset)
		if cerr != nil {
			fmt.Println("L235 Error convColNumToName col=", weekNdx, " err=", cerr)
		} else {
			axis := colName + "2"
			f.SetCellValue(sheetName, axis, weekNdx)
		}
	}

	outRowNdx := 3
	for _, proj := range pln.Proj.Items {
		isProg := false
		if proj.TotDirectResourceUnit == 0 && len(proj.Children) > 0 {
			isProg = true
		}

		if isProg && outRowNdx != 2 {
			outRowNdx += 1
		}
		axRow := strconv.Itoa(outRowNdx)
		f.SetCellValue(sheetName, "A"+axRow, proj.Id)
		f.SetCellValue(sheetName, "B"+axRow, leftSlice(proj.Name, 45))
		startWeek := proj.StartWorkDay / 7
		stopWeek := proj.EndWorkDay / 7
		fmt.Println("startDay=", proj.StartWorkDay, "endDay=", proj.EndWorkDay, "startWeek=", startWeek, "stopWeek=", stopWeek)
		for weekNdx := startWeek; weekNdx <= stopWeek; weekNdx++ {
			colName, cerr := excelize.ColumnNumberToName(weekNdx + chartOffset)
			if cerr != nil {
				fmt.Println("L235: Error convColNumToName col=", weekNdx, " err=", cerr)
			} else {
				axis := colName + axRow
				if isProg {
					f.SetCellStyle(sheetName, axis, axis, progStyle)
				} else {
					f.SetCellStyle(sheetName, axis, axis, styleID)
				}
				//f.SetCellValue(sheetName, axis, "X")
			}
		}
		outRowNdx++
	}
	f.SetActiveSheet(sheetNdx)

}
*/
