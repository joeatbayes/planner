package portfolio

import (
	//"encoding/json"
	"fmt"
	//"math"
	"strconv"
	//"strings"
	//"github.com/360EntSecGroup-Skylar/excelize"
)

func (ravail *ResTrack) FindFirstDayWhenAvailable(rneed *ProjectResourceNeed, startDay int) (int, *ResByDay) {
	for ndx, bDay := range ravail.ByDay {
		if ndx < startDay {
			// TODO: Do this better to avoid wasted scan
			continue
		}
		if bDay.Avail >= rneed.MinUnitsPerDay {
			return ndx, &bDay
		}
	}
	return -1, nil
}

// Find the first Day when All Resources are Available and return
// a that days offset into the future. Otherwise if no day can be
// found return -1.
func (resAll ResAllTrack) FindFirstDayWhenAllResourcesAvailable(proj *Project, startDay int) int {
	maxNdx := -1
	if len(proj.Resources) != len(resAll.AllRes) {
		panic("L136: FATAL ERROR Number of Resources proj.Resources " + strconv.Itoa(len(proj.Resources)) + " must equal len resAll.AllRes " + strconv.Itoa(len(resAll.AllRes)))
	}
	for ndx, rneed := range proj.Resources {
		ravail := resAll.AllRes[ndx]
		fndx, _ := ravail.FindFirstDayWhenAvailable(rneed, startDay)
		if fndx == -1 {
			return -1 // At least one resource never available
		}
		maxNdx = MaxInt(maxNdx, fndx)
	}
	return maxNdx
}

// Pretend we are working the project to find the resource that will take the longest to
// deliver based on working MaxUnitsPerDay once the project is scheduled.
// If it's time is less than the projects specified MinDuration then return
// the longer time.
func (pln *Planner) FindLongestResource(proj *Project, startDay int) (float32, *ProjectResourceNeed) {
	//rall = pln.Res.MaxResourceDay()
	worstRes := proj.Resources[0]
	worstDur := float32(0.0)
	for _, rgrp := range proj.Resources {
		minDur := rgrp.Units / rgrp.MaxUnitsPerDay
		if minDur > worstDur {
			worstDur = minDur
			worstRes = rgrp
		}

		// TODO: Need to scan forward in days from specified day
		// to see how long it will take to deliver desired work
		// because a higher priority project may have consumed
		// enough resources to prevent delivery of full hours.

	}
	if float32(proj.MinDuration) > worstDur {
		worstDur = float32(proj.MinDuration)
	}
	return worstDur, worstRes
}

const MaxWorkRecDepth = 12

