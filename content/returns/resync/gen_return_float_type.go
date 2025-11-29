package resync

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/args/argsync"
	"github.com/VastBlast/gonode/tools"
	"strings"
)

// Generate handler body
func GenHandleReturnFloatCode(method string, args []string, endCode string) string {
	code := tools.FormatCodeIndentLn(`float wg_res_ = `+method+`(`+strings.Join(args, ",")+`);`, 2)
	code += endCode
	code += tools.FormatCodeIndentLn(`return Number::New(wg_env, wg_res_);`, 2)
	return code
}

// Generate return code for float type
func GenReturnFloatTypeCode(export config.Export) string {
	methodName := export.Name
	args := export.Args

	code := `
// ---------- GenCode ---------- 
Value _` + methodName + `(const CallbackInfo& wg_info) {`
	code += tools.FormatCodeIndentLn(`Env wg_env = wg_info.Env();`, 2)

	c, argNames, endCode := argsync.GenArgCode(args)
	code += c

	code += GenHandleReturnFloatCode(methodName, argNames, endCode)

	code += tools.FormatCodeIndentLn(`}`, 0)
	return code
}
