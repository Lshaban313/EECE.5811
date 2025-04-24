package raid

import (
	"fmt"
)

// RAID is the interface for all RAID levels.
type RAID interface {
	Write(blockNum int64, data []byte) error
	Read(blockNum int64) ([]byte, error)
}

// Open constructs a RAID instance of the given level (0,1,4,5).
// files should be 5 paths, blockSize in bytes.
func Open(level int, files []string, blockSize int) (RAID, error) {
	if len(files) != 5 {
		return nil, fmt.Errorf("need exactly 5 disk files, got %d", len(files))
	}
	disks := make([]*Disk, 5)
	for i, path := range files {
		d, err := NewDisk(path, blockSize)
		if err != nil {
			return nil, err
		}
		disks[i] = d
	}
	switch level {
	case 0:
		return &RAID0{disks}, nil
	case 1:
		return &RAID1{disks}, nil
	case 4:
		return &RAID4{disks}, nil
	case 5:
		return &RAID5{disks}, nil
	default:
		return nil, fmt.Errorf("unsupported RAID level %d", level)
	}
}
