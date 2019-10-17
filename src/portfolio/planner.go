package portfolio

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/joeatbayes/goutil/jutil"
	//"github.com/360EntSecGroup-Skylar/excelize"
)

type Planner struct {
	Cfg      *Config
	Proj     *Projects
	Res      *ResourceGrps
	ResAvail *ResAllTrack
	files    *ExcelFiManager

	// Generic stuff from command line parser
	Start             float64
	Perf              *jutil.PerfMeasure
	Pargs             *jutil.ParsedCommandArgs
	LoopDelay         float32
	ProjectFileNames  []string
	ResourceFileNames []string
	ConfigFileName    string
	OutFileName       string
	GenUsageRep       bool
}

func (cfg *Planner) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func MakePlannerDirect(configFiName string, resourceFiNames []string, projectsFiNames []string, outFiName string) *Planner {
	pln := new(Planner)
	pln.Pargs = nil
	pln.Start = jutil.Nowms()
	pln.Perf = jutil.MakePerfMeasure(25000)
	fmt.Println("L50: MakePlannerDirect() configFiName=", configFiName, "resGrpFiNames=", resourceFiNames, " projListFiNames=", projectsFiNames, " outExcelName=", outFiName)

	pln.GenUsageRep = false
	pln.LoopDelay = -1
	pln.ConfigFileName = configFiName
	pln.ProjectFileNames = projectsFiNames
	pln.ResourceFileNames = resourceFiNames
	pln.OutFileName = outFiName
	pln.files = MakeExeclFiManager(pln.OutFileName)
	if jutil.Exists(pln.ConfigFileName) == false {
		log.Fatal("L58: FATAL ERROR: Config File ", configFiName, " does not exist")
	}

	for _, fname := range pln.ProjectFileNames {
		if jutil.Exists(fname) == false {
			log.Fatal("L58: FATAL ERROR: Project files ", fname, " does not exist")
		}
	}

	for _, fname := range pln.ResourceFileNames {
		if jutil.Exists(fname) == false {
			log.Fatal("L69: FATAL ERROR: resources files ", pln.ResourceFileNames, " does not exist")
		}
	}

	jutil.EnsurDir(pln.OutFileName)
	return pln
}

func MakePlannerFromArgs(parms *jutil.ParsedCommandArgs) *Planner {
	projFiles := strings.Split(parms.Sval("projects", ""), " ")
	resFiles := strings.Split(parms.Sval("resources", ""), " ")
	outFile := parms.Sval("out", "")
	configFile := parms.Sval("config", "")
	pln := MakePlannerDirect(configFile, resFiles, projFiles, outFile)
	pln.GenUsageRep = parms.Bval("usageRep", false)
	pln.LoopDelay = parms.Fval("loopdelay", -1)
	return pln
}

func (pln *Planner) LoadPlannerFiles() {
	pln.Cfg = LoadConfig(pln.ConfigFileName)
	pln.Res = LoadResourceGroups(pln.ResourceFileNames, pln.Cfg)
	pln.Proj = pln.MakeAnLoadProjectLists(pln.ProjectFileNames)
	pln.ResAvail = pln.ProduceResourceByDayMatrix()
	pln.GenUsageRep = false
	//pln.Proj.AdjustUnitsRelativeMax(pln)
}

// find highest priority project that is not complete
// and return it and the ndx it was located at. If none
// found then return -1, nil to indicate no work remaining
func (allProj *Projects) FindHighestPriorityIncomplete(startNdx int) (int, *Project, int) {
	for ndx := startNdx; ndx < len(allProj.Items); ndx++ {
		proj := allProj.Items[ndx]
		//fmt.Println("L49: FindHighestPriorityIncomplete startndx=", startNdx, " ndx=", ndx, " proj.Id=", proj.Id, " proj=", proj)
		if proj.IsComplete == false && proj.NotWorkable == false {
			fmt.Println("L51: FindHighestPriorityIncomplete Found a project ndx=", ndx, "proj.Id=", proj.Id)
			return ndx, proj, proj.EndWorkDay
		}
	}
	return -1, nil, -1 // no project is available
}

// Update start stop of parent projects so they
// start as early as first child and finish as
// late as the last child.  Return true if anything changed
func (pln *Planner) UpdateParentStartStop() bool {
	pall := pln.Proj
	changed := false
	for _, prj := range pall.Items {
		if len(prj.Children) < 1 {
			// No children available
			continue
		}
		for _, childId := range prj.Children {
			childProj, found := pall.ItemsById[childId]

			if found != true {
				fmt.Println("L73: Cound not find childId ", childId, " for proj", prj.ToJSON(false))
				continue
			}
			if prj.StartWorkDay < 0 && childProj.StartWorkDay > -1 {
				prj.StartWorkDay = childProj.StartWorkDay
				changed = true
			}
			if childProj.StartWorkDay < prj.StartWorkDay && childProj.StartWorkDay > -1 {
				prj.StartWorkDay = childProj.StartWorkDay
				changed = true
			}
			if childProj.EndWorkDay > prj.EndWorkDay {
				prj.EndWorkDay = childProj.EndWorkDay
				changed = true
			}
			if childProj.IsComplete == false {
				prj.IsComplete = false
				changed = true
			}
		}
		if prj.IsComplete == false {
			// If not complete by the time we get here then
			// one of our children could not be worked
			// or one of our resources were not available.
			prj.NotWorkable = true
		}
	}
	return changed
}