// Iterate through days until all the resources
// needed for the project have been delivered.
// record the usage in resource avaiable so we do not
// use it again.  Returns last day project was worked
// if we could not find all the resources before running
// out of days.
func (rall *ResAllTrack) WorkProject(proj *Project, startDay int, plan *Planner, recDepth int) int {
	fmt.Println("L80: WorkProject id=", proj.Id, "startDay=", startDay, "recDepth=", recDepth)
	maxDayNdx := rall.MaxResourceDay()
	maxDayUsedThisProj := startDay
	if proj.TotDirectResourceUnit <= 0.0 {
		return 0
	}

	if proj.IsComplete {
		return proj.EndWorkDay
	}

	if recDepth > MaxWorkRecDepth {
		// Do not allow circular references to loop forever
		fmt.Println("L92: ERROR Max Child Depth reached breakout.  Check circular dependencies")
		return 0
	}

	// Work all our precursors before we attempt
	// to work ourself
	lastAllPrecDay := 0
	for _, precProj := range proj.Precurs {
		// TODO: Problem with this one is that a child of lower priority
		// could get worked before a child of higher priority.
		// Should build a sort arry of effectivePriority, precursors
		// sort it by priority and work in that order.
		lastPrecDay := rall.WorkProject(plan.Proj.ItemsById[precProj], 0, plan, recDepth+1)
		if lastPrecDay > lastAllPrecDay {
			lastAllPrecDay = lastPrecDay
		}
	}
	if lastAllPrecDay > startDay {
		startDay = lastAllPrecDay
	}

	minDur, worstRes := plan.FindLongestResource(proj, startDay)

	// Add our desired records for our worst resources
	// so the GAP shows up in our project reports
	for dayNdx := 0; dayNdx <= startDay; dayNdx++ {
		worstNdx := worstRes.Position
		rtrack := rall.AllRes[worstNdx]
		rday := &rtrack.ByDay[dayNdx]
		rday.UsedBy[proj.Id] = ResUsed{ProjId: proj.Id, Units: 0, DesiredUnits: worstRes.MinUnitsPerDay}
	}

	fmt.Println("L85: projId=", proj.Id, " minDur=", minDur, " worstRes=", worstRes.ToJSON(false))
	proj.MostDemandResource = worstRes
	proj.MinDurationWorstResource = minDur
	// TODO: If The SDE project had a diferent MaxUnits per Day then
	// could deliver faster.  Figure out a report to show this
	// value.

	//worstResGrp := plan.Res.ItemsById[worstResNeed.Id]
	foundEnoughThisProj := true // Assume project will be sucessful until we find a resource we can not furnish.
	for pndx, res := range proj.Resources {
		foundEnoughThisResource := false
		rtrack := rall.AllRes[pndx]
		if res.MaxUnitsPerDay < 0 {
			fmt.Println("L94: ERROR MaxUnitsPerDay can not be negative id=", res.Id, " maxUnits=", res.MaxUnitsPerDay)
			continue
		}

		for dayNdx := startDay; dayNdx <= maxDayNdx; dayNdx++ {
			unitsNeed := res.Units - res.UnitsDelivered
			if unitsNeed <= 0 {
				foundEnoughThisResource = true
				break
			}

			maxDayUsedThisProj = MaxInt(dayNdx, maxDayUsedThisProj) // keep max day from all resources
			rday := &rtrack.ByDay[dayNdx]
			desiredUnits := res.MaxUnitsPerDay
			if desiredUnits > unitsNeed {
				desiredUnits = unitsNeed
			}

			rgrp := plan.Res.ItemsById[res.Id]

			fmt.Println("L140: dayNdx=", dayNdx, " projId=", proj.Id, "resId=", res.Id, "desiredUnits=", desiredUnits,
				"res.MaxUnitsPerDay=", res.MaxUnitsPerDay,
				"res.MinUnitsPerDay=", res.MinUnitsPerDay,
				"maxDayUsedThisProj=", maxDayUsedThisProj, "unitsNeed=", unitsNeed, "res.Units=", res.Units,
				"proj.StartWorkDay=", proj.StartWorkDay, " endWorkDay=", proj.EndWorkDay,
				"unitsDelivered=", res.UnitsDelivered, "units=", res.Units, "avail=", rday.Avail, "used=", rday.Used,
				"total=", rday.Total, "avail+used=", rday.Used+rday.Avail)

			if rgrp.UsageModel == UseModelEven {
				// Min Duration. When working a project we must look at the Resources needed by the project
				// between current day and end of project. Then we must look at duration in number of days
				// required finish the project and dividing the resource remaining by the resource need to
				// get the maxResource use for that resource per day.  If that number is below max available
				// we only take that amount.
				daysWorked := float32(dayNdx) - float32(startDay)
				daysRemain := float32(minDur) - float32(daysWorked)
				if daysRemain <= 1 {
					// No days left so don't adjust desiredUnits
					// to fit time
				} else {
					unitsNeedPerDay := float32(unitsNeed) / daysRemain
					if unitsNeedPerDay < res.MinUnitsPerDay {
						unitsNeedPerDay = 0
					}
					if unitsNeedPerDay < desiredUnits && unitsNeedPerDay > 0 {
						// reducing our demand to deliver a resource leveled
						// duration.
						desiredUnits = unitsNeedPerDay
					}
				}
			}

			// TODO: Consider LATE consumer where we should
			//  delay work until the last possible day based on
			//  known minDuration and the max we can work when
			//  project is almost done.
			fmt.Println("L178: dayNdx=", dayNdx, " projId=", proj.Id, "resId=", res.Id, "desiredUnits=", desiredUnits,
				"res.MaxUnitsPerDay=", res.MaxUnitsPerDay,
				"res.MinUnitsPerDay=", res.MinUnitsPerDay,
				"maxDayUsedThisProj=", maxDayUsedThisProj, "unitsNeed=", unitsNeed, "res.Units=", res.Units,
				"proj.StartWorkDay=", proj.StartWorkDay, " endWorkDay=", proj.EndWorkDay,
				"unitsDelivered=", res.UnitsDelivered)

			if desiredUnits > unitsNeed {
				desiredUnits = unitsNeed
			}

			tunits := desiredUnits
			if desiredUnits <= 0 {
				// do not do anything when we do not need any more work.
				// can not break because we may need more on a future day
				// when leveling.
				continue
			}

			if rday.Avail <= 0 {
				// This is a bottleneck area and should be recorded as such
				// Still have to finish the loop to record the days
				// we are waiting for resources.
				//continue // Nothing to use so look at next day
			}

			if tunits > rday.Avail {
				// TODO: Add Error catch Here since we did not have
				// sufficient resources then this this resource for
				// this day is a bottleneck.  Add it to it to a data
				// structure to output latter.  This should be by
				// project by resource by day so we can generate
				// a report of bottlenecks.
				tunits = rday.Avail
			}

			if tunits < 0 {
				fmt.Println("L125: ERROR Can Not Work Negative hours id=", res.Id)
				tunits = 0
			}

			if tunits > 0 || desiredUnits > 0 {
				fmt.Println("L226: desiredUnits=", desiredUnits, "tunits=", tunits)
				res.UnitsDelivered += tunits
				rday.Used += tunits
				rday.Avail -= tunits
				rday.UsedBy[proj.Id] = ResUsed{ProjId: proj.Id, Units: tunits, DesiredUnits: desiredUnits}

				if tunits > 0 {
					if dayNdx > plan.Proj.LastDayWorked {
						plan.Proj.LastDayWorked = dayNdx
					}

					if proj.StartWorkDay == -1 {
						proj.StartWorkDay = dayNdx
					}
				}
			}

		} // for dayNdx
		if foundEnoughThisResource == false {
			// If we could not find enough resources for this resource
			// then we also have to fail the project.
			foundEnoughThisProj = false
		}
	} // for proj.resources
	proj.EndWorkDay = maxDayUsedThisProj
	if foundEnoughThisProj == true {
		proj.IsComplete = true
	}
	return maxDayUsedThisProj
}

// Return any resources previously claimed by this project
// to the available pool so another project can use them.
// This is normally used when we tried to work a project and
// ran out of days before we could fully resource it.
func (rall *ResAllTrack) ReturnResources(proj *Project) {

}
