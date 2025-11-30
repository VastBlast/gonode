package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnObjectTypeCode() string {
	return tools.FormatCodeIndentLn(`const char* wg_raw_res_ = static_cast<char*>(wg_data);
  string wg_res_ = wg_raw_res_ ? wg_raw_res_ : "";`, 2)
}

func GenAsyncCallReturnObjectTypeCode(methodName string, argNames []string) string {
	code := `
  // -------- genHandlerCode
  const char* wg_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);`
	return code
}

func GenAsyncCallbackArgObjectTypeCode() string {
	code := `Object wg_obj = wg_string_to_object(wg_res_, wg_env);
    napi_value wg_argv[] = { wg_obj };`
	return code
}
