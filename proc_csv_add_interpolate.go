package psn

import (
	"fmt"
	"sort"

	humanize "github.com/dustin/go-humanize"
)

// Combine combines a list Proc and returns one combined Proc.
// Field values are estimated. UnixNanosecond is reset 0.
// And UnixSecond and other fields that cannot be averaged are set
// with the field value in the last element. This is meant to be
// used to combine Proc rows with duplicate unix second timestamps.
func Combine(procs ...Proc) Proc {
	if len(procs) < 1 {
		return Proc{}
	}
	if len(procs) == 1 {
		return procs[0]
	}

	lastProc := procs[len(procs)-1]
	combined := lastProc
	combined.UnixNanosecond = 0

	// calculate the average
	var (
		// for PSEntry
		voluntaryCtxtSwitches    uint64
		nonVoluntaryCtxtSwitches uint64
		cpuNum                   float64
		vmRSSNum                 uint64
		vmSizeNum                uint64

		// for LoadAvg
		loadAvg1Minute                   float64
		loadAvg5Minute                   float64
		loadAvg15Minute                  float64
		runnableKernelSchedulingEntities int64
		currentKernelSchedulingEntities  int64

		// for DSEntry
		readsCompleted       uint64
		sectorsRead          uint64
		writesCompleted      uint64
		sectorsWritten       uint64
		timeSpentOnReadingMs uint64
		timeSpentOnWritingMs uint64

		// for DSEntry delta
		readsCompletedDelta  uint64
		sectorsReadDelta     uint64
		writesCompletedDelta uint64
		sectorsWrittenDelta  uint64

		// for NSEntry
		receivePackets   uint64
		transmitPackets  uint64
		receiveBytesNum  uint64
		transmitBytesNum uint64

		// for NSEntry delta
		receivePacketsDelta   uint64
		transmitPacketsDelta  uint64
		receiveBytesNumDelta  uint64
		transmitBytesNumDelta uint64
	)

	for _, p := range procs {
		// for PSEntry
		voluntaryCtxtSwitches += p.PSEntry.VoluntaryCtxtSwitches
		nonVoluntaryCtxtSwitches += p.PSEntry.NonvoluntaryCtxtSwitches
		cpuNum += p.PSEntry.CPUNum
		vmRSSNum += p.PSEntry.VMRSSNum
		vmSizeNum += p.PSEntry.VMSizeNum

		// for LoadAvg
		loadAvg1Minute += p.LoadAvg.LoadAvg1Minute
		loadAvg5Minute += p.LoadAvg.LoadAvg5Minute
		loadAvg15Minute += p.LoadAvg.LoadAvg15Minute
		runnableKernelSchedulingEntities += p.LoadAvg.RunnableKernelSchedulingEntities
		currentKernelSchedulingEntities += p.LoadAvg.CurrentKernelSchedulingEntities

		// for DSEntry
		readsCompleted += p.DSEntry.ReadsCompleted
		sectorsRead += p.DSEntry.SectorsRead
		writesCompleted += p.DSEntry.WritesCompleted
		sectorsWritten += p.DSEntry.SectorsWritten
		timeSpentOnReadingMs += p.DSEntry.TimeSpentOnReadingMs
		timeSpentOnWritingMs += p.DSEntry.TimeSpentOnWritingMs

		// for DSEntry delta
		readsCompletedDelta += p.ReadsCompletedDelta
		sectorsReadDelta += p.SectorsReadDelta
		writesCompletedDelta += p.WritesCompletedDelta
		sectorsWrittenDelta += p.SectorsWrittenDelta

		// for NSEntry
		receivePackets += p.NSEntry.ReceivePackets
		transmitPackets += p.NSEntry.TransmitPackets
		receiveBytesNum += p.NSEntry.ReceiveBytesNum
		transmitBytesNum += p.NSEntry.TransmitBytesNum

		// for NSEntry delta
		receivePacketsDelta += p.ReceivePacketsDelta
		transmitPacketsDelta += p.TransmitPacketsDelta
		receiveBytesNumDelta += p.ReceiveBytesNumDelta
		transmitBytesNumDelta += p.TransmitBytesNumDelta
	}

	pN := len(procs)

	// for PSEntry
	combined.PSEntry.VoluntaryCtxtSwitches = uint64(voluntaryCtxtSwitches) / uint64(pN)
	combined.PSEntry.NonvoluntaryCtxtSwitches = uint64(nonVoluntaryCtxtSwitches) / uint64(pN)
	combined.PSEntry.CPUNum = float64(cpuNum) / float64(pN)
	combined.PSEntry.CPU = fmt.Sprintf("%3.2f %%", combined.PSEntry.CPUNum)
	combined.PSEntry.VMRSSNum = uint64(vmRSSNum) / uint64(pN)
	combined.PSEntry.VMRSS = humanize.Bytes(combined.PSEntry.VMRSSNum)
	combined.PSEntry.VMSizeNum = uint64(vmSizeNum) / uint64(pN)
	combined.PSEntry.VMSize = humanize.Bytes(combined.PSEntry.VMSizeNum)

	// for LoadAvg
	combined.LoadAvg.LoadAvg1Minute = float64(loadAvg1Minute) / float64(pN)
	combined.LoadAvg.LoadAvg5Minute = float64(loadAvg5Minute) / float64(pN)
	combined.LoadAvg.LoadAvg15Minute = float64(loadAvg15Minute) / float64(pN)
	combined.LoadAvg.RunnableKernelSchedulingEntities = int64(loadAvg15Minute) / int64(pN)
	combined.LoadAvg.CurrentKernelSchedulingEntities = int64(loadAvg15Minute) / int64(pN)

	// for DSEntry
	combined.DSEntry.ReadsCompleted = uint64(readsCompleted) / uint64(pN)
	combined.DSEntry.SectorsRead = uint64(sectorsRead) / uint64(pN)
	combined.DSEntry.WritesCompleted = uint64(writesCompleted) / uint64(pN)
	combined.DSEntry.SectorsWritten = uint64(sectorsWritten) / uint64(pN)
	combined.DSEntry.TimeSpentOnReadingMs = uint64(timeSpentOnReadingMs) / uint64(pN)
	combined.DSEntry.TimeSpentOnReading = humanizeDurationMs(combined.DSEntry.TimeSpentOnReadingMs)
	combined.DSEntry.TimeSpentOnWritingMs = uint64(timeSpentOnWritingMs) / uint64(pN)
	combined.DSEntry.TimeSpentOnWriting = humanizeDurationMs(combined.DSEntry.TimeSpentOnWritingMs)
	combined.ReadsCompletedDelta = uint64(readsCompletedDelta) / uint64(pN)
	combined.SectorsReadDelta = uint64(sectorsReadDelta) / uint64(pN)
	combined.WritesCompletedDelta = uint64(writesCompletedDelta) / uint64(pN)
	combined.SectorsWrittenDelta = uint64(sectorsWrittenDelta) / uint64(pN)

	// for NSEntry
	combined.NSEntry.ReceiveBytesNum = uint64(receiveBytesNum) / uint64(pN)
	combined.NSEntry.TransmitBytesNum = uint64(transmitBytesNum) / uint64(pN)
	combined.NSEntry.ReceivePackets = uint64(receivePackets) / uint64(pN)
	combined.NSEntry.TransmitPackets = uint64(transmitPackets) / uint64(pN)
	combined.NSEntry.ReceiveBytes = humanize.Bytes(combined.NSEntry.ReceiveBytesNum)
	combined.NSEntry.TransmitBytes = humanize.Bytes(combined.NSEntry.TransmitBytesNum)
	combined.ReceivePacketsDelta = uint64(receivePacketsDelta) / uint64(pN)
	combined.TransmitPacketsDelta = uint64(transmitPacketsDelta) / uint64(pN)
	combined.ReceiveBytesNumDelta = uint64(receiveBytesNumDelta) / uint64(pN)
	combined.ReceiveBytesDelta = humanize.Bytes(combined.ReceiveBytesNumDelta)
	combined.TransmitBytesNumDelta = uint64(transmitBytesNumDelta) / uint64(pN)
	combined.TransmitBytesDelta = humanize.Bytes(combined.TransmitBytesNumDelta)

	return combined
}

