# Requirements,  Actions and Future Features 

## Actions:

* Double Check that change is work allocated to program changes roadmap appropriately.

* Add Max Duration to special fields.

* Add a True Child containment output for Roadmap slide.  EG: List parents first then their children then their children rather than priority order.

* Add a notion of objectives which is separate that precursor or parent.

* Support the notion of reserved capacity so you can state plan at X% of capacity

* Produce something like an Agile board with sprints across top,   Teams in Cells going down with projects listed in appropriate cell. Where possible draw red line from dependencies.

* Add Allocation report as a separate axis.  Eg:  Capitalized,  Expensed or Tech Deb vs Feature

* For Request report add a known total capacity to show when we exceed total capacity.

* Worst Resource for blocking start may not be the same as the resource blocking project start.  We should track them separate and modify usage report so blocking start is in grey and blocking once started is in yellow. 

* DONE:JOE:2018-08-31: Add By Project By Resource By Time unit showing Ask,  Avail,  GAP (hours needed that were not available to sell) show those with a gap before the project starts as normal a gap after projects starts in zero.

  * DONE:JOE:2018-07-29: Add Resource Requested to the work allocation tracking  when we record how much given to the project on that day.

* DONE:JOE:2019-08-01: Add Resource Request by Project By Resource so we can double check parser to be sure everything was loaded correctly.

* Modify  computation of the worst resource times to use working time rather than simple computed time because the project may take longer if some other resource has consumed hours leaving more than minimum hours and less than max hours available during the period when the project is being worked.   A more common issue would be that a project is competing with other projects for a heavily used resource and a higher priority project has already consumed most of the available resources for that resource.   May need to introduce a second concept which is  HowLong project takes to complete from the planning start date for the most limiting resource and a second of how long it takes to complete once started.  We already have how long it takes to complete once started.  Both should be shown on the repWorstResByProj.

* DONE:2018-08-01: When same projectId is encountered for second time then add a counter to it and make it unique.

* On roadmap and other reports include the priority number.

  Add support to parse version 1 project file input.

  * Add ability to change team or resource group size by a time unit other than month. For Agile planning it  would likely be by week. Add ability to interpret # days in the resource count matrix to a unit other than month.  For sprint planning you may need two weeks.
  * Add a Sprint # to all time divisor reports.  Allow sprint length to be  set in config file.
  * Add Ability to automatically create a  new program level task the first time we encounter a program description.  Assign it an automatic ID and then Then set any task using the same program level feature to use that parent ID as it's parent.  Set the target date equal to the farthest target date of any encountered children for the automatically created programs.
  * Allow feature to automatically set the target date of the parent to the farthest target date of any of it's children.  Or at least generate a report where the children target date is further out than the program target date. 
  * Add support for Agile Size to Hours Matrix in specified in configuration.  This will convert the lookup value to units by parsing the input string in a agile scaling size it could be 0=10 1=30 2=60 3=200 4=500 6=750 where 0 is tiny,  1= very small and should be estimated at 30 hours, etc.

* 

* In Roadmap view display the portion of projects that will not complete before their target date it red.   To do this we must compute a starting date or if it is not specified in config file use current date. Then we must subtract current date from target date to obtain a maximum duration.  Then we can generate the warning reports from the max duration. We can also just add the maxDuration as an option for in the input with the rule that when target date is specified it overrides maxDuration.

  - Add Max Duration as optional project  line input field.
  - Modify roadmap view to display items that exceed their MaxDuration when MaxDuration is set.
  - Compute MaxDuration from Target date when target date is specified. 

* Add minCompletionDays to augment minDurationDays.  MinCompletion is relative to start of planning while min duration is relative to start of the project.   

* Implement logic to support MinCompletionDays.

* DONE:JOE:20190731 Show projects running longer than MaxCompletionDays to the  roadmap tab.

* Add support for the Input Units to Resource units Mutliplier.  For more mature teams they often have a simple conversion where each point value is worth about 30 hours to a 3 point card is worth 90 hours.

