package check

import (
	"fmt"
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func CheckBaseConfig(config config.Config) error {
	//defer func() {
	//	if err := recover(); err != nil {
	//		println(err.(string))
	//	}
	//}()

	clog.Info("Start checking the configured ...")

	// Check whether the parameter lists conflict
	exports := config.Exports

	var goApiList = make([]string, 0)
	var JsCallApiList = make([]string, 0)
	// "arraybuffer",
	var allowArgsList = []string{"int", "int32", "int64", "uint32", "float", "double", "boolean", "string", "array", "object", "callback"}
	// "arraybuffer",
	var returnAllowArgsList = []string{"int", "int32", "int64", "uint32", "float", "double", "boolean", "string", "array", "object"}

	for _, export := range exports {
		// Check for duplicate Go export names or JS call names
		if tools.InSlice(goApiList, export.Name) {
			return fmt.Errorf("The export Name \"%s\" already exists", export.Name)
		}
		goApiList = append(goApiList, export.Name)

		if tools.InSlice(JsCallApiList, export.JsCallName) {
			return fmt.Errorf("The export JsCallName \"%s\" already exists", export.JsCallName)
		}
		JsCallApiList = append(JsCallApiList, export.JsCallName)

		// Ensure parameter types are valid and names are not duplicated
		var curArgsList = make([]string, 0)
		for _, arg := range export.Args {
			if !tools.InSlice(allowArgsList, arg.Type) {
				return fmt.Errorf("The arguments type \"%s\" of the exported [%s] is not supported", arg.Type, export.Name)
			}

			if tools.InSlice(curArgsList, arg.Name) {
				return fmt.Errorf("The export parameter \"%s\" of [%s] already exists", arg.Name, export.Name)
			}
			curArgsList = append(curArgsList, arg.Name)
		}

		// Validate return types (callback not allowed)
		if !tools.InSlice(returnAllowArgsList, export.ReturnType) {
			return fmt.Errorf("The return type \"%s\" of the exported [%s] is not supported", export.ReturnType, export.Name)
		}
	}
	clog.Success("Checking the configured done ~")
	return nil
}

func CheckAsyncCorrectnessConfig(config config.Config) error {
	for _, export := range config.Exports {
		if export.JsCallMode == "async" {
			checked := false
			for _, arg := range export.Args {
				if arg.Type == "callback" {
					checked = true
					break
				}
			}
			if !checked {
				return fmt.Errorf("The export [%s] is missing the \"callback\" parameter", export.Name)
			}
		}
	}
	return nil
}

func CheckExportApiWithSourceFile(config config.Config) bool {

	// @todo Compare Golang source files with json configuration for consistency
	// @todo Avoid compilation failures caused by mismatched parameters or types; to be completed later ...
	//clog.Info("Start checking the API of Golang export ...")
	//clog.Success("Checking the API of Golang export done ~")

	return true
}