// Interpolate returns the missing, estimated 'Proc's if any.
// It assumes that 'upper' Proc is later than 'lower'.
// And UnixSecond and other fields that cannot be averaged are set
// with the field value in the last element.
func Interpolate(lower, upper Proc) (procs []Proc, err error) {
	if upper.UnixSecond <= lower.UnixSecond {
		return nil, fmt.Errorf("lower unix second %d >= upper unix second %d", lower.UnixSecond, upper.UnixSecond)
	}

	// min unix second is 5, max is 7
	// then the expected row number is 7-5+1=3
	expectedRowN := upper.UnixSecond - lower.UnixSecond + 1
	if expectedRowN == 2 {
		// no need to interpolate
		return
	}

	// calculate the delta
	var (
		// for PSEntry
		voluntaryCtxtSwitches    = (upper.PSEntry.VoluntaryCtxtSwitches - lower.PSEntry.VoluntaryCtxtSwitches) / uint64(expectedRowN)
		nonVoluntaryCtxtSwitches = (upper.PSEntry.NonvoluntaryCtxtSwitches - lower.PSEntry.NonvoluntaryCtxtSwitches) / uint64(expectedRowN)
		cpuNum                   = (upper.PSEntry.CPUNum - lower.PSEntry.CPUNum) / float64(expectedRowN)
		vmRSSNum                 = (upper.PSEntry.VMRSSNum - lower.PSEntry.VMRSSNum) / uint64(expectedRowN)
		vmSizeNum                = (upper.PSEntry.VMSizeNum - lower.PSEntry.VMSizeNum) / uint64(expectedRowN)

		// for LoadAvg
		loadAvg1Minute                   = (upper.LoadAvg.LoadAvg1Minute - lower.LoadAvg.LoadAvg1Minute) / float64(expectedRowN)
		loadAvg5Minute                   = (upper.LoadAvg.LoadAvg5Minute - lower.LoadAvg.LoadAvg5Minute) / float64(expectedRowN)
		loadAvg15Minute                  = (upper.LoadAvg.LoadAvg15Minute - lower.LoadAvg.LoadAvg15Minute) / float64(expectedRowN)
		runnableKernelSchedulingEntities = (upper.LoadAvg.RunnableKernelSchedulingEntities - lower.LoadAvg.RunnableKernelSchedulingEntities) / int64(expectedRowN)
		currentKernelSchedulingEntities  = (upper.LoadAvg.RunnableKernelSchedulingEntities - lower.LoadAvg.RunnableKernelSchedulingEntities) / int64(expectedRowN)

		// for DSEntry
		readsCompleted       = (upper.DSEntry.ReadsCompleted - lower.DSEntry.ReadsCompleted) / uint64(expectedRowN)
		sectorsRead          = (upper.DSEntry.SectorsRead - lower.DSEntry.SectorsRead) / uint64(expectedRowN)
		writesCompleted      = (upper.DSEntry.WritesCompleted - lower.DSEntry.WritesCompleted) / uint64(expectedRowN)
		sectorsWritten       = (upper.DSEntry.SectorsWritten - lower.DSEntry.SectorsWritten) / uint64(expectedRowN)
		timeSpentOnReadingMs = (upper.DSEntry.TimeSpentOnReadingMs - lower.DSEntry.TimeSpentOnReadingMs) / uint64(expectedRowN)
		timeSpentOnWritingMs = (upper.DSEntry.TimeSpentOnWritingMs - lower.DSEntry.TimeSpentOnWritingMs) / uint64(expectedRowN)

		// for DSEntry delta
		readsCompletedDelta  = (upper.ReadsCompletedDelta - lower.ReadsCompletedDelta) / uint64(expectedRowN)
		sectorsReadDelta     = (upper.SectorsReadDelta - lower.SectorsReadDelta) / uint64(expectedRowN)
		writesCompletedDelta = (upper.WritesCompletedDelta - lower.WritesCompletedDelta) / uint64(expectedRowN)
		sectorsWrittenDelta  = (upper.SectorsWrittenDelta - lower.SectorsWrittenDelta) / uint64(expectedRowN)

		// for NSEntry
		receivePackets   = (upper.NSEntry.ReceivePackets - lower.NSEntry.ReceivePackets) / uint64(expectedRowN)
		transmitPackets  = (upper.NSEntry.TransmitPackets - lower.NSEntry.TransmitPackets) / uint64(expectedRowN)
		receiveBytesNum  = (upper.NSEntry.ReceiveBytesNum - lower.NSEntry.ReceiveBytesNum) / uint64(expectedRowN)
		transmitBytesNum = (upper.NSEntry.TransmitBytesNum - lower.NSEntry.TransmitBytesNum) / uint64(expectedRowN)

		// for NSEntry delta
		receivePacketsDelta   = (upper.ReceivePacketsDelta - lower.ReceivePacketsDelta) / uint64(expectedRowN)
		transmitPacketsDelta  = (upper.TransmitPacketsDelta - lower.TransmitPacketsDelta) / uint64(expectedRowN)
		receiveBytesNumDelta  = (upper.ReceiveBytesNumDelta - lower.ReceiveBytesNumDelta) / uint64(expectedRowN)
		transmitBytesNumDelta = (upper.TransmitBytesNumDelta - lower.TransmitBytesNumDelta) / uint64(expectedRowN)
	)

	procs = make([]Proc, expectedRowN-2)
	for i := range procs {
		procs[i] = upper

		procs[i].UnixNanosecond = 0
		procs[i].UnixSecond = lower.UnixSecond + int64(i+1)

		// for PSEntry
		procs[i].PSEntry.VoluntaryCtxtSwitches = lower.PSEntry.VoluntaryCtxtSwitches + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].PSEntry.NonvoluntaryCtxtSwitches = lower.PSEntry.NonvoluntaryCtxtSwitches + uint64(i+1)*nonVoluntaryCtxtSwitches
		procs[i].PSEntry.CPUNum = lower.PSEntry.CPUNum + float64(i+1)*cpuNum
		procs[i].PSEntry.CPU = fmt.Sprintf("%3.2f %%", procs[i].PSEntry.CPUNum)
		procs[i].PSEntry.VMRSSNum = lower.PSEntry.VMRSSNum + uint64(i+1)*vmRSSNum
		procs[i].PSEntry.VMRSS = humanize.Bytes(procs[i].PSEntry.VMRSSNum)
		procs[i].PSEntry.VMSizeNum = lower.PSEntry.VMSizeNum + uint64(i+1)*vmSizeNum
		procs[i].PSEntry.VMSize = humanize.Bytes(procs[i].PSEntry.VMSizeNum)

		// for LoadAvg
		procs[i].LoadAvg.LoadAvg1Minute = lower.LoadAvg.LoadAvg1Minute + float64(i+1)*loadAvg1Minute
		procs[i].LoadAvg.LoadAvg5Minute = lower.LoadAvg.LoadAvg5Minute + float64(i+1)*loadAvg5Minute
		procs[i].LoadAvg.LoadAvg15Minute = lower.LoadAvg.LoadAvg15Minute + float64(i+1)*loadAvg15Minute
		procs[i].LoadAvg.RunnableKernelSchedulingEntities = lower.LoadAvg.RunnableKernelSchedulingEntities + int64(i+1)*runnableKernelSchedulingEntities
		procs[i].LoadAvg.CurrentKernelSchedulingEntities = lower.LoadAvg.CurrentKernelSchedulingEntities + int64(i+1)*currentKernelSchedulingEntities

		// for DSEntry
		procs[i].DSEntry.ReadsCompleted = lower.DSEntry.ReadsCompleted + uint64(i+1)*readsCompleted
		procs[i].DSEntry.SectorsRead = lower.DSEntry.SectorsRead + uint64(i+1)*sectorsRead
		procs[i].DSEntry.WritesCompleted = lower.DSEntry.WritesCompleted + uint64(i+1)*writesCompleted
		procs[i].DSEntry.SectorsWritten = lower.DSEntry.SectorsWritten + uint64(i+1)*sectorsWritten
		procs[i].DSEntry.TimeSpentOnReadingMs = lower.DSEntry.TimeSpentOnReadingMs + uint64(i+1)*timeSpentOnReadingMs
		procs[i].DSEntry.TimeSpentOnReading = humanizeDurationMs(procs[i].DSEntry.TimeSpentOnReadingMs)
		procs[i].DSEntry.TimeSpentOnWritingMs = lower.DSEntry.TimeSpentOnWritingMs + uint64(i+1)*timeSpentOnWritingMs
		procs[i].DSEntry.TimeSpentOnWriting = humanizeDurationMs(procs[i].DSEntry.TimeSpentOnWritingMs)
		procs[i].ReadsCompletedDelta = lower.ReadsCompletedDelta + uint64(i+1)*readsCompletedDelta
		procs[i].SectorsReadDelta = lower.SectorsReadDelta + uint64(i+1)*sectorsReadDelta
		procs[i].WritesCompletedDelta = lower.WritesCompletedDelta + uint64(i+1)*writesCompletedDelta
		procs[i].SectorsWrittenDelta = lower.SectorsWrittenDelta + uint64(i+1)*sectorsWrittenDelta

		// for NSEntry
		procs[i].NSEntry.ReceiveBytesNum = uint64(receiveBytesNum) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].NSEntry.TransmitBytesNum = uint64(transmitBytesNum) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].NSEntry.ReceivePackets = uint64(receivePackets) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].NSEntry.TransmitPackets = uint64(transmitPackets) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].NSEntry.ReceiveBytes = humanize.Bytes(procs[i].NSEntry.ReceiveBytesNum)
		procs[i].NSEntry.TransmitBytes = humanize.Bytes(procs[i].NSEntry.TransmitBytesNum)
		procs[i].ReceivePacketsDelta = uint64(receivePacketsDelta) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].TransmitPacketsDelta = uint64(transmitPacketsDelta) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].ReceiveBytesNumDelta = uint64(receiveBytesNumDelta) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].ReceiveBytesDelta = humanize.Bytes(procs[i].ReceiveBytesNumDelta)
		procs[i].TransmitBytesNumDelta = uint64(transmitBytesNumDelta) + uint64(i+1)*voluntaryCtxtSwitches
		procs[i].TransmitBytesDelta = humanize.Bytes(procs[i].TransmitBytesNumDelta)
	}

	return
}

