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

		fileMap, err := lib.GenerateErrorWrappersFromPackage(*filePath, "main", "must-")
		if err != nil {
			panic(err)
		}

		for fp, file := range fileMap {
			dirPath := filepath.Dir(fp)
			fileName := filepath.Base(fp)
			newFilePath := filepath.Join(dirPath, "must-"+fileName)
			if err := lib.WriteAstFile(newFilePath, file); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mustifyCmd)
	filePath = mustifyCmd.Flags().String("file", "", "target file path")
	outFilePath = mustifyCmd.Flags().String("out", "", "file path to save output")
}