* Provide ability to source different columns from different files for project lists.  read multiple files for project lists but pull some columns from one project file and other columns from other project files provided they match on ID. 
  - Demo reading multiple files with multiple column inputs for different resource types to allow use of existing files without adding new total columns.
  - Different file types such as product management may contain some columns in secondary files than say app dev. Need to support this notion to get a wholistic data view without needing to modify input files that can not be changed due to integration with other systems.
  - In config file consider adding  a section called dataRowsStart.  Also add a section that allows several different file formats to be defined so each column for a project could be sourced from a different column from a different file.  You just need to add the column for which file and a way to map different files eg:  fileType=prodm, col=3.  then you need a way to map which files are which type and the only way I know to do that would be a regex expression eg if the file   is of that type.
  
* Add notion of seminal field cell to config.  If that field does not contain the expected value or label string then we can assume somebody did something crazy like inserting or deleting columns and should abort for that file right then.

* Test Private Demo file with new changes.

* Test file reading multiple input project files for private demo

* Produce tab in output for warnings when more than one tab re-uses same project #

* Update roadmap generator to show projects finishing past target date with those cells in red.

* Convert Project # to upper case to avoid missing matchup

* Log Errors into a special error file where users can see why system refused to plan easier.

* Parse and Load the Config Basics TAB
  * Add to basic config the length of typical month. Replace any location in code that uses the 22* to use this value.     Also do the same for length of typical week. Replace any location using 5 days to use this value.  Otherwise will not work for retailers open 24X7.
  * Load Config Resources, Projects, Basics using same excel handle rather than re-opening which forces a reparse.
  
* Implement MustStartDate and MustStopDate support 

* Produce resource blocking Report by time unit   Loop across all days for each project for each resource type. If the detail line is as the resource type then can show NoVal if no resource needed,  Yellow if minimum resource need delivered,  Green if Max Resource need delivered,  Red if no resources available when it was needed.  If using the time divisor then show the worst color for the days in that period.  To do this we will need to add a memory structure that shows what the project thought it needed on a given day vs what was given to it.   To make this useful would need to know what the project thought it needed after min duration and resource usage type was applied.  To make that work we need to know which resource was actually the most blocking resource that caused time to extend. Provide a report which shows which resources are blocking the ability to do work and finish the project. Add ability to track which resources where not available at the requested rate by project as the system works the project.   Do this in a way that allows the system to generate a report that shows the bottleneck by resource by day.  Internally this may be modeled as  a array of resource days by resource type for the work project where each item in array includes  targetUnits,  AvailableUnits.  From this we can compute GAP Units show the most to least GAP by resource type in the bar graph by color. 

* Produce a report that rolls up child projects to show each layer cost per per month.  Produce variant of the costs reports sum the cost of parent program and then show the children under that parent.  Only display direct programs if they have not parent. To keep ordering any children that are independent of the parent and have no parent should be listed if the are higher priority than the parent.    Will need to support indenting the child labels as nested children to an arbitrary depth.

* Add support for resourceType to handle hardware purchase.  This must support the target date override.

* Add total by second time unit divisor to all reports by time unit.   This is to allow totals by Year

* Add report for resource usage by project by time unit to show  Show resources with incomplete usage in different color.   

* Add support for Usage Model.  Default used  **Even** - spread evenly through the project,  Late - Bill no hours until as late in project as possible,   First - Bill all hours for this resource on days first do not use other resources until all first resources have been delivered,  Last - Bill all hours after all other use models are complete.

* Write a test that uses multiple layers of parent containment as nested programs.

* Produce a report that shows program  / Parent cost based on resource  costs for children.   Write a report showing the containment hierarchy rather than natural order. We can use the output priority which is only Loosely set from the effective priority as part of its recursive I walk to of parents and children that way you could do a simple sort for the report itself and avoid having to do the recursive work inside the report.

* Add Support for the Apply Approved Filter Flag which also requires adding the Approved Filter Flag support in project list.  Implement support for the approved column.  Includes ability to override it from config file prior to operation.

* Develop a more efficient way to locate the first record which contains available resources for a resource type.   This may be as easy as recording the earliest day where not all resources have been used or where more than than default minimum resource for the type is available. Find a more efficient way to find the first day when all the resources needed are available to start a project.  Current approach of linear search works but is not very optimal. 

* Double check estimated total cost for projects.

* Produce a utility that generates a spreadsheet nearly identical to format needed for  private Client to demonstrate ability to generation format in form they desire.

* Produce a clone utility that will read one spreadsheet and update another mapping columns between them.   Should only update the columns that are mapped leaving others in place to support editing specific detail in one without touching the source sheet. 

* * 

* Modify costs report to show the costs delivered prior to starting this planning cycle before the first time unit costs. 

