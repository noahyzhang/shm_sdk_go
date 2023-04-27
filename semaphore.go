package shm_sdk_go

import "fmt"

type Semaphore struct {
	isCreate bool
	semId    int
}

func (s *Semaphore) Create(semKey uint32) error {
	semId, err := shmGet(int(semKey))
	if err != nil {
		return fmt.Errorf("[Semaphore::Create] shmGet err: %s", err.Error())
	}
	s.semId = semId
	s.isCreate = true
	return nil
}

func (s *Semaphore) Lock(isWait bool) error {
	if !s.isCreate {
		return fmt.Errorf("[Semaphore::Lock] no create Semaphore")
	}
	return semLock(s.semId, isWait)
}

func (s *Semaphore) Unlock() error {
	if !s.isCreate {
		return fmt.Errorf("[Semaphore::Unlock] no create Semaphore")
	}
	return semUnLock(s.semId)
}

func (s *Semaphore) Destroy() error {
	if !s.isCreate {
		return fmt.Errorf("[Semaphore::Destroy] no create Semaphore")
	}
	return semDestroy(s.semId)
}
