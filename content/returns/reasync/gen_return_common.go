package reasync

import (
	"fmt"
	"github.com/VastBlast/gonode/tools"
)

func GenAsyncReturnCommonCode(resultStructName string) string {
	return tools.FormatCodeIndentLn(fmt.Sprintf(`%s* wg_async_res = (%s*)wg_data;
  bool wg_is_error = wg_async_res != NULL && wg_async_res->is_error;
  const char* wg_err_msg = (wg_async_res != NULL) ? wg_async_res->err : NULL;
  std::string wg_err_str = wg_err_msg != NULL ? std::string(wg_err_msg) : std::string();`, resultStructName, resultStructName), 2)
}

func GenAsyncFreeResultWrapperCode(returnType string, resultStructName string) string {
	freeData := ""
	if returnType == "string" || returnType == "array" || returnType == "object" {
		freeData = `
      if (wg_async_res_free->data != NULL) {
        free(wg_async_res_free->data);
      }`
	} else if returnType == "boolean" || returnType == "int" || returnType == "int32" || returnType == "int64" || returnType == "uint32" || returnType == "float" || returnType == "double" {
		freeData = `
      if (wg_async_res_free->data != NULL) {
        free(wg_async_res_free->data);
      }`
	} else {
		freeData = ""
	}

	return tools.FormatCodeIndent(fmt.Sprintf(`%s* wg_async_res_free = (%s*)wg_res_;
      if (wg_async_res_free != NULL) {%s
        if (wg_async_res_free->err != NULL) {
          free(wg_async_res_free->err);
        }
        free(wg_async_res_free);
      }`, resultStructName, resultStructName, freeData), 6)
}
