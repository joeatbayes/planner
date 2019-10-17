package portfolio

import (
	//"encoding/json"
	//"fmt"

	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

/*********
 ****  TODO MODIFY TO FIT STATED INTENT
 *********/

func (pln *Planner) RequestByProjByResReport(sheetName string) {
	pln.Proj.SortByOutputPriority()

	f := pln.files.GetWithLastSeg("request", false)
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
	chartOffset := 4
	dataStartRowNdx := 3
	numRow := len(pln.Proj.Items)

	numRes := len(pln.Cfg.Proj.ResNeeds)
	totOffset := chartOffset + (numRes * 3)
	numCol := totOffset + 2
	postTotOffset := totOffset + 3

	numFmt, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(3) + `}`)
	colHeadFmt, _ := f.NewStyle(`{"alignment":{"horizontal":"right","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)
	wrapFmt, _ := f.NewStyle(`{"alignment":{"horizontal":"right","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)
	sheetNdx := f.NewSheet(sheetName)
	lastNumCell := totOffset + 25
	lastNumCol, _ := excelize.ColumnNumberToName(lastNumCell)

	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Name")
	//f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "A", "A", 10)
	f.SetColWidth(sheetName, "B", "B", 50)
	f.SetColWidth(sheetName, "C", "C", 10)
	f.SetColWidth(sheetName, "D", lastNumCol, 8)
	//f.SetColWidth(sheetName, "V", "V", 8)
	//f.SetColWidth(sheetName, "C", "C", 6)
	f.SetCellValue(sheetName, "C2", "Priority")
	preTotColNarrow, _ := excelize.ColumnNumberToName(totOffset)
	f.SetColWidth(sheetName, preTotColNarrow, preTotColNarrow, 1.2)
	postTotColNarrow, _ := excelize.ColumnNumberToName(postTotOffset)
	f.SetColWidth(sheetName, postTotColNarrow, postTotColNarrow, 1.2)
	postTotStartCol, _ := excelize.ColumnNumberToName(postTotOffset + 1)
	postTotStopCol, _ := excelize.ColumnNumberToName(postTotOffset + 7)
	f.SetColWidth(sheetName, postTotStartCol, postTotStopCol, 12)
	postTotStopColx2, _ := excelize.ColumnNumberToName(postTotOffset + 14)
	f.SetColWidth(sheetName, postTotStopCol, postTotStopColx2, 11)

	f.SetCellValue(sheetName, MakeAxisCNR(totOffset+1, 2), "Proj Total")
	f.SetCellValue(sheetName, MakeAxisCNR(totOffset+2, 2), "Run Total")
	f.SetCellStyle(sheetName, "C2", MakeAxisCNR(lastNumCell, 2), colHeadFmt)
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+1, 2), "is program")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+2, 2), "Target Date")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+3, 2), "Min Duration")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+4, 2), "Start Work Day")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+5, 2), "End Work Day")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+6, 2), "Parent")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+7, 2), "Precurs")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+8, 2), "Precurs For")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+9, 2), "Need By")
	f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+10, 2), "Children")
	f.SetCellStyle(sheetName, "D3", MakeAxisCNR(numCol, numRow+22), numFmt)
	f.SetCellStyle(sheetName, MakeAxisCNR(postTotOffset+5, 3), MakeAxisCNR(postTotOffset+10, numRow+22), wrapFmt)

	//lastCol := numCell + chartOffset
	// Set Col Widths for enough cells to take the Weeks
	//lastColStr, cerr := excelize.ColumnNumberToName(lastCol)
	//if cerr != nil {
	//	f.SetColWidth(sheetName, "C", lastColStr, 3)
	//} else {
	//	f.SetColWidth(sheetName, "C", "ZZ", 3)
	//}

	// Add # Headers

	currOffset := chartOffset
	for _, cfgRes := range pln.Cfg.Proj.ResNeeds {
		axisRes := MakeAxisCNR(currOffset+1, 2)
		axisRunTot := MakeAxisCNR(currOffset+2, 2)
		f.SetCellValue(sheetName, axisRes, cfgRes.Id)
		f.SetCellValue(sheetName, axisRunTot, "Run Total "+cfgRes.Id)
		f.SetCellStyle(sheetName, axisRes, axisRunTot, colHeadFmt)
		currOffset += 3
	}

	//lastOffset := currOffset
	currOffset = chartOffset
	outRowNdx := dataStartRowNdx
	resTotArr := make([]float32, numRes)
	runTotal := float32(0)
	for _, proj := range pln.Proj.Items {
		isProg := "NO"
		if proj.TotDirectResourceUnit == 0 && len(proj.Children) > 0 {
			isProg = "YES"
		}
		totUnits := float32(0)
		f.SetCellValue(sheetName, MakeAxis("A", outRowNdx), proj.Id)
		f.SetCellValue(sheetName, MakeAxis("B", outRowNdx), leftSlice(proj.Name, 150))
		f.SetCellValue(sheetName, MakeAxis("C", outRowNdx), proj.Priority)
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+1, outRowNdx), isProg)
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+2, outRowNdx), proj.TargetDate)

		if proj.MinDuration != -1 {
			f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+3, outRowNdx), proj.MinDuration)
		}

		if proj.StartWorkDay != -1 {
			f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+4, outRowNdx), proj.StartWorkDay)
		}

		if proj.EndWorkDay != -1 {
			f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+5, outRowNdx), proj.EndWorkDay)
		}

		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+6, outRowNdx), proj.Parent)
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+7, outRowNdx), strings.Join(proj.Precurs, "\n"))
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+8, outRowNdx), strings.Join(proj.PrecursFor, "\n"))
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+9, outRowNdx), strings.Join(proj.NeedBy, "\n"))
		f.SetCellValue(sheetName, MakeAxisCNR(postTotOffset+10, outRowNdx), strings.Join(proj.Children, "\n"))

		for resNdx, resNeed := range proj.Resources {
			if resNeed.Units <= 0.01 {
				continue
			}
			colNdx := currOffset + (resNdx * 3)
			f.SetCellValue(sheetName, MakeAxisCNR(colNdx+1, outRowNdx), resNeed.Units)
			resTotArr[resNdx] += resNeed.Units
			totUnits += resNeed.Units
			//resGroup := pln.Res.ItemsById[resNeed.Id]
			f.SetCellValue(sheetName, MakeAxisCNR(colNdx+2, outRowNdx), resTotArr[resNdx])
			if resNeed == proj.MostDemandResource {
				//f.SetCellValue(sheetName, MakeAxis("G", outRowNdx), "BLOCKER")
				// Change Color for this resource
			}
			narrowColNdx := colNdx
			narrowColName, _ := excelize.ColumnNumberToName(narrowColNdx)
			f.SetColWidth(sheetName, narrowColName, narrowColName, 1.2)

			//				workRemain := resNeed.Units - workDone

			/*				wrkStyle := numFmt
							if requested == awarded && awarded > 0.001 {
								wrkStyle = numFmtGreen
							} else if gap > 0 && awarded < 0.001 {
								wrkStyle = numFmtGrey
							} else if gap > 0 {
								wrkStyle = numFmtYellow
							}
							f.SetCellStyle(sheetName, MakeAxisCNR(cellOffset, outRowNdx), MakeAxisCNR(cellOffset, outRowNdx+7), wrkStyle)
			*/

		}
		runTotal += totUnits
		f.SetCellValue(sheetName, MakeAxisCNR(totOffset+1, outRowNdx), totUnits)
		f.SetCellValue(sheetName, MakeAxisCNR(totOffset+2, outRowNdx), runTotal)
		outRowNdx += 1
	}

	f.SetActiveSheet(sheetNdx)
}
