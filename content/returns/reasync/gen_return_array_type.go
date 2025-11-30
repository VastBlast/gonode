package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnArrayTypeCode() string {
	return tools.FormatCodeIndentLn(`const char* wg_raw_res_ = static_cast<char*>(wg_data);
  string wg_res_ = wg_raw_res_ ? wg_raw_res_ : "";`, 2)
}

func GenAsyncCallReturnArrayTypeCode(methodName string, argNames []string) string {
	code := `
  // -------- genHandlerCode
  const char* wg_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);`
	return code
}

func GenAsyncCallbackArgArrayTypeCode() string {
	code := `Array wg_arr_ = wg_string_to_array(wg_res_, wg_env);
    napi_value wg_argv[] = { wg_arr_ };`
	return code
}
