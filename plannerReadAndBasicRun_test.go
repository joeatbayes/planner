package planner_test

import (
	"fmt"
	//"os"
	"testing"

	//github.com/joeatbayes/portfolio-roadmap-planner/src/portfolio"
	"portfolio"
)

// NOTE: When testing local change the "github.com//joeatbayes/portfolio-roadmap-planner/src/portfolio/"
//  to "portfolio" to allow testing without pushing changes to git server.

func TestPlannerRead(t *testing.T) {

	fmt.Println("Test Loading Config")
	configFiName := "data/examples/sample-01/config-sample-01.xlsx"
	resGrpFiNames := []string{"data/examples/sample-01/resources/teams-s1.xlsx"}
	projListFiNames := []string{"data/examples/sample-01/projects/projects-s1.xlsx"}
	outExcelName := "out/test-1.xlsx"

	pln := portfolio.PlannerBasicRun(configFiName, resGrpFiNames, projListFiNames, outExcelName)
	projects := pln.Proj

	// TODO: Add additional checks to see if data is in locations we expect
	//  commented out for now while the underlying structures are rapidly
	//  evolving.

	//if pln.ResAvail.AllRes[0].ResId != "SDE" {
	//	fmt.Println("expected pln.ResAvail.AllRes[0].ResId = SDE but it equaled ", pln.ResAvail.AllRes[0].ResId)
	//	t.FailNow()
	//}

	//if pln.ResAvail.AllRes[0].ByDay[0].Avail != 84 {
	//	fmt.Println("expected pln.ResAvail.AllRes[0].ByDay[0].Avail = 84 but it equaled ", pln.ResAvail.AllRes[0].ByDay[0].Avail)
	//	t.FailNow()
	//}

	if projects.Items[0].Id != "W001" {
		fmt.Println("Expected ID for  project should be W001 and is ", projects.Items[0].Id, " fiName=", projListFiNames)
		t.FailNow()
	}

	//if projects.Items[0].Resources[0].Id != "SDE" {
	//	fmt.Println("ID of first resource for first project should be SDE and is ", projects.Items[0].Resources[0].Id, " fiName=", projListFiNames)
	//	t.FailNow()
	//}

	fmt.Println("\n\n\n\n\nplan=", pln.ToJSON(true))

	//pln.PrintBasicRoadMap()

}
