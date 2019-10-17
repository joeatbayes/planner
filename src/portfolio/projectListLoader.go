package portfolio

import (
	//"encoding/json"
	"fmt"

	"path/filepath"
	"strconv"
	"strings"

	"github.com/joeatbayes/goutil/jutil"

	"github.com/360EntSecGroup-Skylar/excelize"
	//"github.com/joeatbayes/goutil/jutil"
)

func (p *Projects) UpdateCrossRef() {
	return
}

// Created to make it easier to split the parser into smaller
// functions while making it easy to pass in the required context
// without rebuidling over again.
type ProjParserContext struct {
	NumResources         int
	NdxProjNum           int
	NdxPriority          int
	NdxProjName          int
	NdxTargetDate        int
	NdxParent            int
	NdxPrecursors        int
	NdxApproved          int
	NdxNeedBy            int
	NdxNPV               int
	NdxPercComplete      int
	NdxMinDuration       int
	NdxMinCompletionDays int
	NdxMaxCompletionDays int
	NdxTeamName          int
	NdxAdditFields       int
}

// Parse values from the project config and other parameters
// use-used by various parts of the parsing proces.  Needed this
// object to make it easier to break the parsing into discrete
// functions that needed some shared context.
func (plan *Planner) MakeProjParserContext() *ProjParserContext {
	pcfg := plan.Cfg.Proj
	nx := new(ProjParserContext)
	nx.NumResources = len(pcfg.ResNeeds)
	nx.NdxProjNum = GetColNdx(pcfg.ColProjNum, -1)
	nx.NdxPriority = GetColNdx(pcfg.ColPriority, -1)
	nx.NdxProjName = GetColNdx(pcfg.ColProjName, -1)
	nx.NdxTargetDate = GetColNdx(pcfg.ColTargetDate, -1)
	nx.NdxParent = GetColNdx(pcfg.ColParent, -1)
	nx.NdxPrecursors = GetColNdx(pcfg.ColPrecursors, -1)
	nx.NdxApproved = GetColNdx(pcfg.ColApproved, -1)
	nx.NdxNeedBy = GetColNdx(pcfg.ColNeedBy, -1)
	nx.NdxNPV = GetColNdx(pcfg.ColNPV, -1)
	nx.NdxPercComplete = GetColNdx(pcfg.ColPortionComplete, -1)
	nx.NdxMinDuration = GetColNdx(pcfg.ColMinDuration, -1)
	//nx.//ndxMaxDuration := GetColNdx(pcfg.ColMaxDuration, -1)
	nx.NdxMinCompletionDays = GetColNdx(pcfg.ColMinCompletionDays, -1)
	nx.NdxMaxCompletionDays = GetColNdx(pcfg.ColMaxCompletionDays, -1)
	nx.NdxAdditFields = GetColNdx(pcfg.ColAdditFields, -1)

	nx.NdxTeamName = GetColNdx(pcfg.ColTeamName, -1)
	return nx
}

// Scan project less based on predecessor and promote
// any lower priority projects that are considered
// predecessor by higher priority items.
func (p *Projects) UpdateEffectivePriority() {
	for _, proj := range p.Items {
		for _, prec := range proj.Precurs {
			preProj, found := p.ItemsById[prec]
			if found {
				if preProj.Priority > proj.Priority {
					// Change priority of prec Priority so it
					// will run before the current project.
					preProj.EffectivePriority = proj.Priority - 0.001
				}
				preProj.PrecursFor = append(preProj.PrecursFor, proj.Id)
				preProj.NeedBy = append(preProj.NeedBy, proj.Id)
			} else {
				// Item not found so need to update precursorOrphens
				p.PrecursOrphens = append(p.PrecursOrphens, prec)
				fmt.Println("WARN project ", proj.Id, " precursor ", prec, " can not be found")
			}
		}
	}
	p.PrecursOrphens = UniqueStrArr(p.PrecursOrphens)
}

func (p *Projects) UpdateOrphens() {
}

