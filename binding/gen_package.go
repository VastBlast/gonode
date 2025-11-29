package binding

import (
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
)

func GenPackageFile(cfgs config.Config, packageName string) bool {
	code := `{
  "name": "` + cfgs.Name + `",
  "version": "0.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "install": "node ./build.js",
    "build": "node ./build.js",
    "build:debug": "node ./build.js --debug",
    "build:release": "node ./build.js"
  },
  "author": "",
  "license": "ISC",
  "dependencies": {
    "bindings": "^1.5.0",
    "node-addon-api": "^8.0.0"
  },
  "devDependencies": {},
  "gypfile": true,
  "gyp": true
}
`
	// "bindings": "^1.5.0",
	// "node-addon-api": "^5.0.0"
	writePackageFile(code, packageName, cfgs.OutPut)
	return true
}

func writePackageFile(content string, filename string, outPath string) {
	outputDir := tools.FormatDirPath(outPath)
	tools.WriteFile(content, outputDir, filename)
}
