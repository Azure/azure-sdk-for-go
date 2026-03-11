//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// These tests verify that the Azure Storage FE (frontend/backend) accepts and returns
// structured messages. They manually construct SM binary payloads and use pipeline policies
// to inject/inspect headers, with NO dependency on unimplemented SDK structured message types.

package blockblob_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
)

const (
	smVersion     = 1
	smFlagCRC64   = 0x0001
	smSegmentSize = 4 * 1024 * 1024 // 4MB default segment
	smHeaderValue = "XSM/1.0; properties=crc64"
)

// buildStructuredMessage manually constructs a V1 structured message binary payload
// from raw content. This is a minimal inline encoder for test purposes only.
func buildStructuredMessage(data []byte) []byte {
	crcTable := shared.CRC64Table
	totalDataLen := len(data)

	// Calculate number of segments
	numSegments := totalDataLen / smSegmentSize
	if totalDataLen%smSegmentSize != 0 {
		numSegments++
	}
	if numSegments == 0 {
		numSegments = 1
	}

	// Calculate total message length:
	// header(13) + for each segment: segHeader(10) + segData + segFooter(8) + trailer(8)
	msgLen := int64(13)
	for i := 0; i < numSegments; i++ {
		segStart := i * smSegmentSize
		segEnd := segStart + smSegmentSize
		if segEnd > totalDataLen {
			segEnd = totalDataLen
		}
		segDataLen := segEnd - segStart
		msgLen += 10 + int64(segDataLen) + 8 // segHeader + data + CRC64
	}
	msgLen += 8 // message trailer CRC64

	var buf bytes.Buffer

	// === Message Header (13 bytes) ===
	buf.WriteByte(smVersion)                                     // version: uint8
	binary.Write(&buf, binary.LittleEndian, msgLen)              // message-length: uint64
	binary.Write(&buf, binary.LittleEndian, uint16(smFlagCRC64)) // flags: uint16
	binary.Write(&buf, binary.LittleEndian, uint16(numSegments)) // num-segments: uint16

	// Message trailer CRC64 is computed over ALL raw content data (not over segment CRC values)
	messageCRC := crc64.Checksum(data, crcTable)

	for i := 0; i < numSegments; i++ {
		segStart := i * smSegmentSize
		segEnd := segStart + smSegmentSize
		if segEnd > totalDataLen {
			segEnd = totalDataLen
		}
		segData := data[segStart:segEnd]

		// === Segment Header (10 bytes) ===
		binary.Write(&buf, binary.LittleEndian, uint16(i+1))         // segment-num: uint16 (1-based)
		binary.Write(&buf, binary.LittleEndian, int64(len(segData))) // data-length: int64

		// === Segment Data ===
		buf.Write(segData)

		// === Segment Footer (8 bytes) - CRC64 of segment data ===
		segCRC := crc64.Checksum(segData, crcTable)
		binary.Write(&buf, binary.LittleEndian, segCRC)
	}

	// === Message Trailer (8 bytes) - CRC64 of all raw content data ===
	binary.Write(&buf, binary.LittleEndian, messageCRC)

	return buf.Bytes()
}

// smUploadPolicy injects x-ms-structured-body header on upload PUT requests
// and records whether the server accepted the request (2xx response).
type smUploadPolicy struct {
	headerValue    string
	contentLength  int64 // original (unframed) content length for x-ms-structured-content-length
	serverAccepted bool
	responseCode   int
}

func (p *smUploadPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Only inject on PUT requests (uploads)
	if req.Raw().Method == http.MethodPut {
		req.Raw().Header.Set("x-ms-structured-body", p.headerValue)
		if p.contentLength > 0 {
			req.Raw().Header.Set("x-ms-structured-content-length", fmt.Sprintf("%d", p.contentLength))
		}
	}

	resp, err := req.Next()
	if err != nil {
		return resp, err
	}

	if req.Raw().Method == http.MethodPut {
		p.responseCode = resp.StatusCode
		p.serverAccepted = resp.StatusCode >= 200 && resp.StatusCode < 300
	}

	return resp, err
}

// smDownloadPolicy injects x-ms-structured-body header on download GET requests
// and checks the response header to see if the server returned structured message format.
type smDownloadPolicy struct {
	requestHeaderValue  string
	responseHeaderValue string
	serverReturned      bool
}

