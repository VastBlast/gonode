package buildtask

import (
	"path/filepath"

	"github.com/wenlng/gonacli/binding"
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/content"
	"github.com/wenlng/gonacli/tools"
)

func generateAddonBridge(cfgs config.Config) bool {

	cppName := cfgs.Name + ".cc"
	bindingName := "binding.gyp"
	indexJsName := "index.js"
	indexDTsName := "index.d.ts"
	packageName := "package.json"
	buildScriptName := "build.js"
	defName := cfgs.Name + ".def"
	platformHelper := "platform.js"

	// Remove previously generated artifacts
	outputDir := tools.FormatDirPath(cfgs.OutPut)
	paths := []string{
		filepath.Join(outputDir, cppName),
		filepath.Join(outputDir, bindingName),
		filepath.Join(outputDir, indexJsName),
		filepath.Join(outputDir, indexDTsName),
		filepath.Join(outputDir, packageName),
		filepath.Join(outputDir, buildScriptName),
		filepath.Join(outputDir, defName),
	}
	_ = tools.RemoveFiles(paths)

	buildCfg := cfgs
	buildCfg.OutPut = outputDir

	// Generate addon C/C++ code
	if g := content.GenCode(buildCfg, cppName); !g {
		//clog.Warning("Please check whether the \"goaddon\" configuration file is correct.")
		return false
	}

	// Generate node-gyp build configuration
	if y := binding.GenGypFile(buildCfg, bindingName); !y {
		return false
	}

	// Generate JS call API to index.js
	if i := binding.GenJsCallIndexFile(cfgs, indexJsName); !i {
		return false
	}

	// Generate JS call API typings to index.d.ts
	if t := binding.GenJsCallDeclareIndexFile(cfgs, indexDTsName); !t {
		return false
	}

	// Generate npm package template
	if p := binding.GenPackageFile(cfgs, packageName); !p {
		return false
	}

	if h := binding.GenPlatformHelper(platformHelper, cfgs); !h {
		return false
	}

	// Generate Windows def template
	if d := binding.GenDefFile(cfgs, defName); !d {
		return false
	}

	// Generate npm build script
	moduleRoot := findGoModuleRoot(cfgs.Sources)
	if b := binding.GenBuildScriptFile(buildCfg, buildScriptName, moduleRoot); !b {
		return false
	}

	return true
}

func findGoModuleRoot(sources []string) string {
	for _, source := range sources {
		srcPath := source
		if !filepath.IsAbs(srcPath) {
			srcPath = filepath.Join(tools.GetPWD(), source)
		}

		dir := filepath.Dir(srcPath)
		for {
			if tools.Exists(filepath.Join(dir, "go.mod")) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}
