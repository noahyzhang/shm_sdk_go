package shm_sdk_go

import "testing"

var semKey uint32 = 0xcc9d

func TestSemaphore_Create(t *testing.T) {
	var sem Semaphore
	err := sem.Create(semKey)
	if err != nil {
		t.Fatalf("create sem err: %s", err.Error())
	}
	err = sem.Lock(true)
	if err != nil {
		t.Fatalf("lock sem err: %s", err.Error())
	}
	err = sem.Unlock()
	if err != nil {
		t.Fatalf("unlock sem err: %s", err.Error())
	}
	err = sem.Destroy()
	if err != nil {
		t.Fatalf("destroy sem err: %s", err.Error())
	}
}

func TestSemaphore_Lock(t *testing.T) {
	var sem Semaphore
	err := sem.Create(semKey)
	if err != nil {
		t.Fatalf("create sem err: %s", err.Error())
	}
	err = sem.Lock(true)
	if err != nil {
		t.Fatalf("lock sem err: %s", err.Error())
	}
	err = sem.Lock(true)
	if err != nil {
		t.Logf("Repeated locking, timed out. res: %s", err.Error())
	}
}

func TestSemaphore_Unlock(t *testing.T) {
	var sem Semaphore
	err := sem.Create(semKey)
	if err != nil {
		t.Fatalf("create sem err: %s", err.Error())
	}
	err = sem.Unlock()
	if err != nil {
		t.Fatalf("unlock sem err: %s", err.Error())
	}
	err = sem.Unlock()
	if err != nil {
		t.Logf("Repeated unlock. no problem")
	}
}
