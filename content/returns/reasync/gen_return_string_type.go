package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnStringTypeCode(resultStructName string) (string, string) {
	_ = resultStructName
	code := tools.FormatCodeIndentLn(`string wg__res_ = "";
  if (!wg_is_error && wg_async_res != NULL) {
    const char* wg_raw_res_ = static_cast<const char*>(wg_async_res->data);
    wg__res_ = wg_raw_res_ ? wg_raw_res_ : "";
    if (wg_raw_res_ != NULL) {
      free((void*)wg_raw_res_);
    }
  }
  if (wg_async_res != NULL && wg_async_res->err != NULL) {
    free((void*)wg_async_res->err);
  }
  if (wg_async_res != NULL) {
    free(wg_async_res);
  }`, 2)

	pCode := ""
	return code, pCode
}

func GenAsyncCallReturnStringTypeCode(methodName string, argNames []string, cleanupLabel string, resultStructName string) string {
	code := `
  // -------- genHandlerCode
  ` + resultStructName + `* wg_async_res_success = NULL;
  const void* wg_res_src_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  const char* wg_src_str_ = static_cast<const char*>(wg_res_src_);
  size_t wg_src_len_ = wg_src_str_ ? strlen(wg_src_str_) : 0;
  char* wg_res_buf_ = (char*)malloc(wg_src_len_ + 1);
  if (wg_res_buf_ == NULL) {
    wg_free_cstring(wg_src_str_);
    wg_send_async_error("alloc async result");
    goto ` + cleanupLabel + `;
  }
  if (wg_src_len_ > 0) {
    memcpy(wg_res_buf_, wg_src_str_, wg_src_len_);
  }
  wg_res_buf_[wg_src_len_] = '\0';
  wg_free_cstring(wg_src_str_);
  wg_async_res_success = (` + resultStructName + `*)malloc(sizeof(*wg_async_res_success));
  if (wg_async_res_success == NULL) {
    free(wg_res_buf_);
    wg_send_async_error("alloc async result wrapper");
    goto ` + cleanupLabel + `;
  }
  wg_async_res_success->is_error = false;
  wg_async_res_success->data = (void*)wg_res_buf_;
  wg_async_res_success->err = NULL;
  wg_res_ = (void*)wg_async_res_success;`
	return code
}

func GenAsyncCallbackArgStringTypeCode() string {
	return `napi_value wg_result = wg_env_scope.Null();
    if (!wg_is_error) {
      wg_result = String::New(wg_env_scope, wg__res_);
    }`
}
