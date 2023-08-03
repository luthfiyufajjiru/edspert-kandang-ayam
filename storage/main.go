package storage

import (
	"strconv"
	"sync"
)

type SerialFloat struct {
	data []float64
	mtx  sync.RWMutex
}

func (s *SerialFloat) Append(data string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	dataF, err := strconv.ParseFloat(data, 64)
	if err == nil {
		s.data = append(s.data, dataF)
	}
}

func (s *SerialFloat) Len() (res int) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	res = len(s.data)
	return
}

func (s *SerialFloat) Range(start, end int) (res []float64) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	ln := len(s.data)
	if ln == 0 {
		return
	}
	if start < 0 {
		start = 0
	}
	if end > ln {
		end = ln
	}
	res = s.data[start:end]
	return
}
