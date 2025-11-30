package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnArrayTypeCode() string {
	return tools.FormatCodeIndentLn(`const char* wg_raw_res_ = static_cast<char*>(wg_data);
  string wg_res_ = wg_raw_res_ ? wg_raw_res_ : "";
  if (wg_raw_res_ != NULL) {
    free((void*)wg_raw_res_);
  }`, 2)
}

func GenAsyncCallReturnArrayTypeCode(methodName string, argNames []string) string {
	code := `
  // -------- genHandlerCode
  const char* wg_src_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  size_t wg_src_len_ = wg_src_res_ ? strlen(wg_src_res_) : 0;
  char* wg_res_ = (char*)malloc(wg_src_len_ + 1);
  if (wg_res_ == NULL) {
    wg_catch_err_bg(napi_generic_failure, "alloc async array result");
    return;
  }
  if (wg_src_len_ > 0) {
    memcpy(wg_res_, wg_src_res_, wg_src_len_);
  }
  wg_res_[wg_src_len_] = '\0';`
	return code
}

func GenAsyncCallbackArgArrayTypeCode() string {
	code := `Array wg_arr_ = wg_string_to_array(wg_res_, wg_env);
    napi_value wg_argv[] = { wg_arr_ };`
	return code
}