* Add the notion of cost center to resource type.   Produce a report by Month of the costs per resource type.  Assume that hardware costs will be all accrued at the first of the project.      This leads to the notion of resource types.  Some will be accrued at first of the project.  Others will be spread across the life of the project.  others will wait until the end of the project and then be accrued.  Hardware will be easy as a one off cost on a given day either early or late.   Labor is more complex since for late in project cost should be deferred using resourced until latter in the project.    This leads to the notion of specifying the type of allocation approach in the configuration spreadsheet.

* Add Support for multiple resource files Modify ResourceList Loader to  handle loading more files after initial file is done. Also modify to accept a list of files.   Also modify to check each filename to see if it is a directory and if it is then process all xlsx files in that directory.

* Optimizer to deliver  minimized total # hours. 

* Fix Precursor automatically when Precursor has a reference to it's self to ignore the redundant

* Fix Precursor to handle properly when there is a recursive reference bath to it's self especially when separated by several layers of children.

* Double check not workable is sent when we run out of max days and have not been able to mark a project as done.

* 



## *Requirements*

- Produce a hierarchy indented view based on named parent projects.

- Allow parent projects to suppress display of detailed projects in the constructed roadmap.  Allow this detail to be changed.

- Show projects that will miss a Require by date in a separate color. 

- 

- Produce program report showing prioritized projects grouped by program rather than strict execution view.

- Report that shows programs or parent projects that violate target date due to resource constraints of child projects.

- Report showing adjusted priority based on precursor projects

  



## Extensions

- Write an optimizer that changes priority levels and attempts to deliver greater business value. Report on what it finds.  A proxy for business value until we develop one is maximum projects.   Also write optimizer that changes priorities as little as possible to allow projects with target dates to be delivered on or before that date.   
- Write an optimizer that incrementally changes resource group sizes to maximize business value delivered and / or maximize project counts delivered before due date.   A proxy for business value can be project count or could be to minimize idle resources.
- Support a notion that  assigned time from all resource groups allocated per day so the last hour for the smallest resource demand is delivered at the same time as the most demanded resource.  EG: If you need 10 hours of testing and 100 hours of SDE then the testers should only be assigned at a max capacity of 0.1 compared to the testers. Otherwise the testers would finish all their work too soon which we all know would not work since they need to test as code is finished.
- Allow a group of resources to be managed as a team. Allocated as a team even if all the resources in that team do not exactly fit the resource profile needed for the project.
- Modify the notion of resource groups to allow an alternative approach where every resource is named and includes with a list of skills.  Then allow the system to skill match.  This should allow more aggressive load balancing.  For example I may be the CTO but I also have skills as a Solution Architect,  Enterprise Architect.   I can generally support 4 hours per week as a Enterprise architect so this should be available for the portfolio planner to use as  a resource to match that skill.    Otherwise I would not generally be considered in the group of enterprise architects.



# Features & Work Completed

