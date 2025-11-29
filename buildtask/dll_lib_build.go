package buildtask

import (
	"github.com/wenlng/gonacli/binding"
	"github.com/wenlng/gonacli/clog"
	"github.com/wenlng/gonacli/cmd"
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
	"path/filepath"
)

func buildToDll(cfgs config.Config) bool {
	rootPath := tools.FormatDirPath(cfgs.OutPut)
	path := tools.FormatDirPath(filepath.Join(cfgs.OutPut, "prebuild"))
	buildCfg := cfgs
	buildCfg.OutPut = path

	// Check whether "gonacli generate" has been run
	if !tools.Exists(filepath.Join(rootPath, cfgs.Name+".cc")) {
		clog.Error("You need to run \"gonacli generate\" generate c/c++ bridge code.")
		return false
	}

	// Check whether "gonacli build" has been run
	if !tools.Exists(filepath.Join(path, cfgs.Name+".a")) {
		clog.Error("You need to run \"gonacli build\" build golang lib.")
		return false
	}

	libFile := cfgs.Name + ".a"
	defFile := cfgs.Name + ".def"
	targetLibName := cfgs.Name + ".dll"
	targetLibName2 := cfgs.Name + ".dll.a"

	// Remove previously generated artifacts
	paths := []string{
		filepath.Join(path, defFile),
		filepath.Join(path, targetLibName),
		filepath.Join(path, targetLibName2),
	}
	_ = tools.RemoveFiles(paths)

	clog.Info("Start build library ...")
	// Generate def file
	if e := binding.GenDefFile(buildCfg, defFile); !e {
		return false
	}

	// Build dll
	if d := buildDll(path, defFile, libFile, targetLibName, targetLibName2); !d {
		clog.Warning("Please check whether the \"gcc\" command is executed correctly.")
		return false
	}

	return true
}

func buildDll(rootPath string, defName string, libName string, dllName string, dllAName string) bool {
	msg, err := cmd.RunCommand(
		"./",
		"cd "+rootPath+" && gcc "+defName+" "+libName+" -shared -lwinmm -lWs2_32 -o "+dllName+" -Wl,--out-implib,"+dllAName,
	)
	// gcc goaddon.def goaddon.a -shared -lwinmm -lWs2_32 -o goaddon.dll -Wl,--out-implib,goaddon.dll.a
	if err != nil {
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	return true
}

func buildToMSVCLib(cfgs config.Config, useVS bool, msvc32Vs bool) bool {
	libFile := cfgs.Name + ".lib"
	defFile := cfgs.Name + ".def"
	targetLibName := cfgs.Name + ".dll"

	outputDir := tools.FormatDirPath(filepath.Join(cfgs.OutPut, "build"))
	paths := []string{
		filepath.Join(outputDir, libFile),
	}
	_ = tools.RemoveFiles(paths)

	if useVS {
		if s := buildMSVCLibWithVSTool(outputDir, defFile, targetLibName, libFile, msvc32Vs); !s {
			//clog.Warning("Please check whether the \"lib.exe\" command exists.")
			return false
		}
	}

	if r := buildMSVCLibWithMinGWTool(outputDir, defFile, targetLibName, libFile); !r {
		//clog.Warning("Please check whether the \"dlltool.exe\" command exists.")
		return false
	}
	return true
}

func buildMSVCLibWithMinGWTool(rootPath string, defName string, dllName string, libName string) bool {
	msg, err := cmd.RunCommand(
		"./",
		"cd "+rootPath+" && dlltool -d "+defName+" -D "+dllName+" -l "+libName,
	)
	// dlltool -d goaddon.def -D goaddon.dll -l goaddon.lib
	if err != nil {
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	return true
}

func buildMSVCLibWithVSTool(rootPath string, defName string, dllName string, libName string, msvc32Vs bool) bool {
	bit := "64"
	if msvc32Vs {
		bit = "32"
	}

	msg, err := cmd.RunCommand(
		"./",
		"cd "+rootPath+" && lib /def:"+defName+" /name:"+dllName+" /out:"+libName+" /MACHINE:X"+bit,
	)
	// lib /def:goaddon.def /name:goaddon.dll /out:goaddon.lib /MACHINE:X64
	if err != nil {
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	return true
}
