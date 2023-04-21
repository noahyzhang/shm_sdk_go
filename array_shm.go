package shm_sdk_go

import (
	"fmt"
	"time"
	"unsafe"
)

const globalShmVersion uint32 = 0xFFFFFF01

type ArrayShmHeader struct {
	version      uint32
	curNodeCount uint32
	maxNodeCount uint32
	headerCRCVal uint32
	timeNs       uint64
}

type CArrayShm struct {
	CShm
	isInit      bool
	arrayHeader ArrayShmHeader
}

func (cas *CArrayShm) Init(shmKey uint32, maxNodeCount uint32, isCreate bool) error {
	if cas.isInit == false {
		return fmt.Errorf("[CArrayShm::init] Already initialized, can't reinitialized")
	}
	cas.arrayHeader.version = globalShmVersion
	cas.arrayHeader.maxNodeCount = maxNodeCount
	cas.arrayHeader.curNodeCount = 0
	err := cas.init(shmKey, maxNodeCount*cas.getNodeSize(), isCreate)
	if err != nil {
		return err
	}
	cas.isInit = true
	return nil
}

func (cas *CArrayShm) Insert(arr []ShmNode) (uint32, error) {
	if cas.isInit == false {
		return 0, fmt.Errorf("[CArrayShm::insert] init might be mistaken")
	}
	var curNodeCount uint32 = 0
	for i := 0; i < len(arr) && i < int(cas.arrayHeader.maxNodeCount); i++ {
		node, err := cas.getNodeByPos(i)
		if err == nil {
			continue
		}
		copy([]byte(node), arr[i])
		curNodeCount++
	}
	cas.arrayHeader.curNodeCount = curNodeCount
	cas.setHeader()
	return curNodeCount, nil
}

func (cas *CArrayShm) setHeader() error {
	if cas.arrayHeader.maxNodeCount == 0 {
		return fmt.Errorf("[CArrayShm::setHeader] input maxNodeCount invalid")
	}
	cas.arrayHeader.headerCRCVal = 0
	cas.arrayHeader.timeNs = uint64(time.Now().UnixNano())
	headerSize := uint32(unsafe.Sizeof(ArrayShmHeader{}))
	headerBytes := &struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(unsafe.Pointer(&cas.arrayHeader)), int(headerSize), int(headerSize)}
	crc, _ := CalcCRCVal(*((*[]uint8)(unsafe.Pointer(headerBytes))), headerSize)
	cas.arrayHeader.headerCRCVal = crc
	return cas.doSetHeader(cas.arrayHeader)
}

func (cas *CArrayShm) parseHeader(header *ArrayShmHeader) (uint32, error) {
	var version = header.version
	if version != globalShmVersion {
		return 0, fmt.Errorf("[CArrayShm::parseHeader] version check error")
	}
	copy(cas.arrayHeader, header)
}

func (cas *CArrayShm) Traverse(eachFunc TraverseMethodFunc) error {

}

func (cas *CArrayShm) GetHeader() (ArrayShmHeader, error) {

}
