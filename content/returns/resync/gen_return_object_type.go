package resync

import (
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/content/args/argsync"
	"github.com/VastBlast/gonode/tools"
	"strings"
)

// Generate handler body
func GenHandleReturnObjectCode(method string, args []string, preCode string) string {
	code := tools.FormatCodeIndentLn(`const char* wg_raw_res_ = `+method+`(`+strings.Join(args, ",")+`);`, 2)
	code += tools.FormatCodeIndentLn(`string wg_res_ = wg_raw_res_ ? wg_raw_res_ : "";`, 2)
	code += tools.FormatCodeIndentLn(`wg_free_cstring(wg_raw_res_);`, 2)

	code += tools.FormatCodeIndentLn(`Object wg_obj_ = wg_string_to_object(wg_res_, wg_env);`, 2)
	code += preCode
	code += tools.FormatCodeIndentLn(`return wg_obj_;`, 2)
	return code
}

// Generate return code for object type
func GenReturnObjectTypeCode(export config.Export) string {
	methodName := export.Name
	args := export.Args

	body, argNames, endCode := argsync.GenArgCode(args)
	body += GenHandleReturnObjectCode(methodName, argNames, endCode)

	code := `
// ---------- GenCode ---------- 
Value _` + methodName + `(const CallbackInfo& wg_info) {`
	code += tools.FormatCodeIndentLn(`Env wg_env = wg_info.Env();`, 2)
	code += tools.FormatCodeIndentLn(`#ifdef NAPI_CPP_EXCEPTIONS`, 0)
	code += tools.FormatCodeIndentLn(`  try {`, 0)
	code += tools.FormatCodeIndentLn(`#endif`, 0)
	code += body
	code += tools.FormatCodeIndentLn(`#ifdef NAPI_CPP_EXCEPTIONS`, 0)
	code += tools.FormatCodeIndentLn(`  } catch (const Error& wg_ex) {`, 0)
	code += tools.FormatCodeIndentLn(`    wg_ex.ThrowAsJavaScriptException();`, 0)
	code += tools.FormatCodeIndentLn(`    return wg_env.Null();`, 0)
	code += tools.FormatCodeIndentLn(`  } catch (const std::exception& wg_ex) {`, 0)
	code += tools.FormatCodeIndentLn(`    napi_throw_error(wg_env, NULL, wg_ex.what());`, 0)
	code += tools.FormatCodeIndentLn(`    return wg_env.Null();`, 0)
	code += tools.FormatCodeIndentLn(`  } catch (...) {`, 0)
	code += tools.FormatCodeIndentLn(`    napi_throw_error(wg_env, NULL, "native exception");`, 0)
	code += tools.FormatCodeIndentLn(`    return wg_env.Null();`, 0)
	code += tools.FormatCodeIndentLn(`  }`, 0)
	code += tools.FormatCodeIndentLn(`#endif`, 0)

	code += tools.FormatCodeIndentLn(`}`, 0)
	return code
}
