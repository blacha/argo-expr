package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/argoproj/argo-workflows/v3/util/template"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

type InputJson struct {
	Template string            `json:"template"`
	Values   map[string]string `json:"values"`
}

var Version string

// Attempt to fill in the Version var if its not already filled in
// CI fills it using `-ldflags="-X main.Version=$(git describe --tags --always --match 'v*')"`
// doing a go install will use `debug.ReadBuildInfo`
func init_version() {
	if Version != "" {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		Version = "unknown"
		return
	}
	Version = info.Main.Version
}

func create_command() *cobra.Command {

	var values map[string]string
	var from_file string
	var output_to_json bool
	var quiet bool

	rootCmd := &cobra.Command{
		Use:     "argo-expr",
		Short:   "Testing argo expressions",
		Version: Version,
		Long: `Testing argo expression expansions, useful for debugging work expressions without submitting a job to argo
		
Examples:

  Directly convert a input value from a template  
  $ argo-expr "{{=input.parameters.name}}" --value input.parameters.name="hello world" # hello world  
  
  Using Sprig functions  
  $ argo-expr "{{=sprig.trim(input.parameters.name)}}" --value input.parameters.name=" hello world " # hello world  
 
  Convert input to a integer and use math  
  $ argo-expr "{{=asInt(input.parameters.name) + 1}}" --value input.parameters.name="1" # 2`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Preload the argo replacements
			replacement_map := map[string]string{}
			var input_template string

			if from_file != "" {
				jsonFile, err := os.Open(from_file)
				if err != nil {
					panic(err)
				}
				defer jsonFile.Close()

				byteValue, _ := ioutil.ReadAll(jsonFile)
				var jsonInputData InputJson
				json.Unmarshal(byteValue, &jsonInputData)

				if jsonInputData.Template != "" {
					input_template = jsonInputData.Template
				}
				if jsonInputData.Values != nil {
					maps.Copy(replacement_map, jsonInputData.Values)
				}
			}

			if len(args) > 0 {
				if !quiet && input_template != "" {
					cmd.PrintErrf("Replacing template from:'%s' to:'%s'\n", input_template, args[0])
				}
				input_template = args[0]
			}

			// inject all of the --value parameters into the argo replacements
			maps.Copy(replacement_map, values)

			// Convert the template into a JSON object so it can be used by argo
			template_raw := map[string]string{
				"result": input_template,
			}

			template_json, err := json.Marshal(template_raw)
			if err != nil {
				panic(err)
			}

			// Replace the values in the template
			s, err := template.Replace(string(template_json), replacement_map, false)
			if err != nil {
				cmd.Println(err.Error())
				os.Exit(1)
			}

			var replaced_data map[string]interface{}
			err = json.Unmarshal([]byte(s), &replaced_data)
			if err != nil {
				panic(err)
			}

			if output_to_json {
				output_json := map[string]interface{}{
					"template": input_template,
					"values":   replacement_map,
					"result":   replaced_data["result"],
				}
				output, err := json.Marshal(output_json)
				if err != nil {
					panic(err)
				}
				cmd.Println(string(output))
				return
			}
			cmd.Println(replaced_data["result"])
		},
	}
	rootCmd.Flags().StringToStringVarP(&values, "value", "v", map[string]string{}, "Key value pairs of inputs")
	rootCmd.Flags().BoolVar(&output_to_json, "json", false, "Output as a JSON object")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Do not print messages to stderr")
	rootCmd.Flags().StringVarP(&from_file, "from-file", "f", "", "Load parameters from a file")
	return rootCmd
}

func main() {
	init_version()

	var rootCmd = create_command()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
