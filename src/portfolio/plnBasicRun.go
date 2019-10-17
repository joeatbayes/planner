package portfolio

import (
	"fmt"
	//"os"
	"github.com/joeatbayes/goutil/jutil"
	//github.com/joeatbayes/portfolio-roadmap-planner/src/portfolio"
	"log"
)

// NOTE: When testing local change the "github.com//joeatbayes/portfolio-roadmap-planner/src/portfolio/"
//  to "portfolio" to allow testing without pushing changes to git server.

func (pln *Planner) BasicRun() {

	start := jutil.Nowms()

	pln.files.RemoveDerivedFiles()
	jutil.Elap("L18: Finished Remove files", start, jutil.Nowms())
	pln.LoadPlannerFiles()
	jutil.Elap("L20: Finished Load", start, jutil.Nowms())

	cfg := pln.Cfg
	if pln.Cfg == nil {
		log.Fatal("L23: FATAL ERROR: Failed to load config file " + pln.ConfigFileName)
	}

	if cfg.Proj == nil {
		fmt.Println("projListFiNames=", pln.ConfigFileName)
		panic("L30: FATAL ERROR: project Config failed to load")
	}

	if cfg.Res == nil {
		fmt.Println("resGrpFiNames=", pln.ResourceFileNames)
		panic("L35: FATAL ERROR: Resource Config failed to load")
	}

	if pln.Res == nil || len(pln.Res.Items) < 1 {
		fmt.Println(":Failed load Resources resGrpFiNames=", pln.ResourceFileNames)
		panic("L39: FATAL ERROR resource groups failed to load ")
	}

	if pln.Proj == nil {
		fmt.Println("Proj can not be nil fiName=", pln.ProjectFileNames)
		panic("L41: FATAL ERROR proj can not be nil")
	}

	if len(pln.Proj.Items) < 1 {
		fmt.Println("Project.Items has no items fiName=", pln.ProjectFileNames)
		panic("L45: FATAL ERROR project.Items contains no project lines")
	}

	pln.Analyze()
	pln.GenUsageRep = true

	//fmt.Println("\n\n\n\n\nplan=", pln.ToJSON(true))
	jutil.Elap("Finished Analyze", start, jutil.Nowms())
	pln.Proj.SortByOutputPriority()
	pln.PrintBasicRoadMap()
	pln.RequestByProjByResReport("request")
	pln.GenBasicRoadMap("ProjStartStop")
	//pln.GenChartRoadMap("ProjRoadMap")
	pln.GenChartRoadMapFlexUnit("RoadmapWeek", 5, "Week")
	pln.GenChartRoadMapFlexUnit("RoadmapSprint", 10, "Sprint")
	pln.GenChartRoadMapFlexUnit("RoadmapMonth", 22, "Month")
	pln.GenChartRoadMapFlexUnit("RoadmapQuarter", 22*3, "Quarter")
	pln.GenChartRoadMapFlexUnit("RoadmapYear", 22*12, "Year")

	//pln.GenCostsPerPeriodReport("CostByDay", 1, 22*12, "Day")
	pln.CostsByProjByPeriodReport("CostByWeek", 5, 22*12, "Week")
	pln.CostsByProjByPeriodReport("CostByMonth", 22, 22*12, "Month")
	pln.CostsByProjByPeriodReport("CostByQuarter", 22*3, 22*12, "Quarter")
	pln.CostsByProjByPeriodReport("CostByYear", 22*12, 22*12, "Year")

	pln.RepResAvailVsCapacy("resCapWeek", "Resource Capacity By Week", 5, 22*12, "week")
	pln.RepResAvailVsCapacy("resCapSprint", "Resource Capacity By Sprint", 10, 22*12, "sprint")
	pln.RepResAvailVsCapacy("resCapMonth", "Resource Capacity By Month", 22, 22*12, "month")
	pln.RepResAvailVsCapacy("resCapQuarter", "Resource Capacity By Quarter", 22*3, 22*12, "Quarter")

	if pln.GenUsageRep == true {
		pln.UsageByProjByResByTimeUnitRep("UsageWeek", 5, 22, "week", "month")
		pln.UsageByProjByResByTimeUnitRep("UsageSprint", 10, 22, "sprint", "month")
		pln.UsageByProjByResByTimeUnitRep("UsageMonth", 22, 22*3, "Month", "Quarter")
	}

	pln.RepWorstResContraintByProject("BlockingResource")
	jutil.Elap("L83: Finished Report Gen", start, jutil.Nowms())
	pln.files.FlushAll()
	jutil.Elap("L85: Finished Flush Gen", start, jutil.Nowms())
}

func PlannerBasicRun(configFiName string, resGrpFiNames []string, projListFiNames []string, outExcelName string) *Planner {

	fmt.Println("L16: PlannerBasicRun() configFiName=", configFiName, "resGrpFiNames=", resGrpFiNames, " projListFiNames=", projListFiNames, " outExcelName=", outExcelName)

	pln := MakePlannerDirect(configFiName, resGrpFiNames, projListFiNames, outExcelName)
	pln.BasicRun()

	return pln
}
