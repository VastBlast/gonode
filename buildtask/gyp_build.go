package buildtask

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/cmd"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func makeToAddon(cfgs config.Config, args string) bool {

	rootPath := tools.FormatDirPath(cfgs.OutPut)
	goBuildPath := tools.FormatDirPath(filepath.Join(cfgs.OutPut, "prebuild"))
	buildOutputPath := tools.FormatDirPath(filepath.Join(cfgs.OutPut, "build"))
	platformId, err := getPlatformIdentifier(rootPath)
	if err != nil {
		clog.Error(err)
		return false
	}
	targetDir := filepath.Join(rootPath, "prebuilds", platformId)

	if hasExistingBuild(targetDir) {
		clog.Info("Existing build detected for platform, skip rebuild.")
		return true
	}

	// Check whether "gonode generate" has been run
	if !tools.Exists(filepath.Join(rootPath, cfgs.Name+".cc")) {
		clog.Error("You need to run \"gonode generate\" generate c/c++ bridge code.")
		return false
	}

	// Check whether "gonode build" has been run
	if !tools.Exists(filepath.Join(goBuildPath, cfgs.Name+".a")) {
		clog.Error("You need to run \"gonode build\" build golang lib.")
		return false
	}

	// Check whether "gonode install" has been run
	if !tools.Exists(filepath.Join(rootPath, "node_modules")) {
		clog.Error("You need to run \"gonode install\" install dependencies.")
		return false
	}

	// Remove previously generated artifacts
	files := []string{
		filepath.Join(rootPath, "package-lock.json"),
	}
	_ = tools.RemoveFiles(files)

	clog.Info("Starting make addon ...")
	msg, err := cmd.RunCommand(
		"./",
		"cd "+rootPath+" && node-gyp configure --build-dir="+buildOutputPath+" && node-gyp build "+args,
	)
	if err != nil {
		//clog.Warning("Please check whether the \"node-gyp\" command is executed correctly.")
		clog.Error(err)
		return false
	}
	clog.Info("Make addon done ~")
	clog.Info(msg)

	if err := copyBuiltArtifacts(goBuildPath, buildOutputPath, targetDir, cfgs.Name); err != nil {
		clog.Error(err)
		return false
	}

	_ = os.RemoveAll(goBuildPath)
	_ = os.RemoveAll(buildOutputPath)

	return true
}

func getPlatformIdentifier(rootPath string) (string, error) {
	cmdStr := `node -e "process.stdout.write(require('./platform').platformIdentifier())"`
	out, err := cmd.RunCommand(rootPath, cmdStr)
	if err != nil {
		return "", fmt.Errorf("failed to get platform identifier via node: %v", err)
	}
	if len(out) == 0 {
		return "", fmt.Errorf("empty platform identifier from node")
	}
	return out, nil
}

func hasExistingBuild(targetDir string) bool {
	matches, err := filepath.Glob(filepath.Join(targetDir, "*.node"))
	if err != nil {
		return false
	}
	return len(matches) > 0
}

func moveDllNearNodeBinary(prebuildPath string, name string, buildOutputPath string) bool {
	return true
}

func copyBuiltArtifacts(goBuildPath string, buildOutputPath string, targetDir string, name string) error {
	if err := tools.EnsureDir(targetDir); err != nil {
		return err
	}

	nodeDirs := []string{
		filepath.Join(buildOutputPath, "Release"),
		filepath.Join(buildOutputPath, "Debug"),
	}

	for _, dir := range nodeDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if filepath.Ext(f.Name()) == ".node" {
				src := filepath.Join(dir, f.Name())
				dst := filepath.Join(targetDir, f.Name())
				if copyErr := tools.CopyFile(src, dst); copyErr != nil {
					return copyErr
				}
			}
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err = tools.EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err = io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}

	return out.Close()
}
