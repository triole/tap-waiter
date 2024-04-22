package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	lg = logseal.Init("debug", nil, true, false)
)

var CLI struct {
	Path        string `help:"path to scan, default it current dir" arg:"" optional:"" default:"${curdir}"`
	Threads     int    `help:"max threads to run, default no of avail. cpu threads" short:"t" default:"${proc}"`
	Limit       int    `help:"max of requests to execute" short:"l" default:"0"`
	UA          string `help:"user agent" short:"u" default:"${ua}"`
	OmdbBaseURL string `help:"omdb base url, for firing requests" short:"b" default:"${omdbBaseURL}"`
	Force       bool   `help:"force overwrite because usually episode list are skipped when existent" short:"f"`
	Verbose     bool   `help:"verbose mode" short:"v"`
	DryRun      bool   `help:"dry run, just print don't do" short:"n"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level" default:"info" enum:"trace,debug,info,error"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	LogJSON     bool   `help:"enable json log, instead of text one"`
	VersionFlag bool   `help:"display version" short:"V"`
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
			"curdir":      curdir,
			"proc":        strconv.Itoa(runtime.NumCPU()),
			"ua":          "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
			"omdbBaseURL": "http://www.omdbapi.com/?plot=full",
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
