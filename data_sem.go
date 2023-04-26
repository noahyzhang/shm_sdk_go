package shm_sdk_go

import (
	"fmt"
)

type Semaphore struct {
	semSet *semaphoreSet
}

func (s *Semaphore) Create(semKey uint32) error {
	semSet, err := getSemSet(int64(semKey), 1, &semSetFlags{
		Create:    false,
		Exclusive: false,
		Perms:     0666,
	})
	if err == nil {
		s.semSet = semSet
		return nil
	}
	semSet, err = getSemSet(int64(semKey), 1, &semSetFlags{
		Create:    true,
		Exclusive: true,
		Perms:     0666,
	})
	if err != nil {
		return fmt.Errorf("[Semaphore:Create] getSemSet err: %s", err.Error())
	}
	s.semSet = semSet
	if err = s.semSet.setVal(0, 1); err != nil {
		return fmt.Errorf("[Semaphore:Create] setVal err: %s", err.Error())
	}
	return nil
}

func (s *Semaphore) Lock(isWait bool) error {
	flag := semOpFlags{DontWait: !isWait, SemUnDo: true}
	ops := newSemOps()
	if err := ops.decrement(0, 1, &flag); err != nil {
		return fmt.Errorf("[Semaphore::Lock] Desrement err: %s", err.Error())
	}
	return s.semSet.run(ops, -1)
}

func (s *Semaphore) Unlock() error {
	flag := semOpFlags{DontWait: false, SemUnDo: true}
	ops := newSemOps()
	if err := ops.increment(0, 1, &flag); err != nil {
		return fmt.Errorf("[Semaphore::Unlock] increment err: %s", err.Error())
	}
	return s.semSet.run(ops, -1)
}

func (s *Semaphore) Destroy() error {
	if err := s.semSet.remove(); err != nil {
		return fmt.Errorf("[Semaphore::Destroy] remove err: %s", err.Error())
	}
	return nil
}
