// +build windows
package sharedlib

import (
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/0xrawsec/golang-win32/win32"
	"github.com/soyum2222/editPE"
)

const (
	IMAGE_DIRECTORY_ENTRY_EXPORT = 0
)

var (
	// Importing imagehlp.h & dbghelp.h functions though imagehlp.dll.
	imgHelperDLL                  = syscall.NewLazyDLL("imagehlp.dll")
	dbghelpDLL                    = syscall.NewLazyDLL("dbghelp.dll")
	procMapAndLoad                = imgHelperDLL.NewProc("MapAndLoad")
	procUnMapAndLoad              = imgHelperDLL.NewProc("UnMapAndLoad")
	procImageDirectoryEntryToData = dbghelpDLL.NewProc("ImageDirectoryEntryToData") // imgHelperDLL.NewProc("ImageDirectoryEntryToData")
	procImageRvaToVa              = dbghelpDLL.NewProc("ImageRvaToVa")
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
	ModuleName    *uint8
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
	Signature      uint32
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER64
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

type IMAGE_OPTIONAL_HEADER64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               []IMAGE_DATA_DIRECTORY
}

type IMAGE_DATA_DIRECTORY struct {
	VirtualAddress uint32
	Size           uint32
}

//
func stringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}

func ListExportedFunctions(dllName string) (Symbols, error) {
	peFile := editPE.PE{}
	funcNames := Symbols{}
	f, err := ioutil.ReadFile(dllName)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}

	peFile.Parse(f)

	for _, element := range peFile.GetExportFunc().FuncName {
		funcNames = append(funcNames, string(element.Name))
	}

	return funcNames, nil
}

func SearchExportedFunctions(dllName string, funcName string) (bool, error) {
	peFile := editPE.PE{}
	f, err := ioutil.ReadFile(dllName)
	if err != nil {
		return false, fmt.Errorf("opening file: %v", err)
	}

	peFile.Parse(f)

	for _, element := range peFile.GetExportFunc().FuncName {
		if string(element.Name) == funcName {
			return true, nil
		}
	}

	return false, nil
}
