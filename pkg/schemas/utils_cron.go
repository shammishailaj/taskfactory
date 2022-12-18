package schemas

type ScheduleCommandOutput struct {
	CPUUsage    int64
	MemoryUsage uint64
	LoadAverage float64
}
