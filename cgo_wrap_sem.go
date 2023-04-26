package shm_sdk_go

/*
#include <stdlib.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/sem.h>
#include <sys/msg.h>
int semget(key_t key, int nsems, int semflg);
int semtimedop(int semid, struct sembuf *sops, size_t nsops, const struct timespec *timeout);
union arg4 {
	int             val;
	struct semid_ds *buf;
	unsigned short  *array;
};
int semctl_noarg(int semid, int semnum, int cmd) {
	return semctl(semid, semnum, cmd);
};
int semctl_buf(int semid, int cmd, struct semid_ds *buf) {
	union arg4 arg;
	arg.buf = buf;
	return semctl(semid, 0, cmd, arg);
};
int semctl_arr(int semid, int cmd, unsigned short *arr) {
	union arg4 arg;
	arg.array = arr;
	return semctl(semid, 0, cmd, arg);
};
int semctl_val(int semid, int semnum, int cmd, int value) {
	union arg4 arg;
	arg.val = value;
	return semctl(semid, semnum, cmd, arg);
};
*/
import "C"
import (
	"fmt"
	"time"
)

// sem_id 和 信号量个数
type semaphoreSet struct {
	id    int64
	count uint
}

// 创建或者获取信号量
func getSemSet(key, count int64, flags *semSetFlags) (*semaphoreSet, error) {
	rc, err := C.semget(C.key_t(key), C.int(count), C.int(flags.flags()))
	if rc == -1 {
		return nil, fmt.Errorf("[getSemSet] semget err: %v", err)
	}
	return &semaphoreSet{int64(rc), uint(count)}, nil
}

// 对信号量进行操作
func (ss *semaphoreSet) run(ops *semOps, timeout time.Duration) error {
	var cto *C.struct_timespec = nil
	if timeout >= 0 {
		cto = &C.struct_timespec{
			tv_sec:  C.__time_t(timeout / time.Second),
			tv_nsec: C.__syscall_slong_t(timeout % time.Second),
		}
	}

	var opptr *C.struct_sembuf
	if len(*ops) > 0 {
		opptr = &(*ops)[0]
	}

	rc, err := C.semtimedop(C.int(ss.id), opptr, C.size_t(len(*ops)), cto)
	if rc == -1 {
		return err
	}
	return nil
}

// 获取某个信号量的值
func (ss *semaphoreSet) getVal(num uint16) (int, error) {
	val, err := C.semctl_noarg(C.int(ss.id), C.int(num), C.GETVAL)
	if val == -1 {
		return -1, fmt.Errorf("[semaphoreSet::getVal] semctl err: %v", err)
	}
	return int(val), nil
}

// 设置某个信号量的值
func (ss *semaphoreSet) setVal(num uint16, value int) error {
	val, err := C.semctl_val(C.int(ss.id), C.int(num), C.SETVAL, C.int(value))
	if val == -1 {
		return fmt.Errorf("[semaphoreSet::setVal] semctl err: %v", err)
	}
	return nil
}

// 获取信号量集中的所有信号量的值
func (ss *semaphoreSet) getAll() ([]uint16, error) {
	carr := make([]C.ushort, ss.count)

	rc, err := C.semctl_arr(C.int(ss.id), C.GETALL, &carr[0])
	if rc == -1 {
		return nil, fmt.Errorf("[semaphoreSet::getAll] semctl_arr err: %v", err)
	}

	results := make([]uint16, ss.count)
	for i, ci := range carr {
		results[i] = uint16(ci)
	}
	return results, nil
}

// 设置信号量集中的所有信号量的值
func (ss *semaphoreSet) setAll(values []uint16) error {
	if uint(len(values)) != ss.count {
		return fmt.Errorf("[semaphoreSet::setAll] wrong number of values for setAll")
	}

	carr := make([]C.ushort, ss.count)
	for i, val := range values {
		carr[i] = C.ushort(val)
	}

	rc, err := C.semctl_arr(C.int(ss.id), C.SETALL, &carr[0])
	if rc == -1 {
		return fmt.Errorf("[semaphoreSet::setAll] semctl_arr err: %v", err)
	}
	return nil
}

// 获取最后一个操作某个信号量的进程 ID
func (ss *semaphoreSet) getPid(num uint16) (int, error) {
	rc, err := C.semctl_noarg(C.int(ss.id), C.int(num), C.GETPID)
	if rc == -1 {
		return 0, fmt.Errorf("[semaphoreSet::getPid] semctl err: %v", err)
	}
	return int(rc), nil
}

// 获取当前等待该信号量的值增长的进程数
func (ss *semaphoreSet) getNCnt(num uint16) (int, error) {
	rc, err := C.semctl_noarg(C.int(ss.id), C.int(num), C.GETNCNT)
	if rc == -1 {
		return 0, err
	}
	return int(rc), nil
}

// 获取当前等待该信号量的值变为 0 的进程数
func (ss *semaphoreSet) getZCnt(num uint16) (int, error) {
	rc, err := C.semctl_noarg(C.int(ss.id), C.int(num), C.GETZCNT)
	if rc == -1 {
		return 0, fmt.Errorf("[semaphoreSet::getZCnt] semctl err: %v", err)
	}
	return int(rc), nil
}

