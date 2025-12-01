package cli

import (
	"flag"
	"fmt"
	"github.com/VastBlast/gonode/buildtask"
	"log"
	"os"
)

type CLI struct{}

func printUsage() {
	fmt.Println("Usage:")

	fmt.Println("\tversion -- Get version")
	fmt.Println("\thelp -- Help")
	fmt.Println("\tall -- Run clean, generate, build")
	fmt.Println("\tclean -- Clean output directory")
	fmt.Println("\tgenerate -- Generate napi c/c++ code of golang and addon bridge")
	fmt.Println("\tbuild -- Compile the golang source file of the export api")
	fmt.Println("\tmake -- Compile addon bindings of nodejs")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run(name string, version string) {
	isValidArgs()

	// gonode build
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	// gonode generate
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	// gonode make
	makeCmd := flag.NewFlagSet("make", flag.ExitOnError)
	// gonode version
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	// gonode help
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	// gonode clean
	cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	// gonode all
	allCmd := flag.NewFlagSet("all", flag.ExitOnError)

	// gonode build --config xxx.json
	buildCofig := buildCmd.String("config", "goaddon.json", "Addon api export configuration file")
	// gonode build --args '-ldflags "-s -w"'
	buildArg := buildCmd.String("args", "-ldflags \"-s -w\"", "Golang compilation arguments")
	// gonode build --upx
	//buildUpx := buildCmd.Bool("upx", false, "Call the upx compression command")
	// gonode build --xgo
	//buildXgo := buildCmd.Bool("xgo", false, "Call the xgo compression command")
	// gonode generate --config xxx.json
	generateConfig := generateCmd.String("config", "goaddon.json", "Addon api export configuration file")
	// gonode make --args "xxx"
	makeArg := makeCmd.String("args", "", "Nodegyp compilation arguments")
	// gonode make --config xxx.json
	makeConfig := makeCmd.String("config", "goaddon.json", "Addon api export configuration file")
	//makeMpn := makeCmd.Bool("npm-i", false, "Install npm dependencies")
	// gonode clean --config xxx.json
	cleanConfig := cleanCmd.String("config", "goaddon.json", "Addon api export configuration file")
	// gonode all
	allConfig := allCmd.String("config", "goaddon.json", "Addon api export configuration file")
	allBuildArg := allCmd.String("build-args", "-ldflags \"-s -w\"", "Golang compilation arguments")

	switch os.Args[1] {
	case "build":
		err := buildCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "generate":
		err := generateCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "make":
		err := makeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "version":
		err := versionCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "help":
		err := helpCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "clean":
		err := cleanCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "all":
		err := allCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if helpCmd.Parsed() {
		printUsage()
		return
	}

	if versionCmd.Parsed() {
		fmt.Println(name, version)
		return
	}

	if buildCmd.Parsed() {
		if ok := buildtask.RunBuildTask(*buildCofig, *buildArg); !ok {
			os.Exit(1)
		}
		return
	}

	if generateCmd.Parsed() {
		if ok := buildtask.RunGenerateTask(*generateConfig); !ok {
			os.Exit(1)
		}
		return
	}

	if makeCmd.Parsed() {
		if ok := buildtask.RunMakeTask(*makeConfig, *makeArg); !ok {
			os.Exit(1)
		}
		return
	}

	if cleanCmd.Parsed() {
		if ok := buildtask.RunCleanTask(*cleanConfig); !ok {
			os.Exit(1)
		}
		return
	}

	if allCmd.Parsed() {
		if ok := buildtask.RunAllTask(*allConfig, *allBuildArg); !ok {
			os.Exit(1)
		}
		return
	}
}
