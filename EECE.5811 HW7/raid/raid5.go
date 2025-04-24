package raid

// RAID5: rotating parity across all disks.
type RAID5 struct {
	disks []*Disk
}

func (r *RAID5) Write(n int64, data []byte) error {
	D := int64(len(r.disks))
	dataDisks := D - 1
	stripe := n / dataDisks
	parityDisk := stripe % D
	dataIdx := n % dataDisks

	// figure out which real disk holds this data
	var diskIdx int64
	if dataIdx < parityDisk {
		diskIdx = dataIdx
	} else {
		diskIdx = dataIdx + 1
	}
	if err := r.disks[diskIdx].WriteBlock(stripe, data); err != nil {
		return err
	}

	// compute parity
	parity := make([]byte, r.disks[0].BlockSize)
	for i := int64(0); i < D; i++ {
		if i == parityDisk {
			continue
		}
		block, err := r.disks[i].ReadBlock(stripe)
		if err != nil {
			return err
		}
		if parity == nil {
			parity = block
		} else {
			xorBytes(parity, block)
		}
	}
	return r.disks[parityDisk].WriteBlock(stripe, parity)
}

func (r *RAID5) Read(n int64) ([]byte, error) {
	D := int64(len(r.disks))
	dataDisks := D - 1
	stripe := n / dataDisks
	parityDisk := stripe % D
	dataIdx := n % dataDisks

	// map to real disk
	var diskIdx int64
	if dataIdx < parityDisk {
		diskIdx = dataIdx
	} else {
		diskIdx = dataIdx + 1
	}
	data, err := r.disks[diskIdx].ReadBlock(stripe)
	if err == nil {
		return data, nil
	}
	// reconstruct
	res := make([]byte, r.disks[0].BlockSize)
	for i := int64(0); i < D; i++ {
		if i == diskIdx {
			continue
		}
		block, e2 := r.disks[i].ReadBlock(stripe)
		if e2 != nil {
			continue
		}
		xorBytes(res, block)
	}
	return res, nil
}
