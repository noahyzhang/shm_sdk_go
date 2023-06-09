package shm_sdk_go

// System call constants.
const (
	sysShmAt  = 228
	sysShmCtl = 229
	sysShmDt  = 230
	sysShmGet = 231
)

const (
	// Mode bits for `shmget`.

	// Create key if key does not exist.
	IPC_CREAT = 01000
	// Fail if key exists.
	IPC_EXCL = 02000
	// Return error on wait.
	IPC_NOWAIT = 04000

	// Special key values.

	// Private key.
	IPC_PRIVATE = 0

	// Flags for `shmat`.

	// Attach read-only access.
	SHM_RDONLY = 010000
	// Round attach address to SHMLBA.
	SHM_RND = 020000
	// Take-over region on attach.
	SHM_REMAP = 040000
	// Execution access.
	SHM_EXEC = 0100000

	// Commands for `shmctl`.

	// Lock segment (root only).
	SHM_LOCK = 1
	// Unlock segment (root only).
	SHM_UNLOCK = 12

	// Control commands for `shmctl`.

	// remove identifier.
	IPC_RMID = 0
	// set `ipc_perm` options.
	IPC_SET = 1
	// Get `ipc_perm' options.
	IPC_STAT = 2

	SEM_UNDO   = 0x1000
	SEM_SETVAL = 16
	SEM_GETVAL = 12
)

// Perm is used to pass permission information to IPC operations.
type Perm struct {
	// Creator's user ID.
	Cuid uint32
	// Creator's group ID.
	Cgid uint32
	// Owner's user ID.
	Uid uint32
	// Owner's group ID.
	Gid uint32
	// Read/write permission.
	Mode uint16
	// Sequence number.
	Seq uint16
	// Padding.
	PadCgo0 [4]byte
	// Key.
	Key int64
}

// IdDs describes shared memory segment.
type IdDs struct {
	// Operation permission struct.
	Perm Perm
	// Size of segment in bytes.
	SegSz uint64
	// Pid of last shmat/shmdt.
	Lpid int32
	// Pid of creator.
	Cpid int32
	// Number of current attaches.
	Nattch int32
	// Padding.
	PadCgo0 [4]byte
	// Last attach time.
	Atime int64
	// Last detach time.
	Dtime int64
	// Last change time.
	Ctime int64
}