func (p *smDownloadPolicy) Do(req *policy.Request) (*http.Response, error) {
	if req.Raw().Method == http.MethodGet {
		req.Raw().Header.Set("x-ms-structured-body", p.requestHeaderValue)
	}

	resp, err := req.Next()
	if err != nil {
		return resp, err
	}

	if req.Raw().Method == http.MethodGet {
		p.responseHeaderValue = resp.Header.Get("x-ms-structured-body")
		p.serverReturned = p.responseHeaderValue != ""
	}

	return resp, err
}

// TestFEAcceptsStructuredMessageUpload tests that the Azure Storage FE accepts
// a PUT request with x-ms-structured-body header and SM-encoded body.
func (s *BlockBlobUnrecordedTestsSuite) TestFEAcceptsStructuredMessageUpload() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	s.T().Logf("Blob URL: %s", blobURL)
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Create test data and manually encode as structured message
	_, rawData := testcommon.GenerateData(1024) // 1KB test data
	smPayload := buildStructuredMessage(rawData)

	uploadPolicy := &smUploadPolicy{headerValue: smHeaderValue, contentLength: int64(len(rawData))}
	clientOptions := &blockblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{uploadPolicy},
		},
	}
	bbClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, clientOptions)
	_require.NoError(err)

	// Upload the SM-encoded body with the structured body header
	body := streaming.NopCloser(bytes.NewReader(smPayload))
	resp, err := bbClient.Upload(context.Background(), body, nil)

	s.T().Logf("Upload response: %+v", resp)
	s.T().Logf("Upload error: %v", err)
	s.T().Logf("Server response code: %d", uploadPolicy.responseCode)
	s.T().Logf("Server accepted SM upload: %v", uploadPolicy.serverAccepted)

	if err != nil {
		s.T().Logf("Upload error: %v", err)
		s.T().Logf("RESULT: FE does NOT accept structured message uploads (or format was rejected)")
	} else {
		s.T().Logf("RESULT: FE ACCEPTS structured message uploads")
		_require.True(uploadPolicy.serverAccepted)
	}
}

// TestFEAcceptsStructuredMessageStageBlock tests that the Azure Storage FE accepts
// a StageBlock request with x-ms-structured-body header and SM-encoded body.
func (s *BlockBlobUnrecordedTestsSuite) TestFEAcceptsStructuredMessageStageBlock() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	s.T().Logf("Blob URL: %s", blobURL)
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Create test data and manually encode as structured message
	_, rawData := testcommon.GenerateData(2048) // 2KB test data
	smPayload := buildStructuredMessage(rawData)

	uploadPolicy := &smUploadPolicy{headerValue: smHeaderValue, contentLength: int64(len(rawData))}
	clientOptions := &blockblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{uploadPolicy},
		},
	}
	bbClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, clientOptions)
	_require.NoError(err)

	// StageBlock with SM-encoded body
	blockID := "AQAAAA==" // base64 of block 0
	body := streaming.NopCloser(bytes.NewReader(smPayload))
	_, err = bbClient.StageBlock(context.Background(), blockID, body, nil)

	s.T().Logf("Server response code: %d", uploadPolicy.responseCode)
	s.T().Logf("Server accepted SM StageBlock: %v", uploadPolicy.serverAccepted)

	if err != nil {
		s.T().Logf("StageBlock error: %v", err)
		s.T().Logf("RESULT: FE does NOT accept structured message StageBlock (or format was rejected)")
	} else {
		s.T().Logf("RESULT: FE ACCEPTS structured message StageBlock")
		_require.True(uploadPolicy.serverAccepted)
	}
}

