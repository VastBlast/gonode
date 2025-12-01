package buildtask

import (
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/cmd"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
	"os"
	"path/filepath"
	"strings"
)

func buildGoToLibrary(cfgs config.Config, args string) bool {
	libName := cfgs.Name + ".a"
	libHName := cfgs.Name + ".h"

	// Remove previously generated artifacts
	outputDir := tools.FormatDirPath(filepath.Join(cfgs.OutPut, "prebuild"))
	paths := []string{
		filepath.Join(outputDir, libName),
		filepath.Join(outputDir, libHName),
	}
	_ = tools.RemoveFiles(paths)

	clog.Info("Start build library ...")
	tempFreePath := ensureTempFreeCString(cfgs.Sources)
	sourceFiles := genBuildFile(cfgs, tempFreePath)
	moduleRoot := findGoModuleRoot(cfgs.Sources)
	workDir := "./"
	if len(moduleRoot) > 0 {
		workDir = moduleRoot
	}
	if d := buildLibrary(sourceFiles, libName, outputDir, args, workDir); !d {
		return false
	}

	return true
}

// Build the list of Go source files
func genBuildFile(config config.Config, extra string) string {
	files := ""
	for _, source := range config.Sources {
		files += " " + source
	}
	if len(extra) > 0 {
		files += " " + extra
	}
	return files
}

func buildLibrary(sourceFiles string, libName string, outPath string, args string, workDir string) bool {
	oPath := outPath + libName
	msg, err := cmd.RunCommand(
		workDir,
		"go build -buildmode c-archive "+args+" -o "+oPath+sourceFiles,
	)
	if err != nil {
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	return true
}

// Ensure a FreeCString implementation exists by generating a temp helper when absent.
func ensureTempFreeCString(sources []string) string {
	if len(sources) == 0 {
		return ""
	}
	first := sources[0]
	if !filepath.IsAbs(first) {
		first = filepath.Join(tools.GetPWD(), first)
	}
	baseDir := filepath.Dir(first)
	if len(baseDir) == 0 {
		return ""
	}
	tempPath := filepath.Join(baseDir, "temp_gonode_helpers.go")
	if tools.Exists(tempPath) {
		return tempPath
	}
	if hasUserFreeCString(sources) {
		return ""
	}
	pkg := detectPackageName(first)
	if pkg == "" {
		pkg = "main"
	}
	content := `package ` + pkg + `

// #include <stdlib.h>
import "C"
import "unsafe"

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}
`
	if err := os.WriteFile(tempPath, []byte(content), 0644); err != nil {
		clog.Warning("failed to write temp FreeCString helper: " + err.Error())
		return ""
	}
	return tempPath
}

func hasUserFreeCString(sources []string) bool {
	for _, src := range sources {
		abs := src
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(tools.GetPWD(), src)
		}
		data, err := os.ReadFile(abs)
		if err != nil {
			continue
		}
		if strings.Contains(string(data), "FreeCString") {
			return true
		}
	}
	return false
}

func detectPackageName(source string) string {
	path := source
	if !filepath.IsAbs(path) {
		path = filepath.Join(tools.GetPWD(), path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "main"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "package "))
		}
	}
	return "main"
}
