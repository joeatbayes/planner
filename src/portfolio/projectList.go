package portfolio

import (
	"encoding/json"
	"fmt"

	//"os"
	//"strconv"
	"sort"
	"strings"
	//	"github.com/360EntSecGroup-Skylar/excelize"
)

type ProjectResourceNeed struct {
	Id          string
	Units       float32 // normally hours but left as Units to be more generic
	MaxAssigned float32 // Max # resources that can be assigned to this resource
	MinAssigned float32 // Min # resources that can be assigned to this resource
	Position    int

	/* Took another approach so do not use this */
	//UnitsRelativeToMax float32 // If we assume that it is desirable to allocate
	// all resources to they are completed by the time the most demanding
	// resource has been delivered then we would Adjust the less demanding
	// resource allocated per day demand less per day. Otherwise if we need
	// hundred hours of programmer time and 10 hours of scrum master then
	// need to deliver 1 hour of scrum master time for every hour of
	// of programmer time.

	// Contains all work delivered from working this project for this
	// resources need.
	UnitsDelivered float32

	// We take MaxResources Assigned and multiply
	// by the number of hours per day per resource to figure out maximum
	// hours we could deliver in one day.
	MaxUnitsPerDay float32

	MinUnitsPerDay float32 // Take Min Resources Assigned and multiply
	// by the number of u

}

type Project struct {
	Id                string
	Alias             []string
	Priority          float32
	EffectivePriority float32
	OutputPriority    float32
	Name              string

	// Target date when this project should be completed.
	// normally specified by project sponsor due to business
	// goals or to meet specific legal compliance needs.
	// System will modify some output reports to show in
	// warning color when target date can not be met.
	TargetDate string

	// The Target Date day number is set when target date
	// is populated.
	TargetDateDayNum int

	// Some reports will be modified to show in warning color when
	// project completion is after specified.  The system can
	/// also use max duration to optimize or recomend changes
	// that can reduce max duration.  This is set automatically
	// from the target date specified but can be set manually.
	// in liue of the target date.  Since target dates shift all
	// the time based on when funding is recieved completion days
	// can be more flexible than hard coded dates.
	MaxCompletionDays float32

	// min Completion Days to augment min Duration Days.
	// Min Completion is relative to start of planning while
	// min duration is relative to start of the project.
	// When we know that we want to end a project exactly 365
	// days from the start of planning then it is easier to
	// use min completion and still let start time vary.
	// If you set min duration and min completion to the
	// same value then you can gaurante a start and end
	// time relative to project start.
	MinCompletionDays float32

	Parent          string
	Precurs         []string
	PrecursComplete bool
	PrecursFor      []string // updated from all precursors in other projects
	Approved        bool
	NeedBy          []string
	NPV             float32

	// The percentage of the project already complete before this planning
	// session is started. This amount can be added to totalUnitsDelivered
	// for both project and ProjectResourceNeed.
	PortionCompleteAtStart float32

	// The Minimum duration in days that a project can be allowed to
	// run. This is used to stretch projects like KTLO out over an
	// extended period with consistent resource utilization.
	MinDuration float32

	// Total resource units directly requested for all resource columns.
	// but does not include units requested by any children.
	TotDirectResourceUnit float32

	// Total Resource units delivered from either working the project
	// or from % complete at the start of the project.  Does not include
	// work delivered to children.
	TotDirectUnitsDelivered float32

	Children []string // updated from all Parent in other projects
	// the same order as specified in the config
	// file project resources here and in resource
	// available to allow numeric indexing.

	Needs []string // updated from all needBy in other projects

	StartWorkDay int
	EndWorkDay   int
	IsComplete   bool

	NotWorkable bool // set to true when we already know we can not
	// find all necessary resources on any future day to work on this
	// project.

	Resources []*ProjectResourceNeed // This is created in exactly
	// the same sequence as defined in the configuration file.  To allow
	// fast lookup.
	MostDemandResource *ProjectResourceNeed // This is a fast pointer to
	// the most demanding resource because we want to meter delivery
	// of other resources based on what we can deliver of the most
	// demanding resource.
	MinDurationWorstResource float32

	// Allows a single field like comments or description to have additional
	// data encoded in the form --=fieldname=fieldval  eg:  --=minDuration=62
	AdditFlds    StrStrMap
	AdditFldsStr string

	Log []string
}

type Projects struct {
	FiNames      []string
	Items        []*Project
	ItemsById    map[string]*Project
	ItemsByAlias map[string]*Project

	NeedByOrphens  []string // list of projects listed in NeedBy but not defined
	PrecursOrphens []string // list of precursors listed but not defined
	ParentOrphens  []string // list of parents listed but not defined
	LastDayWorked  int      // Last day in plan work was done. Used to figure out rendering
	lastError      error    // last error encountered. Should be nil most of time
	EmptyIdCnt     int
}

