// +build windows
package sharedlib

import (
	"fmt"
	"syscall"
)

var (
	// Importing imagehlp.h & dbghelp.h functions though imagehlp.dll.
	imgHelperDLL                  = syscall.NewLazyDLL("imagehlp.dll")
	procMapAndLoad                = imgHelperDLL.NewProc("MapAndLoad")
	procUnMapAndLoad              = imgHelperDLL.NewProc("UnMapAndLoad")
	procImageDirectoryEntryToData = imgHelperDLL.NewProc("ImageDirectoryEntryToData")
	procImageRvaToVa              = imgHelperDLL.NewProc("ImageRvaToVa")
)

func Dummy() {
	fmt.Println("Shared Lib")
}
