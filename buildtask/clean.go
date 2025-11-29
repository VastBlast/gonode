package buildtask

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func cleanOutput(cfgs config.Config) bool {
	output := strings.TrimSpace(cfgs.OutPut)
	if len(output) == 0 {
		clog.Error("Output path is empty.")
		return false
	}

	outputDir := tools.FormatDirPath(output)
	if filepath.Clean(outputDir) == string(filepath.Separator) {
		clog.Error("Output path resolves to the root directory, aborting clean.")
		return false
	}

	if !tools.Exists(outputDir) {
		clog.Info("Output directory does not exist, skip clean ~")
		return true
	}

	if !tools.IsDir(outputDir) {
		clog.Error("The output path is not a directory.")
		return false
	}

	clog.Info("Cleaning output directory ...")
	if err := tools.RemoveDirContents(outputDir); err != nil {
		clog.Error(err)
		return false
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		clog.Error(err)
		return false
	}
	if len(files) == 0 {
		if err := os.Remove(outputDir); err != nil {
			clog.Warning("Output directory is empty but could not be removed:", err)
		} else {
			clog.Info("Output directory removed because it was empty.")
		}
	}

	return true
}
