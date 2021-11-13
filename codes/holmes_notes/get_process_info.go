package main


func getUsageNormal() (float64, float64, int, int, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPercent, err := p.Percent(time.Second)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// The default percent is from all cores, multiply by runtime.NumCPU()
	// but it's inconvenient to calculate the proper percent
	// here we divide by core number, so we can set a percent bar more intuitively
	cpuPercent = cpuPercent / float64(runtime.NumCPU())

	mem, err := p.MemoryPercent()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	gNum := runtime.NumGoroutine()

	tNum := getThreadNum()

	return cpuPercent, float64(mem), gNum, tNum, nil
}

func getUsageCGroup() (float64, float64, int, int, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPercent, err := p.Percent(time.Second)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPeriod, err := readUint("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuQuota, err := readUint("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	if err != nil {
		return 0, 0, 0, 0, err
	}
	cpuCore := float64(cpuQuota) / float64(cpuPeriod)

	// the same with physical machine
	// need to divide by core number
	cpuPercent = cpuPercent / cpuCore
	mem, err := p.MemoryInfo()
	p.MemoryPercent()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	memLimit, err := getCGroupMemoryLimit()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// mem.RSS / cgroup limit in bytes
	memPercent := float64(mem.RSS) * 100 / float64(memLimit)

	gNum := runtime.NumGoroutine()

	tNum := getThreadNum()

	return cpuPercent, memPercent, gNum, tNum, nil
}

func getCGroupMemoryLimit() (uint64, error) {
	usage, err := readUint("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		return 0, err
	}
	machineMemory, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	limit := uint64(math.Min(float64(usage), float64(machineMemory.Total)))
	return limit, nil
}

func readUint(path string) (uint64, error) {
	v, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUint(strings.TrimSpace(string(v)), 10, 64)
}

func parseUint(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

func getThreadNum() int {
	return pprof.Lookup("threadcreate").Count()
}

func main() {
	// in machine use getUsageNormal
	fmt.Println(getUsageNormal())
	// in docker us getUsageCGroup
	fmt.Println(getUsageCGroup())
	
}