func (p *Projects) UpdateChildren() {
	for _, proj := range p.Items {
		if proj.Parent > " " {
			parProj, found := p.ItemsById[proj.Parent]
			if found {
				parProj.Children = append(parProj.Children, proj.Id)

			} else {
				// Item not found so need to update ParentOrphens
				p.ParentOrphens = append(p.ParentOrphens, proj.Parent)
				fmt.Println("WARN project ", proj.Id, " parent ", proj.Parent, " can not be found")
			}
		}
	}
	p.ParentOrphens = UniqueStrArr(p.ParentOrphens)
}

// Update needs array for all items named in any projects
// needBy column.  This allows us to determine a total set
// of needs.
func (p *Projects) UpdateNeeds() {
	for _, proj := range p.Items {
		for _, tneed := range proj.NeedBy {
			needProj, found := p.ItemsById[tneed]
			if found {
				needProj.Needs = append(needProj.Needs, proj.Id)
			} else {
				// Item not found so need to update Orphens
				p.NeedByOrphens = append(p.NeedByOrphens, tneed)
				fmt.Println("WARN project ", proj.Id, " needs ", tneed, " can not be found")
			}
		}
	}
	p.PrecursOrphens = UniqueStrArr(p.PrecursOrphens)

	// Remove any duplicates from the needs since it updated
	// from both the precursor and needBy
	for _, proj := range p.Items {
		proj.Needs = UniqueStrArr(proj.Needs)
	}
}

// Load  Set of ProjectsList files into a single unified set
// May be used many times to safely load more files but caller
// must call DoPostLoadAdjustments again after they are all loaded
func (oin *Projects) LoadProjectLists(fiNames []string, plan *Planner) {
	fmt.Println("L138: LoadProjectLists fiNames=", fiNames, " oin=", oin)
	for _, fiName := range fiNames {
		fmt.Println("L93: LoadProjectLists fiName=", fiName)
		if IsDir(fiName) == false {
			fmt.Println("L94: isDir=false")
			oin.LoadProjectListDispatch(fiName, plan)
		} else {
			gp := fiName + "/*.xlsx"
			dirFiles, err := filepath.Glob(gp)
			fmt.Println("L100: LoadProjectLists asDir name=", fiName, " glob=", gp, " dirFiles=", dirFiles, "err=", err)
			if err != nil {
				fmt.Println("L102: Error reading globpath ", gp, " err=", err)
			} else {
				for _, dfname := range dirFiles {
					fmt.Println("L105: LoadProjectLists loading files from dir dfname=", dfname)
					if strings.HasPrefix(dfname, "~") {
						fmt.Println("L161: skip ", dfname, " because starts with ~")
						continue
					}
					oin.LoadProjectListDispatch(dfname, plan)
				}
			}
		}
		fmt.Println("L106: fiName=", fiName, " numItems=", len(oin.Items))
	}
}

func (oin *Projects) LoadProjectListDispatch(fiName string, plan *Planner) {
	fmt.Println("L162: LoadProjectListDispatch fiName=", fiName, " fileType=", plan.Cfg.Proj.FileType)
	if plan.Cfg.Proj.FileType == "version1" {
		oin.LoadProjectListVersion1(fiName, plan)
	} else {
		oin.LoadProjectListColumns(fiName, plan)
	}
	jutil.Elap("L168: LoadProjectListDispatch done", plan.Start, jutil.Nowms())
}

func (oin *Projects) DoPostLoadAdjustments() {
	oin.UpdateEffectivePriority()
	oin.UpdateOrphens()
	oin.UpdateChildren()
	oin.UpdateNeeds()
	oin.SortByEffectivePriority()
}

/*
func (proj *Projects) ParseMainProjSection(row []string) {

}
*/

