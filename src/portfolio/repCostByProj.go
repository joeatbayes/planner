package portfolio

import (
	//"encoding/json"
	"fmt"

	"strconv"
	//"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// produce a array of projects with costs based on what was consumed
// of each resource types by time unit.  The second dimension is a time divisor
// such as 5 days which would represent a work week.  This is intended
// to simplify production of reports by day.
func (pln *Planner) MakeProjCostByDay(timeDivisor int) [][]float32 {
	numDays := pln.Proj.LastDayWorked + 1
	numCells := (numDays / timeDivisor)
	//numRes := pln.Cfg.Proj.NumResourceNeeds
	allResItems := pln.Res.Items
	projItems := pln.Proj.Items
	numProj := len(projItems)
	allResAvail := pln.ResAvail.AllRes
	// Build a summary array of costs by period.
	tot_arr := make([][]float32, numProj)
	for projNdx := range tot_arr {
		d2 := make([]float32, numCells+1) // Add one to contain the total
		tot_arr[projNdx] = d2
	}

	for projNdx, proj := range projItems {
		projId := proj.Id
		for resGrpNdx, resGrp := range allResItems {
			resAvail := allResAvail[resGrpNdx]
			for dayNdx := 0; dayNdx < numDays; dayNdx++ {
				cellNdx := dayNdx / timeDivisor
				resByDay := resAvail.ByDay[dayNdx]
				resUsed, found := resByDay.UsedBy[projId]
				if found == true {
					// Lookup the cost per hour for this resource
					// and add it to the total for this project on this day
					tot_arr[projNdx][cellNdx] += resUsed.Units * resGrp.AvgCostPerUnit
					tot_arr[projNdx][numCells] += resUsed.Units * resGrp.AvgCostPerUnit
				}
			}
		}
	}
	return tot_arr
}

func (pln *Planner) CostsByProjByPeriodReport(sheetName string, timeDivisor int, timeDivisorYear int, divLabel string) {
	pln.Proj.SortByOutputPriority()

	f := pln.files.GetWithLastSeg("cost", false)
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

	numFmt, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(3) + `}`)

	colHeadFmt, _ := f.NewStyle(`{"alignment":{"horizontal":"right","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)

	totalsArr := pln.MakeProjCostByDay(timeDivisor)

	dataStartRowNdx := 3
	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	//f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "B", "B", 50)

	chartOffset := 3
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
	for cellNdx := 0; cellNdx < numCell; cellNdx++ {
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
	colName, _ := excelize.ColumnNumberToName(numCell + chartOffset)
	axis := colName + "2"
	f.SetCellValue(sheetName, axis, "Total")

	outRowNdx := dataStartRowNdx
	for projNdx, proj := range pln.Proj.Items {
		isProg := false

		if proj.TotDirectUnitsDelivered == 0 && len(proj.Children) > 0 {
			isProg = true
		}

		if isProg && outRowNdx != 2 {
			outRowNdx += 1
		}

		axRow := strconv.Itoa(outRowNdx)
		f.SetCellValue(sheetName, "A"+axRow, proj.Id)
		f.SetCellValue(sheetName, "B"+axRow, leftSlice(proj.Name, 50))

		for cellNdx := 0; cellNdx <= numCell; cellNdx++ {
			colName, cerr := excelize.ColumnNumberToName(cellNdx + chartOffset)
			if cerr != nil {
				fmt.Println("L583: Error convColNumToName col=", cellNdx, " err=", cerr)
			} else {
				axis := colName + axRow
				// TODO: Format with numeric ,
				if totalsArr[projNdx][cellNdx] != 0 {
					f.SetCellValue(sheetName, axis, totalsArr[projNdx][cellNdx])
				}
				f.SetCellStyle(sheetName, axis, axis, numFmt)
				//if isProg {
				//	f.SetCellStyle(sheetName, axis, axis, progStyle)
				//} else {
				//		f.SetCellStyle(sheetName, axis, axis, styleID)
				//}
				//f.SetCellValue(sheetName, axis, "X")
			}
		}
		outRowNdx++
	}
	f.SetActiveSheet(sheetNdx)
}
