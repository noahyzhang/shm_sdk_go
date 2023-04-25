package shm_sdk_go

//type NodeType struct {
//	dummy01 int32
//	dummy02 uint32
//}
//
//func (nt NodeType) getSize() uint32 {
//	return uint32(unsafe.Sizeof(NodeType{}))
//}
//
//func (nt NodeType) copy(node ShmNode) error {
//	val, ok := node.(NodeType)
//	if !ok {
//		fmt.Errorf("[NodeType::copy] ShmNode not implement NodeType")
//	}
//	nt.dummy01 = val.dummy01
//	nt.dummy02 = val.dummy02
//	return nil
//}
//
//func TestCArrayShm_Insert(t *testing.T) {
//	var arrayShm CArrayShm
//	var shmKey uint32 = 0x5c8f
//	var maxNodeCount uint32 = 100
//	err := arrayShm.Init(shmKey, maxNodeCount, true)
//	if err != nil {
//		t.Fatalf("init err: %s", err.Error())
//	}
//	//var nodes []NodeType
//	//nodes = append(nodes, NodeType{dummy01: 10, dummy02: "hello"}, NodeType{dummy01: 20, dummy02: "world"})
//	//count, err := arrayShm.Insert(*(*[]ShmNode)(unsafe.Pointer(&nodes)))
//	//if err != nil {
//	//	t.Fatalf("insert err: %s", err.Error())
//	//}
//	//t.Logf("insert count: %d", count)
//
//	//arrayShm.Traverse(func(node *ShmNode) bool {
//	//	val, ok := (*node).(NodeType)
//	//	if !ok {
//	//		fmt.Errorf("[Traverse] ShmNode not implement NodeType")
//	//	}
//	//	t.Logf("traverse node, dummy01: %d, dummy02: %s", val.dummy01, val.dummy02)
//	//	return true
//	//})
//}