func MakeEmptyProject(pcx *ProjParserContext, plan *Planner) *Project {
	aProj := new(Project)
	aProj.Children = make([]string, 0, 4)
	aProj.Resources = make([]*ProjectResourceNeed, pcx.NumResources)
	aProj.Needs = make([]string, 0, 4)
	aProj.PrecursFor = make([]string, 0, 4)
	aProj.Log = make([]string, 0, 10)
	for mndx, configNeed := range plan.Cfg.Proj.ResNeeds {
		projNeed := new(ProjectResourceNeed)
		projNeed.Id = configNeed.Id
		projNeed.Position = mndx
		aProj.Resources[mndx] = projNeed
		// Note the resource need still needs to
		// have min, Max units and assigned in latter
		// parsing steps
	}
	return aProj

}

func MakeEmptyFullyInitializedProject(pctx *ProjParserContext, plan *Planner) *Project {
	proj := MakeEmptyProject(pctx, plan)
	proj.Id = ""
	proj.Priority = 999999
	proj.EffectivePriority = proj.Priority
	proj.Alias = EmptyStrArr()
	proj.Name = ""
	proj.Parent = ""
	proj.Precurs = EmptyStrArr()
	proj.Approved = true
	proj.NeedBy = EmptyStrArr()
	proj.NPV = -1
	proj.PortionCompleteAtStart = 0
	proj.MinDuration = -1
	proj.MaxCompletionDays = -1
	proj.MinCompletionDays = -1
	proj.StartWorkDay = -1
	proj.EndWorkDay = -1
	proj.IsComplete = false
	proj.NotWorkable = false
	return proj
}

// Parse the constant portion of each project.  Separated from
// rest of the project row parser to make it easier to re-use
// from other custom parsers.
func (plan *Planner) ParseProjectHeaderSection(aRow []string, pcx *ProjParserContext) *Project {
	projIds := GetRowValStr(aRow, pcx.NdxProjNum, "")
	projName := GetRowValStr(aRow, pcx.NdxProjName, "")
	if projIds <= " " || projName <= " " {
		fmt.Println("L246 both ProjIds and projName empty so stop arow=", aRow)
		return nil // empty project ID and name so End stop processing
	}

	aProj := MakeEmptyProject(pcx, plan)

	// Convert a project number array that may contain
	// more than one number into an array. The first
	// record we keep
	if projIds == "NOID" || projIds == "need" || projIds <= " " {
		plan.Proj.EmptyIdCnt += 1
		aProj.Id = "NEED" + strconv.Itoa(plan.Proj.EmptyIdCnt)
	} else {
		idarr := strings.Fields(projIds)
		if len(idarr) == 1 {
			aProj.Id = idarr[0]
		} else if len(idarr) > 1 {
			aProj.Id = idarr[0]
			aProj.Alias = idarr[1:]
		}
	}

	//fmt.Println("L163: proj row=", aRow)

	aProj.Priority = GetRowValFloat32(aRow, pcx.NdxPriority, 9999)
	aProj.EffectivePriority = aProj.Priority
	aProj.Name = aRow[pcx.NdxProjName]
	aProj.Name = strings.Split(aProj.Name, ":")[0]

	// TODO: Parse this as a date rather than string
	//  excelizer has date parsing built in.
	aProj.TargetDate = GetRowValStr(aRow, pcx.NdxTargetDate, "")
	aProj.Parent = GetRowValStr(aRow, pcx.NdxParent, "")
	aProj.Precurs = GetRowValStrArr(aRow, pcx.NdxPrecursors, EmptyStrArr())
	aProj.Approved = GetRowValBool(aRow, pcx.NdxApproved, true)
	aProj.NeedBy = GetRowValStrArr(aRow, pcx.NdxNeedBy, EmptyStrArr())
	aProj.NPV = GetRowValFloat32(aRow, pcx.NdxNPV, -1)
	aProj.PortionCompleteAtStart = GetRowValFloat32(aRow, pcx.NdxPercComplete, 0)
	aProj.MinDuration = GetRowValFloat32(aRow, pcx.NdxMinDuration, -1)
	aProj.MaxCompletionDays = GetRowValFloat32(aRow, pcx.NdxMaxCompletionDays, -1)
	aProj.MinCompletionDays = GetRowValFloat32(aRow, pcx.NdxMinCompletionDays, -1)
	additFldsStr := GetRowValStr(aRow, pcx.NdxAdditFields, "")
	aProj.AdditFldsStr = additFldsStr
	aProj.AdditFlds = ParseSpecialFields(additFldsStr)
	aProj.StartWorkDay = -1
	aProj.EndWorkDay = -1
	aProj.IsComplete = false
	aProj.NotWorkable = false
	return aProj
}

