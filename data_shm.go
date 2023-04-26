package shm_sdk_go

import (
	"fmt"
	"time"
	"unsafe"
)

const globalShmDataVersion uint32 = 0xFFFFFF01

type TraverseNodeFunc func(*ShmDataNode) bool

type ShmDataHeader struct {
	version      uint32
	curNodeCount uint32
	maxNodeCount uint32
	headerCRCVal uint32
	timeNs       uint64
}

func (sdh *ShmDataHeader) getSize() uint32 {
	return uint32(unsafe.Sizeof(ShmDataHeader{}))
}

func (sdh *ShmDataHeader) convertIntegerArr() []uint8 {
	headerSize := sdh.getSize()
	headerBytes := &struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(sdh)), int(headerSize), int(headerSize)}
	return *((*[]uint8)(unsafe.Pointer(headerBytes)))
}

func (sdh *ShmDataHeader) copy(other ShmDataHeader) {
	sdh.version = other.version
	sdh.curNodeCount = other.curNodeCount
	sdh.maxNodeCount = other.maxNodeCount
	sdh.headerCRCVal = other.headerCRCVal
	sdh.timeNs = other.timeNs
}

type ShmDataNode struct {
	Tid           uint32
	ArenaId       uint32
	AllocatedKB   uint32
	DeallocatedKB uint32
}

func (sdn *ShmDataNode) copy(other ShmDataNode) {
	sdn.Tid = other.Tid
	sdn.ArenaId = other.ArenaId
	sdn.AllocatedKB = other.AllocatedKB
	sdn.DeallocatedKB = other.DeallocatedKB
}

type ShmData struct {
	isInit        bool
	isNeedCreate  bool
	isAttach      bool
	isSetCallback bool

	shmKey          uint32
	shmLength       uint32
	shmHeaderLength uint32
	shmBodyLength   uint32

	localHeader ShmDataHeader

	pShm      unsafe.Pointer
	shmHeader *ShmDataHeader
	shmBody   *[]ShmDataNode
}

func (sd *ShmData) getHeaderSize() uint32 {
	return uint32(unsafe.Sizeof(ShmDataHeader{}))
}

func (sd *ShmData) getNodeSize() uint32 {
	return uint32(unsafe.Sizeof(ShmDataNode{}))
}

func (sd *ShmData) Init(shmKey uint32, maxNodeCount uint32, isCreate bool) error {
	if sd.isInit == true {
		return fmt.Errorf("[CArrayShm::init] Already initialized, can't reinitialized")
	}
	shmBodyLen := maxNodeCount * sd.getNodeSize()
	if isCreate == false {
		shmBodyLen = 0
	}
	sd.isNeedCreate = isCreate

	sd.shmKey = shmKey
	sd.shmHeaderLength = sd.getHeaderSize()
	sd.shmBodyLength = shmBodyLen
	sd.shmLength = sd.shmHeaderLength + sd.shmBodyLength

	sd.localHeader.version = globalShmDataVersion
	sd.localHeader.maxNodeCount = maxNodeCount
	sd.localHeader.curNodeCount = 0

	headerPoint, err := sd.doAttach(sd.shmHeaderLength)
	if err == nil && headerPoint != nil {
		// 挂载成功，说明已经创建了，无需再创建
		sd.isNeedCreate = false
	}
	// 如果无需创建
	if sd.isNeedCreate == false {
		length, err := sd.parseHeader(*(*ShmDataHeader)(headerPoint))
		_ = sd.doDetach(headerPoint)
		if err != nil || length == 0 {
			return fmt.Errorf("[ShmData::init] parseHeader err: %s", err.Error())
		}
		sd.shmLength = length
		_ = sd.attach()
	} else {
		// 如果需要创建
		if err = sd.create(); err != nil {
			return fmt.Errorf("[ShmData::init] create err: %s", err.Error())
		}
		if err = sd.setHeader(); err != nil {
			return fmt.Errorf("[ShmData::init] setHeader err: %s", err.Error())
		}
	}
	sd.isInit = true
	return nil
}