// Return the computed array of PrecursFor as a string delimited
// by the specified delimiter.  Useful to inlcude the entire output
// into a single report cell.
func (proj *Project) PrecursForAsStr(delim string) string {
	return strings.Join(proj.PrecursFor, " ")
}

// Return the computed needs by field with all values in the array
// concatonated into a string delimited by the specified delimiter
// useful to include the entire set in a single report cell.
func (proj *Project) NeedByAsStr(delim string) string {
	return strings.Join(proj.NeedBy, " ")
}

// Return the computed Children field with all values in the array
// concatonated into a string delimited by the specified delimiter
// useful to include the entire set in a single report cell.
func (proj *Project) ChildrenAsStr(delim string) string {
	return strings.Join(proj.Children, " ")
}

// Make a new Project List and initialize components
// so they are ready to start loading from a file.
func MakeProjects() *Projects {
	tres := new(Projects)
	tres.FiNames = make([]string, 0, 100)
	tres.Items = make([]*Project, 0, 1000)
	tres.ItemsById = make(map[string]*Project)
	tres.ItemsByAlias = make(map[string]*Project)
	tres.NeedByOrphens = make([]string, 0, 100)
	tres.PrecursOrphens = make([]string, 0, 100)
	tres.ParentOrphens = make([]string, 0, 100)
	return tres
}

// return the number of resources allocated for our first project
func (oin *Projects) NumResources() int {
	if len(oin.Items) == 0 {
		return 0
	}
	return len(oin.Items[0].Resources)
}

// Create a set Projects Instance and load the contents of
// many files into project List
func (plan *Planner) MakeAnLoadProjectLists(fiNames []string) *Projects {
	plan.Proj = MakeProjects()
	plan.Proj.LoadProjectLists(fiNames, plan)
	plan.Proj.DoPostLoadAdjustments()
	return plan.Proj
}

func (oin *Projects) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(oin, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(oin)
		return string(slcB)
	}
}

func (oin *Project) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(oin, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(oin)
		return string(slcB)
	}
}

func (oin *ProjectResourceNeed) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(oin, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(oin)
		return string(slcB)
	}
}

func (oin *Project) ArePrecursComplete(allProj *Projects) (bool, int) {
	lastComplete := 0
	if len(oin.Precurs) == 0 {
		// No precursors to find so just return.
		return true, 0
	}
	for _, precId := range oin.Precurs {
		pproj, found := allProj.ItemsById[precId]
		if found == false {
			fmt.Println("L97: Could not find Precursor ", precId, " for project", oin.ToJSON(false))
		} else if pproj.IsComplete == false {
			// At least one precursor is not complete
			// so we can exit with false
			return false, 0
		} else if lastComplete < pproj.EndWorkDay {
			lastComplete = pproj.EndWorkDay
		}
	}
	return true, lastComplete
}

func (oin *Project) GetMaxResourceUnits() float32 {
	return oin.MostDemandResource.Units
}

/* TOOD A different approach that makes this unecessary
// Adjust units relative to max for all resources needed
// for this project.  This is used to allow metered Delivery
// of resources where demand is smaller relative to larger demand
// eg: Need 600 hours of programmers but only 100 hours of Lead.
func (oin *Project) AdustUnitsRelativeMax(rgrps *ResourceGrps) {
	tmax := oin.GetMaxResourceUnits()
	for _, res := range oin.Resources {
		if res == nil {
			panic("L167: FATAL ERROR: res can not be nil when ieterating *Proj.Resources")
		}
		if res.Units > 0 {
			res.UnitsRelativeToMax = res.Units / tmax
			_, found := rgrps.ItemsById[res.Id]
			if found != true {
				fmt.Println("L93: ERROR Could not find matching resource Group id=", res.Id, "for project=", oin.Id)
			}
		} else {
			res.UnitsRelativeToMax = 0
		}
	}
}
*/

func (oin *Projects) SortByOutputPriority() {
	items := oin.Items
	for _, witem := range items {
		witem.OutputPriority = witem.Priority
	}
	for pass := 0; pass < 5; pass++ {
		// run multiple passes to handle priority changes
		// from nested children
		for _, witem := range items {
			for _, cid := range witem.Children {
				childItem := oin.ItemsById[cid]
				if childItem.OutputPriority < witem.OutputPriority {
					witem.OutputPriority = childItem.OutputPriority - 0.0001
				}
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].OutputPriority < items[j].OutputPriority
	})
}

func (oin *Projects) SortByEffectivePriority() {
	items := oin.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].EffectivePriority < items[j].EffectivePriority
	})
}

/*
// Adjust Units relative to max for all projects
func (oin *Projects) AdjustUnitsRelativeMax(pln *Planner) {
	for _, proj := range oin.Items {
		proj.AdustUnitsRelativeMax(pln.Res)
	}
}*/

func (proj *Project) AddLog(sender string, str string) {
	// TODO Add Timems to each line
	outStr := sender + "\t" + str
	proj.Log = append(proj.Log, outStr)
	fmt.Println(outStr, " FROM: Project.AddLog")
}