func (rneed *ProjectResourceNeed) InitFromResourceGroup(rgrp *ResourceGrp, configResNeed *ResNeedConfig) {
	rneed.Id = configResNeed.Id
	rneed.MaxAssigned = rgrp.MaxPerProj
	rneed.MinAssigned = rgrp.MinPerProj
	if rneed.MinAssigned == -1 {
		rneed.MinAssigned = 1
	}
	if rneed.MaxAssigned == -1 {
		rneed.MaxAssigned = 1
	}
	rneed.MaxUnitsPerDay = rneed.MaxAssigned * rgrp.UnitsPerDay
	rneed.MinUnitsPerDay = rneed.MinAssigned * rgrp.UnitsPerDay
}

// For each resource column defined in the configuration we will access the
// column specified in the configuration and parse that value to use as
// the units needed for that resource and use it to populate the resources
// array for that project.
func (proj *Project) ParseColumnBasedResourceNeeds(row []string, pctx *ProjParserContext, plan *Planner) {

	rGroups := plan.Res
	proj.TotDirectResourceUnit = 0.0
	for mndx, configResNeed := range plan.Cfg.Proj.ResNeeds {
		//fmt.Println("L115: mndx=", mndx, " configResNeed=", configResNeed)
		fmt.Println("L305: ParseColumnBasedResourceNeeds mndx=", mndx, "len(proj.Resources)=", len(proj.Resources))
		rneed := proj.Resources[mndx]
		rgrp, rfnd := rGroups.ItemsById[configResNeed.Id]
		rneed.InitFromResourceGroup(rgrp, configResNeed)
		rneed.Id = configResNeed.Id

		// Note: We must add the resource group even if we can not find a matching
		// one because we use simple indexing in latter loops and they can not
		// tolerate a nil record.  Add the record here to allow simple shortcut
		// continue on error conditions.
		if rfnd == false {
			fmt.Println("L188: WARN proj", proj.Id, " needs resourceGroup ", configResNeed.Id, " could not be found in resourceGroups")
			// We are depending on the default go initializers here to ensure that
			// a rneed values are zeroed.
			continue
		}
		//fmt.Println("L197: needId=", rneed.Id, " ndxsForGrp=", configResNeed.NdxsForGroup, "aRow=", aRow, "")

		// Iterate the list of columns and sum their value to get total units
		// for this resource type.
		for _, ndxForGrp := range configResNeed.NdxsForGroup {
			tval := GetRowValFloat32(row, ndxForGrp, 0)
			rneed.Units += tval
			//fmt.Println("L236: ndxForGrp=", ndxForGrp, " value=", tval, " rneed.Units=", rneed.Units)
		}
		// Handle % Complete if specified to reduce amount of
		// work required by marking that work as already delivered.
		// the system will only have to schedule the remaining work
		if proj.PortionCompleteAtStart > 0 {
			rneed.UnitsDelivered = rneed.Units * proj.PortionCompleteAtStart
		}

		// Read Max Assigned from spreadsheet. Fallback to what is
		// specified in resource gorup if not defined.   Fall back
		// to 1 if can find  nothing.
		rneed.MaxAssigned = GetRowValFloat32(row, configResNeed.NdxForMaxAssigned, rneed.MaxAssigned)

		// Read Min Assigned from spreadsheet. Fallback to what is
		// specified in resource gorup if not defined.   Fall back
		// to 1 if can find  nothing.
		rneed.MinAssigned = GetRowValFloat32(row, configResNeed.NdxForMinAssigned, rneed.MinAssigned)

		// overrite MaxUnitsPerDay from
		// using any values loaded from the max, min specified on project line
		// most projects spreadsheets do not define this so this ends up as a no-op
		// with same answer computed during init from the resGroup most of the time
		rneed.MaxUnitsPerDay = rneed.MaxAssigned * rgrp.UnitsPerDay
		rneed.MinUnitsPerDay = rneed.MinAssigned * rgrp.UnitsPerDay
		proj.TotDirectResourceUnit += rneed.Units

		// We do not set MostDemandResources here because we do it when working the project
		// after other resources may have reduced available hours and elongate fulfillment
		// time for some resources that may not initially seem the most demanding.
	}
}

