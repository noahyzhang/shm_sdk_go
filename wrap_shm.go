package shm_sdk_go

import (
	"syscall"
	"unsafe"
)

// Constants.
const (
	// Mode bits for `shmget`.

	// Create key if key does not exist.
	IPC_CREAT = 01000
	// Fail if key exists.
	IPC_EXCL = 02000
	// Return error on wait.
	IPC_NOWAIT = 04000

	// Special key values.

	// Private key.
	IPC_PRIVATE = 0

	// Flags for `shmat`.

	// Attach read-only access.
	SHM_RDONLY = 010000
	// Round attach address to SHMLBA.
	SHM_RND = 020000
	// Take-over region on attach.
	SHM_REMAP = 040000
	// Execution access.
	SHM_EXEC = 0100000

	// Commands for `shmctl`.

	// Lock segment (root only).
	SHM_LOCK = 1
	// Unlock segment (root only).
	SHM_UNLOCK = 12

	// Control commands for `shmctl`.

	// Remove identifier.
	IPC_RMID = 0
	// Set `ipc_perm` options.
	IPC_SET = 1
	// Get `ipc_perm' options.
	IPC_STAT = 2
)

func GetShm(key int, size int, shmFlag int) (int, error) {
	id, _, err := syscall.Syscall(sysShmGet, uintptr(int32(key)), uintptr(int32(size)), uintptr(int32(shmFlag)))
	if int(id) == -1 {
		return -1, err
	}
	return int(id), nil
}

func AttachShm(shmId int, shmAddr uintptr, shmFlag int) (unsafe.Pointer, error) {
	addr, _, err := syscall.Syscall(sysShmAt, uintptr(int32(shmId)), shmAddr, uintptr(int32(shmFlag)))
	if int(addr) == -1 {
		return nil, err
	}
	return unsafe.Pointer(addr), nil
	//length, err2 := getShmSize(shmId)
	//if err2 != nil {
	//	syscall.Syscall(sysShmDt, addr, 0, 0)
	//	return nil, err2
	//}
	//var b = struct {
	//	addr uintptr
	//	len  int
	//	cap  int
	//}{addr, int(length), int(length)}
	//data := (*[]byte)(unsafe.Pointer(&b))
	//return data, nil
}

func DetachShm(data unsafe.Pointer) error {
	res, _, err := syscall.Syscall(sysShmDt, uintptr(data), 0, 0)
	if int(res) == -1 {
		return err
	}
	return nil
}

func ctlShm(shmId int, cmd int, buf *IdDs) (int, error) {
	res, _, err := syscall.Syscall(sysShmCtl, uintptr(int32(shmId)), uintptr(int32(cmd)), uintptr(unsafe.Pointer(buf)))
	if int(res) == -1 {
		return -1, err
	}
	return int(res), nil
}

func DestroyShm(shmId int) error {
	_, err := ctlShm(shmId, IPC_RMID, nil)
	return err
}

func getShmSize(shmId int) (int64, error) {
	var idDs IdDs
	_, err := ctlShm(shmId, IPC_STAT, &idDs)
	if err != nil {
		return 0, err
	}
	return int64(idDs.SegSz), nil
}
