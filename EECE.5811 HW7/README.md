# EECE.5811: HW7

> **Group Members**: Vina Dang, Layann Shaban  
> **Due Date**: April 24, 2025  

## Overview

This assignment implements a RAID simulator in Go (RAID-0, RAID-1, RAID-4 and RAID-5) using plain files to represent disks, plus a benchmarking tool that measures per‐block write and read latencies under each level. A companion Python script parses the CSV output and produces bar charts of write‐ and read‐performance for comparison against textbook predictions.



## File Descriptions

- **go.mod**  
  Declares the module path (`raid-sim`) and Go version.

- **raid/disk.go**  
  Implements `type Disk` with `ReadBlock(idx)` and `WriteBlock(idx, data)` methods that `fsync` after each write to simulate real disk delays.

- **raid/raid.go**  
  Defines the `RAID` interface (`Write(blockNum int64, data []byte) error` and `Read(blockNum int64) ([]byte, error)`) and a constructor `Open(level int, files []string, blockSize int)`.

- **raid/raid0.go**, **raid1.go**, **raid4.go**, **raid5.go**  
  Concrete implementations of the `RAID` interface:
  - RAID-0: simple modulo‐based striping
  - RAID-1: full mirroring to all disks
  - RAID-4: single parity disk at fixed position
  - RAID-5: rotating parity disk; parity computed via XOR of data blocks

- **cmd/benchmark/main.go**  
  CLI tool that:
  1. Parses flags:  
     - `-bs` (block size in bytes, default 4096)  
     - `-size` (total bytes to write, default 100_000_000)  
  2. Opens five disk files (`disk0.dat`…`disk4.dat`) via `raid.Open(...)`  
  3. Writes the specified number of blocks (fsync’ing each) and times the operation  
  4. Reads back the same blocks and times the reads  
  5. Emits `results.csv` with columns:  
     ```
     level,write_ns_per_block,read_ns_per_block
     0,XXXX,YYYY
     1,XXXX,YYYY
     4,XXXX,YYYY
     5,XXXX,YYYY
     ```

- **plot_results.py**  
  A standalone Python 3 script that:
  1. Reads `results.csv` into pandas  
  2. Generates two bar charts using matplotlib:  
     - **write_performance.png** (ns or ms per block vs RAID level)  
     - **read_performance.png**  
  3. Saves the PNGs in the working directory

## How to Run

1. **Prepare disk files** (for 100 MiB each):  
  ```
   for i in 0 1 2 3 4; do
     dd if=/dev/zero of=disk$i.dat bs=1M count=100
   done


