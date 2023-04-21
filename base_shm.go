package shm_sdk_go

type ShmHeader interface {
}

type ShmNode interface {
}

type TraverseMethodFunc func(*ShmNode) bool

type Shm interface {
	setHeader() bool
	parseHeader(*ShmHeader) uint32
	traverse(*TraverseMethodFunc) bool
}

type CShm struct {
}

func (cs *CShm) init(shmKey uint32, shmBodySize uint32, isCreate bool) error {
}

func (cs *CShm) getNodeByPos(pos int) (*ShmNode, error) {

}

func (cs *CShm) doSetHeader(header ShmHeader) error {

}

func (cs *CShm) doGetHeader(header *ShmHeader) error {

}

func (cs *CShm) getHeaderSize() uint32 {
}

func (cs *CShm) getNodeSize() uint32 {

}
