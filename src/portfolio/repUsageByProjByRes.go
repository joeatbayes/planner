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

// produce a array of projects Showing by project by resource by time unit
//   0=to units for resource for day, 1=Units remaining available for resources for day,
//   2=units Requested by project for day, 3=Units Awarded to projec for day,
//   4=Gap between requested used.
func (pln *Planner) MakeUsageByProjByResourceByTimeUnit(timeDivisor int) [][][][]float32 {
	numDays := pln.Proj.LastDayWorked + 1
	numCells := (numDays / timeDivisor)
	numMeasures := 6
	numRes := pln.Cfg.Proj.NumResourceNeeds
	//allResItems := pln.Res.Items
	projItems := pln.Proj.Items
	numProj := len(projItems)
	allResAvail := pln.ResAvail.AllRes
	//numArrEle := numProj * NumDay * NumCells

	// Build a summary array of costs by period.
	// Go does not directly support nested arrays dynamically
	// allocated so we do it by nesting the initializers
	tot_arr := make([][][][]float32, numProj+1)
	for projNdx := range tot_arr {
		d2 := make([][][]float32, numRes+1)
		tot_arr[projNdx] = d2
		for resNdx := range d2 {
			d3 := make([][]float32, numCells+1)
			d2[resNdx] = d3
			for cellNdx := range d3 {
				d4 := make([]float32, numMeasures+1)
				d3[cellNdx] = d4
			}
		}
	}

	// Update with totals
	for projNdx, proj := range projItems {
		projId := proj.Id
		for resGrpNdx, _ := range proj.Resources {
			resAvail := allResAvail[resGrpNdx]
			for dayNdx := 0; dayNdx < numDays; dayNdx++ {
				cellNdx := dayNdx / timeDivisor
				resByDay := resAvail.ByDay[dayNdx]
				resUsed, found := resByDay.UsedBy[projId]
				tot_arr[projNdx][resGrpNdx][cellNdx][0] += resByDay.Total
				tot_arr[projNdx][resGrpNdx][cellNdx][1] += resByDay.Avail
				if found == true {
					tot_arr[projNdx][resGrpNdx][cellNdx][2] += resUsed.DesiredUnits
					tot_arr[projNdx][resGrpNdx][cellNdx][3] += resUsed.Units
				}

				gap := resUsed.DesiredUnits - resUsed.Units
				if gap < 0 {
					gap = 0
				}
				tot_arr[projNdx][resGrpNdx][cellNdx][4] += gap

			}
		}
	}
	return tot_arr
}

func (pln *Planner) UsageByProjByResByTimeUnitRep(sheetName string, timeDivisor int, timeDivisor2 int, divLabel string, label2 string) {
	pln.Proj.SortByOutputPriority()

	f := pln.files.GetWithLastSeg("usage", false)
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
	numFmtYellowStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#FFFF66"],"pattern":1}}`
	numFmtYellow, _ := f.NewStyle(numFmtYellowStr)
	numFmtGreyStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#D3D3D3"],"pattern":1}}`
	numFmtGrey, _ := f.NewStyle(numFmtGreyStr)
	numFmtGreenStr := `{"number_format": ` + strconv.Itoa(3) + `,"fill":{"type":"pattern","color":["#00FF00"],"pattern":1}}`
	numFmtGreen, _ := f.NewStyle(numFmtGreenStr)

	totalsArr := pln.MakeUsageByProjByResourceByTimeUnit(timeDivisor)
	dataStartRowNdx := 3
	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	//f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "B", "B", 8)
	f.SetColWidth(sheetName, "C", "C", 8)
	f.SetColWidth(sheetName, "D", "E", 6)
	f.SetColWidth(sheetName, "F", "F", 22)

	chartOffset := 8
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

		if proj.TotDirectResourceUnit == 0 && len(proj.Children) > 0 {
			isProg = true
		}

		if isProg && outRowNdx != 2 {
			outRowNdx += 1
		}

		f.SetCellValue(sheetName, MakeAxis("A", outRowNdx), proj.Id)
		f.SetCellValue(sheetName, MakeAxis("B", outRowNdx), proj.Priority)
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx), leftSlice(proj.Name, 50))

		outRowNdx += 1
		for resNdx, resNeed := range proj.Resources {
			if resNeed.Units <= 0.01 {
				continue
			}
			f.SetCellValue(sheetName, MakeAxis("D", outRowNdx), resNeed.Id)
			//f.SetCellValue(sheetName, MakeAxis("D", outRowNdx), resNeed.Name)
			resGroup := pln.Res.ItemsById[resNeed.Id]
			f.SetCellValue(sheetName, MakeAxis("E", outRowNdx), resGroup.GroupName)
			if resNeed == proj.MostDemandResource {
				f.SetCellValue(sheetName, MakeAxis("G", outRowNdx), "BLOCKER")
			}
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+1), "tot work")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+2), "tot work delivered")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+3), "Tot Res ")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+4), "Res still avail")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+5), "proj request")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+6), "proj awarded")
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx+7), "proj gap")
			workDone := float32(0.0)
			exitCnt := 0
			for cellNdx := 0; cellNdx <= numCell; cellNdx++ {
				cellOffset := cellNdx + chartOffset
				total := totalsArr[projNdx][resNdx][cellNdx][0]
				avail := totalsArr[projNdx][resNdx][cellNdx][1]
				requested := totalsArr[projNdx][resNdx][cellNdx][2]
				awarded := totalsArr[projNdx][resNdx][cellNdx][3]
				gap := totalsArr[projNdx][resNdx][cellNdx][4]
				workDone += awarded
				workRemain := resNeed.Units - workDone
				if workRemain <= 0.01 {
					// once a project is done then no need to keep output
					exitCnt += 1
					if exitCnt > 3 {
						break
					}
				}
				f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+1), resNeed.Units)
				f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+2), workDone)
				f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+3), total)
				f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+4), avail)
				if requested > 0 {
					f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+5), requested)
					f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+6), awarded)
					f.SetCellValue(sheetName, MakeAxisCNR(cellOffset, outRowNdx+7), gap)
				}

				wrkStyle := numFmt
				if requested == awarded && awarded > 0.001 {
					wrkStyle = numFmtGreen
				} else if gap > 0 && awarded < 0.001 {
					wrkStyle = numFmtGrey
				} else if gap > 0 {
					wrkStyle = numFmtYellow
				}
				f.SetCellStyle(sheetName, MakeAxisCNR(cellOffset, outRowNdx), MakeAxisCNR(cellOffset, outRowNdx+7), wrkStyle)

			}
			outRowNdx += 9
		}
	}
	f.SetActiveSheet(sheetNdx)
}