func (sd *ShmData) Insert(arr []ShmDataNode) (uint32, error) {
	if sd.isInit == false {
		return 0, fmt.Errorf("[CArrayShm::insert] init might be mistaken")
	}
	var curNodeCount = len(arr)
	if curNodeCount > int(sd.localHeader.maxNodeCount) {
		arr = arr[:sd.localHeader.maxNodeCount]
		curNodeCount = len(arr)
	}
	for i := 0; i < len(arr) && i < int(sd.localHeader.maxNodeCount); i++ {
		node, err := sd.getNodeByPos(i)
		if node == nil || err != nil {
			return 0, fmt.Errorf("[CArrayShm::insert] getNodeByPos err: %s", err.Error())
		}
		(*node).copy(arr[i])
	}
	sd.localHeader.curNodeCount = uint32(curNodeCount)
	if err := sd.setHeader(); err != nil {
		return 0, fmt.Errorf("[ShmData::Insert] setHeader err: %s", err.Error())
	}
	return uint32(curNodeCount), nil
}

func (sd *ShmData) Traverse(eachFunc TraverseNodeFunc) error {
	if sd.isInit == false {
		return fmt.Errorf("[CArrayShm::Traverse] init might be mistaken")
	}
	var pNode *ShmDataNode = nil
	var err error = nil
	for i := 0; i < int(sd.localHeader.curNodeCount); i++ {
		pNode, err = sd.getNodeByPos(i)
		if err != nil || pNode == nil {
			return fmt.Errorf("[CArrayShm::Traverse] Failed to get node")
		}
		if !eachFunc(pNode) {
			return fmt.Errorf("[CArrayShm::Traverse] callback TRAVERSE_METHOD function return false")
		}
	}
	return nil
}

func (sd *ShmData) GetHeader() (ShmDataHeader, error) {
	if sd.isInit == false {
		return ShmDataHeader{}, fmt.Errorf("[CArrayShm::Traverse] init might be mistaken")
	}
	header, err := sd.doGetHeader()
	if err != nil {
		return ShmDataHeader{}, fmt.Errorf("[ShmData::GetHeader] doGetHeader err: %s", err.Error())
	}
	return header, nil
}

func (sd *ShmData) create() error {
	if !sd.isNeedCreate || sd.shmLength == 0 {
		return fmt.Errorf("[CShm::create] isNeedCreate is false or shmLength is 0")
	}
	flag := 0666 | IPC_CREAT
	pShm, err := sd.getShmPointWithCreate(int(sd.shmKey), int(sd.shmLength), flag)
	if err != nil || pShm == nil {
		return fmt.Errorf("[CShm::create] getShmPointWithCreate err: %s", err.Error())
	}
	sd.pShm = pShm
	sd.shmHeader = (*ShmDataHeader)(pShm)
	var shmBodyPoint = (uintptr)(pShm) + uintptr(sd.shmHeaderLength)
	dummy := &struct {
		addr uintptr
		len  int
		cap  int
	}{shmBodyPoint, int(sd.localHeader.maxNodeCount), int(sd.localHeader.maxNodeCount)}
	sd.shmBody = (*[]ShmDataNode)(unsafe.Pointer(dummy))

	sd.isAttach = true
	return nil
}

func (sd *ShmData) setHeader() error {
	if sd.localHeader.maxNodeCount == 0 {
		return fmt.Errorf("[ShmData::setHeader] maxNodeCount is 0")
	}
	sd.localHeader.headerCRCVal = 0
	sd.localHeader.timeNs = uint64(time.Now().UnixNano())
	crc, _ := CalcCRCVal(sd.localHeader.convertIntegerArr(), sd.localHeader.getSize())
	sd.localHeader.headerCRCVal = crc
	sd.shmHeader.copy(sd.localHeader)
	return nil
}

func (sd *ShmData) parseHeader(shmHeader ShmDataHeader) (uint32, error) {
	if shmHeader.version != globalShmDataVersion {
		return 0, fmt.Errorf("[ShmData::parseHeader] version check error")
	}
	sd.localHeader.copy(shmHeader)
	sd.localHeader.headerCRCVal = 0
	crc, _ := CalcCRCVal(sd.localHeader.convertIntegerArr(), sd.localHeader.getSize())
	sd.localHeader.headerCRCVal = shmHeader.headerCRCVal
	if crc != shmHeader.headerCRCVal {
		return 0, fmt.Errorf("[ShmData::parseHeader] CRC calibration error")
	}
	return sd.localHeader.maxNodeCount*sd.getNodeSize() + sd.getHeaderSize(), nil
}

