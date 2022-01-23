//nolint
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}
func PrettyPrintAsJson(v interface{}) (Json string, err error) {
	JsonBuffer, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	Json = string(JsonBuffer)
	return
}

type CliGlobals struct {
	Config   string      `help:"Location of client config files" type:"path" default:""`
	Debug    bool        `short:"D" help:"Enable debug mode"`
	LogLevel string      `short:"l" help:"Set the logging level (debug|info|warn|error|fatal)" default:"info"`
	Version  VersionFlag `name:"version" help:"Print version information and quit"`
}

type RegisterTargetCmd struct {
	Name string `required help:"Named alias of the Seagate Enclosure target."`
	Url  string `required help:"URL of the Seagate enclosure target."`
	User string `required help:"user name to use to authenticate with the target."`
	Pass string `help:"password to use to authenticate with the target"`
}

func (aCmd *RegisterTargetCmd) Run(globals *CliGlobals) error {
	fmt.Println("Running target command")
	fmt.Printf("Config: %s\n", globals.Config)
	fmt.Printf("URL : %v\n", aCmd.Name)
	fmt.Printf("URL : %v\n", aCmd.Url)
	fmt.Printf("User: %v\n", aCmd.User)
	fmt.Printf("Pass: %v\n", aCmd.Pass)
	return nil
}
func (aCmd *RegisterTargetCmd) AsJson() (jsonStr string, err error) {
	prettyJSON, err := json.MarshalIndent(aCmd, "", "  ")
	if err != nil {
		return
	}
	jsonStr = string(prettyJSON)
	return
}

type CvtRawCmd struct {
	Target []string `required help:"named enclosure target" short:"t"`
	Cmd    []string `arg passthrough:"" required help:"Pass through any command string to the targeted enclosure"`
}

func (aCmd *CvtRawCmd) Run(globals *CliGlobals) error {
	fmt.Println("Running raw command")
	fmt.Println("Target: ", aCmd.Target)
	fmt.Println("Cmd:", strings.Join(aCmd.Cmd[:], " "))
	fmt.Println("Cmd.len:", len(aCmd.Cmd))
	return nil
}
func (aCmd *CvtRawCmd) AsJson() (jsonStr string, err error) {
	prettyJSON, err := json.MarshalIndent(aCmd, "", "  ")
	if err != nil {
		return
	}
	jsonStr = string(prettyJSON)
	return
}

type CvtShowCmd struct {
	Target []string `required help:"named enclosure target" short:"t"`
	Disks  struct {
	} `cmd help:"Generates a report of disk present in the enclosure"`
	DiskGroups struct {
	} `cmd help:"Generates a report of defined disk groups."`

	AlertConditionHistory struct {
	} `cmd help:"show alert condition history"`
	AdvancedSettings struct {
	} `cmd help:"show advanced settings"`
	Certificates struct {
		Json bool `negateable:"" optional:"" default:"false" help:"Flag to output certificates as json"`
	} `cmd help:"show enclosure https certificates"`
	Volumes struct {
	} `cmd help:"show volumes"`
	Tester struct {
		Fred struct {
			Fred string `arg required help:"set fred=SomeValue"`
		} `arg help:"fred is a nested argument, try fred=Argument"`
		//Barney string `arg help:"set barney=SomeValue"`
		Wilma bool `negatable:"" short:"w" help:"wilma is a nested flag, try --wilma"`
	} `cmd help:"tester cmd"`
}

func (aCmd *CvtShowCmd) AsJson() (jsonStr string, err error) {
	prettyJSON, err := json.MarshalIndent(aCmd, "", "  ")
	if err != nil {
		return
	}
	jsonStr = string(prettyJSON)
	return
}

func (aCmd *CvtShowCmd) Run(globals *CliGlobals, kCtx *kong.Context) error {
	subCmdName := kCtx.Selected().Name
	fmt.Printf("Sending \"show %s to %s\n", subCmdName, strings.Join(aCmd.Target, ","))
	//fmt.Println(aCmd.AsJson())
	fmt.Println(PrettyPrintAsJson(aCmd))
	switch subCmdName {
	case "disks":
		fmt.Println("Show Disks command found")
	case "advanced-settings":
		fmt.Println("Show AdvancedSettings command found")
	case "alert-condition-history":
		fmt.Println("Show AlertConditionHistory command found")
	case "certificates":
		fmt.Println("Show Certificates command found")
	case "volumes":
		fmt.Println("Show Volumes command found")
	}
	return nil
}

type CLI struct {
	CliGlobals
	Register RegisterTargetCmd `cmd help:"Register an enclosure target to manage." short:"R"`
	Show     CvtShowCmd        `cmd help:"Show commands"`
	Raw      CvtRawCmd         `cmd help:"Send Raw command string to target enclosure"`
}
