package resync

import (
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/content/args/argsync"
	"github.com/wenlng/gonacli/tools"
	"strings"
)

// Generate handler body
func GenHandleReturnArrayCode(method string, args []string, preCode string) string {
	code := tools.FormatCodeIndentLn(`string wg_res_ = `+method+`(`+strings.Join(args, ",")+`);`, 2)

	// Convert to array
	code += tools.FormatCodeIndentLn(`Array wg_arr_ = wg_string_to_array(wg_res_, wg_env);`, 2)

	code += preCode
	code += tools.FormatCodeIndentLn(`return wg_arr_;`, 2)
	return code
}

// Generate return code for array type
func GenReturnArrayTypeCode(export config.Export) string {
	methodName := export.Name
	args := export.Args

	code := `
// ---------- GenCode ---------- 
Value _` + methodName + `(const CallbackInfo& wg_info) {`
	code += tools.FormatCodeIndentLn(`Env wg_env = wg_info.Env();`, 2)

	c, argNames, endCode := argsync.GenArgCode(args)
	code += c

	code += GenHandleReturnArrayCode(methodName, argNames, endCode)

	code += tools.FormatCodeIndentLn(`}`, 0)
	return code
}
