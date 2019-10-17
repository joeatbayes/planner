package portfolio

import (
	"encoding/json"
	"fmt"

	//"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ResNeedConfig struct {
	Name              string
	Id                string
	CostCenter        string
	ColsForGroup      []string
	ColForMaxAssigned string
	ColForMinAssigned string
	ColForUsageModel  string
	ColResType        string // labor, purchFirst, purchLast
	NdxsForGroup      []int
	NdxForMaxAssigned int
	NdxForMinAssigned int
	NdxForUsageModel  int
	NdxForResType     int
	Position          int // This items index in array.  Used to allow fast indexing after lookup by ID in other objects.

	// Stores the parsed out values for
	// agile sizing for the specified team.
	AgileSizes map[string]float32
}

type ProjectConfig struct {
	SheetNames   []string
	StartDataRow int

	// Used to allow invokation of different kinds of parsers.  Many files can be handled without our default
	// parser but with Jira and V1 exports they list one project per line with one estimate and do not break
	// estimates out by column.  We have to drive some of this with a slightly different parser
	FileType           string
	ColProjNum         string
	ColPriority        string
	ColProjName        string
	ColTargetDate      string
	ColParent          string
	ColPrecursors      string
	ColApproved        string
	ColNeedBy          string
	ColNPV             string
	ColPortionComplete string
	ColDateWorkStarted string
	ColDateWorkStoped  string
	ColMinDuration     string
	ColMaxDuration     string
	ColMustStartDate   string
	ColMustStopDate    string
	// Min Completion Days is similar to minimum duration
	// but rather than measuring days from the start of
	// project Min completion measures days from the
	// start of planning.  It includes any days of delay
	// before the project starts.  In is a good way to
	// state that maintenance activity must complete 264
	// working days from the start of planning or after
	// one full year.  Min completion can require that
	// a project start is delayed even if the project
	// could have been started sooner based on priorities
	// When minimum duration and minimum completion are both
	// set it allows precise control such as setting both
	// to 264 will cause the project to last exactly one work
	// year and will reduce consumption of resources to
	// fit.  When min completion is set without minimum
	// duration then the start date may be moved out
	// to make the project end on that date even if it
	// could have been completed sooner.
	// Min Completion is set automatically based on
	// starting date of planning when must Stop date
	// is set.
	//
	ColMinCompletionDays string
	ColMaxCompletionDays string
	ColTeamName          string
	NumResourceNeeds     int // short hand access to len of resource needs.
	ResNeedsById         map[string]*ResNeedConfig
	ResNeeds             []*ResNeedConfig
	ResNeedsByName       map[string]*ResNeedConfig
	ColAgileSizeEst      string
	ColAdditFields       string
}

func (cfg *ProjectConfig) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func (rgrp *ResNeedConfig) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(rgrp, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(rgrp)
		return string(slcB)

	}
}

