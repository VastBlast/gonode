package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnBooleanTypeCode(resultStructName string) string {
	_ = resultStructName
	return tools.FormatCodeIndentLn(`bool wg_res_ = false;
  if (!wg_is_error && wg_async_res != NULL) {
    bool* wg_res_ptr_ = (bool*)wg_async_res->data;
    wg_res_ = wg_res_ptr_ ? *wg_res_ptr_ : false;
    if (wg_res_ptr_ != NULL) {
      free(wg_res_ptr_);
    }
  }
  if (wg_async_res != NULL && wg_async_res->err != NULL) {
    free((void*)wg_async_res->err);
  }
  if (wg_async_res != NULL) {
    free(wg_async_res);
  }`, 2)
}

func GenAsyncCallReturnBooleanTypeCode(methodName string, argNames []string, cleanupLabel string, resultStructName string) string {
	code := `
  // -------- genHandlerCode
  const bool wg_tmp_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  bool* wg_res_ptr_ = (bool*)malloc(sizeof(bool));
  if (wg_res_ptr_ == NULL) {
    wg_send_async_error("alloc async result");
    goto ` + cleanupLabel + `;
  }
  *wg_res_ptr_ = wg_tmp_res_;
  ` + resultStructName + `* wg_async_res_success = (` + resultStructName + `*)malloc(sizeof(*wg_async_res_success));
  if (wg_async_res_success == NULL) {
    free(wg_res_ptr_);
    wg_send_async_error("alloc async result wrapper");
    goto ` + cleanupLabel + `;
  }
  wg_async_res_success->is_error = false;
  wg_async_res_success->data = (void*)wg_res_ptr_;
  wg_async_res_success->err = NULL;
  wg_res_ = (void*)wg_async_res_success;`

	return code
}

func GenAsyncCallbackArgBooleanTypeCode() string {
	return `napi_value wg_result = wg_env_scope.Null();
    if (!wg_is_error) {
      Boolean wg_bool_ = Boolean::New(wg_env_scope, wg_res_);
      wg_result = wg_bool_;
    }`
}
