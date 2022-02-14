package basic

// TFill 不同的字段排序会影响"填充字节"的大小,如下所示
type TFill struct {
	b byte
	i int64
	u uint16
}

// SFill 与 TFill 不同的字段排序
type SFill struct {
	b byte
	u uint16
	i int64
}

// go1.17/src/runtime/mstats.go:22
// 自动填充事例，通常我们会通过空标识符来进行主动填充
/*
type mstats struct {
	// General statistics.
	alloc       uint64 // bytes allocated and not yet freed
	total_alloc uint64 // bytes allocated (even if freed)
	sys         uint64 // bytes obtained from system (should be sum of xxx_sys below, no locking, approximate)
	nlookup     uint64 // number of pointer lookups (unused)
	nmalloc     uint64 // number of mallocs
	nfree       uint64 // number of frees

	// Statistics about malloc heap.
	// Updated atomically, or with the world stopped.
	//
	// Like MemStats, heap_sys and heap_inuse do not count memory
	// in manually-managed spans.
	heap_sys      sysMemStat // virtual address space obtained from system for GC'd heap
	heap_inuse    uint64     // bytes in mSpanInUse spans
	heap_released uint64     // bytes released to the os

	// heap_objects is not used by the runtime directly and instead
	// computed on the fly by updatememstats.
	heap_objects uint64 // total number of allocated objects

	// Statistics about stacks.
	stacks_inuse uint64     // bytes in manually-managed stack spans; computed by updatememstats
	stacks_sys   sysMemStat // only counts newosproc0 stack in mstats; differs from MemStats.StackSys

	// Statistics about allocation of low-level fixed-size structures.
	// Protected by FixAlloc locks.
	mspan_inuse  uint64 // mspan structures
	mspan_sys    sysMemStat
	mcache_inuse uint64 // mcache structures
	mcache_sys   sysMemStat
	buckhash_sys sysMemStat // profiling bucket hash table

	// Statistics about GC overhead.
	gcWorkBufInUse           uint64     // computed by updatememstats
	gcProgPtrScalarBitsInUse uint64     // computed by updatememstats
	gcMiscSys                sysMemStat // updated atomically or during STW

	// Miscellaneous statistics.
	other_sys sysMemStat // updated atomically or during STW

	// Statistics about the garbage collector.

	// Protected by mheap or stopping the world during GC.
	last_gc_unix    uint64 // last gc (in unix time)
	pause_total_ns  uint64
	pause_ns        [256]uint64 // circular buffer of recent gc pause lengths
	pause_end       [256]uint64 // circular buffer of recent gc end times (nanoseconds since 1970)
	numgc           uint32
	numforcedgc     uint32  // number of user-forced GCs
	gc_cpu_fraction float64 // fraction of CPU time used by GC
	enablegc        bool
	debuggc         bool

	// Statistics about allocation size classes.

	by_size [_NumSizeClasses]struct {
		size    uint32
		nmalloc uint64
		nfree   uint64
	}

	// Add an uint32 for even number of size classes to align below fields
	// to 64 bits for atomic operations on 32 bit platforms.
	_ [1 - _NumSizeClasses%2]uint32

	last_gc_nanotime uint64 // last gc (monotonic time)
	last_heap_inuse  uint64 // heap_inuse at mark termination of the previous GC

	// heapStats is a set of statistics
	heapStats consistentHeapStats

	// _ uint32 // ensure gcPauseDist is aligned

	// gcPauseDist represents the distribution of all GC-related
	// application pauses in the runtime.
	//
	// Each individual pause is counted separately, unlike pause_ns.
	gcPauseDist timeHistogram
}
*/