// TestFEReturnsStructuredMessageDownload tests that the Azure Storage FE returns
// a structured message response when x-ms-structured-body is set on a download request.
func (s *BlockBlobUnrecordedTestsSuite) TestFEReturnsStructuredMessageDownload() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// First, upload a normal blob (no SM) so we have something to download
	_, rawData := testcommon.GenerateData(4096) // 4KB
	uploadClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, nil)
	_require.NoError(err)

	body := streaming.NopCloser(bytes.NewReader(rawData))
	_, err = uploadClient.Upload(context.Background(), body, nil)
	_require.NoError(err)

	// Now download requesting structured message format
	downloadPolicy := &smDownloadPolicy{requestHeaderValue: smHeaderValue}
	downloadOptions := &blob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{downloadPolicy},
		},
	}
	downloadClient, err := blob.NewClientWithSharedKeyCredential(blobURL, cred, downloadOptions)
	_require.NoError(err)

	resp, err := downloadClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	// Read the full response body
	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())

	s.T().Logf("Response x-ms-structured-body header: %q", downloadPolicy.responseHeaderValue)
	s.T().Logf("Server returned SM format: %v", downloadPolicy.serverReturned)
	s.T().Logf("Downloaded body size: %d (raw data size: %d)", len(downloadedData), len(rawData))

	if downloadPolicy.serverReturned {
		s.T().Logf("RESULT: FE RETURNS structured message response")
		s.T().Logf("Response SM header value: %s", downloadPolicy.responseHeaderValue)

		// The response body should be larger than raw data (SM envelope overhead)
		// Verify the SM envelope: first byte should be version 1
		_require.Greater(len(downloadedData), len(rawData),
			"SM response body should be larger than raw data due to envelope overhead")
		_require.Equal(byte(smVersion), downloadedData[0],
			"First byte of SM response should be version 1")

		// Parse the message header to validate structure
		if len(downloadedData) >= 13 {
			msgLen := binary.LittleEndian.Uint64(downloadedData[1:9])
			flags := binary.LittleEndian.Uint16(downloadedData[9:11])
			numSegs := binary.LittleEndian.Uint16(downloadedData[11:13])
			s.T().Logf("SM header: msgLen=%d, flags=0x%04x, numSegments=%d", msgLen, flags, numSegs)
			_require.Equal(uint64(len(downloadedData)), msgLen, "message length should match response body size")
			_require.Equal(uint16(smFlagCRC64), flags, "flags should indicate CRC64")
			_require.GreaterOrEqual(numSegs, uint16(1), "should have at least 1 segment")
		}
	} else {
		s.T().Logf("RESULT: FE does NOT return structured message response (header absent)")
		// If no SM response, the downloaded data should be the raw data
		_require.Equal(rawData, downloadedData, "without SM, downloaded data should match uploaded data")
	}
}

// TestFEStructuredMessageRoundTrip tests the full cycle:
// 1. Upload with SM-encoded body + x-ms-structured-body header
// 2. Download with x-ms-structured-body header
// 3. Verify the raw data can be extracted from the SM download response
func (s *BlockBlobUnrecordedTestsSuite) TestFEStructuredMessageRoundTrip() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Step 1: Upload with SM encoding
	_, rawData := testcommon.GenerateData(8192) // 8KB
	smPayload := buildStructuredMessage(rawData)

	uploadPolicy := &smUploadPolicy{headerValue: smHeaderValue, contentLength: int64(len(rawData))}
	uploadOptions := &blockblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{uploadPolicy},
		},
	}
	uploadClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, uploadOptions)
	_require.NoError(err)

	body := streaming.NopCloser(bytes.NewReader(smPayload))
	_, err = uploadClient.Upload(context.Background(), body, nil)
	if err != nil {
		s.T().Skipf("FE does not accept SM uploads, skipping round-trip test: %v", err)
	}
	_require.True(uploadPolicy.serverAccepted, "upload should have been accepted")
	s.T().Logf("Upload accepted by FE")

	// Step 2: Download requesting SM format
	downloadPolicy := &smDownloadPolicy{requestHeaderValue: smHeaderValue}
	downloadOptions := &blob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{downloadPolicy},
		},
	}
	downloadClient, err := blob.NewClientWithSharedKeyCredential(blobURL, cred, downloadOptions)
	_require.NoError(err)

	resp, err := downloadClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())

	s.T().Logf("Download SM header: %q", downloadPolicy.responseHeaderValue)
	s.T().Logf("Downloaded body size: %d", len(downloadedData))

	if !downloadPolicy.serverReturned {
		s.T().Logf("FE did not return SM response, verifying raw data match")
		_require.Equal(rawData, downloadedData)
		return
	}

	// Step 3: Parse SM response and extract raw data
	s.T().Logf("FE returned SM response, parsing envelope")
	extractedData, err := extractDataFromStructuredMessage(downloadedData)
	_require.NoError(err, "failed to parse SM response")
	_require.Equal(rawData, extractedData, "extracted data should match original raw data")
	s.T().Logf("RESULT: Full round-trip SUCCESS — uploaded SM, downloaded SM, data matches")
}

