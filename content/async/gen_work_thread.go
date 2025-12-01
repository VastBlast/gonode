package async

import (
	"fmt"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/args/argasync"
)

// Generate thread entry code
func genWorkThreadCode(
	export config.Export,
	workName string,
	workCompleteName string,
	executeworkName string,
	jsCallbackName string,
	structDataName string,
) string {
	argc := len(export.Args)
	cbArgIndex := 0

	for i, arg := range export.Args {
		if arg.Type == "callback" {
			cbArgIndex = i
			break
		}
	}

	inputArgCode := ""

	for index, arg := range export.Args {
		cCode, _ := argasync.GenParseInputArgCode(arg, index)
		inputArgCode += cCode
	}

	code := `
// ---------- genworkThreadCode
static napi_value ` + workName + `(napi_env wg_env, napi_callback_info wg_info) {
  const size_t wg_expected_argc = ` + fmt.Sprintf("%d", argc) + `;
  size_t wg_argc = wg_expected_argc;
  size_t wg_cb_arg_index = ` + fmt.Sprintf("%d", cbArgIndex) + `;
  napi_value wg_args[` + fmt.Sprintf("%d", argc) + `] = {0};
  napi_value wg_work_name;
  napi_status wg_sts;
  ` + structDataName + `* wg_addon = (` + structDataName + `*)malloc(sizeof(*wg_addon));
  if (wg_addon == NULL) {
    napi_throw_error(wg_env, NULL, "alloc async data");
    return NULL;
  }
  wg_addon->work = NULL;
  wg_addon->argc = wg_expected_argc;
  for (size_t i = 0; i < wg_expected_argc; i++) {
    wg_addon->args[i] = NULL;
  }
  napi_value wg_undefined;
  wg_catch_err(wg_env, napi_get_undefined(wg_env, &wg_undefined));
  auto wg_cleanup = [&]() {
    for (size_t i = 0; i < wg_expected_argc; i++) {
      if (wg_addon->args[i] != NULL && wg_addon->args[i]->type == 1) {
        WgAddonArgInfo* info = (WgAddonArgInfo*)wg_addon->args[i];
        delete [] (char *)info->value;
      }
      if (wg_addon->args[i] != NULL) {
        free(wg_addon->args[i]);
        wg_addon->args[i] = NULL;
      }
    }
    free(wg_addon);
  };
#ifdef NAPI_CPP_EXCEPTIONS
  try {
#endif
  wg_sts = napi_get_undefined(wg_env, &wg_undefined);
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    wg_cleanup();
    return NULL;
  }
  wg_sts = napi_get_cb_info(wg_env, wg_info, &wg_argc, wg_args, NULL, NULL);
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    wg_cleanup();
    return NULL;
  }
  for (size_t i = wg_argc; i < wg_expected_argc; i++) {
    wg_args[i] = wg_undefined;
  }
  napi_value wg_js_cb = wg_args[wg_cb_arg_index];` + inputArgCode + `
  assert(wg_addon->work == NULL && "Only one work item must exist at a time");
  wg_sts = napi_create_string_utf8(wg_env, "N-API Thread-safe Call from Async Work Item", NAPI_AUTO_LENGTH, &wg_work_name);
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    wg_cleanup();
    return NULL;
  }
  wg_sts = napi_create_threadsafe_function(wg_env, wg_js_cb, NULL, wg_work_name, 0, 1, NULL, NULL, NULL, ` + jsCallbackName + `, &(wg_addon->tsfn));
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    wg_cleanup();
    return NULL;
  }
  wg_sts = napi_create_async_work(wg_env, NULL, wg_work_name, ` + executeworkName + `, ` + workCompleteName + `, wg_addon, &(wg_addon->work));
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    napi_release_threadsafe_function(wg_addon->tsfn, napi_tsfn_abort);
    wg_cleanup();
    return NULL;
  }
  wg_sts = napi_queue_async_work(wg_env, wg_addon->work);
  if (wg_sts != napi_ok) {
    wg_catch_err(wg_env, wg_sts);
    napi_release_threadsafe_function(wg_addon->tsfn, napi_tsfn_abort);
    wg_catch_err(wg_env, napi_delete_async_work(wg_env, wg_addon->work));
    wg_cleanup();
    return NULL;
  }
  return NULL;
#ifdef NAPI_CPP_EXCEPTIONS
  } catch (const Error& wg_ex) {
    wg_ex.ThrowAsJavaScriptException();
  } catch (const std::exception& wg_ex) {
    napi_throw_error(wg_env, NULL, wg_ex.what());
  } catch (...) {
    napi_throw_error(wg_env, NULL, "native exception");
  }
  wg_cleanup();
  return NULL;
#endif
}`

	return code
}
