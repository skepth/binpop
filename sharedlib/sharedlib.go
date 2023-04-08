// +build windows
package sharedlib

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/0xrawsec/golang-win32/win32"
)

var (
	// Importing imagehlp.h & dbghelp.h functions though imagehlp.dll.
	imgHelperDLL                  = syscall.NewLazyDLL("imagehlp.dll")
	procMapAndLoad                = imgHelperDLL.NewProc("MapAndLoad")
	procUnMapAndLoad              = imgHelperDLL.NewProc("UnMapAndLoad")
	procImageDirectoryEntryToData = imgHelperDLL.NewProc("ImageDirectoryEntryToData")
	procImageRvaToVa              = imgHelperDLL.NewProc("ImageRvaToVa")
)

type Symbols []string

// ref: https://microsoft.github.io/windows-docs-rs/doc/windows/Win32/System/SystemServices/struct.IMAGE_EXPORT_DIRECTORY.html
type IMAGE_EXPORT_DIRECTORY struct {
	Characteristics       uint32
	TimeDateStamp         uint32
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}

// ref: https://microsoft.github.io/windows-docs-rs/doc/windows/Win32/System/Diagnostics/Debug/struct.LOADED_IMAGE.html
type LOADED_IMAGE struct {
	// ModuleName    win32.PWSTR //pstr?
	hFile         win32.HANDLE
	MappedAddress *uint8
	FileHeader    *IMAGE_NT_HEADERS64
	// LastRvaSection *IMAGE_SECTION_HEADER
	NumberOfSections uint32
	// Sections *IMAGE_SECTION_HEADER
	Characteristics uint32
	fSystemImage    bool
	fDOSImage       bool
	fReadOnly       bool
	Version         uint8
	// Links LIST_ENTRY
	SizeOfImage uint32
}

type IMAGE_NT_HEADERS64 struct {
	Signature  uint32
	FileHeader IMAGE_FILE_HEADER
	// OptionalHeader IMAGE_OPTIONAL_HEADER64
}

type IMAGE_FILE_HEADER struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

//
func stringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}

func ListExportedFunctions(dllName string) (Symbols, error) {

	// var dNameRVAs *win32.DWORD
	// imgExpDir := new(IMAGE_EXPORT_DIRECTORY)
	// var cDirSize uint64 // unsigned long?
	loadedImg := new(LOADED_IMAGE)
	// var sName string
	// sListOfFunctions := new(Symbols)

	// MapAndLoad the DLL in Question
	// ref: https://learn.microsoft.com/en-us/windows/win32/api/imagehlp/nf-imagehlp-mapandload
	result, _, _ := procMapAndLoad.Call(
		uintptr(unsafe.Pointer(stringToCharPtr(dllName))),
		win32.NULL,
		uintptr(unsafe.Pointer(loadedImg)),
		uintptr(win32.TRUE),
		uintptr(win32.TRUE),
	)

	if int(result) == int(win32.TRUE) {
		fmt.Println("Success")
	}

	// fmt.Printf("Call returned: %v\n", (win32.BOOL)(success))
	fmt.Printf("Loaded Image: %v\n", loadedImg)

	// if err != nil {
	// 	return nil, fmt.Errorf("procMapAndLoad call failed: %v", err)
	// }

	fmt.Println("In Load")

	// fmt.Printf("Call returned: %v\n", success)
	// fmt.Printf("Loaded Image: %v\n", loadedImg)

	// Call ImageDirectoryEntryToData to get pointer to AddressOfNames & NumberOfNames

	// Call ImageRvaToVa over AddressOfNames

	// Iterate over NumberOfNames

	// |_> Call ImageRvaToVa to get the func names from RVA

	return nil, nil
}

func Dummy() {
	fmt.Println("Shared Lib")
}
