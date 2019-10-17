package main

// portAnalyze.go  - Command Line Interface for Portfolio Analyzer / Interpolator

import (
	//"bufio"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	//"path/filepath"
	//m2h "mdtohtml"
	//"path"
	//"path/filepath"
	"portfolio"
	//"regexp"
	//"sort"
	//s "strings"
	//"time"

	"github.com/joeatbayes/goutil/jutil"
)

func PrintHelp() {
	fmt.Println(`Sample Usage	
analyzePortfolio -config=data/examples/sample-01/config-sample-01.xlsx -resources=data/examples/sample-01/resources  -projects=data/examples/sample-01/projects -out=out/sample-1.xlsx

 
  -resources = A set of one or more paths to files or 
        directories containing resource files.  Pathnames
		are separated by comma.  If path is to file it will
		be read.  If path is to a directory then all excel
		files in that directory will be read and combined
		as if they had been defined in a single file. 
		
  -projects = A set of one or more paths to files or directories
         containing project files. Paths are separated by
		 comma.  If path is to directory it will read all excel
		 files in that directory. If to a specific file it will
		 only read that file.   All paths listed will be read
		 and combined as a set of projects as if they had been
		 defined in a single file.  
		
  -config = Name of a excel file the system should read that 
         contains the column mappings for both resources
		 files and project files.  The config file can also
		 specify some items that will modify operation of the
		 system. 
		
  -out = A path specifying the excel file that should be 
         generated.  The system will generate several files
		 at this location.  For example if you sepecity
		 -out=out/fil1.xlsx it will generate out/file1.cost.xlsx
		 and several others. 
		
  -loopDelay - When set the system will process the input.  
         Sleep for a number of seconds and then re-process.
		 This is intended to keep a generated file available to 
		 easily reload.  eg:  -loopDelay=30 will cause the system to 
		 reprocess the input files once every 30 seconds.
		
  -usageRep - When set to Y will generate usage Report. Otherwise
         will skip.  The ussage report can generate very large excel
		 files that cause out of memory errors on very large portfolios.
		
   
	
   We made a design decsion to name each of the input files so you can mix and match.
   It is common to run the same set of projects against several different resource files
   when doing what-if analysis.
	
	-`)
}

func main() {
	startms := jutil.Nowms()

	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}

	if parms.Exists("out") == false ||
		parms.Exists("projects") == false ||
		parms.Exists("resources") == false ||
		parms.Exists("config") == false {
		fmt.Println("Missing Mandatory command line parameter")
		fmt.Println(parms.String())
		PrintHelp()
		return
	}

	fmt.Println(parms.String())
	plan := portfolio.MakePlannerFromArgs(parms)
	plan.BasicRun()
	jutil.Elap("Finished Basic Run", startms, jutil.Nowms())

}
