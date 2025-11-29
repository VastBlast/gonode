package binding

import (
	"github.com/wenlng/gonacli/config"
	"github.com/wenlng/gonacli/tools"
)

func GenGypFile(cfgs config.Config, bindingName string) bool {
	code := `{
    "targets": [
        {
            "target_name": "` + cfgs.Name + `",
            "cflags": [ "-O3" ],
            "cflags_cc": [ "-O3" ],
            "cflags!": [ "-fno-exceptions" ],
            "cflags_cc!": [ "-fno-exceptions" ],
            "sources": [ "` + cfgs.Name + `.cc" ],
            "include_dirs": [
                "<!@(node -p \"require('node-addon-api').include\")",
                "<(module_root_dir)/prebuild"
            ],
            "defines": [ "NAPI_CPP_EXCEPTIONS", "NAPI_EXPERIMENTAL" ],
            "conditions": [
                [ 'OS=="linux"', {
                    "cflags": [ "-O3", "-fdata-sections", "-ffunction-sections" ],
                    "cflags_cc": [ "-O3", "-fdata-sections", "-ffunction-sections" ],
                    "ldflags": [ "-Wl,--gc-sections" ],
                    "libraries": [
                        "<(module_root_dir)/prebuild/` + cfgs.Name + `.a"
                    ]
                }],
                [ 'OS=="mac"', {
                    "cflags": [ "-O3", "-fdata-sections", "-ffunction-sections" ],
                    "cflags_cc": [ "-O3", "-fdata-sections", "-ffunction-sections" ],
                    "xcode_settings": {
                        "OTHER_LDFLAGS": [ "-Wl,-dead_strip" ]
                    },
                    "libraries": [
                        "<(module_root_dir)/prebuild/` + cfgs.Name + `.a"
                    ]
                }],
                [ 'OS=="win"', {
                    "msvs_settings": {
                        "VCCLCompilerTool": {
                            "Optimization": "2",
                            "InlineFunctionExpansion": "2",
                            "FavorSizeOrSpeed": "1",
                            "StringPooling": "true",
                            "MinimalRebuild": "false",
                            "BufferSecurityCheck": "false"
                        },
                        "VCLinkerTool": {
                            "OptimizeReferences": "2",
                            "EnableCOMDATFolding": "2"
                        }
                    },
                    "libraries": [
                        "<(module_root_dir)/prebuild/` + cfgs.Name + `.lib"
                    ]
                }],
                [ 'OS!="win" and OS!="linux" and OS!="mac"', {
                    "cflags": [ "-O3" ],
                    "cflags_cc": [ "-O3" ],
                    "libraries": [
                        "<(module_root_dir)/prebuild/` + cfgs.Name + `.a"
                    ]
                }]
            ]
        }
    ]
}`

	writeGypFile(code, bindingName, cfgs.OutPut)
	return true
}

func writeGypFile(content string, filename string, outPath string) {
	outputDir := tools.FormatDirPath(outPath)
	tools.WriteFile(content, outputDir, filename)
}
