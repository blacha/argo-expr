package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/argoproj/argo-workflows/v3/util/template"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func main() {
	var var_map map[string]string
	var raw bool

	rootCmd := &cobra.Command{
		Use:   "argo-expr",
		Short: "Testing argo expressions",
		Long:  ``,
		Example: `  
  Directly convert a input value from a template
  
  $ argo-expr "{{=input.parameters.name}}" --value input.parameters.name="hello world" --raw # hello world

  Using Sprig functions

  $ argo-expr "{{=sprig.trim(input.parameters.name)}}" --value input.parameters.name=" hello world " --raw # hello world

  Convert input to a integer and use math

  $ argo-expr "{{=asInt(input.parameters.name) + 1}}" --value input.parameters.name="1" --raw # 2

		`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input_template := args[0]
			// Preload the argo replacements
			base_map := map[string]string{}

			// inject all of the --value parameters into the argo replacements
			maps.Copy(base_map, var_map)

			// Convert the template into a JSON object so it can be used by argo
			template_raw := map[string]string{
				"result": input_template,
			}

			template_json, err := json.Marshal(template_raw)
			if err != nil {
				panic(err)
			}

			// Replace the values in the template
			s, err := template.Replace(string(template_json), base_map, true)
			if err != nil {
				panic(err)
			}
			var replaced_data map[string]interface{}
			err = json.Unmarshal([]byte(s), &replaced_data)
			if err != nil {
				panic(err)
			}

			if raw {
				fmt.Println(replaced_data["result"])
				return
			}

			output_json := map[string]interface{}{
				"input":  input_template,
				"values": base_map,
				"result": replaced_data["result"],
			}
			output, err := json.Marshal(output_json)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(output))
		},
	}
	rootCmd.Flags().StringToStringVarP(&var_map, "value", "v", map[string]string{}, "Key value pairs")
	rootCmd.Flags().BoolVarP(&raw, "raw", "r", false, "output only the raw value")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
