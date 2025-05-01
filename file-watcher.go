package main

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	bufferSize   = 4096
	notifyFilter = windows.FILE_NOTIFY_CHANGE_FILE_NAME |
		windows.FILE_NOTIFY_CHANGE_DIR_NAME |
		windows.FILE_NOTIFY_CHANGE_ATTRIBUTES |
		windows.FILE_NOTIFY_CHANGE_SIZE |
		windows.FILE_NOTIFY_CHANGE_LAST_WRITE
)

var (
	modkernel32               = syscall.NewLazyDLL("kernel32.dll")
	procReadDirectoryChangesW = modkernel32.NewProc("ReadDirectoryChangesW")
)

func WatchDirectory(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf("UTF16PtrFromString: %v", err)
	}

	handle, err := windows.CreateFile(
		pathPtr,
		windows.FILE_LIST_DIRECTORY,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return fmt.Errorf("CreateFile: %v", err)
	}
	defer windows.CloseHandle(handle)

	buffer := make([]byte, bufferSize)

	for {
		var bytesReturned uint32

		r1, _, err := procReadDirectoryChangesW.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0, // FALSE = no recursivo
			uintptr(notifyFilter),
			uintptr(unsafe.Pointer(&bytesReturned)),
			0,
			0,
		)

		if r1 == 0 {
			return fmt.Errorf("ReadDirectoryChangesW failed: %v", err)
		}

		parseEvents(buffer[:bytesReturned])
	}
}

type fileNotifyInformation struct {
	NextEntryOffset uint32
	Action          uint32
	FileNameLength  uint32
}

func parseEvents(buf []byte) {
	offset := 0
	for {
		header := (*fileNotifyInformation)(unsafe.Pointer(&buf[offset]))

		nameBytes := buf[offset+12 : offset+12+int(header.FileNameLength)]
		u16 := make([]uint16, len(nameBytes)/2)
		for i := 0; i < len(u16); i++ {
			u16[i] = uint16(nameBytes[i*2]) | uint16(nameBytes[i*2+1])<<8
		}
		name := string(utf16.Decode(u16))

		switch header.Action {
		case windows.FILE_ACTION_ADDED:
			fmt.Println("[+] Archivo creado:", name)
		case windows.FILE_ACTION_REMOVED:
			fmt.Println("[-] Archivo eliminado:", name)
		case windows.FILE_ACTION_MODIFIED:
			fmt.Println("[*] Archivo modificado:", name)
		case windows.FILE_ACTION_RENAMED_OLD_NAME:
			fmt.Println("[>] Renombrado desde:", name)
		case windows.FILE_ACTION_RENAMED_NEW_NAME:
			fmt.Println("[<] Renombrado a:", name)
		default:
			fmt.Println("[?] Acción desconocida")
		}

		if header.NextEntryOffset == 0 {
			break
		}
		offset += int(header.NextEntryOffset)
	}
}
