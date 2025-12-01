package async

import (
	"github.com/VastBlast/gonode/config"
)

// Generate async work execution code
func genExecuteWorkCode(export config.Export, executeWorkName string, structDataName string, resultStructName string) string {
	cleanupLabel := executeWorkName + "_cleanup"
	handlerCode, endCode, sendFailureCleanup := genHandlerCode(export, cleanupLabel, resultStructName)

	code := `
// -------- genExecuteworkCode
static void ` + executeWorkName + `(napi_env wg_env, void* wg_data) {
  ` + structDataName + `* wg_addon = (` + structDataName + `*)wg_data;
  napi_status wg_sts = napi_acquire_threadsafe_function(wg_addon->tsfn);
  if (wg_sts != napi_ok) {
    wg_catch_err_bg(wg_sts, "acquire threadsafe function");
    return;
  }
  auto wg_send_async_error = [&](const char* wg_msg_) {
    ` + resultStructName + `* wg_async_res_err = (` + resultStructName + `*)malloc(sizeof(*wg_async_res_err));
    if (wg_async_res_err == NULL) {
      wg_catch_err_bg(napi_generic_failure, "alloc async error result");
      return;
    }
    wg_async_res_err->is_error = true;
    wg_async_res_err->data = NULL;
    wg_async_res_err->err = NULL;
    if (wg_msg_ != NULL) {
      size_t wg_err_len_ = strlen(wg_msg_);
      wg_async_res_err->err = (char*)malloc(wg_err_len_ + 1);
      if (wg_async_res_err->err != NULL) {
        memcpy(wg_async_res_err->err, wg_msg_, wg_err_len_);
        wg_async_res_err->err[wg_err_len_] = '\0';
      }
    }
    napi_status wg_sts_err = napi_call_threadsafe_function(wg_addon->tsfn, (void*)(wg_async_res_err), napi_tsfn_blocking);
    if (wg_sts_err != napi_ok) {
      if (wg_async_res_err->err != NULL) {
        free((void*)wg_async_res_err->err);
      }
      free(wg_async_res_err);
      wg_catch_err_bg(wg_sts_err, "call threadsafe function");
    }
  };
  void* wg_res_ = NULL;
#ifdef NAPI_CPP_EXCEPTIONS
  try {
#endif` + handlerCode + `
  if (wg_res_ != NULL) {
    wg_sts = napi_call_threadsafe_function(wg_addon->tsfn, (void*)(wg_res_), napi_tsfn_blocking);
    if (wg_sts != napi_ok) {
      wg_catch_err_bg(wg_sts, "call threadsafe function");` + sendFailureCleanup + `
      goto ` + cleanupLabel + `;
    }
  }
` + cleanupLabel + `:
  wg_sts = napi_release_threadsafe_function(wg_addon->tsfn, napi_tsfn_release);
  wg_catch_err_bg(wg_sts, "release threadsafe function");` + endCode + `
#ifdef NAPI_CPP_EXCEPTIONS
  } catch (const std::exception& wg_ex) {
    wg_send_async_error(wg_ex.what());
  } catch (...) {
    wg_send_async_error("native exception");
  }
#endif
}`

	return code
}
