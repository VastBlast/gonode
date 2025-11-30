package async

import (
	"github.com/VastBlast/gonode/config"
)

// Generate async work execution code
func genExecuteWorkCode(export config.Export, executeWorkName string, structDataName string) string {
	handlerCode, endCode := genHandlerCode(export)

	code := `
// -------- genExecuteworkCode
static void ` + executeWorkName + `(napi_env wg_env, void* wg_data) {
  ` + structDataName + `* wg_addon = (` + structDataName + `*)wg_data;
  napi_status wg_sts = napi_acquire_threadsafe_function(wg_addon->tsfn);
  wg_catch_err_bg(wg_sts, "acquire threadsafe function");` + handlerCode + `
  wg_sts = napi_call_threadsafe_function(wg_addon->tsfn, (void*)(wg_res_), napi_tsfn_blocking);
  wg_catch_err_bg(wg_sts, "call threadsafe function");
  wg_sts = napi_release_threadsafe_function(wg_addon->tsfn, napi_tsfn_release);
  wg_catch_err_bg(wg_sts, "release threadsafe function");` + endCode + `
}`

	return code
}
