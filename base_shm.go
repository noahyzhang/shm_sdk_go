package shm_sdk_go

//type ShmHeader interface {
//	getSize() uint32
//	copy(ShmHeader) error
//	convertIntegerArr() []uint8
//}
//
//type ShmNode interface {
//	getSize() uint32
//	copy(ShmNode) error
//}
//
//type TraverseMethodFunc func(*ShmNode) bool
//
////type Shm interface {
////	setHeader() bool
////	parseHeader(*ShmHeader) uint32
////	traverse(*TraverseMethodFunc) bool
////}
//
//type CShm struct {
//	//ShmNode
//	isInit        bool
//	isNeedCreate  bool
//	isAttach      bool
//	isSetCallback bool
//
//	shmKey          uint32
//	shmLength       uint32
//	shmHeaderLength uint32
//	shmBodyLength   uint32
//
//	pShm      unsafe.Pointer
//	shmHeader *ShmHeader
//	shmBody   *[]ShmNode
//}
//
//func (cs *CShm) init(shmKey uint32, shmHeaderLength uint32, shmBodyLength uint32, isNeedCreate bool) error {
//	if !isNeedCreate {
//		shmBodyLength = 0
//	}
//	cs.shmKey = shmKey
//	cs.isNeedCreate = isNeedCreate
//	cs.shmHeaderLength = shmHeaderLength
//	cs.shmBodyLength = shmBodyLength
//	cs.shmLength = cs.shmHeaderLength + cs.shmBodyLength
//
//	cs.isInit = false
//	if cs.isNeedCreate && cs.shmBodyLength == 0 {
//		return fmt.Errorf("[CShm::init] The specifying length is invalid (==0) when creating SHM")
//	}
//	err := cs.attach()
//	if err == nil {
//		cs.isNeedCreate = false
//	}
//	if cs.isNeedCreate {
//		if err = cs.create(); err != nil {
//			return fmt.Errorf("[CShm::init] create shm err: %s", err.Error())
//		}
//	}
//	cs.isInit = true
//	return nil
//}
//
//func (cs *CShm) create() error {
//	if !cs.isNeedCreate || cs.shmLength == 0 {
//		return fmt.Errorf("[CShm::create] isNeedCreate is false or shmLength is 0")
//	}
//	flag := 0666 | IPC_CREAT
//	fmt.Printf("create shm, key: %d, shm-length: %d, flag: %d\n", cs.shmKey, cs.shmLength, flag)
//	pShm, err := cs.getShmPointWithCreate(int(cs.shmKey), int(cs.shmLength), flag)
//	if err != nil {
//		return fmt.Errorf("[CShm::create] getShmPointWithCreate err: %s", err.Error())
//	}
//	cs.pShm = pShm
//	cs.shmHeader = (*ShmHeader)(pShm)
//
//	var shmBodyPoint = (uintptr)(pShm) + uintptr(cs.shmHeaderLength)
//	var nodeCount = cs.shmBodyLength / cs.getNodeSize()
//	dummy := &struct {
//		addr uintptr
//		len  int
//		cap  int
//	}{shmBodyPoint, int(nodeCount), int(nodeCount)}
//	cs.shmBody = (*[]ShmNode)(unsafe.Pointer(dummy))
//
//	cs.isAttach = true
//	return nil
//}
//
//func (cs *CShm) insertNodeArr(nodeArr []ShmNode) error {
//	copy(*cs.shmBody, nodeArr)
//	return nil
//}
//
//func (cs *CShm) getHeaderSize() uint32 {
//	return (*cs.shmHeader).getSize()
//}
//
//func (cs *CShm) getNodeSize() uint32 {
//	return uint32(unsafe.Sizeof(cs.ShmNode.getSize()))
//}
//
//func (cs *CShm) attach() error {
//	if cs.isNeedCreate || cs.shmLength == 0 {
//		return fmt.Errorf("[CShm::attach] shm not initialized")
//	}
//	pShm, err := cs.doAttach(cs.shmLength)
//	if err != nil || pShm == nil {
//		return fmt.Errorf("[CShm::attach] doAttach err: %s", err.Error())
//	}
//	cs.pShm = pShm
//	cs.shmHeader = (*ShmHeader)(pShm)
//	var shmBodyPoint = (uintptr)(pShm) + uintptr(cs.shmHeaderLength)
//	var nodeCount = cs.shmBodyLength / cs.getNodeSize()
//	dummy := &struct {
//		addr uintptr
//		len  int
//		cap  int
//	}{shmBodyPoint, int(nodeCount), int(nodeCount)}
//	cs.shmBody = (*[]ShmNode)(unsafe.Pointer(dummy))
//	cs.isAttach = true
//	return nil
//}
//
//func (cs *CShm) detach() error {
//	if !cs.isAttach {
//		return fmt.Errorf("[CShm::detach] not attach shm")
//	}
//	if err := cs.doDetach(cs.pShm); err != nil {
//		return fmt.Errorf("[CShm::detach] doDetach err: %s", err.Error())
//	}
//	return nil
//}
//
//func (cs *CShm) doAttach(length uint32) (unsafe.Pointer, error) {
//	var flag = 0666
//	pShm, err := cs.getShmPoint(int(cs.shmKey), int(length), flag)
//	if err != nil || pShm == nil {
//		return nil, fmt.Errorf("[CShm::doAttach] getShmPoint err: %s", err.Error())
//	}
//	return pShm, nil
//}
//
//func (cs *CShm) doDetach(pShm unsafe.Pointer) error {
//	if err := DetachShm(pShm); err != nil {
//		return fmt.Errorf("[CShm::doDetach] detach err: %s", err.Error())
//	}
//	return nil
//}
//
//func (cs *CShm) doSetHeader(header ShmHeader, headerSize uint32) error {
//	if false == cs.isAttach {
//		return fmt.Errorf("[CShm::doSetHeader] shm is not attach")
//	}
//	src := &struct {
//		addr uintptr
//		len  int
//		cap  int
//	}{uintptr(unsafe.Pointer(&header)), int(headerSize), int(headerSize)}
//	dst := &struct {
//		addr uintptr
//		len  int
//		cap  int
//	}{uintptr(unsafe.Pointer(cs.shmHeader)), int(headerSize), int(headerSize)}
//
//	copy(*((*[]uint8)(unsafe.Pointer(dst))), *((*[]uint8)(unsafe.Pointer(src))))
//	//if err := (*cs.shmHeader).copy(header); err != nil {
//	//	return fmt.Errorf("[CShm::doSetHeader] copy err: %s", err.Error())
//	//}
//	return nil
//}
//
//func (cs *CShm) doGetHeader() (ShmHeader, error) {
//	if false == cs.isAttach {
//		return nil, fmt.Errorf("[CShm::doGetHeader] shm is not attach")
//	}
//	var header ShmHeader
//	if err := header.copy(*cs.shmHeader); err != nil {
//		return nil, fmt.Errorf("[CShm::doGetHeader] copy err: %s", err.Error())
//	}
//	return header, nil
//}
//
//func (cs *CShm) getNodeByPos(pos int) (*ShmNode, error) {
//	if false == cs.isAttach {
//		return nil, fmt.Errorf("[CShm::getNodeByPos] shm is not attach")
//	}
//	return &((*cs.shmBody)[pos]), nil
//}
//
//func (cs *CShm) getShmPoint(shmKey int, shmSize int, flag int) (unsafe.Pointer, error) {
//	if shmKey == 0 {
//		return nil, fmt.Errorf("[CShm::getShmPoint] Param shmKey is 0")
//	}
//	shmId, err := GetShm(shmKey, shmSize, flag)
//	if err != nil || shmId < 0 {
//		return nil, fmt.Errorf("[CShm::getShmPoint] GetShm err: %s", err.Error())
//	}
//	pShm, err := AttachShm(shmId, 0, 0)
//	if err != nil {
//		return nil, fmt.Errorf("[CShm::getShmPoint] AttachShm err: %s", err.Error())
//	}
//	return pShm, nil
//}
//
//func (cs *CShm) getShmPointWithCreate(shmKey int, shmSize int, flag int) (unsafe.Pointer, error) {
//	if shmKey == 0 {
//		return nil, fmt.Errorf("[CShm::getShmPointWithCreate] Param shmKey is 0")
//	}
//	pShm, err := cs.getShmPoint(shmKey, shmSize, flag&(^IPC_CREAT))
//	if err != nil || pShm == nil {
//		if (flag & IPC_CREAT) == 0 {
//			return nil, fmt.Errorf("[CShm::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
//				flag&(^IPC_CREAT), err.Error())
//		}
//		pShm, err = cs.getShmPoint(shmKey, shmSize, flag)
//		if err != nil || pShm == nil {
//			return nil, fmt.Errorf("[CShm::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
//				flag, err.Error())
//		}
//		return pShm, nil
//	}
//	return pShm, nil
//}
