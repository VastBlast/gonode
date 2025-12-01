package buildtask

import (
	"fmt"
	"github.com/VastBlast/gonode/check"
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
	"github.com/jinzhu/configor"
	"path/filepath"
	"strings"
)

var cfgs config.Config

func parseConfig(config string) bool {
	configPath := "./goaddon.json"
	if len(config) > 0 {
		configPath = config
	}

	cpath := tools.FormatDirPath(configPath)

	if !tools.Exists(configPath) {
		clog.Error("The json configuration file of addon does not exist.")
		return false
	}

	path := filepath.Join(cpath)
	if err := configor.Load(&cfgs, path); err != nil {
		clog.Error(err)
		return false
	}

	return true
}

func parseAndCheck(config string) bool {
	if ok := parseConfig(config); !ok {
		return false
	}

	if c := checkConfigure(cfgs); !c {
		return false
	}

	return true
}

func normalizeArgs(args string) string {
	if len(args) == 0 {
		return args
	}

	if strings.HasPrefix(args, "'") || strings.HasPrefix(args, "\"") {
		args = args[1:]
		if strings.HasSuffix(args, "'") || strings.HasSuffix(args, "\"") {
			args = args[:len(args)-1]
		}
	}

	return args
}

func checkConfigure(c config.Config) bool {
	// Check configuration file
	if err := check.CheckBaseConfig(c); err != nil {
		clog.Error(err)
		return false
	}
	if err := check.CheckAsyncCorrectnessConfig(c); err != nil {
		clog.Error(err)
		return false
	}
	if c := check.CheckExportApiWithSourceFile(c); !c {
		return false
	}

	return true
}

// Build golang library files
// gonode build => go build -buildmode c-archive -o xxx.a xxx.go xxx1.go xxx2.go ...
func RunBuildTask(config string, args string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runBuildStep(args)
}

// Clean output directory
func RunCleanTask(config string) bool {
	if ok := parseConfig(config); !ok {
		return false
	}

	if done := cleanOutput(cfgs); !done {
		clog.Error("Fail clean output directory!")
		return false
	}

	clog.Success("Successfully cleaned output directory ~")
	fmt.Println("")
	return true
}

// Generate bridge C/C++ code
func RunGenerateTask(config string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runGenerateStep()
}

// Compile node addon
func RunMakeTask(config string, args string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runMakeStep(args)
}

// Run all steps: clean -> generate -> build
func RunAllTask(config string, buildArgs string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	buildArgs = normalizeArgs(buildArgs)

	if done := cleanOutput(cfgs); !done {
		clog.Error("Fail clean output directory!")
		return false
	}

	if ok := runGenerateStep(); !ok {
		return false
	}

	if ok := runBuildStep(buildArgs); !ok {
		return false
	}

	clog.Success("Successfully completed all tasks ~")
	fmt.Println("")
	return true
}

func runGenerateStep() bool {
	if done := generateAddonBridge(cfgs); !done {
		clog.Error("Fail generated bridge code!")
		return false
	}

	clog.Success("Successfully generate the Addon bridge c/c++ code of Nodejs ~")
	fmt.Println("")
	return true
}

func runBuildStep(args string) bool {
	args = normalizeArgs(args)
	if ok := runNpmInstall(cfgs, args); !ok {
		clog.Error("Fail build addon via npm install!")
		return false
	}

	clog.Success("Successfully built addon via npm install ~")
	fmt.Println("")
	return true
}

func runMakeStep(args string) bool {
	args = normalizeArgs(args)

	if done := makeToAddon(cfgs, args); !done {
		clog.Error("Fail make addon!")
		return false
	}

	clog.Success("Successfully make the addon of Nodejs ~")
	fmt.Println("")
	return true
}

// Windows-specific build (msvc) removed; rely on build.js/npm install for windows builds.
