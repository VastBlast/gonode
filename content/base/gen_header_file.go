package base

// Generate header file content
func GenHeaderFileCode(headerFile string) string {
	var code = `#define NAPI_EXPERIMENTAL
#include <napi.h>
#include <string>
#include <assert.h>
#include <functional>

#include "` + headerFile + `"

using namespace Napi;
using namespace std;
`
	return code
}
