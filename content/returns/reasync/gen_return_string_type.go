package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnStringTypeCode() (string, string) {
	code := tools.FormatCodeIndentLn(`const char* wg_raw_res_ = static_cast<char*>(wg_data);
  string wg__res_ = wg_raw_res_ ? wg_raw_res_ : "";
  if (wg_raw_res_ != NULL) {
    free((void*)wg_raw_res_);
  }`, 2)

	pCode := tools.FormatCodeIndentLn(`napi_value wg_res_ = String::New(wg_env, wg__res_);`, 4)
	//endCode := tools.FormatCodeIndentLn(`delete [] _res_`, 0)
	return code, pCode
}

func GenAsyncCallReturnStringTypeCode(methodName string, argNames []string) string {
	code := `
  // -------- genHandlerCode
  const void* wg_res_src_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  const char* wg_src_str_ = static_cast<const char*>(wg_res_src_);
  size_t wg_src_len_ = wg_src_str_ ? strlen(wg_src_str_) : 0;
  char* wg_res_ = (char*)malloc(wg_src_len_ + 1);
  if (wg_res_ == NULL) {
    wg_catch_err_bg(napi_generic_failure, "alloc async result");
    return;
  }
  if (wg_src_len_ > 0) {
    memcpy(wg_res_, wg_src_str_, wg_src_len_);
  }
  wg_res_[wg_src_len_] = '\0';`
	return code
}

func GenAsyncCallbackArgStringTypeCode() string {
	//return `napi_value wg_string_ = String::New(wg_env, wg_res_);
	//napi_value wg_argv[] = { wg_string_ };`
	return `napi_value wg_argv[] = { wg_res_ };`
}
