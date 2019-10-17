package portfolio

import (
	//"encoding/json"
	"fmt"

	"strconv"
	//"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

/*********
 ****  TODO MODIFY TO FIT STATED INTENT
 *********/

// Produce a multi dimensional array showing by resource
// how much of their capacity is available vs consumed
// produce a array of projects with costs based on what was consumed
// by Resource Add extra element at end of each dimension
// for totals.
// Output array [resourceNdx][DayNdx][0]float32 = available
// Output array [resourceNdx][DayNdx][1]float32 = used
func (pln *Planner) MakeResUsedFreeCapacyByTimeUnit(timeDivisor int) [][][]float32 {
	numDays := pln.Proj.LastDayWorked + 1
	numCells := (numDays / timeDivisor) + 1
	numRes := pln.Cfg.Proj.NumResourceNeeds + 1
	allResItems := pln.Res.Items
	allResAvail := pln.ResAvail.AllRes
	// At lowest level we have a two elements [0] = capacity,  [1] used

	// Build a summary array of costs by period.
	// Go does not directly support nested arrays dynamically
	// allocated so we do it by nesting the initializers
	tot_arr := make([][][]float32, numRes)
	for projNdx := range tot_arr {
		d2 := make([][]float32, numCells)
		tot_arr[projNdx] = d2
		for resNdx := range d2 {
			d3 := make([]float32, 2)
			d2[resNdx] = d3
		}
	}

	// Update with totals
	for resGrpNdx, _ := range allResItems {
		resAvail := allResAvail[resGrpNdx]
		for dayNdx := 0; dayNdx < numDays; dayNdx++ {
			cellNdx := dayNdx / timeDivisor
			resForDay := resAvail.ByDay[dayNdx]
			tot_arr[resGrpNdx][cellNdx][0] += resForDay.Avail
			tot_arr[resGrpNdx][cellNdx][1] += resForDay.Used
		}
	}

	return tot_arr
}

// Produce a report showing resource use versus capacity by time unit.
// Designed to illustrate hotspots and idle capacity.
// shows one main line for each resource with 3 rows capacity,
// consumption, Idle, Total,  Percentage Idle
// Colorize those grey where Idle capacity is under 1%
// Colorize those yellow where Idle capcity is greater than 15%
// Colorize those red where Idle capcity is greater than 30%
func (pln *Planner) RepResAvailVsCapacy(sheetName string, title string, timeDivisor int, timeDivisorYear int, divLabel string) {

	allResItems := pln.Res.Items

	f := pln.files.GetWithLastSeg("capacity", false)
	//format := `{"fill":{"type":"pattern","color":["#94d3a2"],"pattern":1}}`
	//styleID, styleErr := f.NewStyle(format)
	//if styleErr != nil {
	//	fmt.Println("L223: Error creating style for format ", format)
	//}

	//progFormat := `{"fill":{"type":"pattern","color":["#0403a2"],"pattern":1}}`
	//progStyle, progStyleErr := f.NewStyle(progFormat)
	//if progStyleErr != nil {
	//	fmt.Println("L223: Error creating style for format ", progFormat)
	//}

	// TODO: Something wrong the set cell style below is not working.
	numFmtNorm, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(3) + `}`)
	numFmtGreyStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#D3D3D3"],"pattern":1}}`
	numFmtGrey, _ := f.NewStyle(numFmtGreyStr)
	numFmtRedStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#B22222"],"pattern":1}}`
	numFmtRed, _ := f.NewStyle(numFmtRedStr)
	numFmtYellowStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#FFFF66"],"pattern":1}}`
	numFmtYellow, _ := f.NewStyle(numFmtYellowStr)
	numFmtGreenStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#32CD32"],"pattern":1}}`
	numFmtGreen, _ := f.NewStyle(numFmtGreenStr)
	fmt.Println("grey=", numFmtGreyStr)
	fmt.Println("yellow=", numFmtYellowStr)
	fmt.Println("red=", numFmtRedStr)

	colHeadFmt, _ := f.NewStyle(`{"alignment":{"horizontal":"right","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)

	totalsArr := pln.MakeResUsedFreeCapacyByTimeUnit(timeDivisor)

	dataStartRowNdx := 3
	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", title)
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	//f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "B", "B", 50)

	chartOffset := 4
	numCell := pln.Proj.LastDayWorked / timeDivisor

	//lastCol := numCell + chartOffset
	// Set Col Widths for enough cells to take the Weeks
	//lastColStr, cerr := excelize.ColumnNumberToName(lastCol)
	//if cerr != nil {
	//	f.SetColWidth(sheetName, "C", lastColStr, 3)
	//} else {
	//	f.SetColWidth(sheetName, "C", "ZZ", 3)
	//}

	// Add # Headers
	for cellNdx := 0; cellNdx <= numCell; cellNdx++ {
		colName, cerr := excelize.ColumnNumberToName(cellNdx + chartOffset)
		if cerr != nil {
			fmt.Println("L556: Error convColNumToName col=", cellNdx, " err=", cerr)
		} else {
			axis := colName + "2"
			label := divLabel + " " + strconv.Itoa(cellNdx+1)
			f.SetCellValue(sheetName, axis, label)
			f.SetCellStyle(sheetName, axis, axis, colHeadFmt)
		}
	}
	axis := MakeAxisCNR(chartOffset+numCell+1, 2)
	f.SetCellValue(sheetName, axis, "Total")

	outRowNdx := dataStartRowNdx
	for resGrpNdx, resGrp := range allResItems {
		totAvail := float32(0.0)
		totUsed := float32(0.0)
		f.SetCellValue(sheetName, MakeAxis("A", outRowNdx), resGrp.Id)
		f.SetCellValue(sheetName, MakeAxis("B", outRowNdx), leftSlice(resGrp.GroupName, 50))
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx+1), "avail")
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx+2), "used")
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx+3), "total")
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx+4), "% Idle")

		for cellNdx := 0; cellNdx <= numCell; cellNdx++ {
			cellCol := cellNdx + chartOffset
			avail := totalsArr[resGrpNdx][cellNdx][0]
			used := totalsArr[resGrpNdx][cellNdx][1]
			total := avail + used
			totAvail += avail
			totUsed += used
			percIdle := (avail / total) * 100.0
			availAxis := MakeAxisCNR(cellCol, outRowNdx+1)
			usedAxis := MakeAxisCNR(cellCol, outRowNdx+2)
			totAxis := MakeAxisCNR(cellCol, outRowNdx+3)
			percIdleAxis := MakeAxisCNR(cellCol, outRowNdx+4)

			f.SetCellValue(sheetName, availAxis, avail)
			f.SetCellValue(sheetName, usedAxis, used)
			f.SetCellValue(sheetName, totAxis, total)
			f.SetCellValue(sheetName, percIdleAxis, percIdle)

			f.SetCellStyle(sheetName, availAxis, availAxis, numFmtNorm)
			f.SetCellStyle(sheetName, usedAxis, usedAxis, numFmtNorm)
			f.SetCellStyle(sheetName, totAxis, totAxis, numFmtNorm)
			f.SetCellStyle(sheetName, percIdleAxis, percIdleAxis, numFmtNorm)

			format := numFmtGreen
			if percIdle > 30 {
				format = numFmtRed
			} else if percIdle > 20 {
				format = numFmtYellow
			} else if percIdle < 5 {
				format = numFmtGrey
			}
			f.SetCellStyle(sheetName, percIdleAxis, percIdleAxis, format)
			//if isProg {
			//	f.SetCellStyle(sheetName, axis, axis, progStyle)
			//} else {
			//		f.SetCellStyle(sheetName, axis, axis, styleID)
			//}
			//f.SetCellValue(sheetName, axis, "X")

		}
		totCol := chartOffset + numCell + 1
		totTot := totAvail + totUsed
		f.SetCellValue(sheetName, MakeAxisCNR(totCol, outRowNdx+1), totAvail)
		f.SetCellValue(sheetName, MakeAxisCNR(totCol, outRowNdx+2), totUsed)
		f.SetCellValue(sheetName, MakeAxisCNR(totCol, outRowNdx+3), totTot)
		if totTot > 0 {
			percIdle := (totAvail / totTot) * 100.0
			f.SetCellValue(sheetName, MakeAxisCNR(totCol, outRowNdx+4), percIdle)
		}

		outRowNdx += 5
	}
	f.SetActiveSheet(sheetNdx)
}
