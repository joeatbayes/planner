package main

import "fmt"
import "regexp"
import "strings"

var ParmMatch, ParmErr = regexp.Compile(`\.=\w+=[\.\w-]+`)

func parseSpecialFields(str string) map[string]string {
	ms := ParmMatch.FindAllIndex([]byte(str), -1)
	tout := make(map[string]string)
	for _, m := range ms {
		subStr := str[m[0]+2 : m[1]]
		tarr := strings.SplitN(subStr, "=", 2)
		if len(tarr) > 1 {
			tout[tarr[0]] = tarr[1]
		} else {
			tout[subStr] = ""
		}
	}
	return tout
}

func main() {
	str := `
	   i ama test
	   .=minDuration=99 and I am a character
	   test1 .=minhack=jimbo test .=floatTest=1918.31  .=negFloatTest=-1818.21
	.=jackz=   timbber
	`
	tout := parseSpecialFields(str)
	fmt.Println("Generated MAP")
	for key, val := range tout {
		fmt.Println("\tkey=", key, "\tval=", val)
	}
}
