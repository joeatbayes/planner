package portfolio

import (
	//"encoding/json"
	"fmt"
	"os"
	//	"strconv"
	//"strings"
	"path/filepath"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Sheet Cache Manager Takes a input main file names manages a set of files
// with modified extensions.  Can also manage an arbitrary set of files.
// Also provides a cache mechanism so if a file is already opened
// as a excelize.File then it will be returned without reopening.
// The caching is intended to help avoid the need to reparse the same
// file multiple times. It is also useful to allow easy adding of
// new tabs or sheets that are all written to disk during a single
// flush operation.
type ExcelFiManager struct {
	// main output excel file
	mainName string

	// Name of output excel file without the extension.
	// Saved here to make easy to compute derived names
	baseName string
	fiExt    string

	// We store the cached excelize files once opened
	// so various utilities can access them without
	// reopening and reloading them.  Very helpful
	// when generating reports onto many tabs of the
	// same file.
	files map[string]*excelize.File

	// Used by some Utilities to compute how many rows
	// we have conumed and provide the next available
	// counter.
	rowCnts map[string]int
}

func MakeExeclFiManager(mainFiName string) *ExcelFiManager {
	mgr := new(ExcelFiManager)
	mgr.mainName = mainFiName
	mgr.fiExt = filepath.Ext(mainFiName)
	mgr.baseName = mgr.mainName[0 : len(mgr.mainName)-len(mgr.fiExt)]
	mgr.files = make(map[string]*excelize.File)
	mgr.rowCnts = make(map[string]int)
	fmt.Println("L51:MakeExeclFiManager  mainFiName=", mainFiName, " baseName=", mgr.baseName, "ext=", mgr.fiExt)
	return mgr
}

func (mgr *ExcelFiManager) NameWithLastSeg(newSeg string) string {
	if newSeg <= " " {
		return mgr.mainName
	} else {
		return mgr.baseName + "." + newSeg + mgr.fiExt
	}
}

// Fetch the file with modified last segment inserted before the mainName segment.
// A new file will be created and replace the existing file when useExistFi is false.
// otherwise will open and read the existing file if it exists.   If the file is already
// opened the cached handle will be returned and the useExistFi will be ignored.
func (mgr *ExcelFiManager) GetWithLastSeg(newSeg string, useExistFi bool) *excelize.File {
	return mgr.Get(mgr.NameWithLastSeg(newSeg), useExistFi)
}

// Return the file specified in mainFiName of the ExcelFiManager
// A new file will be created and replace the existing file when useExistFi is false.
// otherwise will open and read the existing file if it exists.   If the file is already
// opened the cached handle will be returned and the useExistFi will be ignored.
func (mgr *ExcelFiManager) GetMainFi(useExistFi bool) *excelize.File {
	return mgr.Get(mgr.mainName, useExistFi)
}

// Returns the excelize.File handle for the fiName specified. If the file is
// already opened the cached handle will be returned.
// if the file is not open then  new file will be created and replace the existing file
// when useExistFi is false. Otherwise will open and read the existing file if it exists.
// If the file is already opened the cached handle will be returned and the useExistFi
// will be ignored.
func (mgr *ExcelFiManager) Get(fiName string, useExistFi bool) *excelize.File {
	file, found := mgr.files[fiName]
	if found {
		// File is open and file handle already cached.
		return file
	} else {
		if useExistFi == false || IsFile(fiName) == false {
			// We can just create an new file handle and
			// when it flushes it will overwrite the existing
			// file.
			f := excelize.NewFile()
			mgr.files[fiName] = f
			return f
		} else {
			// We know the user wants to read the existing file
			// and we know the existing file does exist so we
			// can try to open it.
			f, err := excelize.OpenFile(fiName)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			return f

		}
	}
}

func (mgr *ExcelFiManager) FlushFi(fiName string) *excelize.File {
	file, found := mgr.files[fiName]
	if found {
		err := file.SaveAs(fiName)
		if err != nil {
			fmt.Println("L128: Error saving file ", fiName, "err=", err)
			return nil
		}
	} else {
		return nil
	}
	return file
}

func (mgr *ExcelFiManager) FlushAll() {
	for name, _ := range mgr.files {
		mgr.FlushFi(name)
	}
}

func (mgr *ExcelFiManager) Evict(fiName string) bool {
	_, found := mgr.files[fiName]
	if found {
		mgr.FlushFi(fiName)
		delete(mgr.files, fiName)
		return true
	} else {
		return false
	}
}

func (mgr *ExcelFiManager) EvictAll() {
	for name, _ := range mgr.files {
		mgr.Evict(name)
	}
}

// since we may have many files with names derived from
// the base name we need to delete them all to ensure a clean
// run.
func (mgr *ExcelFiManager) RemoveDerivedFiles() {
	gp := mgr.baseName + "*"
	dirFiles, err := filepath.Glob(gp)
	if err != nil {
		fmt.Println("Error reading globpath ", gp, " err=", err)
		return
	}
	for _, dfname := range dirFiles {
		err := os.Remove(dfname)
		if err != nil {
			fmt.Println("L64: ERROR failed to delete file=", dfname, " err=", err)
		}
	}
}

// Find a counter for the sheetName and segment name.  Return the counter
// located or zero if it does not exists.
func (mgr *ExcelFiManager) GetCount(segName string, sheetName string) int {
	key := segName + sheetName
	val, found := mgr.rowCnts[key]
	if found {
		return val
	} else {
		return 0
	}
}

// Find a counter for the sheetName and segment name.  Increment the value
// and return it if it exists.  If not exist return 0 and create the counter
// for next time.
func (mgr *ExcelFiManager) NextRow(segName string, sheetName string) int {
	key := segName + sheetName
	val, found := mgr.rowCnts[key]
	if found {
		val++
		mgr.rowCnts[key] = val
		return val
	} else {
		mgr.rowCnts[key] = 0
		return 0
	}
}

// Set value for the counter located at sheetname, segname
// to the specified value.  If the counter does not exist then
func (mgr *ExcelFiManager) SetRow(segName string, sheetName string, valIn int) {
	key := segName + sheetName
	mgr.rowCnts[key] = valIn
}

// Set the value for the column at the current row for the specified file to the value specified.
func (mgr *ExcelFiManager) SetCellVal(segName string, sheetName string, colName string, tval string) {
	f := mgr.GetWithLastSeg(segName, false)
	row := mgr.GetCount(segName, sheetName)
	axis := MakeAxis(colName, row)
	f.SetCellValue(sheetName, axis, tval)
}