- DONE:JOE:2019-07-28: Produce report showing cost by project by resource by time unit
- DONE:JOE-2019-07-28: Add support for minimum duration to stretch projects out over time by reducing labor consumed.
- DONE:2019-07-20: Allow entry of a directory for input files rather than single file.   Process all input task files as if they were a single file which requires a sort and merge process to produce a single list from the merged priorities.
- DONE:JOE: 2019-07-20: Report showing hours consumed by month by resource group once the portfolio analyzer has determined what can be executed when.
- DONE:JOE:2019-07-20: Consume resources to resource projects by priority order.  Working and complete maximum projects as soon as they can reasonably be worked with available resources.
- DONE:JOE:2019-07-20: Produce roadmaps automatically based on project priority and available resources. 
- DONE:JOE:2019-07-28: Write a Tutorial on how to use approached by a person just getting started. Include the advanced concepts like how to use percent complete to minimize future planning overhead.   Include how to configure but recommend they use the default templates adding rows until they have reached some success.
- DONE:JOE:2017-07-27: Modify basic Demo config to demonstrate the % complete functionality.
- DONE:JOE:2019-07-28: Separate the Roadmap,  Cost and Capacity tabs into separate files.
- DONE:JOE:2019-07-29: Create a Cache Manager that writes files by a grouping string into different files.   When a file for a segment is needed the first time it is either loaded or created and then cached to make adding more tabs easy. 
- DONE:JOE:2019-07-28: Double Check resource allocation to ensure we actually delivered the resources requested.
- DONE:JOE:2019-07-28: We capture the most Demanding Resource and a computed minimum duration for it during workProject we should generate a report showing those resources that are are controlling completion times. 
- DONE:JOE:2019-07-27: Modify default allocation model to be Even so unless one of the alternate resource models is chosen the system figures out which resource will take the longest to satisfy at the max resource level and then spread the other resources out to consume some resources every day over that time.  Need to figure out a better way of supporting one resource such as an architect who could be spread across 10 projects but would only allocate 10% per day per project.  The Unit count assigned would be 0.1 architect or if they have 6 hours per day each project would get 0.6 hours of time. Compute a Ratio of hours deliver per hour delivered of the most demanding resources.  EG:  If I need 10 hours of Dev time and 1 hour of IA time then I for every 1 hour of Dev time I need to deliver 10/1 = 0.1 hours of IA time.   Use this ratio for the smaller demand resources to meter delivery of the other resources.
- DONE:JOE:2019-07-27 Support Portion Complete to only work the portion of the project remaining.  Otherwise people would have to recompute remaining hours every run when the system can do it for them. 
- DONE:JOE:2019-07-27 Support for minimum duration to force the system to consume resources at a rate slower than max available.
- DONE:JOE:2019-07-27: Consolidate the two main test drivers to a common driver that is used with two different sets of files or a single larger module called by a smaller test module.
- DONE:JOE:2019-07-25: Project Resource needs loaded from Config is an array. At current time project resources is getting created as the same length as project resources. But when a resource need is set it leaves a project resource item that is nil which is not OK because we use simple numeric indexing to allow faster lookup and iteration.
- DONE:JOE:2019-07-26: Produce report showing costs by project by time unit
- DONE:JOE:2019-07-26: Produce report showing cost by resource by time unit
- DONE:JOE:2019-07-27: Add support for column for resource group to support multiple columns and sum the values read from project list rather than a single column read.
- DONE:JOE:2019-07-27: Add support for Percentage  Completed so the system can mathematically deduct that from the total resources and only plan for the remaining.  This would allow it to adjust the roadmap as things evolve over time.  It would also be used in conjunction with DateStarted to force a accurate graphical view.
- NOTDO:JOE:2019-07-27:  With the percent complete functionality we already have the ability to deduct completed work from total units and other tools are used to track actual hours consumed by project so no need to provide that function.    Add support for Date Completed so the system can skip planning those that have already been finished so the system can show them in the past roadmap while retaining the original estimates.  They would only render graphically other than acting as a filter and not planning them in consumption.
- DONE:JOE:2019-07-20: Modify resource file to vary the resource count by month to test staff ramp up and ramp down.
- DONE:JOE:2019-07-23:Modify ProjectList Loader to  handle loading more files after initial file is done. Also modify to accept a list of files.   Also modify to check each filename to see if it is a directory and if it is then process all xlsx files in that directory.
- DONE:JOE:2019-07-22: Add Color Background to cells with X for active during week to  ProjRoadMap of output file
- DONE:JOE:2019-07-22: Add a blank row above projects that have no direct resources but do have children in ProjRoadMap of output file
- DONE:JOE:2019-07-22:Projections of roadmap showing 2 char wide cells with label on left and X in filled in for Weeks active on that project.  Make an extension of this that shows X or delivered for days when each resource type delivers time.
- DONE:JOE:2019-07-22:Fully Implement the minimum resources assigned per group with proper logic to require minResource to get into assignment but use upto max resource when available.   Override the max, min in resourceGroup when supplied in the project.
- DONE:JOE:2019-07-22:Change All references to Hours to Units to allow something other than Humans to be managed. 
- DONE:JOE:2019-07-21: Produce basic Excel output showing basic generate Roadmap.
- DONE:JOE:2019-07-22: Forcing project to start after last precursor is complete needs to be completed. 
- DONE:JOE:2019-07-22: Fully support Precursor projects blocking execution of project that depends on them.
- DONE:JOE:2019-07-22: Output Order should be so that a project with children is  moved to a output priority slightly lower than any of it's children.  Start with prima
- DONE:JOE:2019-07-22: Find a way to output groups before the children they contain.   Needs to support nested containment.
- DONE:JOE:2019-07-22: Add Sort function for projects based on EffectivePriority
- 
- 