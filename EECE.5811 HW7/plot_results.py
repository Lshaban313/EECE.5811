import matplotlib
matplotlib.use('Agg')
import sys
import pandas as pd
import matplotlib.pyplot as plt

if len(sys.argv) != 2:
    print("Usage: python3 plot_results.py results.csv")
    sys.exit(1)

df = pd.read_csv(sys.argv[1])
print("DataFrame:\n", df, "\n")   # for debug

# Detect which format we have:
if 'write_ns_per_block' in df.columns:
    # Old nanoseconds‐based CSV
    df['write_ms_per_block'] = df['write_ns_per_block'] / 1e6
    df['read_ms_per_block']  = df['read_ns_per_block']  / 1e6
    x = df['level']
    yw = df['write_ms_per_block']
    yr = df['read_ms_per_block']
    y_label = 'Time per Block (ms)'
    title_w = 'RAID Write Time per Block'
    title_r = 'RAID Read Time per Block'
else:

    if 'WriteThroughput(MB/s)' in df.columns:
        x = df['RAIDLevel']
        yw = df['WriteThroughput(MB/s)']
        yr = df['ReadThroughput(MB/s)']
        y_label = 'Throughput (MB/s)'
        title_w = 'RAID Write Throughput'
        title_r = 'RAID Read Throughput'
    else:
        # fallback to seconds per total
        x = df['RAIDLevel']
        yw = df['WriteTime(s)']
        yr = df['ReadTime(s)']
        y_label = 'Time (s)'
        title_w = 'RAID Write Time'
        title_r = 'RAID Read Time'

# Plot write
plt.figure()
plt.bar(x.astype(str), yw)
plt.xlabel('RAID Level')
plt.ylabel(y_label)
plt.title(title_w)
plt.tight_layout()
plt.savefig('write_performance.png')
print("Saved write_performance.png")

# Plot read
plt.figure()
plt.bar(x.astype(str), yr)
plt.xlabel('RAID Level')
plt.ylabel(y_label)
plt.title(title_r)
plt.tight_layout()
plt.savefig('read_performance.png')
print("Saved read_performance.png")
