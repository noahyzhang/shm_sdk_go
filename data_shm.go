package shm_sdk_go

import (
	"fmt"
	"time"
	"unsafe"
)

const globalShmDataVersion uint32 = 0xFFFFFF01

type TraverseNodeFunc func(*ShmDataNode) bool

type ShmDataHeader struct {
	Version      uint32
	CurNodeCount uint32
	MaxNodeCount uint32
	HeaderCRCVal uint32
	TimeNs       uint64
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
	sdh.Version = other.Version
	sdh.CurNodeCount = other.CurNodeCount
	sdh.MaxNodeCount = other.MaxNodeCount
	sdh.HeaderCRCVal = other.HeaderCRCVal
	sdh.TimeNs = other.TimeNs
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
		return fmt.Errorf("[ShmData::Init] Already initialized, can't reinitialized")
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

	sd.localHeader.Version = globalShmDataVersion
	sd.localHeader.MaxNodeCount = maxNodeCount
	sd.localHeader.CurNodeCount = 0

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
			return fmt.Errorf("[ShmData::Init] parseHeader err: %s", err.Error())
		}
		sd.shmLength = length
		_ = sd.attach()
	} else {
		// 如果需要创建
		if err = sd.create(); err != nil {
			return fmt.Errorf("[ShmData::Init] create err: %s", err.Error())
		}
		if err = sd.setHeader(); err != nil {
			return fmt.Errorf("[ShmData::Init] setHeader err: %s", err.Error())
		}
	}
	sd.isInit = true
	return nil
}

func (sd *ShmData) Insert(arr []ShmDataNode) (uint32, error) {
	if sd.isInit == false {
		return 0, fmt.Errorf("[ShmData::Insert] init might be mistaken")
	}
	var curNodeCount = len(arr)
	if curNodeCount > int(sd.localHeader.MaxNodeCount) {
		arr = arr[:sd.localHeader.MaxNodeCount]
		curNodeCount = len(arr)
	}
	for i := 0; i < len(arr) && i < int(sd.localHeader.MaxNodeCount); i++ {
		node, err := sd.getNodeByPos(i)
		if node == nil || err != nil {
			return 0, fmt.Errorf("[ShmData::Insert] getNodeByPos err: %s", err.Error())
		}
		(*node).copy(arr[i])
	}
	sd.localHeader.CurNodeCount = uint32(curNodeCount)
	if err := sd.setHeader(); err != nil {
		return 0, fmt.Errorf("[ShmData::Insert] setHeader err: %s", err.Error())
	}
	return uint32(curNodeCount), nil
}

func (sd *ShmData) Traverse(eachFunc TraverseNodeFunc) error {
	if sd.isInit == false {
		return fmt.Errorf("[ShmData::Traverse] init might be mistaken")
	}
	header, err := sd.GetHeader()
	if err != nil {
		return fmt.Errorf("[ShmData::Traverse] GetHeader err: %s", err.Error())
	}
	if _, err = sd.parseHeader(header); err != nil {
		return fmt.Errorf("[ShmData::Traverse] parseHeader err: %s", err.Error())
	}
	var pNode *ShmDataNode = nil
	for i := 0; i < int(sd.localHeader.CurNodeCount); i++ {
		pNode, err = sd.getNodeByPos(i)
		if err != nil || pNode == nil {
			return fmt.Errorf("[ShmData::Traverse] Failed to get node")
		}
		if !eachFunc(pNode) {
			return fmt.Errorf("[ShmData::Traverse] callback TraverseNodeFunc function return false")
		}
	}
	return nil
}

func (sd *ShmData) GetHeader() (ShmDataHeader, error) {
	if sd.isInit == false {
		return ShmDataHeader{}, fmt.Errorf("[ShmData::GetHeader] init might be mistaken")
	}
	header, err := sd.doGetHeader()
	if err != nil {
		return ShmDataHeader{}, fmt.Errorf("[ShmData::GetHeader] doGetHeader err: %s", err.Error())
	}
	return header, nil
}

func (sd *ShmData) create() error {
	if !sd.isNeedCreate || sd.shmLength == 0 {
		return fmt.Errorf("[ShmData::create] isNeedCreate is false or shmLength is 0")
	}
	flag := 0666 | IPC_CREAT
	pShm, err := sd.getShmPointWithCreate(int(sd.shmKey), int(sd.shmLength), flag)
	if err != nil || pShm == nil {
		return fmt.Errorf("[ShmData::create] getShmPointWithCreate err: %s", err.Error())
	}
	sd.pShm = pShm
	sd.shmHeader = (*ShmDataHeader)(pShm)
	var shmBodyPoint = (uintptr)(pShm) + uintptr(sd.shmHeaderLength)
	dummy := &struct {
		addr uintptr
		len  int
		cap  int
	}{shmBodyPoint, int(sd.localHeader.MaxNodeCount), int(sd.localHeader.MaxNodeCount)}
	sd.shmBody = (*[]ShmDataNode)(unsafe.Pointer(dummy))

	sd.isAttach = true
	return nil
}

