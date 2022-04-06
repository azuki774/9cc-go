/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/azuki774/9cc-go/internal/compiler"
	"github.com/azuki774/9cc-go/internal/parser"
	"github.com/spf13/cobra"
)

var (
	OutputFileName string
	SourceFileName string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "9cc-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return fmt.Errorf("invalid argument")
		}

		OutputFileName, _ = cmd.Flags().GetString("output")
		// fmt.Printf("Output file name : %s\n", OutputFileName)
		SourceFileName = args[0]
		// fmt.Printf("Source file name : %s\n", SourceFileName)

		tf, _ := cmd.Flags().GetBool("tokenize")
		stf, _ := cmd.Flags().GetBool("show-tokenize")
		nmain, _ := cmd.Flags().GetBool("no-main")
		if tf {
			err = compiler.TokenizeOnly(SourceFileName)
			return err
		}

		if nmain {
			parser.NoMain = true
		}

		err = compiler.CompileMain(OutputFileName, SourceFileName, stf)
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("output", "o", "output.s", "output file name")
	rootCmd.Flags().BoolP("tokenize", "", false, "only tokenize")
	rootCmd.Flags().BoolP("show-tokenize", "", false, "show tokenize")
	rootCmd.Flags().BoolP("no-main", "n", false, "no main")
}