// 后去信号量集的信息
func (ss *semaphoreSet) stat() (*semSetInfo, error) {
	sds := C.struct_semid_ds{}

	rc, err := C.semctl_buf(C.int(ss.id), C.IPC_STAT, &sds)
	if rc == -1 {
		return nil, fmt.Errorf("[semaphoreSet::stat] semctl err: %v", err)
	}

	ssinf := semSetInfo{
		Perms: Perm{
			Uid:  uint32(sds.sem_perm.uid),
			Gid:  uint32(sds.sem_perm.gid),
			Cuid: uint32(sds.sem_perm.cuid),
			Cgid: uint32(sds.sem_perm.cgid),
			Mode: uint16(sds.sem_perm.mode),
		},
		LastOp:     time.Unix(int64(sds.sem_otime), 0),
		LastChange: time.Unix(int64(sds.sem_ctime), 0),
		Count:      uint(sds.sem_nsems),
	}
	return &ssinf, nil
}

// 设置信号量集的信息
func (ss *semaphoreSet) set(ssi *semSetInfo) error {
	sds := &C.struct_semid_ds{
		sem_perm: C.struct_ipc_perm{
			uid:  C.__uid_t(ssi.Perms.Uid),
			gid:  C.__gid_t(ssi.Perms.Gid),
			mode: C.ushort(ssi.Perms.Mode & 0x1FF),
		},
	}

	rc, err := C.semctl_buf(C.int(ss.id), C.IPC_SET, sds)
	if rc == -1 {
		return fmt.Errorf("[semaphoreSet::set] semctl err: %v", err)
	}
	return nil
}

// 删除信号量集，将会唤醒所有被阻塞的进程
func (ss *semaphoreSet) remove() error {
	rc, err := C.semctl_noarg(C.int(ss.id), 0, C.IPC_RMID)
	if rc == -1 {
		return fmt.Errorf("[semaphoreSet::remove] semctl err: %v", err)
	}
	return nil
}

// 用于操作信号量的集合
type semOps []C.struct_sembuf

// 创建一个操作信号量的集合
func newSemOps() *semOps {
	sops := semOps(make([]C.struct_sembuf, 0))
	return &sops
}

// 信号量操作：增加信号量值
func (so *semOps) increment(num uint16, by int16, flags *semOpFlags) error {
	if by < 0 {
		return fmt.Errorf("[semOps::increment] param by must be > 0, use desrement")
	} else if by == 0 {
		return fmt.Errorf("[semOps::increment] param by must be > 0, use waitZero")
	}

	*so = append(*so, C.struct_sembuf{
		sem_num: C.ushort(num),
		sem_op:  C.short(by),
		sem_flg: C.short(flags.flags()),
	})
	return nil
}

// 信号量操作：一直阻塞直到信号量的值为 0
func (so *semOps) waitZero(num uint16, flags *semOpFlags) error {
	*so = append(*so, C.struct_sembuf{
		sem_num: C.ushort(num),
		sem_op:  C.short(0),
		sem_flg: C.short(flags.flags()),
	})
	return nil
}

// 信号量操作：减小信号量值
func (so *semOps) decrement(num uint16, by int16, flags *semOpFlags) error {
	if by <= 0 {
		return fmt.Errorf("[semOps::decrement] param by must be > 0, use waitZero or increment")
	}

	*so = append(*so, C.struct_sembuf{
		sem_num: C.ushort(num),
		sem_op:  C.short(-by),
		sem_flg: C.short(flags.flags()),
	})
	return nil
}

// 信号量集的元数据
type semSetInfo struct {
	Perms      Perm
	LastOp     time.Time
	LastChange time.Time
	Count      uint
}

// 信号量获取的标识位
type semSetFlags struct {
	// 是否创建不存在的信号量集
	Create bool

	// 如果信号量集存在，则使用 Exclusive 会失败
	// 仅仅在和 Create 一起使用时有效
	Exclusive bool

	// 信号量集权限，只有和 Create 一起使用时有效
	Perms int
}

// 获取标识位
func (sf *semSetFlags) flags() int64 {
	if sf == nil {
		return 0
	}

	var f int64 = int64(sf.Perms) & 0777
	if sf.Create {
		f |= int64(C.IPC_CREAT)
	}
	if sf.Exclusive {
		f |= int64(C.IPC_EXCL)
	}

	return f
}

// 信号量操作的标识位
type semOpFlags struct {
	// 操作是否阻塞，true 为不阻塞
	DontWait bool

	// 当操作的进程退出后，该进程对 sem 进行的操作将被取消
	SemUnDo bool
}

// 返回标识位
func (so *semOpFlags) flags() int64 {
	if so == nil {
		return 0
	}

	var f int64
	if so.DontWait {
		f |= int64(C.IPC_NOWAIT)
	}
	if so.SemUnDo {
		f |= int64(C.SEM_UNDO)
	}

	return f
}
