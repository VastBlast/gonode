package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnArrayBufferTypeCode() string {
	return tools.FormatCodeIndentLn(`const void* wg_res_ = (void*)wg_data;`, 2)
}

func GenAsyncCallReturnArrayBufferTypeCode(methodName string, argNames []string) string {
	code := `
  // -------- genHandlerCode
  const void* wg_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);`
	return code
}

func GenAsyncCallbackArgArrayBufferTypeCode() string {
	code := `char *wg_ab_ = (char*) wg_res_;
    size_t wg_ab_length_ = strlen(wg_ab_);
    ArrayBuffer wg_arr_buffer_ = ArrayBuffer::New(wg_env, wg_ab_, wg_ab_length_);
    napi_value wg_argv[] = { wg_arr_buffer_ };`
	return code
}
