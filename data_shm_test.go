package shm_sdk_go

import "testing"

var shmKey uint32 = 0x5c9f
var maxNodeCount uint32 = 100

func TestShmData_GetHeader(t *testing.T) {
	var shmData ShmData
	err := shmData.Init(shmKey, maxNodeCount, true)
	if err != nil {
		t.Fatalf("Init failed, err: %s\n", err.Error())
	}
	header, err := shmData.GetHeader()
	if err != nil {
		t.Fatalf("Get Header failed, err: %s", err.Error())
	}
	t.Log(header)
}

func TestShmData_Insert(t *testing.T) {
	var shmData ShmData
	err := shmData.Init(shmKey, maxNodeCount, true)
	if err != nil {
		t.Fatalf("Init failed, err: %s\n", err.Error())
	}
	var nodes []ShmDataNode
	var i uint64 = 0
	for ; i < 10; i++ {
		nodes = append(nodes, ShmDataNode{i, i, i, i, i})
	}
	count, err := shmData.Insert(nodes)
	if err != nil {
		t.Fatalf("Insert failed, err: %s\n", err.Error())
	}
	t.Logf("insert count: %d\n", count)
}

func TestShmData_Traverse(t *testing.T) {
	var shmData ShmData
	err := shmData.Init(shmKey, maxNodeCount, true)
	if err != nil {
		t.Fatalf("Init failed, err: %s\n", err.Error())
	}

	err = shmData.Traverse(func(node *ShmDataNode) bool {
		t.Logf("traverse node, a: %d, b: %d, c: %d, d: %d, e: %d\n",
			node.a, node.b, node.c, node.d, node.e)
		return true
	})
	if err != nil {
		t.Logf("Traverse failed, err: %s\n", err.Error())
	}
}