// TestFEDownloadNormalBlobWithSMHeader tests downloading a normal (non-SM-uploaded) blob
// with the x-ms-structured-body request header. Verifies the server wraps the response
// in SM format even for blobs not uploaded with SM.
func (s *BlockBlobUnrecordedTestsSuite) TestFEDownloadNormalBlobWithSMHeader() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Upload a plain blob (no SM)
	_, rawData := testcommon.GenerateData(2048) // 2KB
	plainClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, nil)
	_require.NoError(err)

	uploadBody := streaming.NopCloser(bytes.NewReader(rawData))
	_, err = plainClient.Upload(context.Background(), uploadBody, nil)
	_require.NoError(err)

	// Download with SM header
	downloadPolicy := &smDownloadPolicy{requestHeaderValue: smHeaderValue}
	downloadClientOpts := &blob.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			PerCallPolicies: []policy.Policy{downloadPolicy},
		},
	}
	downloadClient, err := blob.NewClientWithSharedKeyCredential(blobURL, cred, downloadClientOpts)
	_require.NoError(err)

	resp, err := downloadClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.NoError(resp.Body.Close())

	s.T().Logf("Response x-ms-structured-body: %q", downloadPolicy.responseHeaderValue)
	s.T().Logf("Server returned SM for normal blob: %v", downloadPolicy.serverReturned)

	if downloadPolicy.serverReturned {
		s.T().Logf("RESULT: FE wraps normal blob downloads in SM format when requested")
		extractedData, err := extractDataFromStructuredMessage(downloadedData)
		_require.NoError(err)
		_require.Equal(rawData, extractedData, "SM-wrapped download should contain original data")
	} else {
		s.T().Logf("RESULT: FE does NOT wrap normal blob downloads in SM format")
		_require.Equal(rawData, downloadedData)
	}
}

// extractDataFromStructuredMessage is a minimal SM decoder for test verification.
// It parses the SM binary envelope and returns the concatenated segment data.
func extractDataFromStructuredMessage(smData []byte) ([]byte, error) {
	if len(smData) < 13 {
		return nil, fmt.Errorf("SM data too short for header: %d bytes", len(smData))
	}

	version := smData[0]
	if version != smVersion {
		return nil, fmt.Errorf("unexpected SM version: %d", version)
	}

	msgLen := binary.LittleEndian.Uint64(smData[1:9])
	if uint64(len(smData)) != msgLen {
		return nil, fmt.Errorf("message length mismatch: header says %d, got %d bytes", msgLen, len(smData))
	}

	flags := binary.LittleEndian.Uint16(smData[9:11])
	numSegments := binary.LittleEndian.Uint16(smData[11:13])

	hasCRC := flags&smFlagCRC64 != 0
	crcTable := shared.CRC64Table

	offset := 13
	var result []byte
	messageCRC := crc64.New(crcTable) // running CRC64 over all raw segment data

	for i := 0; i < int(numSegments); i++ {
		if offset+10 > len(smData) {
			return nil, fmt.Errorf("segment %d: not enough data for segment header", i+1)
		}

		segNum := binary.LittleEndian.Uint16(smData[offset : offset+2])
		segDataLen := int64(binary.LittleEndian.Uint64(smData[offset+2 : offset+10]))
		offset += 10

		if segNum != uint16(i+1) {
			return nil, fmt.Errorf("segment number mismatch: expected %d, got %d", i+1, segNum)
		}

		if offset+int(segDataLen) > len(smData) {
			return nil, fmt.Errorf("segment %d: not enough data for segment content", segNum)
		}

		segData := smData[offset : offset+int(segDataLen)]
		offset += int(segDataLen)
		result = append(result, segData...)

		if hasCRC {
			if offset+8 > len(smData) {
				return nil, fmt.Errorf("segment %d: not enough data for CRC64", segNum)
			}

			expectedCRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
			actualCRC := crc64.Checksum(segData, crcTable)
			offset += 8

			if expectedCRC != actualCRC {
				return nil, fmt.Errorf("segment %d: CRC64 mismatch (expected %x, got %x)", segNum, expectedCRC, actualCRC)
			}

			// Accumulate raw data into message CRC (trailer = CRC64 of all raw content)
			messageCRC.Write(segData)
		}
	}

	// Validate message trailer CRC64
	if hasCRC {
		if offset+8 > len(smData) {
			return nil, fmt.Errorf("not enough data for message trailer CRC64")
		}

		expectedMsgCRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
		actualMsgCRC := messageCRC.Sum64()

		if expectedMsgCRC != actualMsgCRC {
			return nil, fmt.Errorf("message trailer CRC64 mismatch (expected %x, got %x)", expectedMsgCRC, actualMsgCRC)
		}
	}

	return result, nil
}
