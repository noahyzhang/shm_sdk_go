package shm_sdk_go

import (
	"syscall"
	"unsafe"
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
