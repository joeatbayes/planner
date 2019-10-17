package portfolio

import (
	"encoding/json"
	"fmt"

	//"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ResourceConfig struct {
	SheetNames              []string
	StartDataRow            int
	ColGroupName            string
	ColStartingCount        string
	ColNetUnitsPerResPerDay string
	ColAvgCostPerUnit       string
	ColMatchName            string
	ColMinPerProj           string
	ColMaxPerProj           string

	// Column to start at when looking at resources levels by month
	ColStartByMonth string
	// Column to stop at when looking at resource levels
	ColEndByMonth string
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
	ColUsageModel string

	// Column in project resources to define the kind of resource.
	// labor is scheduled over time,   purchase is not treated as
	// a limited resource since presumably purchase dollars are
	// controlled through another process.
	ColResourceType string
}

func (cfg *ResourceConfig) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func LoadResourceConfig(fiName string) *ResourceConfig {
	f, err := excelize.OpenFile(fiName)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	cfg := new(ResourceConfig)
	// allocate resNeeds as empty slice but with memory pre allocated to store upto 1000 elements.

	// Get value from cell by given worksheet name and axis.
	tabName := "resources"
	cfg.SheetNames = strings.Fields(GetCellStr(f, tabName, "C3", ""))
	cfg.StartDataRow = GetCellInt(f, tabName, "C4", -1)
	cfg.ColGroupName = GetCellStr(f, tabName, "C5", "")
	cfg.ColStartingCount = GetCellStr(f, tabName, "C6", "")
	cfg.ColNetUnitsPerResPerDay = GetCellStr(f, tabName, "C7", "")
	cfg.ColAvgCostPerUnit = GetCellStr(f, tabName, "C8", "")
	cfg.ColMatchName = GetCellStr(f, tabName, "C9", "")
	cfg.ColMinPerProj = GetCellStr(f, tabName, "C6", "")
	cfg.ColMaxPerProj = GetCellStr(f, tabName, "C10", "")
	cfg.ColStartByMonth = GetCellStr(f, tabName, "C11", "")
	cfg.ColEndByMonth = GetCellStr(f, tabName, "C12", "")

	return cfg
}