func LoadProjectConfig(fiName string) *ProjectConfig {
	f, err := excelize.OpenFile(fiName)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	prjCfg := new(ProjectConfig)
	// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.
	prjCfg.ResNeeds = make([]*ResNeedConfig, 0, 1000)
	prjCfg.ResNeedsById = make(map[string]*ResNeedConfig)
	prjCfg.ResNeedsByName = make(map[string]*ResNeedConfig)

	// Get value from cell by given worksheet name and axis.
	namesStr := GetCellStr(f, "projects", "C4", "")
	namesStr = strings.Replace(namesStr, " ", ",", -1)
	namesStr = strings.Replace(namesStr, ";", ",", -1)
	namesStr = strings.Replace(namesStr, ":", ",", -1)
	namesArr := strings.Split(namesStr, ",")
	prjCfg.SheetNames = namesArr
	prjCfg.StartDataRow = GetCellInt(f, "projects", "C5", -1)
	prjCfg.ColProjName = GetCellStr(f, "projects", "C9", "")
	prjCfg.ColProjNum = GetCellStr(f, "projects", "C6", "")
	prjCfg.ColPriority = GetCellStr(f, "projects", "C7", "")
	prjCfg.ColTargetDate = GetCellStr(f, "projects", "C10", "")
	prjCfg.ColTargetDate = GetCellStr(f, "projects", "C10", "")
	prjCfg.ColParent = GetCellStr(f, "projects", "C11", "")
	prjCfg.ColPrecursors = GetCellStr(f, "projects", "C12", "")
	prjCfg.ColApproved = GetCellStr(f, "projects", "C13", "")
	prjCfg.ColNeedBy = GetCellStr(f, "projects", "C14", "")
	prjCfg.ColNPV = GetCellStr(f, "projects", "C15", "")
	prjCfg.ColPortionComplete = GetCellStr(f, "projects", "C16", "")
	prjCfg.ColApproved = GetCellStr(f, "projects", "C17", "")
	prjCfg.ColMustStartDate = GetCellStr(f, "projects", "C19", "")
	prjCfg.ColMustStopDate = GetCellStr(f, "projects", "C19", "")
	prjCfg.ColMinDuration = GetCellStr(f, "projects", "C21", "")
	prjCfg.ColMinCompletionDays = GetCellStr(f, "projects", "C22", "")
	prjCfg.ColMaxCompletionDays = GetCellStr(f, "projects", "C23", "")
	prjCfg.ColTeamName = GetCellStr(f, "projects", "G6", "")
	prjCfg.FileType = strings.ToLower(GetCellStr(f, "projects", "C24", ""))
	prjCfg.ColAdditFields = strings.ToLower(GetCellStr(f, "projects", "G8", ""))
	prjCfg.ColAgileSizeEst = GetCellStr(f, "projects", "G7", "")

	//fmt.Println("Project Cfg=", prjCfg.ToJSON())

	// Parse Resource Needs Array
	// from the Excel file
	startNeedsRow := 27
	rows, err := f.GetRows("projects")
	rlen := len(rows)
	if err != nil {
		fmt.Println("L162: FATAL ERROR: feting rows for projects: ", err)
	} else {

		for currRow := startNeedsRow; currRow < rlen; currRow++ {
			trow := rows[currRow]
			if len(trow) < 5 {
				break // Encountered Empty Row so stop loading resource needs
			}

			//fmt.Println("trow=", trow)
			resId := GetRowValStr(trow, 2, "")
			if resId < " " {
				break // Encountere a empty row so stop loading resource needs
			}
			if resId[0] == '#' {
				continue
			}
			agrp := new(ResNeedConfig)
			agrp.AgileSizes = make(map[string]float32)
			agrp.Id = resId
			agrp.Name = StripNonWordChar(GetRowValStr(trow, 1, ""))
			agrp.CostCenter = GetRowValStr(trow, 0, "")
			agrp.ColsForGroup = strings.Fields(strings.Replace(GetRowValStr(trow, 3, ""), ",", " ", -1))
			agrp.ColForMaxAssigned = GetRowValStr(trow, 4, "")
			agrp.ColForMinAssigned = GetRowValStr(trow, 5, "")
			agrp.ColForUsageModel = GetRowValStr(trow, 6, "")
			agrp.ColResType = GetRowValStr(trow, 7, "")
			agrp.NdxForUsageModel = GetColNdx(agrp.ColForUsageModel, -1)
			agrp.NdxForResType = GetColNdx(agrp.ColResType, -1)
			agrp.NdxForMaxAssigned = GetColNdx(agrp.ColForMaxAssigned, -1)
			agrp.NdxForMinAssigned = GetColNdx(agrp.ColForMinAssigned, -1)
			agrp.NdxsForGroup = make([]int, len(agrp.ColsForGroup))
			for ndx, colForResGrp := range agrp.ColsForGroup {
				agrp.NdxsForGroup[ndx] = GetColNdx(colForResGrp, -1)

			}
			prjCfg.ResNeeds = append(prjCfg.ResNeeds, agrp)
			prjCfg.ResNeedsById[agrp.Id] = agrp
			prjCfg.ResNeedsByName[agrp.Name] = agrp
			agrp.Position = prjCfg.NumResourceNeeds // position in array is before index due to zero based indexing

			// Parse the Agile Sizes so we can lookup
			// the size from the agile sizes exported from
			// tools like version1 and Jira
			agileSizeStr := GetRowValStrCN(trow, "J", "")
			agrp.AgileSizes = ParseStrFloatMap(agileSizeStr)
			fmt.Println("L205: agileSizeStr=", agileSizeStr, " sizes=", agrp.AgileSizes)

			prjCfg.NumResourceNeeds++

			//fmt.Println("resName=", resName, "agrp=", agrp)
			//fmt.Println("L149: config ResNeed = ", agrp.ToJSON(true))

		}
	}

	return prjCfg
}
