package portfolio

import (
	//"encoding/json"
	"fmt"

	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func LoadResourceGroups(fiNames []string, config *Config) *ResourceGrps {
	rcfg := config.Res
	ndxMaxPerProj := GetColNdx(rcfg.ColMaxPerProj, -1)
	ndxMinPerProj := GetColNdx(rcfg.ColMinPerProj, -1)
	ndxGrpId := GetColNdx(rcfg.ColMatchName, -1)
	ndxStartByMonth := GetColNdx(rcfg.ColStartByMonth, -1)
	ndxEndByMonth := GetColNdx(rcfg.ColEndByMonth, -1)
	numMon := ndxEndByMonth - ndxStartByMonth
	ndxGroupName := GetColNdx(rcfg.ColGroupName, -1)
	ndxUnitsPerDay := GetColNdx(rcfg.ColNetUnitsPerResPerDay, -1)
	ndxAvgCostPerUnit := GetColNdx(rcfg.ColAvgCostPerUnit, -1)
	ndxUsageModel := GetColNdx(rcfg.ColUsageModel, -1)
	ndxResType := GetColNdx(rcfg.ColResourceType, -1)

	tres := new(ResourceGrps)
	tres.FiNames = make([]string, 0, 5)
	// make array of empty resource groups
	// Make and initialize my Resource Groups
	// array.   Also create and do basic initialization
	// for each group
	tres.Items = make([]*ResourceGrp, config.Proj.NumResourceNeeds, config.Proj.NumResourceNeeds)
	tres.ItemsById = make(map[string]*ResourceGrp)
	tres.ByName = make(map[string]*ResourceGrp)

	for _, fiName := range fiNames {
		f, err := excelize.OpenFile(fiName)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		tres.FiNames = append(tres.FiNames, fiName)

		for ndx, rneed := range config.Proj.ResNeeds {
			agrp := new(ResourceGrp)
			agrp.Id = rneed.Id
			agrp.GroupName = rneed.Name
			agrp.Position = ndx
			tres.Items[ndx] = agrp
			agrp.CntByMonth = make([]float32, numMon)
			agrp.NumMonth = numMon
		}

		// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.

		// Get value from cell by given worksheet name and axis.
		tabNames := rcfg.SheetNames
		for _, tabName := range tabNames {

			// Parse Resource Needs Array
			// from the Excel file
			startNeedsRow := rcfg.StartDataRow - 1
			rows, err := f.GetRows(tabName)
			rlen := len(rows)
			//fmt.Println("# resource rows in tab=", rlen)
			if err != nil {
				fmt.Println("Err feting rows for projects", err)
			} else {

				for currRowNdx := startNeedsRow; currRowNdx < rlen; currRowNdx++ {

					//rowndxstr := strconv.Itoa(currRowNdx)
					aRow := rows[currRowNdx-1]
					//axis := rcfg.ColMatchName + rowndxstr
					groupId := GetRowValStr(aRow, ndxGrpId, "")
					groupName := GetRowValStr(aRow, ndxGroupName, "")
					if groupId <= " " && groupName <= " " {
						// first group without a ID and name so indicates a blank row so stop processing
						// them.
						break
					} else if groupId <= "" || groupId[0] == '#' {
						// group id is blank or this line is a comment.
						continue
					}

					projConfigResNeed, found := config.Proj.ResNeedsById[groupId]
					if found == false {
						fmt.Println("L170: ERROR NON FATAL: ", "Resourced ID ", groupId, " in Resources is not found in config and will not be planned.")
						continue
					}

					aGrp := tres.Items[projConfigResNeed.Position]

					// Make our array order for resource groups loaded exactly the same
					// as when resources were defined in configuration file.
					// So we can do fast numeric lookup.

					aGrp.Id = groupId
					aGrp.GroupName = GetRowValStr(aRow, ndxGroupName, "ERROR")
					aGrp.GroupName = aGrp.GroupName
					// Index by name to support project types that list team name on single row rather than all on one row.
					aGrp.StartCount = GetRowValFloat32(aRow, ndxStartByMonth, -1)
					aGrp.UnitsPerDay = GetRowValFloat32(aRow, ndxUnitsPerDay, -1)
					aGrp.AvgCostPerUnit = GetRowValFloat32(aRow, ndxAvgCostPerUnit, -1)

					aGrp.MaxPerProj = GetRowValFloat32(aRow, ndxMaxPerProj, -1)
					aGrp.MinPerProj = GetRowValFloat32(aRow, ndxMinPerProj, -1)

					useModelStr := strings.ToLower(GetRowValStr(aRow, ndxUsageModel, ""))
					useModel, ufound := UseModelLook[useModelStr]
					if ufound {
						aGrp.UsageModel = useModel
					}

					resTypeStr := strings.ToLower(GetRowValStr(aRow, ndxResType, ""))
					resType, rfound := ResTypeLook[resTypeStr]
					if rfound {
						aGrp.ResourceType = resType
					}

					// Save our constructed group fo latter use.
					tres.ItemsById[aGrp.Id] = aGrp
					tres.ByName[aGrp.GroupName] = aGrp

					// Copy The variable length staff count by month
					// into the resource record.
					moutndx := 0
					for mndx, sval := range aRow[ndxStartByMonth:ndxEndByMonth] {
						fval, err := strconv.ParseFloat(sval, 32)
						if err == nil {
							tcnt := float32(fval)
							aGrp.CntByMonth[moutndx] = tcnt
							aGrp.TotCnt += tcnt
							aGrp.LastCnt = tcnt
						} else {
							fmt.Println("Error converting ", sval, " to float", " mndx=", mndx)
							aGrp.CntByMonth[moutndx] = -1
						}
						moutndx++
					}

					//fmt.Println("tabName=", tabName, " currRowNdx=", currRowNdx, "aGrp=", aGrp, "tndx=", tndx, "aRow=", aRow)
					fmt.Println("L125: new group=", aGrp.ToJSON(false))

				}
			}
		}

	}

	return tres
}
