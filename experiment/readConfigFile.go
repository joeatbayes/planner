package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ConfigResNeed struct {
	Name              string
	Id                string
	ColForGroup       string
	ColForMaxAssigned string
}

type ProjectConfig struct {
	SheetName        string
	StartDataRow     int
	ColProjectNumber string
	ColPriority      string
	ColProjectName   string
	ResNeeds         []ConfigResNeed
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
	fmt.Println("getCellInt cellVal=", cellVal, "err=", err)
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

func main() {
	f, err := excelize.OpenFile("../data/examples/sample-01/config.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	prjCfg := new(ProjectConfig)
	// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.
	prjCfg.ResNeeds = make([]ConfigResNeed, 0, 1000)
	// Get value from cell by given worksheet name and axis.
	prjCfg.SheetName = GetCellStr(f, "projects", "C4", "")
	prjCfg.StartDataRow = GetCellInt(f, "projects", "C5", -1)
	prjCfg.ColProjectName = GetCellStr(f, "projects", "C9", "")
	prjCfg.ColProjectNumber = GetCellStr(f, "projects", "C9", "")
	prjCfg.ColPriority = GetCellStr(f, "projects", "C6", "")

	fmt.Println("Project Cfg=", prjCfg)

	slcB, _ := json.Marshal(prjCfg)
	fmt.Println(string(slcB))

	// Parse Resource Needs Array
	// from the Excel file
	startNeedsRow := 19
	currRow := startNeedsRow
	rows, err := f.GetRows("projects")
	rlen := len(rows)
	if err != nil {
		fmt.Println("Err feting rows for projects", err)
	} else {

		for {
			if currRow >= rlen {
				break
			}
			trow := rows[currRow]
			if len(trow) < 5 {
				break // Encountered Empty Row so stop loading resource needs
			}

			fmt.Println("trow=", trow)
			resId := strings.TrimSpace(trow[2])
			if resId < " " {
				break // Encountere a empty row so stop loading resource needs
			}
			resName := strings.TrimSpace(trow[1])
			colResGrp := strings.TrimSpace(trow[3])
			colResGrpMax := strings.TrimSpace(trow[4])
			agrp := ConfigResNeed{
				Name:              resName,
				Id:                resId,
				ColForGroup:       colResGrp,
				ColForMaxAssigned: colResGrpMax}
			prjCfg.ResNeeds = append(prjCfg.ResNeeds, agrp)
			fmt.Println("resName=", resName, "agrp=", agrp)
			currRow++
		}
	}
	slcB, _ = json.Marshal(prjCfg)
	fmt.Println("config after reading resource needs", string(slcB))

}
