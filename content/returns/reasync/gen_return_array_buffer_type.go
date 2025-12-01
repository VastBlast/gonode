package reasync

import (
	"github.com/VastBlast/gonode/tools"
	"strings"
)

func GenAsyncReturnArrayBufferTypeCode(resultStructName string) string {
	_ = resultStructName
	return tools.FormatCodeIndentLn(`const void* wg_res_ = NULL;
  if (!wg_is_error && wg_async_res != NULL) {
    wg_res_ = wg_async_res->data;
  }
  if (wg_async_res != NULL && wg_async_res->err != NULL) {
    free((void*)wg_async_res->err);
  }
  if (wg_async_res != NULL) {
    free(wg_async_res);
  }`, 2)
}

func GenAsyncCallReturnArrayBufferTypeCode(methodName string, argNames []string, cleanupLabel string, resultStructName string) string {
	code := `
  // -------- genHandlerCode
  const void* wg_res_buf_ = ` + methodName + `(` + strings.Join(argNames, ",") + `);
  ` + resultStructName + `* wg_async_res_success = (` + resultStructName + `*)malloc(sizeof(*wg_async_res_success));
  if (wg_async_res_success == NULL) {
    wg_send_async_error("alloc async array buffer result wrapper");
    goto ` + cleanupLabel + `;
  }
  wg_async_res_success->is_error = false;
  wg_async_res_success->data = (void*)wg_res_buf_;
  wg_async_res_success->err = NULL;
  wg_res_ = (void*)wg_async_res_success;`
	return code
}

func GenAsyncCallbackArgArrayBufferTypeCode() string {
	code := `napi_value wg_result = wg_env_scope.Null();
    if (!wg_is_error && wg_res_ != NULL) {
      char *wg_ab_ = (char*) wg_res_;
      size_t wg_ab_length_ = strlen(wg_ab_);
      ArrayBuffer wg_arr_buffer_ = ArrayBuffer::New(wg_env_scope, wg_ab_, wg_ab_length_);
      wg_result = wg_arr_buffer_;
    }`
	return code
}
