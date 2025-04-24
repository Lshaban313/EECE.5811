package raid

import (
	"fmt"
	"os"
)

// Disk simulates one physical disk file at a fixed block size.
type Disk struct {
	f         *os.File
	BlockSize int
}

// NewDisk opens (and creates, if needed) the file at path, using blockSize.
func NewDisk(path string, blockSize int) (*Disk, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return &Disk{f: f, BlockSize: blockSize}, nil
}

// ReadBlock reads the block at index idx.
func (d *Disk) ReadBlock(idx int64) ([]byte, error) {
	buf := make([]byte, d.BlockSize)
	if _, err := d.f.Seek(idx*int64(d.BlockSize), 0); err != nil {
		return nil, err
	}
	if _, err := d.f.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteBlock writes data to block idx and fsyncs.
func (d *Disk) WriteBlock(idx int64, data []byte) error {
	if _, err := d.f.Seek(idx*int64(d.BlockSize), 0); err != nil {
		return err
	}
	if _, err := d.f.Write(data); err != nil {
		return err
	}
	return d.f.Sync()
}
