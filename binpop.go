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

	searchSymCmd = &cobra.Command{
		Use:   "search",
		Short: "search exported symbols from shared library",
		RunE:  searchExportSymbols,
	}
)

var (
	dumpBinPath = flag.String("path", "", "path to dumpbin.exe")
	walkPath    = flag.String("walk", "C:\\Windows\\System32", "directory or file path")
	findFunc    = flag.String("function", "", "function to search for")
	walkDLL     = flag.String("dll", "", "dll to look in")
)

func exportSymbols(cmd *cobra.Command, args []string) error {
	var symbols sharedlib.Symbols

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
			if filepath.Ext(entry.Name()) == ".dll" {
				symbols, err = sharedlib.ListExportedFunctions(path)
				if err != nil {
					return fmt.Errorf("ListExportedFunctions: %v", err)
				}
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

	if len(symbols) == 0 {
		fmt.Println("No Exported Symbols Found!")
	}
	fmt.Println(symbols)

	return nil
}

func searchExportSymbols(cmd *cobra.Command, args []string) error {

	path, err := cmd.Flags().GetString("path")
	if err != nil {
		return fmt.Errorf("path flag: %v", err)
	}

	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("path is invalid: %v", err)
	}

	search, err := cmd.Flags().GetString("function")
	if err != nil {
		return fmt.Errorf("path flag: %v", err)
	}

	if search == "" {
		return fmt.Errorf("set function flag")
	}

	err = filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if filepath.Ext(entry.Name()) == ".dll" {
				found, err := sharedlib.SearchExportedFunctions(path, search)
				if err != nil {
					return fmt.Errorf("SearchExportedFunctions: %v", err)
				}

				if found {
					fmt.Printf("Function %s was found in %s\n", search, path)
					os.Exit(0)
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("walking path: %v", err)
	}

	fmt.Printf("Function %s was NOT found in %v\n", search, path)
	return nil
}

func init() {
	rootCmd.AddCommand(exportSymCmd)
	exportSymCmd.AddCommand(searchSymCmd)

	exportSymCmd.PersistentFlags().StringP("path", "p", "C:\\Windows\\System32", "dll or directory path to walk")
	exportSymCmd.PersistentFlags().BoolP("debug", "d", true, "print symbols")
	searchSymCmd.PersistentFlags().StringP("function", "f", "", "symbol name to search")

	searchSymCmd.MarkPersistentFlagRequired("function")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
