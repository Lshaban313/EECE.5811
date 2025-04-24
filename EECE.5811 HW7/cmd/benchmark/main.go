package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "raid-sim/raid"
)

func main() {
    // Command-line flags
    bs := flag.Int("bs", 4096, "block size in bytes")
    totalSize := flag.Int("size", 100000000, "total bytes to write")
    flag.Parse()

    // Compute number of blocks
    blocks := int64(*totalSize) / int64(*bs)
    if blocks < 1 {
        blocks = 1
    }

    // Paths to disk files
    files := []string{"disk0.dat", "disk1.dat", "disk2.dat", "disk3.dat", "disk4.dat"}

    // Prepare one data block
    data := make([]byte, *bs)

    // Create results.csv
    f, err := os.Create("results.csv")
    if err != nil {
        log.Fatalf("failed to create results.csv: %v", err)
    }
    defer f.Close()

    writer := csv.NewWriter(f)
    defer func() {
        writer.Flush()
        if err := writer.Error(); err != nil {
            log.Fatalf("error flushing CSV: %v", err)
        }
    }()

    // CSV header
    if err := writer.Write([]string{"level", "write_ns_per_block", "read_ns_per_block"}); err != nil {
        log.Fatalf("failed to write CSV header: %v", err)
    }

    // Benchmark each RAID level
    for _, lvl := range []int{0, 1, 4, 5} {
        fmt.Printf("Benchmarking RAID %d: %d blocks of %d bytes\n", lvl, blocks, *bs)

        // Open RAID
        r, err := raid.Open(lvl, files, *bs)
        if err != nil {
            log.Fatalf("open RAID level %d: %v", lvl, err)
        }

        // -- Write phase --
        startW := time.Now()
        for i := int64(0); i < blocks; i++ {
            if err := r.Write(i, data); err != nil {
                log.Fatalf("write error, level %d block %d: %v", lvl, i, err)
            }
        }
        writeDur := time.Since(startW).Nanoseconds()
        perWrite := writeDur / blocks

        // -- Read phase --
        startR := time.Now()
        for i := int64(0); i < blocks; i++ {
            if _, err := r.Read(i); err != nil {
                log.Fatalf("read error, level %d block %d: %v", lvl, i, err)
            }
        }
        readDur := time.Since(startR).Nanoseconds()
        perRead := readDur / blocks

        // Write CSV record
        record := []string{
            strconv.Itoa(lvl),
            strconv.FormatInt(perWrite, 10),
            strconv.FormatInt(perRead, 10),
        }
        if err := writer.Write(record); err != nil {
            log.Fatalf("failed to write CSV row for level %d: %v", lvl, err)
        }
    }

    fmt.Println("Done — results.csv written.")
}
