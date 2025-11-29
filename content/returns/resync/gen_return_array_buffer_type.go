package resync

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/args/argsync"
	"github.com/VastBlast/gonode/tools"
	"strings"
)

// Generate handler body
func GenHandleReturnArrayBufferCode(method string, args []string, preCode string) string {
	// Convert to array buffer
	code := `
  void * wg_res_ = ` + method + `(` + strings.Join(args, ",") + `);
  char *wg_ab_ = (char*) wg_res_;
  size_t wg_ab_length_ = strlen(wg_ab_);
  ArrayBuffer wg_arr_buffer_ = ArrayBuffer::New(wg_env, wg_ab_, wg_ab_length_);`

	code += preCode
	code += tools.FormatCodeIndentLn(`return wg_arr_buffer_;`, 2)
	return code
}

// Generate return code for array buffer type
func GenReturnArrayBufferTypeCode(export config.Export) string {
	methodName := export.Name
	args := export.Args

	code := `
// ---------- GenCode ---------- 
Value _` + methodName + `(const CallbackInfo& wg_info) {`
	code += tools.FormatCodeIndentLn(`Env wg_env = wg_info.Env();`, 2)

	c, argNames, endCode := argsync.GenArgCode(args)
	code += c

	code += GenHandleReturnArrayBufferCode(methodName, argNames, endCode)

	code += tools.FormatCodeIndentLn(`}`, 0)
	return code
}
