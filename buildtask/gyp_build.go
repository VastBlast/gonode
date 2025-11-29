package buildtask

import (
	"io"
	"os"
	"path/filepath"

	"github.com/wenlng/gonacli/clog"
	"github.com/wenlng/gonacli/cmd"
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
)

func makeToAddon(cfgs config.Config, args string) bool {

	path := tools.FormatDirPath(cfgs.OutPut)

	// Check whether "gonacli generate" has been run
	if !tools.Exists(filepath.Join(path, cfgs.Name+".cc")) {
		clog.Error("You need to run \"gonacli generate\" generate c/c++ bridge code.")
		return false
	}

	// Check whether "gonacli build" has been run
	if !tools.Exists(filepath.Join(path, cfgs.Name+".a")) {
		clog.Error("You need to run \"gonacli build\" build golang lib.")
		return false
	}

	// Check whether "gonacli install" has been run
	if !tools.Exists(filepath.Join(path, "node_modules")) {
		clog.Error("You need to run \"gonacli install\" install dependencies.")
		return false
	}

	// On Windows, verify whether "gonacli msvc" has been run
	if tools.IsWindowsOs() {
		if !tools.Exists(filepath.Join(path, cfgs.Name+".lib")) {
			clog.Error("You need to run \"gonacli msvc\" build lib on windows OS.")
			return false
		}
		if !tools.Exists(filepath.Join(path, cfgs.Name+".dll")) {
			clog.Error("You need to run \"gonacli msvc\" build dll on windows OS.")
			return false
		}
	}

	// Remove previously generated artifacts
	_ = tools.RemoveDirContents(filepath.Join(path, "build"))
	files := []string{
		filepath.Join(path, "package-lock.json"),
	}
	_ = tools.RemoveFiles(files)

	clog.Info("Starting make addon ...")
	msg, err := cmd.RunCommand(
		"./",
		"cd "+path+" && node-gyp configure && node-gyp build "+args,
	)
	if err != nil {
		//clog.Warning("Please check whether the \"node-gyp\" command is executed correctly.")
		clog.Error(err)
		return false
	}
	clog.Info("Make addon done ~")
	clog.Info(msg)

	if tools.IsWindowsOs() {
		if ok := moveDllNearNodeBinary(path, cfgs.Name); !ok {
			return false
		}
	}

	return true
}

func moveDllNearNodeBinary(rootPath string, name string) bool {
	dllPath := filepath.Join(rootPath, name+".dll")
	if !tools.Exists(dllPath) {
		clog.Error("The dll file is missing, please execute \"gonacli msvc\" first.")
		return false
	}

	nodeBins, err := filepath.Glob(filepath.Join(rootPath, "build", "*", name+".node"))
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
