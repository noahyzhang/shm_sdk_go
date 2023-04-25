package shm_sdk_go

//const globalShmVersion uint32 = 0xFFFFFF01
//
//type ArrayShmHeader struct {
//	version      uint32
//	curNodeCount uint32
//	maxNodeCount uint32
//	headerCRCVal uint32
//	timeNs       uint64
//}
//
//func (ash ArrayShmHeader) getSize() uint32 {
//	return uint32(unsafe.Sizeof(ArrayShmHeader{}))
//}
//
//func (ash ArrayShmHeader) convertIntegerArr() []uint8 {
//	headerSize := ash.getSize()
//	headerBytes := &struct {
//		addr uintptr
//		len  int
//		cap  int
//	}{uintptr(unsafe.Pointer(&ash)), int(headerSize), int(headerSize)}
//	return *((*[]uint8)(unsafe.Pointer(headerBytes)))
//}
//
//func (ash ArrayShmHeader) copy(header ShmHeader) error {
//	arrayHeader, ok := header.(ArrayShmHeader)
//	if !ok {
//		return fmt.Errorf("[ArrayShmHeader::copy] param ShmHeader type is not ArrayShmHeader")
//	}
//	ash.version = arrayHeader.version
//	ash.curNodeCount = arrayHeader.curNodeCount
//	ash.maxNodeCount = arrayHeader.maxNodeCount
//	ash.headerCRCVal = arrayHeader.headerCRCVal
//	ash.timeNs = arrayHeader.timeNs
//	return nil
//}
//
//type CArrayShm struct {
//	CShm
//	isInit      bool
//	arrayHeader ArrayShmHeader
//}
//
//func (cas *CArrayShm) Init(shmKey uint32, nodeType ShmNode, maxNodeCount uint32, isCreate bool) error {
//	if cas.isInit == true {
//		return fmt.Errorf("[CArrayShm::init] Already initialized, can't reinitialized")
//	}
//	cas.arrayHeader.version = globalShmVersion
//	cas.arrayHeader.maxNodeCount = maxNodeCount
//	cas.arrayHeader.curNodeCount = 0
//	var nodeSize = nodeType.getSize()
//	err := cas.init(shmKey, cas.arrayHeader.getSize(), maxNodeCount*nodeSize, isCreate)
//	if err != nil {
//		return err
//	}
//	cas.isInit = true
//	return nil
//}
//
//func (cas *CArrayShm) Insert(arr []ShmNode) (uint32, error) {
//	if cas.isInit == false {
//		return 0, fmt.Errorf("[CArrayShm::insert] init might be mistaken")
//	}
//	var curNodeCount = len(arr)
//	if curNodeCount > int(cas.arrayHeader.maxNodeCount) {
//		arr = arr[:cas.arrayHeader.maxNodeCount]
//		curNodeCount = len(arr)
//	}
//
//	cas.insertNodeArr(arr)
//	//for i := 0; i < len(arr) && i < int(cas.arrayHeader.maxNodeCount); i++ {
//	//	node, err := cas.getNodeByPos(i)
//	//	if node == nil || err != nil {
//	//		return 0, err
//	//	}
//	//	err = (*node).copy(arr[i])
//	//	if err != nil {
//	//		return 0, err
//	//	}
//	//	curNodeCount++
//	//}
//
//	cas.arrayHeader.curNodeCount = uint32(curNodeCount)
//	if err := cas.setHeader(); err != nil {
//		return 0, err
//	}
//	return uint32(curNodeCount), nil
//}
//
//func (cas *CArrayShm) setHeader() error {
//	if cas.arrayHeader.maxNodeCount == 0 {
//		return fmt.Errorf("[CArrayShm::setHeader] input maxNodeCount invalid")
//	}
//	cas.arrayHeader.headerCRCVal = 0
//	cas.arrayHeader.timeNs = uint64(time.Now().UnixNano())
//	crc, _ := CalcCRCVal(cas.arrayHeader.convertIntegerArr(), cas.arrayHeader.getSize())
//	cas.arrayHeader.headerCRCVal = crc
//	return cas.doSetHeader(cas.arrayHeader, cas.arrayHeader.getSize())
//}
//
//func (cas *CArrayShm) parseHeader(header ArrayShmHeader) (uint32, error) {
//	var version = header.version
//	if version != globalShmVersion {
//		return 0, fmt.Errorf("[CArrayShm::parseHeader] version check error")
//	}
//	err := cas.arrayHeader.copy(header)
//	if err != nil {
//		return 0, fmt.Errorf("[CArrayShm::parseHeader] Failed to copy, err: %s", err.Error())
//	}
//	cas.arrayHeader.headerCRCVal = 0
//	crc, _ := CalcCRCVal(cas.arrayHeader.convertIntegerArr(), cas.arrayHeader.getSize())
//	cas.arrayHeader.headerCRCVal = header.headerCRCVal
//	if crc != header.headerCRCVal {
//		return 0, fmt.Errorf("[CArrayShm::parseHeader] CRC calibration error")
//	}
//	return cas.arrayHeader.maxNodeCount*cas.getNodeSize() + cas.arrayHeader.getSize(), nil
//}
//
//func (cas *CArrayShm) Traverse(eachFunc TraverseMethodFunc) error {
//	if cas.isInit == false {
//		return fmt.Errorf("[CArrayShm::Traverse] init might be mistaken")
//	}
//	var pNode *ShmNode = nil
//	var err error = nil
//	for i := 0; i < int(cas.arrayHeader.curNodeCount); i++ {
//		pNode, err = cas.getNodeByPos(i)
//		if err != nil || pNode == nil {
//			return fmt.Errorf("[CArrayShm::Traverse] Failed to get node")
//		}
//		if !eachFunc(pNode) {
//			return fmt.Errorf("[CArrayShm::Traverse] callback TRAVERSE_METHOD function return false")
//		}
//	}
//	return nil
//}
//
//func (cas *CArrayShm) GetHeader() (ArrayShmHeader, error) {
//	if cas.isInit == false {
//		return ArrayShmHeader{}, fmt.Errorf("[CArrayShm::Traverse] init might be mistaken")
//	}
//	var header ShmHeader
//	var err error
//	header, err = cas.doGetHeader()
//	if err != nil {
//		return ArrayShmHeader{}, nil
//	}
//	arrayHeader, ok := header.(ArrayShmHeader)
//	if !ok {
//		return ArrayShmHeader{}, fmt.Errorf("[CArrayShm::Traverse] elem type is not ArrayShmHeader")
//	}
//	return arrayHeader, nil
//}
