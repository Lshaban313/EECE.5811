package raid

// RAID0 implements simple striping.
type RAID0 struct {
	disks []*Disk
}

func (r *RAID0) Write(n int64, data []byte) error {
	idx := n % int64(len(r.disks))
	stripe := n / int64(len(r.disks))
	return r.disks[idx].WriteBlock(stripe, data)
}

func (r *RAID0) Read(n int64) ([]byte, error) {
	idx := n % int64(len(r.disks))
	stripe := n / int64(len(r.disks))
	return r.disks[idx].ReadBlock(stripe)
}
