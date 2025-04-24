package raid

// RAID1 implements mirroring across all disks.
type RAID1 struct {
	disks []*Disk
}

func (r *RAID1) Write(n int64, data []byte) error {
	for _, d := range r.disks {
		if err := d.WriteBlock(n, data); err != nil {
			return err
		}
	}
	return nil
}

func (r *RAID1) Read(n int64) ([]byte, error) {
	// just read from disk 0
	return r.disks[0].ReadBlock(n)
}
