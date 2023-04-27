## 共享内存操作的 SDK

### 一、介绍共享内存

目前实现共享内存的操作，共享内存中的数据格式如下

```go
type ShmDataHeader struct {
	version      uint32
	curNodeCount uint32
	maxNodeCount uint32
	headerCRCVal uint32
	timeNs       uint64
}

type ShmDataNode struct {
    Tid           uint32
    ArenaId       uint32
    AllocatedKB   uint32
    DeallocatedKB uint32
}
```

如下为共享内存中数据的格式：
| ShmDataHeader | ShmDataNode | ShmDataNode | ... | ShmDataNode |

### 二、介绍信号量

为了同步多进程操作共享内存，采用信号量来同步。
使用 cgo 实现了封装信号量

### 三、如何使用

首先 `go get github.com/noahyzhang/shm_sdk_go` 

然后便可方便的操作共享内存了，可以使用信号量同步多进程

```go
package main

import (
	"fmt"
	"github.com/noahyzhang/shm_sdk_go"
	"time"
)

var SHM_KEY uint32 = 0x5c9f
var MAX_SHM_ARR_COUNT uint32 = 500
var SEM_KEY uint32 = 0xcc9f

func main() {
	var shm shm_sdk_go.ShmData
	if err := shm.Init(SHM_KEY, MAX_SHM_ARR_COUNT, true); err != nil {
		fmt.Printf("shm Init err: %s\n", err.Error())
		return
	}
	var sem shm_sdk_go.Semaphore
	if err := sem.Create(SEM_KEY); err != nil {
		fmt.Printf("sem Create err: %s\n", err.Error())
		return
	}

	for {
		// 非阻塞
		if err := sem.Lock(true); err != nil {
			fmt.Printf("sem Lock err: %s\n", err.Error())
			time.Sleep(time.Second)
			continue
		}
		header, err := shm.GetHeader()
		if err != nil {
			_ = sem.Unlock()
			fmt.Printf("shm GetHeader err: %s\n", err.Error())
			continue
		}
		fmt.Printf("header: version: %d, cur_node_count: %d, max_node_count: %d, time_ns: %v, crc: %d\n",
			header.Version, header.CurNodeCount, header.MaxNodeCount, header.TimeNs, header.HeaderCRCVal)
		err = shm.Traverse(func(node *shm_sdk_go.ShmDataNode) bool {
			fmt.Println("node: ", node.Tid, node.ArenaId, node.AllocatedKB, node.DeallocatedKB)
			return true
		})
		if err != nil {
			_ = sem.Unlock()
			fmt.Printf("shm Traverse err: %s\n", err.Error())
			return
		}
		fmt.Println()
		if err = sem.Unlock(); err != nil {
			fmt.Printf("sem Unlock err: %s\n", err.Error())
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
	}
}
```

### 四、注意

1. 注意时间，共享内存中的头部中有"时间字段"，请判断时间，如果和当前时间相差太远，就说明这个数据是过期数据，需要酬情处理