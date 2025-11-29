package buildtask

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wenlng/gonacli/clog"
	"github.com/wenlng/gonacli/cmd"
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
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

	// Check whether "gonacli generate" has been run
	if !tools.Exists(filepath.Join(rootPath, cfgs.Name+".cc")) {
		clog.Error("You need to run \"gonacli generate\" generate c/c++ bridge code.")
		return false
	}

	// Check whether "gonacli build" has been run
	if !tools.Exists(filepath.Join(goBuildPath, cfgs.Name+".a")) {
		clog.Error("You need to run \"gonacli build\" build golang lib.")
		return false
	}

	// Check whether "gonacli install" has been run
	if !tools.Exists(filepath.Join(rootPath, "node_modules")) {
		clog.Error("You need to run \"gonacli install\" install dependencies.")
		return false
	}

	// On Windows, verify whether "gonacli msvc" has been run
	if tools.IsWindowsOs() {
		if !tools.Exists(filepath.Join(goBuildPath, cfgs.Name+".lib")) {
			clog.Error("You need to run \"gonacli msvc\" build lib on windows OS.")
			return false
		}
		if !tools.Exists(filepath.Join(goBuildPath, cfgs.Name+".dll")) {
			clog.Error("You need to run \"gonacli msvc\" build dll on windows OS.")
			return false
		}
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

	if tools.IsWindowsOs() {
		if ok := moveDllNearNodeBinary(goBuildPath, cfgs.Name, buildOutputPath); !ok {
			return false
		}
	}

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
	dllPath := filepath.Join(prebuildPath, name+".dll")
	if !tools.Exists(dllPath) {
		clog.Error("The dll file is missing, please execute \"gonacli msvc\" first.")
		return false
	}

	nodeBins, err := filepath.Glob(filepath.Join(buildOutputPath, "*", name+".node"))
	if err != nil {
		clog.Error(err)
		return false
	}

	if len(nodeBins) == 0 {
		clog.Error("The addon binary was not generated, please confirm node-gyp build output.")
		return false
	}

	for _, nodeBin := range nodeBins {
		targetPath := filepath.Join(filepath.Dir(nodeBin), name+".dll")
		if err = copyFile(dllPath, targetPath); err != nil {
			clog.Error(err)
			return false
		}
	}

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

	dllPath := filepath.Join(goBuildPath, name+".dll")
	if tools.Exists(dllPath) {
		dst := filepath.Join(targetDir, name+".dll")
		if err := tools.CopyFile(dllPath, dst); err != nil {
			return err
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
