package shm_sdk_go

import (
	"fmt"
	"time"
)

type Semaphore struct {
	semSet *SemaphoreSet
}

func (s *Semaphore) Create(semKey uint32) error {
	semSet, err := GetSemSet(int64(semKey), 1, &SemSetFlags{
		Create:    false,
		Exclusive: false,
		Perms:     0666,
	})
	if err == nil {
		s.semSet = semSet
		return nil
	}
	semSet, err = GetSemSet(int64(semKey), 1, &SemSetFlags{
		Create:    true,
		Exclusive: true,
		Perms:     0666,
	})
	if err != nil {
		return fmt.Errorf("[Semaphore:Create] GetSemSet err: %s", err.Error())
	}
	s.semSet = semSet
	if err = s.semSet.Setval(0, 1); err != nil {
		return fmt.Errorf("[Semaphore:Create] Setval err: %s", err.Error())
	}
	return nil
}

func (s *Semaphore) Lock(isWait bool) error {
	flag := SemOpFlags{DontWait: isWait, SemUnDo: true}
	ops := NewSemOps()
	if err := ops.Decrement(0, 1, &flag); err != nil {
		return fmt.Errorf("[Semaphore::Lock] Desrement err: %s", err.Error())
	}
	return s.semSet.Run(ops, time.Second)
}

func (s *Semaphore) Unlock() error {
	flag := SemOpFlags{DontWait: false, SemUnDo: true}
	ops := NewSemOps()
	if err := ops.Increment(0, 1, &flag); err != nil {
		return fmt.Errorf("[Semaphore::Unlock] Increment err: %s", err.Error())
	}
	return s.semSet.Run(ops, -1)
}

func (s *Semaphore) Destroy() error {
	if err := s.semSet.Remove(); err != nil {
		return fmt.Errorf("[Semaphore::Destroy] remove err: %s", err.Error())
	}
	return nil
}
