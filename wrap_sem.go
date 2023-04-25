package shm_sdk_go

//type Semaphore struct {
//	semId int32
//}
//
//func (s *Semaphore) Create(semKey uint32) error {
//	s.semId = -1
//	res, _, err := syscall.Syscall(sysSemGet,
//		uintptr(int(semKey)), uintptr(int(1)), uintptr(int(0666)))
//	if res == -1 {
//		res, _, err = syscall.Syscall(sysSemGet,
//			uintptr(int(semKey)), uintptr(int(1)), uintptr(int(0666|IPC_CREAT|IPC_EXCL)))
//		if res == -1 {
//			return fmt.Errorf("[Semaphore::Create] call semGet err: %s", err.Error())
//		}
//	}
//	res, _, err = syscall.Syscall(sysSemCtl, uintptr(int(1)), uintptr(int(sysSemSetVal)), uintptr(int(1)))
//	if res == -1 {
//		return fmt.Errorf("[Semaphore::Create] call semCtl err: %s", err.Error())
//	}
//	return nil
//}
//
//func (s *Semaphore) Lock(isWait bool) error {
//	if s.semId == -1 {
//		return fmt.Errorf("[Semaphore::Lock] semId is -1")
//	}
//	type semBuf struct {
//		semNum  uint16
//		semOp   int16
//		semFlag int16
//	}
//	var buf semBuf
//	if isWait {
//		buf.semNum = 0
//		buf.semOp = -1
//		buf.semFlag = sysSemUNDO
//	} else {
//		buf.semNum = 0
//		buf.semOp = -1
//		buf.semFlag = IPC_NOWAIT | sysSemUNDO
//	}
//	res, _, err := syscall.Syscall(sysSemOp, uintptr(unsafe.Pointer(&buf)), uintptr(1), 0)
//	if res == -1 {
//		return fmt.Errorf("[Semaphore::Lock] call SemOp err: %s", err.Error())
//	}
//	return nil
//}
//
//func (s *Semaphore) Unlock() error {
//	if s.semId == -1 {
//		return fmt.Errorf("[Semaphore::Unlock] semId is -1")
//	}
//	buf := &struct {
//		semNum  uint16
//		semOp   int16
//		semFlag int16
//	}{0, 1, sysSemUNDO}
//	res, _, err := syscall.Syscall(sysSemOp, uintptr(unsafe.Pointer(buf)), uintptr(1), 0)
//	if res == -1 {
//		return fmt.Errorf("[Semaphore::Unlock] call semOp err: %s", err.Error())
//	}
//	return nil
//}
//
//func (s *Semaphore) Destroy() error {
//	res, _, err := syscall.Syscall(sysSemCtl, uintptr(s.semId), uintptr(0), uintptr(IPC_CREAT))
//	if res == -1 {
//		return fmt.Errorf("[Semaphore::Destroy] call semCtl err: %s", err.Error())
//	}
//	return nil
//}
