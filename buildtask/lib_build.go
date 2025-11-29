package buildtask

import (
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/cmd"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
	"path/filepath"
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
	sourceFiles := genBuildFile(cfgs)
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
func genBuildFile(config config.Config) string {
	files := ""
	for _, source := range config.Sources {
		files += " " + source
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