func (sd *ShmData) setHeader() error {
	if sd.localHeader.MaxNodeCount == 0 {
		return fmt.Errorf("[ShmData::setHeader] MaxNodeCount is 0")
	}
	sd.localHeader.HeaderCRCVal = 0
	sd.localHeader.TimeNs = uint64(time.Now().UnixNano())
	crc, _ := CalcCRCVal(sd.localHeader.convertIntegerArr(), sd.localHeader.getSize())
	sd.localHeader.HeaderCRCVal = crc
	sd.shmHeader.copy(sd.localHeader)
	return nil
}

func (sd *ShmData) parseHeader(shmHeader ShmDataHeader) (uint32, error) {
	if shmHeader.Version != globalShmDataVersion {
		return 0, fmt.Errorf("[ShmData::parseHeader] Version check error")
	}
	sd.localHeader.copy(shmHeader)
	sd.localHeader.HeaderCRCVal = 0
	crc, _ := CalcCRCVal(sd.localHeader.convertIntegerArr(), sd.localHeader.getSize())
	sd.localHeader.HeaderCRCVal = shmHeader.HeaderCRCVal
	if crc != shmHeader.HeaderCRCVal {
		return 0, fmt.Errorf("[ShmData::parseHeader] CRC calibration error")
	}
	return sd.localHeader.MaxNodeCount*sd.getNodeSize() + sd.getHeaderSize(), nil
}

func (sd *ShmData) attach() error {
	if sd.isNeedCreate || sd.shmLength == 0 {
		return fmt.Errorf("[ShmData::attach] shm not initialized")
	}
	pShm, err := sd.doAttach(sd.shmLength)
	if err != nil || pShm == nil {
		return fmt.Errorf("[ShmData::attach] doAttach err: %s", err.Error())
	}
	sd.pShm = pShm
	sd.shmHeader = (*ShmDataHeader)(pShm)
	// go 语言指针加操作，每次只移动一个字节
	var shmBodyPoint = (uintptr)(pShm) + uintptr(sd.shmHeaderLength)
	dummy := &struct {
		addr uintptr
		len  int
		cap  int
	}{shmBodyPoint, int(sd.localHeader.MaxNodeCount), int(sd.localHeader.MaxNodeCount)}
	sd.shmBody = (*[]ShmDataNode)(unsafe.Pointer(dummy))
	sd.isAttach = true
	return nil
}

func (sd *ShmData) detach() error {
	if !sd.isAttach {
		return fmt.Errorf("[ShmData::detach] not attach shm")
	}
	if err := sd.doDetach(sd.pShm); err != nil {
		return fmt.Errorf("[ShmData::detach] doDetach err: %s", err.Error())
	}
	return nil
}

func (sd *ShmData) doAttach(length uint32) (unsafe.Pointer, error) {
	var flag = 0666
	pShm, err := sd.getShmPoint(int(sd.shmKey), int(length), flag)
	if err != nil || pShm == nil {
		return nil, fmt.Errorf("[ShmData::doAttach] getShmPoint err: %s", err.Error())
	}
	return pShm, nil
}

func (sd *ShmData) doDetach(pShm unsafe.Pointer) error {
	if err := DetachShm(pShm); err != nil {
		return fmt.Errorf("[ShmData::doDetach] detach err: %s", err.Error())
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
		return nil, fmt.Errorf("[ShmData::getShmPoint] Param shmKey is 0")
	}
	shmId, err := GetShm(shmKey, shmSize, flag)
	if err != nil || shmId < 0 {
		return nil, fmt.Errorf("[ShmData::getShmPoint] GetShm err: %s", err.Error())
	}
	pShm, err := AttachShm(shmId, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("[ShmData::getShmPoint] AttachShm err: %s", err.Error())
	}
	return pShm, nil
}

func (sd *ShmData) getShmPointWithCreate(shmKey int, shmSize int, flag int) (unsafe.Pointer, error) {
	if shmKey == 0 {
		return nil, fmt.Errorf("[ShmData::getShmPointWithCreate] Param shmKey is 0")
	}
	pShm, err := sd.getShmPoint(shmKey, shmSize, flag&(^IPC_CREAT))
	if err != nil || pShm == nil {
		if (flag & IPC_CREAT) == 0 {
			return nil, fmt.Errorf("[ShmData::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
				flag&(^IPC_CREAT), err.Error())
		}
		pShm, err = sd.getShmPoint(shmKey, shmSize, flag)
		if err != nil || pShm == nil {
			return nil, fmt.Errorf("[ShmData::getShmPointWithCreate] getShmPoint(flag: %d) err: %s",
				flag, err.Error())
		}
		return pShm, nil
	}
	return pShm, nil
}
