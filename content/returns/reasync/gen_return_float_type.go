package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnFloatTypeCode(resultStructName string) string {
	_ = resultStructName
	code := tools.FormatCodeIndentLn(`float wg_res_ = 0.0;
  if (!wg_is_error && wg_async_res != NULL) {
    float* wg_res_ptr_ = (float*)wg_async_res->data;
    wg_res_ = wg_res_ptr_ ? *wg_res_ptr_ : 0.0;
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

	return code
}

func GenAsyncCallReturnFloatTypeCode(methodName string, argNames []string, cleanupLabel string, resultStructName string) string {
	code := `
  // -------- genHandlerCode
  ` + resultStructName + `* wg_async_res_success = NULL;
  const float wg_tmp_res_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  float* wg_res_ptr_ = (float*)malloc(sizeof(float));
  if (wg_res_ptr_ == NULL) {
    wg_send_async_error("alloc async result");
    goto ` + cleanupLabel + `;
  }
  *wg_res_ptr_ = wg_tmp_res_;
  wg_async_res_success = (` + resultStructName + `*)malloc(sizeof(*wg_async_res_success));
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

func GenAsyncCallbackArgFloatTypeCode() string {
	return `napi_value wg_result = wg_env_scope.Null();
    if (!wg_is_error) {
      Number wg_float_ = Number::New(wg_env_scope, wg_res_);
      wg_result = wg_float_;
    }`
}