func (sd *ShmData) attach() error {
	if sd.isNeedCreate || sd.shmLength == 0 {
		return fmt.Errorf("[CShm::attach] shm not initialized")
	}
	pShm, err := sd.doAttach(sd.shmLength)
	if err != nil || pShm == nil {
		return fmt.Errorf("[CShm::attach] doAttach err: %s", err.Error())
	}
	sd.pShm = pShm
	sd.shmHeader = (*ShmDataHeader)(pShm)
	// go 语言指针加操作，每次只移动一个字节
	var shmBodyPoint = (uintptr)(pShm) + uintptr(sd.shmHeaderLength)
	dummy := &struct {
		addr uintptr
		len  int
		cap  int
	}{shmBodyPoint, int(sd.localHeader.maxNodeCount), int(sd.localHeader.maxNodeCount)}
	sd.shmBody = (*[]ShmDataNode)(unsafe.Pointer(dummy))
	sd.isAttach = true
	return nil
}

func (sd *ShmData) detach() error {
	if !sd.isAttach {
		return fmt.Errorf("[CShm::detach] not attach shm")
	}
	if err := sd.doDetach(sd.pShm); err != nil {
		return fmt.Errorf("[CShm::detach] doDetach err: %s", err.Error())
	}
	return nil
}

func (sd *ShmData) doAttach(length uint32) (unsafe.Pointer, error) {
	var flag = 0666
	pShm, err := sd.getShmPoint(int(sd.shmKey), int(length), flag)
	if err != nil || pShm == nil {
		return nil, fmt.Errorf("[CShm::doAttach] getShmPoint err: %s", err.Error())
	}
	return pShm, nil
}

func (sd *ShmData) doDetach(pShm unsafe.Pointer) error {
	if err := DetachShm(pShm); err != nil {
		return fmt.Errorf("[CShm::doDetach] detach err: %s", err.Error())
	}
	return nil
}

func (sd *ShmData) getNodeByPos(pos int) (*ShmDataNode, error) {
	if sd.isAttach == false {
		return nil, fmt.Errorf("[ShmData::getNodeByPos] isAttach is false")
	}
	return &(*sd.shmBody)[pos], nil
}

func (sd *ShmData) doGetHeader() (ShmDataHeader, error) {
	if sd.isAttach == false {
		return ShmDataHeader{}, fmt.Errorf("[ShmData::doGetHeader] isAttach is false")
	}
	var header ShmDataHeader
	header.copy(*sd.shmHeader)
	return header, nil
}

func (sd *ShmData) getShmPoint(shmKey int, shmSize int, flag int) (unsafe.Pointer, error) {
	if shmKey == 0 {
		return nil, fmt.Errorf("[CShm::getShmPoint] Param shmKey is 0")
	}
	shmId, err := GetShm(shmKey, shmSize, flag)
	if err != nil || shmId < 0 {
		return nil, fmt.Errorf("[CShm::getShmPoint] GetShm err: %s", err.Error())
	}
	pShm, err := AttachShm(shmId, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("[CShm::getShmPoint] AttachShm err: %s", err.Error())
	}
	return pShm, nil
}

func (sd *ShmData) getShmPointWithCreate(shmKey int, shmSize int, flag int) (unsafe.Pointer, error) {
	if shmKey == 0 {
		return nil, fmt.Errorf("[CShm::getShmPointWithCreate] Param shmKey is 0")
	}
	pShm, err := sd.getShmPoint(shmKey, shmSize, flag&(^IPC_CREAT))
	if err != nil || pShm == nil {
		if (flag & IPC_CREAT) == 0 {
			return nil, fmt.Errorf("[CShm::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
				flag&(^IPC_CREAT), err.Error())
		}
		pShm, err = sd.getShmPoint(shmKey, shmSize, flag)
		if err != nil || pShm == nil {
			return nil, fmt.Errorf("[CShm::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
				flag, err.Error())
		}
		return pShm, nil
	}
	return pShm, nil
}
