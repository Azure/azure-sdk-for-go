// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package arrow

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/apache/arrow-go/v18/arrow"
	arrowArray "github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

// benchSchema returns a schema with representative fields for benchmarking.
func benchSchema(md *arrow.Metadata) *arrow.Schema {
	fields := []arrow.Field{
		{Name: "Name", Type: arrow.BinaryTypes.String, Nullable: false},
		{Name: "Creation-Time", Type: &arrow.TimestampType{Unit: arrow.Microsecond}, Nullable: true},
		{Name: "Last-Modified", Type: &arrow.TimestampType{Unit: arrow.Microsecond}, Nullable: true},
		{Name: "BlobType", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "Etag", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "Content-Length", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "Content-Type", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "ServerEncrypted", Type: arrow.FixedWidthTypes.Boolean, Nullable: true},
		{Name: "AccessTier", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "LeaseState", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "LeaseStatus", Type: arrow.BinaryTypes.String, Nullable: true},
	}
	return arrow.NewSchema(fields, md)
}

// buildBenchArrowStream creates a synthetic Arrow IPC stream with n blobs.
func buildBenchArrowStream(b *testing.B, n int) []byte {
	b.Helper()
	alloc := memory.DefaultAllocator
	md := arrow.NewMetadata(
		[]string{"NextMarker", "NumberOfRecords"},
		[]string{"", strconv.Itoa(n)},
	)
	schema := benchSchema(&md)

	var buf bytes.Buffer
	w := ipc.NewWriter(&buf, ipc.WithSchema(schema), ipc.WithAllocator(alloc))

	builder := arrowArray.NewRecordBuilder(alloc, schema)
	defer builder.Release()

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("blob-%05d.txt", i)
		builder.Field(0).(*arrowArray.StringBuilder).Append(name)
		builder.Field(1).(*arrowArray.TimestampBuilder).Append(arrow.Timestamp(1700000000000000 + int64(i)*1000000))
		builder.Field(2).(*arrowArray.TimestampBuilder).Append(arrow.Timestamp(1700000000000000 + int64(i)*1000000))
		builder.Field(3).(*arrowArray.StringBuilder).Append("BlockBlob")
		builder.Field(4).(*arrowArray.StringBuilder).Append(fmt.Sprintf("0x%08x", i))
		builder.Field(5).(*arrowArray.Uint64Builder).Append(uint64(1024 + i))
		builder.Field(6).(*arrowArray.StringBuilder).Append("application/octet-stream")
		builder.Field(7).(*arrowArray.BooleanBuilder).Append(true)
		builder.Field(8).(*arrowArray.StringBuilder).Append("Hot")
		builder.Field(9).(*arrowArray.StringBuilder).Append("available")
		builder.Field(10).(*arrowArray.StringBuilder).Append("unlocked")
	}

	rec := builder.NewRecordBatch()
	defer rec.Release()

	if err := w.Write(rec); err != nil {
		b.Fatal(err)
	}
	if err := w.Close(); err != nil {
		b.Fatal(err)
	}

	return buf.Bytes()
}

// buildBenchXMLStream creates a synthetic XML ListBlobsFlatSegmentResponse with n blobs.
func buildBenchXMLStream(b *testing.B, n int) []byte {
	b.Helper()
	type xmlBlobProperties struct {
		XMLName        xml.Name `xml:"Properties"`
		CreationTime   string   `xml:"Creation-Time"`
		LastModified   string   `xml:"Last-Modified"`
		BlobType       string   `xml:"BlobType"`
		Etag           string   `xml:"Etag"`
		ContentLength  int64    `xml:"Content-Length"`
		ContentType    string   `xml:"Content-Type"`
		ServerEncrypted bool   `xml:"ServerEncrypted"`
		AccessTier     string   `xml:"AccessTier"`
		LeaseState     string   `xml:"LeaseState"`
		LeaseStatus    string   `xml:"LeaseStatus"`
	}
	type xmlBlob struct {
		Name       string             `xml:"Name"`
		Properties xmlBlobProperties  `xml:"Properties"`
	}
	type xmlBlobs struct {
		Blob []xmlBlob `xml:"Blob"`
	}
	type xmlResponse struct {
		XMLName xml.Name `xml:"EnumerationResults"`
		Blobs   xmlBlobs `xml:"Blobs"`
	}

	resp := xmlResponse{}
	for i := 0; i < n; i++ {
		resp.Blobs.Blob = append(resp.Blobs.Blob, xmlBlob{
			Name: fmt.Sprintf("blob-%05d.txt", i),
			Properties: xmlBlobProperties{
				CreationTime:   "Mon, 01 Jan 2024 00:00:00 GMT",
				LastModified:   "Mon, 01 Jan 2024 00:00:00 GMT",
				BlobType:       "BlockBlob",
				Etag:           fmt.Sprintf("0x%08x", i),
				ContentLength:  int64(1024 + i),
				ContentType:    "application/octet-stream",
				ServerEncrypted: true,
				AccessTier:     "Hot",
				LeaseState:     "available",
				LeaseStatus:    "unlocked",
			},
		})
	}

	data, err := xml.Marshal(resp)
	if err != nil {
		b.Fatal(err)
	}
	return data
}

func BenchmarkArrowParsing(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 5000} {
		data := buildBenchArrowStream(b, count)
		b.Run(fmt.Sprintf("Arrow_%d_blobs", count), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				resp := makeHTTPResponse(data, nil)
				_, err := HandleFlatListResponse(resp)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkXMLParsing(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 5000} {
		data := buildBenchXMLStream(b, count)
		b.Run(fmt.Sprintf("XML_%d_blobs", count), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var result generated.ListBlobsFlatSegmentResponse
				if err := xml.Unmarshal(data, &result); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
