package storage

import "sync"

// MemoryRecord represent some memory metric about node
type MemoryRecord struct {
	// Total represent the total memory of physical node
	Total int64
	// Allocated represent the memory that has been allocated to pods by k8s
	Allocated int64
	// Capacity represent the memory the memory that be use by pods
	Capacity int64
	// Allocatable represent the memory that can be allocate by k8s
	Allocatable int64  // doest trustable
}



type CacheStorage struct {
	lock sync.Locker
	memInfo map[string]MemoryRecord
}

func NewCacheStorage()*CacheStorage{
	return &CacheStorage{
		lock:    &sync.Mutex{},
		memInfo: make(map[string]MemoryRecord),
	}
}

func(cache *CacheStorage)Update(nodeName string,record MemoryRecord){
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.memInfo[nodeName] = MemoryRecord{
		Total:       record.Total,
		Allocated:   record.Allocated,
		Capacity:    record.Capacity,
		Allocatable: record.Allocatable,
	}
}


func(cache *CacheStorage)DeepCopy()map[string]MemoryRecord{
	cache.lock.Lock()
	defer cache.lock.Unlock()
	minfor := map[string]MemoryRecord{}
	for k,v := range cache.memInfo{
		minfor[k] = v
	}
	return minfor
}