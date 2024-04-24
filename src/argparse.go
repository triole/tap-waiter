package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	BUILDTAGS      string
	appName        = "tyson tap"
	appDescription = "scan folders for toml, yaml, json or markdown files and offer a web server to fetch information about them"
	appMainversion = "0.1"
)

var CLI struct {
	Conf         string `help:"path to scan, default is current dir" arg:"" optional:"" default:"${curdir}"`
	Threads      int    `help:"max threads, default no of avail. cpu threads" short:"t" default:"${proc}"`
	LogFile      string `help:"log file" default:"/dev/stdout"`
	LogLevel     string `help:"log level" default:"info" enum:"trace,debug,info,error"`
	LogNoColors  bool   `help:"disable output colours, print plain text"`
	LogJSON      bool   `help:"enable json log, instead of text one"`
	ValidateConf bool   `help:"validate configuration and pretty print it"`
	VersionFlag  bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir, _ := os.Getwd()
	ctx := kong.Parse(&CLI,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"curdir":  curdir,
			"logfile": path.Join(os.TempDir(), alnum(appName)+".log"),
			"output":  path.Join(curdir, "tyson.json"),
			"proc":    strconv.Itoa(runtime.NumCPU()),
			"filter":  "\\.(md|toml|yaml|json)$",
		},
	)
	_ = ctx.Run()

	if CLI.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}

func alnum(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile("[^a-z0-9_-]")
	return re.ReplaceAllString(s, "-")
}
