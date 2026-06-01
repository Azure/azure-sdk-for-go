// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

// Shared CLI flags exposed by every storage perf test. Defaults preserve
// the historical behavior so existing perf matrices keep working unchanged.
var (
	// commonBlockSize is the SDK BlockSize for chunked upload/download paths.
	// 0 means "use the SDK default".
	commonBlockSize int64

	// commonConcurrency is the SDK Concurrency (parallel blocks per operation)
	// for chunked upload/download paths. 0 means "use the SDK default".
	//
	// Typed as uint16 because the strictest SDK option that consumes this
	// value (blockblob.UploadBufferOptions.Concurrency,
	// blob.DownloadBufferOptions.Concurrency) is uint16. Using uint16 here
	// rules out silent truncation/overflow at the cast sites and lets the
	// CLI fail fast (via uint16Flag.Set) for values outside [0, 65535].
	commonConcurrency uint16

	// uploadMethod selects which upload API the upload test exercises:
	//   single  -> blockblob.Client.Upload       (single REST PUT; default)
	//   buffer  -> blockblob.Client.UploadBuffer (parallel staged blocks)
	//   stream  -> blockblob.Client.UploadStream (chunked from io.Reader)
	uploadMethod string

	// downloadMethod selects which download API the download test exercises:
	//   stream  -> blob.Client.DownloadStream (single GET; default)
	//   buffer  -> blob.Client.DownloadBuffer (parallel ranged GETs)
	downloadMethod string

	// listPageSize is the MaxResults page size used by the list test.
	listPageSize int
)

func init() {
	flag.Int64Var(&commonBlockSize, "block-size", 0, "Block/chunk size in bytes for chunked upload/download (0=SDK default).")
	flag.Var((*uint16Flag)(&commonConcurrency), "concurrency", "Max number of blocks transferred in parallel per chunked upload/download operation (0=SDK default; max 65535).")
	flag.StringVar(&uploadMethod, "upload-method", "single", "Upload method: single|buffer|stream.")
	flag.StringVar(&downloadMethod, "download-method", "stream", "Download method: stream|buffer.")
	flag.IntVar(&listPageSize, "page-size", 5000, "MaxResults page size used by list operations.")
}

// uint16Flag is a flag.Value backed by a uint16 so the CLI rejects values
// outside [0, 65535] at parse time instead of silently truncating at the
// SDK option cast site.
type uint16Flag uint16

func (u *uint16Flag) String() string {
	if u == nil {
		return "0"
	}
	return strconv.FormatUint(uint64(*u), 10)
}

func (u *uint16Flag) Set(value string) error {
	parsed, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid --concurrency value %q (must be a uint16 in [0, 65535]): %w", value, err)
	}
	*u = uint16Flag(parsed)
	return nil
}

// generateRandomBytes returns size bytes of random data, suitable for use as an
// in-memory payload for upload tests.
//
// size is typed as int64 to match the rest of the perf code (e.g.
// uploadTestGlobal.size) and to make overflow on 32-bit platforms (or from a
// JSON workload-config value beyond math.MaxInt) a deterministic, actionable
// error rather than a silent truncation inside make.
func generateRandomBytes(size int64) ([]byte, error) {
	if size < 0 {
		return nil, fmt.Errorf("invalid size %d", size)
	}
	if size > int64(math.MaxInt) {
		return nil, fmt.Errorf("size %d exceeds platform max int (%d); use a streaming method for payloads larger than %d bytes", size, int64(math.MaxInt), int64(math.MaxInt))
	}
	buf := make([]byte, size)
	if size == 0 {
		return buf, nil
	}
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// randomSeedSize is the size of the repeating random seed used to back
// streaming/single-PUT upload payloads. Picked large enough to amortize the
// per-iteration tiling cost without holding the full --size in RAM.
const randomSeedSize = 1 << 20 // 1 MiB

// randomStream is an io.ReadSeekCloser that produces `total` bytes by tiling
// `seed` (a small random buffer) repeatedly. It allows multi-GiB / TiB upload
// payloads without allocating the full payload in memory and without per-byte
// crypto/rand calls on the hot path.
//
// Safe for use by a single goroutine. Each parallel perf-test goroutine should
// own its own randomStream so the read offset is goroutine-local.
type randomStream struct {
	seed  []byte
	total int64
	pos   int64
}

// newRandomStream returns a randomStream that yields `total` bytes by tiling
// the provided seed (which must be non-empty).
func newRandomStream(seed []byte, total int64) *randomStream {
	return &randomStream{seed: seed, total: total}
}

func (r *randomStream) Read(p []byte) (int, error) {
	if r.pos >= r.total {
		return 0, io.EOF
	}
	remaining := r.total - r.pos
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}
	n := 0
	for n < len(p) {
		off := int((r.pos + int64(n)) % int64(len(r.seed)))
		c := copy(p[n:], r.seed[off:])
		n += c
	}
	r.pos += int64(n)
	return n, nil
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.pos + offset
	case io.SeekEnd:
		abs = r.total + offset
	default:
		return 0, fmt.Errorf("randomStream.Seek: invalid whence %d", whence)
	}
	if abs < 0 {
		return 0, fmt.Errorf("randomStream.Seek: negative position %d", abs)
	}
	r.pos = abs
	return abs, nil
}

func (r *randomStream) Close() error { return nil }

// availableSystemMemoryBytes returns the amount of physical memory the kernel
// believes is currently available to user processes, in bytes. Returns 0 if
// the value can't be determined (e.g. non-Linux OS, restricted container).
//
// On Linux it parses /proc/meminfo's MemAvailable line. On other platforms it
// returns 0 so callers should treat 0 as "unknown" rather than "no memory".
func availableSystemMemoryBytes() uint64 {
	if runtime.GOOS != "linux" {
		return 0
	}
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "MemAvailable:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0
		}
		kb, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0
		}
		return kb * 1024
	}
	return 0
}

// checkBufferMemoryBudget returns an error when the caller requests an
// in-memory buffer for upload or download whose total size would exceed a
// safe fraction (80%) of available system memory. This guards against the
// runtime OOM that would otherwise abort the process before any HTTP request
// is made.
//
// flagLabel is the flag-and-value pair that triggered the check (e.g.
// "--upload-method buffer" or "--download-method buffer") and is embedded in
// the error message so the text matches the invoked flag.
//
// sizeBytes is the total number of bytes the caller intends to allocate
// across all goroutines (upload shares one payload across --parallel, so the
// caller passes the per-payload size; download allocates one buffer per
// goroutine, so the caller passes size * parallel).
func checkBufferMemoryBudget(flagLabel string, sizeBytes int64) error {
	if sizeBytes <= 0 {
		return nil
	}
	avail := availableSystemMemoryBytes()
	if avail == 0 {
		return nil // unknown — skip the check rather than guess
	}
	const safetyFraction = 0.8
	budget := uint64(float64(avail) * safetyFraction)
	if uint64(sizeBytes) > budget {
		return fmt.Errorf(
			"%s requires materializing %d bytes (%.2f GiB) in RAM, "+
				"which exceeds %d%% of available system memory (%d bytes, %.2f GiB available). "+
				"Use a streaming method for large payloads, or reduce --size/--parallel",
			flagLabel, sizeBytes, float64(sizeBytes)/(1<<30),
			int(safetyFraction*100), avail, float64(avail)/(1<<30),
		)
	}
	return nil
}