/* Read a single File and add it's contents
to the existing project List */
func (oin *Projects) LoadProjectListColumns(fiName string, plan *Planner) error {
	dupProjIdCntr := 0
	f, err := excelize.OpenFile(fiName)
	if err != nil {
		fmt.Println("ERROR: FAIL: Loading Project List:", err)
		return err
	}
	oin.FiNames = append(oin.FiNames, fiName)
	cfg := plan.Cfg
	//rGroups := plan.Res
	pcfg := cfg.Proj
	pctx := plan.MakeProjParserContext()
	//rcfg := cfg.Res
	//fmt.Println("fiName=", fiName, "SheetCount=", f.SheetCount)

	// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.

	// Get value from cell by given worksheet name and axis.
	tabNames := pcfg.SheetNames
	for _, tabName := range tabNames {

		// Parse Resource Needs Array
		// from the Excel file
		startDataRow := pcfg.StartDataRow
		//rows, err := f.GetRows(tabName)
		//if err != nil {
		//	fmt.Println("ERROR: FAIL: Could not find sheet (", tabName, ") in file ", fiName)
		//	return err
		//}
		currNdx := 0
		rows, _ := f.Rows(tabName)
		//fmt.Println("tabName=", tabName, " NumRows=", len(rows))
		for rows.Next() {
			currNdx += 1
			if currNdx <= startDataRow {
				continue
			}
			row, _ := rows.Columns()

			// Make Index Array to allow fast retrieval
			// of columns from the resource Configuration
			//fmt.Println("L49: startDataRow=", startDataRow, " ndxProjNum=", ndxProjNum, " numResources=", numResources, "ndxNeedBy=", ndxNeedBy)
			//rlen := len(rows)
			//fmt.Println("# resource rows in tab=", rlen)
			//if err != nil {
			//	fmt.Println("Err feting rows for projects", err)
			//} else {

			//for currRowNdx := startDataRow; currRowNdx < rlen; currRowNdx++ {

			//rowndxstr := strconv.Itoa(currRowNdx)
			//row := rows[currRowNdx-1]
			//axis := rcfg.ColMatchName + rowndxstr
			aProj := plan.ParseProjectHeaderSection(row, pctx)

			if aProj == nil {
				continue
				//break // empty project ID and name so End stop processing
			}

			// Check to see if the project already exists
			// in the system and if so modify it's ID so we
			// do not loose duplicates across sheets.
			_, projAlreadyExists := oin.ItemsById[aProj.Id]
			if projAlreadyExists == true {
				dupProjIdCntr++
				aProj.Id = aProj.Id + "-DUP" + strconv.Itoa(dupProjIdCntr)
			}

			aProj.ParseColumnBasedResourceNeeds(row, pctx, plan)

			oin.Items = append(oin.Items, aProj)
			oin.ItemsById[aProj.Id] = aProj

			// Handle % Complete if specified
			if aProj.PortionCompleteAtStart > 0 {
				aProj.TotDirectUnitsDelivered = aProj.TotDirectResourceUnit * aProj.PortionCompleteAtStart
			}
			fmt.Println("L293: new Project=", aProj.ToJSON(false))

			// Index project so they can be found by
			// alias.
			for _, alias := range aProj.Alias {
				if alias != aProj.Id {
					oin.ItemsByAlias[alias] = aProj
				}
			}
			//fmt.Println("tres=", tres, "agrp=", agrp)
		}
		// Now we have loaded the complete set we can
		// update indexing for projects that reference each other
		//

	}

	return nil
}
