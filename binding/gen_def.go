package binding

import (
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
)

func GenDefFile(cfgs config.Config, filename string) bool {
	code := tools.FormatCodeIndent("EXPORTS", 0)

	for _, export := range cfgs.Exports {
		code += tools.FormatCodeIndentLn(export.Name, 2)
	}

	if e := tools.WriteFile(code, tools.FormatDirPath(cfgs.OutPut), filename); e != nil {
		return false
	}
	return true
}
