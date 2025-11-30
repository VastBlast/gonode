package base

// Generate header file content
func GenHeaderFileCode(headerFile string) string {
	var code = `#ifndef NAPI_EXPERIMENTAL
#define NAPI_EXPERIMENTAL
#endif
#include <napi.h>
#include <string>
#include <cstring>
#include <cstdlib>
#include <assert.h>
#include <functional>

#include "` + headerFile + `"

using namespace Napi;
using namespace std;
`
	return code
}
