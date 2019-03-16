package cmd

import (
	"path/filepath"

	"github.com/mpppk/goofy/lib"

	"github.com/spf13/cobra"
)

var filePath *string
var outFilePath *string
var mustifyCmd = &cobra.Command{
	Use:   "mustify",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if *outFilePath == "" {
			base := filepath.Base(*filePath)
			o := filepath.Join(filepath.Dir(*filePath), "must-"+base)
			outFilePath = &o
		}

		file, newDecls, err := lib.GenerateErrorWrappersFromProgram(*filePath)
		if err != nil {
			panic(err)
		}
		file.Decls = newDecls

		if err := lib.WriteAstFile(*outFilePath, file); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mustifyCmd)
	filePath = mustifyCmd.Flags().String("file", "", "target file path")
	outFilePath = mustifyCmd.Flags().String("out", "", "file path to save output")
}
