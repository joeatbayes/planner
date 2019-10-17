package portfolio

import (
	//"encoding/json"
	"fmt"

	//"path/filepath"
	"strconv"
	"strings"
	//"regexp"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// With Version1 The format assumptions are different that the column oriented.
// each line includes a Team column and we have to lookup the team name by string
// and match it to the resource column.  As a result any single project will only contain
// a single resource with all others being set to zero units.
//
// We must read the portfolio unit and create a program item as a project for
// each new unique item in that program. we have to detect the latest date for
// items under that portfolio item and set the portfolio
// due date to most distant date of any of its items.
//
// The Priority column in version1 is not properly exported so we will assign
// the priority based on the row Ndx being read.
//
// The Agile size specification in the configuration file
// must be used to convert Estimate value into Units.
func (oin *Projects) LoadProjectListVersion1(fiName string, plan *Planner) error {
	fmt.Println("L30: LoadProjectListVersion1 for file=", fiName)
	f, err := excelize.OpenFile(fiName)
	if err != nil {
		fmt.Println("L33: ERROR: FAIL: Loading Project List:", err)
		return err
	}
	//fmt.Println("L36: Cfg=", plan.Cfg.ToJSON(true))
	programs := make(map[string]*Project)

	oin.FiNames = append(oin.FiNames, fiName)
	cfg := plan.Cfg
	//rGroups := plan.Res
	pcfg := cfg.Proj
	pctx := plan.MakeProjParserContext()
	//rcfg := cfg.Res
	//fmt.Println("fiName=", fiName, "SheetCount=", f.SheetCount)

	// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.

	// TODO:
	// Make a Lookup Map from Team Names
	// to Team ID and lookup of from Team ID
	// to Index in project.resource needed.

	// Get value from cell by given worksheet name and axis.

	Sheets := f.GetSheetMap()
	for _, tabName := range Sheets {
		fmt.Println("L59: tabName=", tabName)
		// Parse Resource Needs Array
		// from the Excel file
		startDataRow := pcfg.StartDataRow
		//rows, err := f.GetRows(tabName)
		//if err != nil {
		//	fmt.Println("ERROR: FAIL: Could not find sheet (", tabName, ") in file ", fiName)
		//	return err
		//}

		//fmt.Println("tabName=", tabName, " NumRows=", len(rows))

		// Make Index Array to allow fast retrieval
		// of columns from the resource Configuration
		//fmt.Println("L49: startDataRow=", startDataRow, " ndxProjNum=", ndxProjNum, " numResources=", numResources, "ndxNeedBy=", ndxNeedBy)
		//rlen := len(rows)
		//fmt.Println("# resource rows in tab=", rlen)
		//if err != nil {
		//	fmt.Println("Err feting rows for projects", err)
		//} else {
		currRowNdx := 0
		rows, _ := f.Rows(tabName)
		for rows.Next() {
			currRowNdx += 1
			row, _ := rows.Columns()
			fmt.Println("L83: currRowNdx=", currRowNdx, " row=", row)
			if currRowNdx <= startDataRow {
				continue
			}

			//for currRowNdx := startDataRow; currRowNdx < rlen; currRowNdx++ {
			//rowndxstr := strconv.Itoa(currRowNdx)
			//row := rows[currRowNdx-1]
			//axis := rcfg.ColMatchName + rowndxstr
			proj := plan.ParseProjectHeaderSection(row, pctx)
			if proj == nil {
				fmt.Println("L87 empty project found so end run")
				continue // empty project ID and name so End stop processing
			}
			proj.Priority = float32(currRowNdx)
			proj.EffectivePriority = proj.Priority

			// Either update the existing program
			// or create a new program to act as a
			// parent for this item.
			progName := proj.Parent
			progProj, progFound := programs[progName]
			if progFound {
				// parse date item and update it in the program if it is after
				// the one currently in the program
				// TODO: ONLY SET DATE WHEN LATER THAN CURRENT CHILD DATE
				progProj.TargetDate = proj.TargetDate
				proj.Parent = progProj.Id

			} else {
				// Create the Program project the first time
				// we ecounter a new program name
				progProj = MakeEmptyFullyInitializedProject(pctx, plan)
				progProj.Id = "prog-" + strconv.Itoa(len(programs))
				progProj.Name = progName
				progProj.Priority = float32(currRowNdx)
				progProj.EffectivePriority = progProj.Priority
				fmt.Println("L109: New Program progProj=", progProj.ToJSON(false))
				programs[progName] = progProj
				oin.Items = append(oin.Items, progProj)
				oin.ItemsById[progProj.Id] = progProj
				proj.Parent = progProj.Id
			}

			// Lookup Team ID from Team Column in the
			// spreadsheet.  Then figure out what column
			// that resource would go into by doing a match against the
			// name column from the config file.
			teamName := GetRowValStr(row, pctx.NdxTeamName, "")
			teamName = StripNonWordChar(teamName)

			if teamName <= " " {
				msg := "L129: ERROR: team name can not be blank because it is used to map to resource type project=" + proj.ToJSON(true)
				proj.AddLog("LoadProjectListVersion1 L128:", msg)
			} else {
				configResGrp, configResGrpFound := cfg.Proj.ResNeedsByName[teamName]
				if configResGrpFound == false {
					msg := "L129: ERROR: team name can " + teamName + " not found and is required to map to resource type project=" + proj.ToJSON(true)
					proj.AddLog("LoadProjectListVersion1 L128:", msg)
				} else {
					resGroup := plan.Res.ItemsById[configResGrp.Id]
					rneed := proj.Resources[configResGrp.Position]
					rneed.InitFromResourceGroup(resGroup, configResGrp)
					// Now that we know the resource group ID we update Resource Group.
					// Match The Agile Size from the Configuration
					// file and lookup the size from the specified values.
					// for that resource type.
					estNdx := GetColNdx(pcfg.ColAgileSizeEst, -1)     // get column where agile size will be on  project row
					inValStr := GetRowValStr(row, estNdx, "")         // get the string representing agile size for project
					inValStr = strings.Split(inValStr, ".")[0]        // take only the int value we will be looking it up on key match
					units, found := configResGrp.AgileSizes[inValStr] // lookup the effective units for the specified agile size
					if found == false {
						proj.AddLog("LoadProjectListVersion1", "L146: Agile Time Unit="+inValStr+" could not be found in agile size to hours matrix, will be logged with zero units")
						units = 0
					}
					fmt.Println("L146: teamName=", teamName, " id=", "configResGrp.Id", " estNdx=", estNdx, " inValStr=", inValStr, " units=", units)
					rneed.Units = units
					proj.TotDirectResourceUnit += units
					rneed.MaxAssigned = proj.AdditFlds.GetFloat("maxassigned", rneed.MaxAssigned)
					rneed.MinAssigned = proj.AdditFlds.GetFloat("minassigned", rneed.MinAssigned)

					if rneed.MaxAssigned > -1 {
						rneed.MaxUnitsPerDay = rneed.MaxAssigned * resGroup.UnitsPerDay
					}
					if rneed.MinAssigned > -1 {
						rneed.MinUnitsPerDay = rneed.MinAssigned * resGroup.UnitsPerDay
					}

					// Check overrides in special fields.
					proj.MinDuration = proj.AdditFlds.GetFloat("minduration", proj.MinDuration)
					proj.PortionCompleteAtStart = proj.AdditFlds.GetFloat("portioncompleteatstart", proj.PortionCompleteAtStart)
					proj.MinCompletionDays = proj.AdditFlds.GetFloat("mincompletiondays", proj.MinCompletionDays)
					proj.MaxCompletionDays = proj.AdditFlds.GetFloat("maxcompletion", proj.MaxCompletionDays)
				}
			}
			fmt.Println("L147 new project read=", proj.ToJSON(false))
			oin.Items = append(oin.Items, proj)
			oin.ItemsById[proj.Id] = proj

			// Handle % Complete if specified
			if proj.PortionCompleteAtStart > 0 {
				proj.TotDirectUnitsDelivered = proj.TotDirectResourceUnit * proj.PortionCompleteAtStart
			}
			fmt.Println("L125: new Project=", proj.ToJSON(false))

			// Index project so they can be found by
			// alias.
			for _, alias := range proj.Alias {
				if alias != proj.Id {
					oin.ItemsByAlias[alias] = proj
				}
			}
			//fmt.Println("tres=", tres, "agrp=", agrp)
		}
	}
	//}
	return nil
}
