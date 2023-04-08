// +build windows

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/skepth/binpop/sharedlib"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "binpop",
		Short: "binpop is a binary exploration tool",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	exportSymCmd = &cobra.Command{
		Use:   "export-symbols",
		Short: "export symbols from shared library",
		RunE:  exportSymbols,
	}
)

var (
	dumpBinPath = flag.String("path", "", "path to dumpbin.exe")
	walkPath    = flag.String("walk", "C:\\Windows\\System32", "directory to walk")
	findFunc    = flag.String("func", "", "function to search for")
	walkDLL     = flag.String("dll", "", "dll to look in")
)

func exportSymbols(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return fmt.Errorf("path flag: %v", err)
	}

	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("path is invalid: %v", err)
	}

	err = filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			fmt.Println(path)
			if filepath.Ext(entry.Name()) == ".dll" {
				fmt.Println(entry.Name())
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walking path: %v", err)
	}

	// _, err = cmd.Flags().GetString("search")
	// if err != nil {
	// 	return fmt.Errorf("search flag: %v", err)
	// }

	_, err = sharedlib.ListExportedFunctions("kernel32.dll")
	if err != nil {
		return fmt.Errorf("ListExportedFunctions: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(exportSymCmd)

	exportSymCmd.PersistentFlags().StringP("path", "p", "C:\\Windows\\System32", "dll or directory path to walk")
	exportSymCmd.PersistentFlags().StringP("search", "s", "", "symbol name to search")
	exportSymCmd.PersistentFlags().BoolP("debug", "d", true, "print symbols")

	// exportSymCmd.MarkPersistentFlagRequired("path")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
