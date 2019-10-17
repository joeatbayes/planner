package portfolio

import (
	"encoding/json"
	//"fmt"
	//"strconv"
	//"strings"
	//	"github.com/360EntSecGroup-Skylar/excelize"
)

type ResourceType int
type UseModel int

const (
	ResTypeLabor ResourceType = 0
	ResTypePurch ResourceType = 1

	UseModelEven   UseModel = 0
	UseModelEarly  UseModel = 1
	UseModelLate   UseModel = 2
	UseModelBefore UseModel = 3
	UseModelAfter  UseModel = 4
)

var UseModelLook = map[string]UseModel{"even": 0, "early": 1, "late": 2, "before": 3, "after": 4}
var ResTypeLook = map[string]ResourceType{"labor": 0, "purch": 1}

type ResourceGrp struct {
	Id             string
	GroupName      string
	StartCount     float32
	UnitsPerDay    float32 // Units/Hours Per Day Per resource used to compute available hours
	AvgCostPerUnit float32
	MaxPerProj     float32   // Max # resources that can be assigned to project simutaneously
	MinPerProj     float32   // Min # resources needed to allow project to be started
	NumMonth       int       // Number of Months the user specified hours for
	CntByMonth     []float32 // Count of resources by month specified by user
	TotCnt         float32   // Total Count of resources by user
	LastCnt        float32   // Last count of resources specified by user

	// Column in project resources to define default resource usage
	// model.   even, early, late, before, after.
	// even is default and means spread usage evenly across life of project
	// early means use max resources possible early in project
	// late means use no resources at all until using minimum resources
	// would finish by the time the most damanding resource finishes.
	// before  means to use up all this resource before consuming any
	// other resource that is not marked before.
	// after means to use none of this resource until all resources
	// not with use model of after are fulfilled.
	UsageModel UseModel

	// Column in project resources to define the kind of resource.
	// labor is scheduled over time,   purchase is not treated as
	// a limited resource since presumably purchase dollars are
	// controlled through another process.
	ResourceType ResourceType

	// The elements potion in the array
	// Kept here to may reverse lookups fast
	Position int
}

type ResourceGrps struct {
	FiNames   []string
	Items     []*ResourceGrp
	ItemsById map[string]*ResourceGrp
	ByName    map[string]*ResourceGrp
}

func (oin *ResourceGrps) NumGroups() int {
	return len(oin.Items)
}

func (oin *ResourceGrps) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(oin, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(oin)
		return string(slcB)

	}
}

func (oin *ResourceGrp) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(oin, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(oin)
		return string(slcB)

	}
}
