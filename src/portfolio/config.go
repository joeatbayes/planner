package portfolio

import (
	"encoding/json"
	"fmt"

	//"strconv"
	//"strings"

	//"github.com/360EntSecGroup-Skylar/excelize"
)

type Config struct {
	FiName string
	Res    *ResourceConfig
	Proj   *ProjectConfig
}

func (cfg *Config) ToJSON(pretty bool) string {
	if pretty {
		slcB, _ := json.MarshalIndent(cfg, "", "  ")
		return string(slcB)
	} else {
		slcB, _ := json.Marshal(cfg)
		return string(slcB)

	}
}

func LoadConfig(fiName string) *Config {
	if IsFile(fiName) == false {
		panic("L32 FATAL ERROR: Config File " + fiName + " does not exist")
	}
	cfg := new(Config)
	cfg.FiName = fiName
	cfg.Res = LoadResourceConfig(fiName)
	cfg.Proj = LoadProjectConfig(fiName)
	fmt.Println("L38: LoadConfig Sucess()")
	return cfg

}
