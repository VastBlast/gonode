package buildtask

import (
	"github.com/wenlng/gonacli/binding"
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/content"
	"github.com/wenlng/gonacli/tools"
	"path/filepath"
)

func generateAddonBridge(cfgs config.Config) bool {

	cppName := cfgs.Name + ".cc"
	bindingName := "binding.gyp"
	indexJsName := "index.js"
	indexDTsName := "index.d.ts"
	packageName := "package.json"

	// Remove previously generated artifacts
	outputDir := tools.FormatDirPath(cfgs.OutPut)
	paths := []string{
		filepath.Join(outputDir, cppName),
		filepath.Join(outputDir, bindingName),
		filepath.Join(outputDir, indexJsName),
		filepath.Join(outputDir, indexDTsName),
		filepath.Join(outputDir, packageName),
	}
	//_ = tools.RemoveDirContents(outputDir)
	_ = tools.RemoveFiles(paths)

	// Generate addon C/C++ code
	if g := content.GenCode(cfgs, cppName); !g {
		//clog.Warning("Please check whether the \"goaddon\" configuration file is correct.")
		return false
	}

	// Generate node-gyp build configuration
	if y := binding.GenGypFile(cfgs, bindingName); !y {
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

	return true
}
