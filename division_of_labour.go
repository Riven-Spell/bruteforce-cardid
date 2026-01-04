package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/Pashugan/trie"
	"github.com/Riven-Spell/generic/enumerable"
)

type LabourDivider interface {
	GetJob(ThreadID uint) string
}

var MaxValidID = func() uint64 {
	max, _ := strconv.ParseUint(prepareInput("0FFF FFFF FFFF FFFF"), 16, 64)

	return max
}()

type LabourDividerPureRandom struct{}

func (l *LabourDividerPureRandom) GetJob(ThreadID uint) string {
	return fmt.Sprintf("%016X", rand.N(MaxValidID))
}

type LabourDividerCollisionChecker struct {
	t *trie.Trie
}

func (l *LabourDividerCollisionChecker) GetJob(ThreadID uint) string {
	for {
		id := fmt.Sprintf("%016X", rand.N(MaxValidID))
		if l.t.Search(id) != nil {
			continue
		}

		l.t.Insert(id, true)
		return id
	}
}

type LabourDividerStrict struct {
	multiple  uint64
	origins   map[uint]uint64
	threadMap map[uint]*uint64
}

type MachineDefinition struct {
	coreCount      uint
	roamingThreads uint
}

func NewLabourDividerStrict(thisMachineIndex uint, machineCores ...MachineDefinition) LabourDivider {
	out := &LabourDividerStrict{}

	coreSum := enumerable.Sum(enumerable.FromList(machineCores, false), func(i MachineDefinition, o uint) uint {
		return o + (i.coreCount - i.roamingThreads)
	})

	roamsum := enumerable.Sum(enumerable.FromList(machineCores, false), func(i MachineDefinition, o uint) uint {
		return o + i.roamingThreads
	})

	fmt.Printf("%d cores total (%d strict, %d roam)", coreSum+roamsum, coreSum, roamsum)

	out.origins = make(map[uint]uint64)
	out.threadMap = make(map[uint]*uint64)
	out.multiple = MaxValidID / uint64(coreSum)

	thisMachineBegins := enumerable.Sum(enumerable.FromList(machineCores[:thisMachineIndex], false), func(i MachineDefinition, o uint64) uint64 {
		return o + (uint64(i.coreCount-i.roamingThreads) * out.multiple)
	})

	thisMachineCores := uint64(machineCores[thisMachineIndex].coreCount - machineCores[thisMachineIndex].roamingThreads)

	for v := range thisMachineCores {
		var initial = thisMachineBegins + out.multiple*v
		out.origins[uint(v)] = initial

		var zero = uint64(0)
		out.threadMap[uint(v)] = &zero
	}

	fmt.Println(out.origins)

	return out
}

func (l LabourDividerStrict) GetJob(ThreadID uint) string {
	if n, ok := l.threadMap[ThreadID]; !ok {
		return fmt.Sprintf("%016X", rand.N(MaxValidID)) // those remaining can be exploratory
	} else {
		out := *n + l.origins[ThreadID]

		if out == l.multiple {
			return fmt.Sprintf("%016X", rand.N(MaxValidID)) // those remaining can be exploratory
		}

		*n += 1

		return fmt.Sprintf("%016X", out)
	}
}