// Interpolate interpolates missing rows in CSV assuming CSV is to be collected for every second.
// 'Missing' means unix seconds in rows are not continuous.
// It fills in the empty rows by estimating the averages.
// It returns a new copy of CSV. And the new copy sets all unix nanoseconds to 0.,
// since it's now aggregated by the unix "second".
func (c *CSV) Interpolate() (cc *CSV, err error) {
	if c == nil || len(c.Rows) < 2 {
		// no need to interpolate
		return
	}

	// copy the original CSV data
	cc = &(*c)

	// find missing rows, assuming CSV is to be collected every second
	if cc.MinUnixSecond == cc.MaxUnixSecond {
		// no need to interpolate
		return
	}

	// min unix second is 5, max is 7
	// then the expected row number is 7-5+1=3
	expectedRowN := cc.MaxUnixSecond - cc.MinUnixSecond + 1
	secondToAllProcs := make(map[int64][]Proc)
	for _, row := range cc.Rows {
		if _, ok := secondToAllProcs[row.UnixSecond]; ok {
			secondToAllProcs[row.UnixSecond] = append(secondToAllProcs[row.UnixSecond], row)
		} else {
			secondToAllProcs[row.UnixSecond] = []Proc{row}
		}
	}
	if int64(len(cc.Rows)) == expectedRowN && len(cc.Rows) == len(secondToAllProcs) {
		// all rows have distinct unix second
		// and they are all continuous unix seconds
		return
	}

	// interpolate cases
	//
	// case #1. If duplicate rows are found (equal/different unix nanoseconds, equal unix seconds),
	//          combine those into one row with its average.
	//
	// case #2. If some rows are discontinuous in unix seconds, there are missing rows.
	//          Fill in those rows with average estimates.

	// case #1, find duplicate rows!
	// It finds duplicates by unix second! Not by unix nanoseconds!
	secondToProc := make(map[int64]Proc)
	for sec, procs := range secondToAllProcs {
		if len(procs) == 0 {
			return nil, fmt.Errorf("empty row found at unix second %d", sec)
		}

		if len(procs) == 1 {
			secondToProc[sec] = procs[0]
			continue // no need to combine
		}

		// procs conflicted on unix second,
		// we want to combine those into one
		secondToProc[sec] = Combine(procs...)
	}

	// sort and reset the unix second
	rows2 := make([]Proc, 0, len(secondToProc))
	allUnixSeconds := make([]int64, 0, len(secondToProc))
	for _, row := range secondToProc {
		row.UnixNanosecond = 0
		rows2 = append(rows2, row)
		allUnixSeconds = append(allUnixSeconds, row.UnixSecond)
	}
	sort.Sort(ProcSlice(rows2))

	cc.Rows = rows2
	cc.MinUnixNanosecond = rows2[0].UnixNanosecond
	cc.MinUnixSecond = rows2[0].UnixSecond
	cc.MaxUnixNanosecond = rows2[len(rows2)-1].UnixNanosecond
	cc.MaxUnixSecond = rows2[len(rows2)-1].UnixSecond

	// case #2, find missing rows!
	// if unix seconds have discontinued ranges, it's missing some rows!
	missingTS := make(map[int64]struct{})
	for unixSecond := cc.MinUnixSecond; unixSecond <= cc.MaxUnixSecond; unixSecond++ {
		_, ok := secondToProc[unixSecond]
		if !ok {
			missingTS[unixSecond] = struct{}{}
		}
	}
	if len(missingTS) == 0 {
		// now all rows have distinct unix second
		// and there's no missing unix seconds
		return
	}

	missingSeconds := make([]int64, 0, len(missingTS))
	for ts := range missingTS {
		missingSeconds = append(missingSeconds, ts)
	}
	sort.Sort(int64Slice(missingSeconds))

	for i := range missingSeconds {
		second := missingSeconds[i]
		if _, ok := secondToProc[second]; ok {
			return nil, fmt.Errorf("second %d is not supposed to be found at secondToProc but found", second)
		}
	}

	// now we need to estimate the Proc for missingTS
	// fmt.Printf("total %d points available, missing %d points\n", len(allUnixSeconds), len(missingTS))
	bds := buildBoundaries(allUnixSeconds)

	// start from mid, in case missing seconds are continuous (several seconds empty)
	for i := range missingSeconds {
		second := missingSeconds[i]
		if _, ok := secondToProc[second]; ok {
			// already estimated!
			continue
		}

		bd := bds.findBoundary(second)
		if bd.lower == second && bd.upper == second {
			return nil, fmt.Errorf("%d is supposed to be missing but found at index %d", second, bd.lowerIdx)
		}

		// not found at boundaries pool
		// must have been found since it was created with min,max unix second
		if bd.lowerIdx == -1 || bd.upperIdx == -1 {
			return nil, fmt.Errorf("boundary is not found for missing second %d", second)
		}

		procLower, ok := secondToProc[bd.lower]
		if !ok {
			return nil, fmt.Errorf("%d is not found at secondToProc", bd.lower)
		}
		procUpper, ok := secondToProc[bd.upper]
		if !ok {
			return nil, fmt.Errorf("%d is not found at secondToProc", bd.upper)
		}
		missingRows, err := Interpolate(procLower, procUpper)
		if err != nil {
			return nil, err
		}
		for _, mrow := range missingRows {
			secondToProc[mrow.UnixSecond] = mrow

			// now 'mrow.UnixSecond' is not missing anymore
			bds.add(mrow.UnixSecond)
		}
	}

	rows3 := make([]Proc, 0, len(secondToProc))
	for _, row := range secondToProc {
		row.UnixNanosecond = 0
		rows3 = append(rows3, row)
	}
	sort.Sort(ProcSlice(rows3))

	cc.Rows = rows3
	cc.MinUnixNanosecond = rows3[0].UnixNanosecond
	cc.MinUnixSecond = rows3[0].UnixSecond
	cc.MaxUnixNanosecond = rows3[len(rows3)-1].UnixNanosecond
	cc.MaxUnixSecond = rows3[len(rows3)-1].UnixSecond

	return
}

// ConvertUnixNano unix nanoseconds to unix second.
func ConvertUnixNano(unixNano int64) (unixSec int64) {
	return int64(unixNano / 1e9)
}
