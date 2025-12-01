package binding

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func GenDefFile(cfgs config.Config, filename string) bool {
	code := tools.FormatCodeIndent("EXPORTS", 0)

	hasFree := false
	for _, export := range cfgs.Exports {
		code += tools.FormatCodeIndentLn(export.Name, 2)
		if export.Name == "FreeCString" {
			hasFree = true
		}
	}
	if !hasFree {
		code += tools.FormatCodeIndentLn("FreeCString", 2)
	}

	if e := tools.WriteFile(code, tools.FormatDirPath(cfgs.OutPut), filename); e != nil {
		return false
	}
	return true
}