func (pln *Planner) Analyze() {
	pln.Proj.SortByEffectivePriority()
	for {
		// Take Highest Priority Task Not complete
		pndx, highProj, _ := pln.Proj.FindHighestPriorityIncomplete(0)

		if pndx == -1 {
			fmt.Println("L63: No more workable projects available")
			break
		}
		// TODO: Check HighProj predecessors are all complete
		//  before starting this work.

		// Find first Day when all resources are available.
		fmt.Println("L117: Found project to work id=", highProj.Id)
		dayNdx := pln.ResAvail.FindFirstDayWhenAllResourcesAvailable(highProj, 0)
		fmt.Println("L119: firstAvailDay dayNdx =", dayNdx)
		if dayNdx == -1 {
			highProj.NotWorkable = true
			continue
		}

		// Do not set StartDay here because it is set in the work project.
		lastDay := pln.ResAvail.WorkProject(highProj, dayNdx, pln, 0)
		if lastDay != -1 {
			// project fully worked so now can
			highProj.EndWorkDay = lastDay
			highProj.IsComplete = true
		} else {
			// Project can not be worked due to
			// running out of resource days
			pln.ResAvail.ReturnResources(highProj)
			highProj.NotWorkable = true
			continue
		}
	}

	// Run Parent Update multiple times to
	// handle any parents that reference other
	// parents.  Runs multiple times to handle
	// nested changes.
	ucnt := 0
	for {
		changed1 := pln.UpdateParentStartStop()
		changed2 := pln.UpdateParentStartStop()
		if changed1 == false && changed2 == false {
			break
		}
		ucnt++
		if ucnt > 25 {
			break
		}
	}

	// See if all Resources are available
	// If so start work and record resource usage
	// by looping through days until all Resources
	// have been delivered. And Mark Complete on
	// that day.
	//
	//
	// Find next highest priority task not comlete
	//
	//
}

// FUNC PlannerAdjustParentStartTop
// FUNC CalcTotalCost

func (pln *Planner) PrintBasicRoadMap() {
	fmt.Println("Basic Roadmap")
	for pndx, proj := range pln.Proj.Items {
		fmt.Println("ndx=", pndx, " id=", proj.Id, "priority=", proj.Priority, "effectivePriority=", proj.EffectivePriority, "startDay=", proj.StartWorkDay, "endDay=", proj.EndWorkDay, "Complete=", proj.IsComplete, "NotWorkable=", proj.NotWorkable, "name=", proj.Name)
	}
}

func leftSlice(str string, num int) string {
	if len(str) > num {
		return str[0:num]
	} else {
		return str
	}
}
func (pln *Planner) GenBasicRoadMap(sheetName string) {
	pln.Proj.SortByOutputPriority()
	fmt.Println("Basic Roadmap")
	f := pln.files.GetWithLastSeg("roadmap", false)
	f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "Basic Roadmap Start / Stop")
	f.SetCellValue(sheetName, "A2", "ID")
	f.SetCellValue(sheetName, "B2", "Priority")
	f.SetCellValue(sheetName, "C2", "Effective Priority")
	f.SetCellValue(sheetName, "D2", "Start Day")
	f.SetCellValue(sheetName, "E2", "End Day")
	f.SetCellValue(sheetName, "F3", "Name")
	f.SetColWidth(sheetName, "F", "F", 45)
	f.SetCellValue(sheetName, "G2", "Not Workable")

	for pndx, proj := range pln.Proj.Items {
		axRow := strconv.Itoa(pndx + 3)
		f.SetCellValue(sheetName, "A"+axRow, proj.Id)
		f.SetCellValue(sheetName, "B"+axRow, proj.Priority)
		f.SetCellValue(sheetName, "C"+axRow, proj.EffectivePriority)
		f.SetCellValue(sheetName, "D"+axRow, proj.StartWorkDay)
		f.SetCellValue(sheetName, "E"+axRow, proj.EndWorkDay)

		f.SetCellValue(sheetName, "F"+axRow, leftSlice(proj.Name, 45))
		f.SetCellValue(sheetName, "G"+axRow, proj.NotWorkable)
	}
	//f.SetActiveSheet(sheetNdx)

}
