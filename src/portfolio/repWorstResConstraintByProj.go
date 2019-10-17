package portfolio

import (
	//"encoding/json"
	//"fmt"

	"strconv"
	//"strings"

	//"github.com/360EntSecGroup-Skylar/excelize"
)

func (pln *Planner) RepWorstResContraintByProject(sheetName string) {
	pln.Proj.SortByOutputPriority()

	f := pln.files.GetWithLastSeg("constraint", false)

	TitleFormat, _ := f.NewStyle(`{"font":{"bold":true,"italic":true,"family":"Berlin Sans FB Demi","size":36,"color":"#777777","underline":"none"}}`)
	numFmtNorm, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(3) + `}`)
	numFmt2Dig, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(4) + `}`)
	colHeadFmt, _ := f.NewStyle(`{"alignment":{"horizontal":"right","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)
	colHeadFmtLeft, _ := f.NewStyle(`{"alignment":{"horizontal":"left","ident":0,"justify_last_line":true,"reading_order":0,"relative_indent":0,"shrink_to_fit":true,"text_rotation":0,"vertical":"bottom","wrap_text":true}}`)
	numFmtYellow, _ := f.NewStyle(`{"number_format": ` + strconv.Itoa(4) + `,"fill":{"type":"pattern","color":["#FFFF66"],"pattern":1}}`)

	dataStartRowNdx := 3
	sheetNdx := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "Resources Requiring Longest Time By Project")
	f.SetCellStyle(sheetName, "A1", "A1", TitleFormat)
	f.SetCellValue(sheetName, "A2", "Proj ID")
	f.SetCellValue(sheetName, "B2", "Proj Name")
	//f.SetCellValue(sheetName, "C1", divLabel)
	f.SetColWidth(sheetName, "B", "B", 50)
	f.SetColWidth(sheetName, "D", "D", 20)
	//f.SetColWidth(sheetName, "G", "G", 2)
	f.SetColWidth(sheetName, "M", "M", 2)
	f.SetCellValue(sheetName, "C2", "Res ID")
	f.SetCellValue(sheetName, "D2", "Block Res Name")
	f.SetCellValue(sheetName, "F2", "Calc Duration")
	f.SetCellValue(sheetName, "E2", "Specified min Duration")
	f.SetCellValue(sheetName, "G2", "Units Needed")
	f.SetCellValue(sheetName, "H2", "Tot Res Start Count")
	f.SetCellValue(sheetName, "I2", "Max Res Cnt per Proj Allowed")
	f.SetCellValue(sheetName, "J2", "Max Units per day Allowed")
	f.SetCellValue(sheetName, "K2", "Min Assigned Res cnt Allowed")
	f.SetCellValue(sheetName, "L2", "Min Units per day Allowed")
	f.SetCellValue(sheetName, "N2", "Precursor For")
	f.SetCellValue(sheetName, "O2", "Needed By")

	f.SetCellStyle(sheetName, "A2", "P2", colHeadFmt)
	f.SetCellStyle(sheetName, "E2", "P2", colHeadFmt)
	f.SetCellStyle(sheetName, "A2", "D2", colHeadFmtLeft)
	f.SetCellStyle(sheetName, "E3", "E500", numFmtNorm)
	f.SetCellStyle(sheetName, "N3", "P500", colHeadFmt)
	f.SetCellStyle(sheetName, "E3", "L500", numFmt2Dig)

	// Figure this out.  It almost works but the title area is not displayed correctly
	//f.SetPanes(sheetName, `{"freeze":false,"split":true,"x_split":5270,"y_split":1600,"top_left_cell":"D3","active_pane":"bottomRight","panes":[{"sqref":"F6","active_cell":"F6"},{"sqref":"O2","active_cell":"O2","pane":"topRight"},{"sqref":"J60","active_cell":"J60","pane":"bottomLeft"},{"sqref":"O60","active_cell":"O60","pane":"bottomRight"}]}`)
	//f.SetPanes(sheetName, `{"freeze":true,"split":false,"x_split":1,"y_split":1,"top_left_cell":"D3","active_pane":"bottomRight","panes":[{"sqref":"A2","active_cell":"A2","pane":"topLeft"}]}`)

	outRowNdx := dataStartRowNdx
	for _, proj := range pln.Proj.Items {
		//isProg := false

		if proj.TotDirectResourceUnit == 0 {
			continue
		}

		wRes := proj.MostDemandResource
		wResDur := proj.MinDurationWorstResource

		f.SetCellValue(sheetName, MakeAxis("A", outRowNdx), proj.Id)
		f.SetCellValue(sheetName, MakeAxis("B", outRowNdx), leftSlice(proj.Name, 50))

		if float32(proj.MinDuration) > wResDur-0.01 {
			// skip generating data for this resource because the user
			// specifed duration exceeds the computed worst duration
			wRes = nil
			axis := MakeAxis("D", outRowNdx)
			f.SetCellValue(sheetName, axis, "User specified duration is limiter")
			f.SetCellStyle(sheetName, axis, axis, numFmtYellow)
		}

		if wRes != nil {
			rgrp := pln.Res.ItemsById[wRes.Id]
			f.SetCellValue(sheetName, MakeAxis("C", outRowNdx), wRes.Id)
			f.SetCellValue(sheetName, MakeAxis("D", outRowNdx), rgrp.GroupName)
			f.SetCellValue(sheetName, MakeAxis("F", outRowNdx), wResDur)
			f.SetCellValue(sheetName, MakeAxis("G", outRowNdx), wRes.Units)
			f.SetCellValue(sheetName, MakeAxis("H", outRowNdx), rgrp.StartCount)
			f.SetCellValue(sheetName, MakeAxis("I", outRowNdx), wRes.MaxAssigned)
			f.SetCellValue(sheetName, MakeAxis("J", outRowNdx), wRes.MaxUnitsPerDay)
			f.SetCellValue(sheetName, MakeAxis("K", outRowNdx), wRes.MinAssigned)
			f.SetCellValue(sheetName, MakeAxis("L", outRowNdx), wRes.MinUnitsPerDay)
		}
		//f.SetCellValue(sheetName, MakeAxis("D", outRowNdx), wRes.Name)
		if proj.MinDuration > 0 {
			f.SetCellValue(sheetName, MakeAxis("E", outRowNdx), proj.MinDuration)
		} else {
			f.SetCellValue(sheetName, MakeAxis("E", outRowNdx), " ")
		}

		f.SetCellValue(sheetName, MakeAxis("N", outRowNdx), proj.PrecursForAsStr(" "))
		f.SetCellValue(sheetName, MakeAxis("O", outRowNdx), proj.NeedByAsStr(" "))

		outRowNdx++
	}
	f.SetActiveSheet(sheetNdx)

}
