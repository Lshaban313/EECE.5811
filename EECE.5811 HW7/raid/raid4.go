package raid

// RAID4: data on disks 0–3, parity always on disk 4.
type RAID4 struct {
	disks []*Disk
}

func xorBytes(a, b []byte) {
	for i := range a {
		a[i] ^= b[i]
	}
}

func (r *RAID4) Write(n int64, data []byte) error {
	D := int64(len(r.disks))
	dataDisks := D - 1
	stripe := n / dataDisks
	dataIdx := n % dataDisks

	// write data
	if err := r.disks[dataIdx].WriteBlock(stripe, data); err != nil {
		return err
	}

	// recompute parity
	parity := make([]byte, r.disks[0].BlockSize)
	copy(parity, data)
	for i := int64(0); i < dataDisks; i++ {
		if i == dataIdx {
			continue
		}
		block, err := r.disks[i].ReadBlock(stripe)
		if err != nil {
			return err
		}
		xorBytes(parity, block)
	}
	return r.disks[D-1].WriteBlock(stripe, parity)
}

func (r *RAID4) Read(n int64) ([]byte, error) {
	D := int64(len(r.disks))
	dataDisks := D - 1
	stripe := n / dataDisks
	dataIdx := n % dataDisks

	// try read
	data, err := r.disks[dataIdx].ReadBlock(stripe)
	if err == nil {
		return data, nil
	}
	// reconstruct via parity
	res := make([]byte, r.disks[0].BlockSize)
	for i := int64(0); i < D; i++ {
		block, e2 := r.disks[i].ReadBlock(stripe)
		if e2 != nil {
			continue
		}
		xorBytes(res, block)
	}
	return res, nil
}
