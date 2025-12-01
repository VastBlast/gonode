package buildtask

import (
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/cmd"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

// Run npm install inside the generated bindings directory to trigger build.js.
func runNpmInstall(cfgs config.Config, args string) bool {
	bindingsDir := tools.FormatDirPath(cfgs.OutPut)

	cmdStr := "cd " + bindingsDir + " && "
	if len(args) > 0 {
		if tools.IsWindowsOs() {
			cmdStr += "set GO_BUILD_ARGS=" + args + " && "
		} else {
			cmdStr += "GO_BUILD_ARGS='" + args + "' "
		}
	}
	cmdStr += "npm install"

	clog.Info("Running npm install in bindings directory ...")
	msg, err := cmd.RunCommand("./", cmdStr)
	if err != nil {
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	return true
}
