package buildtask

import (
	"github.com/VastBlast/gonode/clog"
	"github.com/VastBlast/gonode/cmd"
	"github.com/VastBlast/gonode/config"
	"github.com/VastBlast/gonode/tools"
)

func installDep(cfgs config.Config) bool {
	path := tools.FormatDirPath(cfgs.OutPut)

	clog.Info("Starting install dependencies ...")
	// "bindings" "node-addon-api"
	msg, err := cmd.RunCommand(
		"./",
		"cd "+path+" && npm install bindings node-addon-api -S",
	)
	if err != nil {
		//clog.Warning("Please check whether the \"npm\" command is executed correctly.")
		clog.Error(err)
		return false
	}
	clog.Info(msg)
	clog.Info("Install dependencies done ~")
	return true
}
