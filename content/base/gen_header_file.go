package base

// Generate header file content
func GenHeaderFileCode(headerFile string) string {
	var code = `#include <napi.h>
#include <string>
#include <assert.h>
#include <functional>
`
	code += `
#include "` + headerFile + `"
`

	code += `
using namespace Napi;
using namespace std;

#define NAPI_EXPERIMENTAL
`
	return code
}
