package shm_sdk_go

import (
	"fmt"
	"syscall"
	"unsafe"
)

type semBuf struct {
	semNum uint16
	semOp  int16
	semFlg int16
}

func shmGet(key int) (int, error) {
	r1, _, err := syscall.Syscall(sysSemGet, uintptr(key), uintptr(1), uintptr(00666))
	if int(r1) < 0 {
		r1, _, err = syscall.Syscall(sysSemGet, uintptr(key),
			uintptr(1), uintptr(IPC_CREAT|IPC_EXCL|00666))
		if int(r1) < 0 {
			return -1, fmt.Errorf("[shmGet] syscall SYS_SEMGET err: %s", err.Error())
		}
		if err2 := semSetVal(int(r1), 1); err2 != nil {
			return -1, fmt.Errorf("[shmGet] semSetVal err: %s", err2.Error())
		}
	}
	return int(r1), nil
}

func semLock(semId int, isWait bool) error {
	var stSemBuf semBuf
	if isWait {
		stSemBuf = semBuf{semNum: 0, semOp: -1, semFlg: SEM_UNDO}
	} else {
		stSemBuf = semBuf{semNum: 0, semOp: -1, semFlg: IPC_NOWAIT | SEM_UNDO}
	}
	r1, _, err := syscall.Syscall(sysSemOp, uintptr(semId), uintptr(unsafe.Pointer(&stSemBuf)), 1)
	if int(r1) < 0 {
		return fmt.Errorf("[semLock] syscall SYS_SEMOP err: %s", err.Error())
	}
	return nil
}

func semUnLock(semId int) error {
	var stSemBuf = semBuf{semNum: 0, semOp: 1, semFlg: SEM_UNDO}
	r1, _, err := syscall.Syscall(sysSemOp, uintptr(semId), uintptr(unsafe.Pointer(&stSemBuf)), 1)
	if int(r1) < 0 {
		return fmt.Errorf("[semUnLock] syscall SYS_SETOP err: %s", err.Error())
	}
	return nil
}

func semSetVal(semId int, val int32) error {
	r1, _, err := syscall.Syscall6(sysSemCtl, uintptr(semId), 0,
		uintptr(SEM_SETVAL), uintptr(val), uintptr(0), uintptr(0))
	if int(r1) < 0 {
		return fmt.Errorf("[setVal] syscall SYS_SEMCTL err: %s", err.Error())
	}
	return nil
}

func semGetVal(semId int) (int, error) {
	r1, _, err := syscall.Syscall(sysSemCtl, uintptr(semId), 0, uintptr(SEM_GETVAL))
	if int(r1) < 0 {
		return -1, fmt.Errorf("[getVal] syscall SYS_GETCTL err: %s", err.Error())
	}
	return int(r1), nil
}

func semDestroy(semId int) error {
	r1, _, err := syscall.Syscall(sysSemCtl, uintptr(semId), 0, uintptr(IPC_RMID))
	if int(r1) < 0 {
		return fmt.Errorf("[destroy] syscall SYS_SEMCTL err: %s", err.Error())
	}
	return nil
}
