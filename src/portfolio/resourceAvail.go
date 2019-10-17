package portfolio

import (
	"encoding/json"
	"fmt"
	//"math"
	"strconv"
	//"strings"
	//"github.com/360EntSecGroup-Skylar/excelize"
)

type ResUsed struct {
	ProjId string

	// Number of units consumed by this Project
	Units float32

	// The # of units desired by this project.
	// Can help show resource shortages when
	// Units is less than Desired Units.
	DesiredUnits float32
}

type ResByDay struct {
	//ResId  string
	// Hours available to use for this resource for
	// this day.
	Avail float32

	// Hours used for this resource for this day
	Used float32

	// Total hours available for this resource for this day
	Total float32

	// The project who used this resource for this Day
	UsedBy map[string]ResUsed
}

// A array of resource usage by day for a single
// resource.
type ResTrack struct {
	ResId string
	ByDay []ResByDay
}

// A matix that contains one array element for each resource
// defined in the system. Each array contains a array
// days each of which measures resources available, used
// and what project used them.
type ResAllTrack struct {
	AllRes []ResTrack // This array is in exactly the same order
	// that we load resource needs for projects which is
	// the same order project resources are listed in config
	// project resources.  This is to allow synchronized
	// searching.
}

func (rall *ResAllTrack) MaxResourceDay() int {
	return len(rall.AllRes[0].ByDay) - 1
}

func (cfg *ResAllTrack) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func (cfg *ResTrack) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func (pln *Planner) ProduceResourceByDayMatrix() *ResAllTrack {
	rall := new(ResAllTrack)
	rgrps := pln.Res
	numGrp := pln.Cfg.Proj.NumResourceNeeds
	rall.AllRes = make([]ResTrack, numGrp, numGrp)
	fmt.Println("L176: num Groups when making all res = ", len(rgrps.Items))
	allRes := rall.AllRes
	gitems := rgrps.Items // resources to supply needs
	for gndx, rgrp := range gitems {
		if rgrp == nil {
			panic("L180 ERROR FATAL resourceGroups.items can not be nil at ndx " + strconv.Itoa(gndx) + " rgrps=" + rgrps.ToJSON(true))
		}
		calcDay := rgrp.NumMonth * WorkDaysPerMonth
		numDay := MaxInt(calcDay, DefaultNumberDaysModel)
		//allRes[gndx] = new(ResTrack)
		wrkByDays := make([]ResByDay, numDay, numDay)
		allRes[gndx].ResId = rgrp.Id
		allRes[gndx].ByDay = wrkByDays

		//WorkDaysPerMonth
		// initialize by day to last day
		// and create it's resource used
		for dndx := 0; dndx < numDay; dndx++ {
			wrkByDays[dndx].UsedBy = make(map[string]ResUsed)
			//wrkByDays[dndx].Id = rgrp.Id
			monthNum := dndx / WorkDaysPerMonth
			resCount := rgrp.LastCnt
			// Lookup the actual number of resources
			// if user specified them. Otherwise just
			// project the last one they did specify
			if monthNum < rgrp.NumMonth {
				// override default with count specified
				// by the user for that month
				resCount = rgrp.CntByMonth[monthNum]
			}
			wrkByDays[dndx].Avail = resCount * rgrp.UnitsPerDay
			wrkByDays[dndx].Total = resCount * rgrp.UnitsPerDay
		}
	}
	return rall
}
