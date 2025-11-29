package buildtask

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/VastBlast/gonode/buildtask/compatible"
	"github.com/VastBlast/gonode/check"
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
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

// Initialize npm install dependencies
func RunInstallTask(config string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runInstallStep()
}

// Compile node addon
func RunMakeTask(config string, args string) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runMakeStep(args)
}

// Run all steps: clean -> generate -> build -> install -> (windows only) msvc -> make
func RunAllTask(config string, buildArgs string, makeArgs string, useVS bool, msvc32Vs bool) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	buildArgs = normalizeArgs(buildArgs)
	makeArgs = normalizeArgs(makeArgs)

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

	if ok := runInstallStep(); !ok {
		return false
	}

	if tools.IsWindowsOs() {
		if ok := runMsvcStep(useVS, msvc32Vs); !ok {
			return false
		}
	} else {
		clog.Info("Skip msvc step on non-windows platform")
	}

	if ok := runMakeStep(makeArgs); !ok {
		return false
	}

	clog.Success("Successfully completed all tasks ~")
	fmt.Println("")
	return true
}

// Windows environment compatibility handling
func RunMsvcTask(config string, useVS bool, msvc32Vs bool) bool {
	if ok := parseAndCheck(config); !ok {
		return false
	}

	return runMsvcStep(useVS, msvc32Vs)
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

	if d := buildGoToLibrary(cfgs, args); !d {
		clog.Error("Fail build golang lib!")
		return false
	}

	clog.Success("Successfully build golang lib ~")
	fmt.Println("")
	return true
}

func runInstallStep() bool {
	if done := installDep(cfgs); !done {
		clog.Error("Fail installed!")
		return false
	}

	clog.Success("Successfully installed ~")
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

func runMsvcStep(useVS bool, msvc32Vs bool) bool {
	if !tools.IsWindowsOs() {
		clog.Error("The \"msvc\" command is only supported on Windows.")
		return false
	}

	clog.Info("Starting fix file ...")
	compatible.FixCGOWithWindow(cfgs)
	clog.Info("Fix file done ~")

	clog.Info("Starting build dll ...")
	if done := buildToDll(cfgs); !done {
		clog.Error("Fail build dll!")
		return false
	}
	clog.Success("Successfully build dll ~")
	fmt.Println("")

	clog.Info("Starting build lib ...")
	if done := buildToMSVCLib(cfgs, useVS, msvc32Vs); !done {
		clog.Error("Fail build lib!")
		return false
	}

	clog.Success("Successfully build lib ~")
	fmt.Println("")

	return true
}
