// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/stretchr/testify/require"
)

// newMockListClient returns a Client whose transport returns body as a single 200 OK List Directories and Files response.
func newMockListClient(t *testing.T, body string) (*directory.Client, *testcommon.MockTransport) {
	transport := testcommon.NewMockTransport(testcommon.MockResponse{StatusCode: http.StatusOK, Body: body})
	client, err := directory.NewClientWithNoCredential("https://account.file.core.windows.net/share/dir", &directory.ClientOptions{
		ClientOptions: testcommon.MockClientOptions(transport),
	})
	require.NoError(t, err)
	return client, transport
}

func listFirstPage(t *testing.T, client *directory.Client, options *directory.ListFilesAndDirectoriesOptions) directory.ListFilesAndDirectoriesResponse {
	pager := client.NewListFilesAndDirectoriesPager(options)
	require.True(t, pager.More())
	resp, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	return resp
}

func TestListFilesAndDirectoriesNFSAllItemTypesMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, nfsListFilesAndDirectoriesXML)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment
	_require.NotNil(segment)
	_require.Equal("9223442405599084545", *resp.DirectoryID)

	// 1 Directory + 4 Files + 1 SymLink + 1 BlockDevice + 1 CharDevice + 1 Fifo + 1 Socket = 10.
	_require.Len(segment.Directories, 1)
	_require.Len(segment.Files, 4)
	_require.Len(segment.SymLinkItems, 1)
	_require.Len(segment.BlockDeviceItems, 1)
	_require.Len(segment.CharDeviceItems, 1)
	_require.Len(segment.FifoItems, 1)
	_require.Len(segment.SocketItems, 1)

	// Directory
	dir := segment.Directories[0]
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.Equal("subdir", *dir.Name)
	_require.Equal("12682206919419625485", *dir.ID)
	_require.Equal(int64(2), *dir.LinkCount)
	_require.Nil(dir.Properties.ContentLength)
	_require.Equal("1000", *dir.Properties.Owner)
	_require.Equal("1000", *dir.Properties.Group)
	_require.Equal("0755", *dir.Properties.FileMode)
	_require.Equal(azcore.ETag(`"0x8DEAF1479E1C087"`), *dir.Properties.ETag)
	_require.NotNil(dir.Properties.CreationTime)
	_require.NotNil(dir.Properties.LastModified)

	// Files
	expectedFiles := []struct {
		name      string
		id        string
		size      int64
		linkCount int64
	}{
		{"hardlink.txt", "13835128424026472451", 80, 2},
		{"regular.txt", "13835128424026472451", 80, 2},
		{"regularClose.txt", "9799903157902508049", 14, 1},
		{"regularOpen.txt", "17293892937847013391", 14, 1},
	}
	for i, expected := range expectedFiles {
		f := segment.Files[i]
		_require.Equal(directory.FileTypeRegular, f.FileType())
		_require.Equal(expected.name, *f.Name)
		_require.Equal(expected.id, *f.ID)
		_require.Equal(expected.size, *f.Properties.ContentLength)
		_require.Equal(expected.linkCount, *f.LinkCount)
		_require.Equal("1000", *f.Properties.Owner)
		_require.Equal("1000", *f.Properties.Group)
		_require.Equal("0644", *f.Properties.FileMode)
		_require.Nil(f.Attributes)
		_require.Nil(f.PermissionKey)
	}

	// SymLink
	symlink := segment.SymLinkItems[0]
	_require.Equal(directory.FileTypeSymLink, symlink.FileType())
	_require.Equal("symlink.txt", *symlink.Name.Content)
	_require.Nil(symlink.Name.Encoded)
	_require.Equal("10376363910205931529", *symlink.FileID)
	_require.Equal(int64(1), *symlink.LinkCount)
	_require.Equal("/mnt/s2/dir2/regular.txt", *symlink.LinkText)
	_require.Equal("1000", *symlink.Properties.Owner)
	_require.Equal("1000", *symlink.Properties.Group)
	_require.Equal("0777", *symlink.Properties.FileMode)

	// BlockDevice
	block := segment.BlockDeviceItems[0]
	_require.Equal(directory.FileTypeBlockDevice, block.FileType())
	_require.Equal("block_device", *block.Name.Content)
	_require.Equal("10952824662509355033", *block.FileID)
	_require.Equal(int64(1), *block.LinkCount)
	_require.Equal(int64(8), *block.DeviceMajor)
	_require.Equal(int64(0), *block.DeviceMinor)
	_require.Equal("0", *block.Properties.Owner)
	_require.Equal("0", *block.Properties.Group)
	_require.Equal("0640", *block.Properties.FileMode)

	// CharDevice
	charDev := segment.CharDeviceItems[0]
	_require.Equal(directory.FileTypeCharacterDevice, charDev.FileType())
	_require.Equal("char_device", *charDev.Name.Content)
	_require.Equal("16717432185543589911", *charDev.FileID)
	_require.Equal(int64(1), *charDev.LinkCount)
	_require.Equal(int64(1), *charDev.DeviceMajor)
	_require.Equal(int64(7), *charDev.DeviceMinor)
	_require.Equal("0", *charDev.Properties.Owner)
	_require.Equal("0", *charDev.Properties.Group)
	_require.Equal("0644", *charDev.Properties.FileMode)

	// Fifo
	fifo := segment.FifoItems[0]
	_require.Equal(directory.FileTypeFifo, fifo.FileType())
	_require.Equal("fifo_pipe", *fifo.Name.Content)
	_require.Equal("14988049928633319435", *fifo.FileID)
	_require.Equal(int64(1), *fifo.LinkCount)
	_require.Equal("1000", *fifo.Properties.Owner)
	_require.Equal("1000", *fifo.Properties.Group)
	_require.Equal("0644", *fifo.Properties.FileMode)

	// Socket
	socket := segment.SocketItems[0]
	_require.Equal(directory.FileTypeSocket, socket.FileType())
	_require.Equal("unix_socket", *socket.Name.Content)
	_require.Equal("16429201809391878183", *socket.FileID)
	_require.Equal(int64(1), *socket.LinkCount)
	_require.Equal("0", *socket.Properties.Owner)
	_require.Equal("0", *socket.Properties.Group)
	_require.Equal("0755", *socket.Properties.FileMode)
}

func TestListFilesAndDirectoriesNFSEncodedNameMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="nfsshare" DirectoryPath="dir" Encoded="true">
	<Entries>
		<BlockDevice>
			<Name Encoded="true">block%EF%BF%BF</Name>
			<Properties />
		</BlockDevice>
		<CharDevice>
			<Name Encoded="true">char%EF%BF%BF</Name>
			<Properties />
		</CharDevice>
		<Fifo>
			<Name Encoded="true">fifo%EF%BF%BF</Name>
			<Properties />
		</Fifo>
		<SymLink>
			<Name Encoded="true">link%EF%BF%BF.txt</Name>
			<Properties />
			<LinkText>target.txt</LinkText>
		</SymLink>
		<Socket>
			<Name Encoded="false">socket%20name</Name>
			<Properties />
		</Socket>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment

	assertDecoded := func(name *directory.StringEncoded, expected string) {
		_require.Equal(expected, *name.Content)
		_require.False(*name.Encoded)
	}
	assertDecoded(segment.BlockDeviceItems[0].Name, "block￿")
	assertDecoded(segment.CharDeviceItems[0].Name, "char￿")
	assertDecoded(segment.FifoItems[0].Name, "fifo￿")
	assertDecoded(segment.SymLinkItems[0].Name, "link￿.txt")
	_require.Equal("target.txt", *segment.SymLinkItems[0].LinkText)
	// A name that is not encoded is returned unchanged.
	assertDecoded(segment.SocketItems[0].Name, "socket%20name")
}

func TestListFilesAndDirectoriesNFSMalformedEncodedNameMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="nfsshare" DirectoryPath="dir" Encoded="true">
	<Entries>
		<SymLink>
			<Name Encoded="true">bad%zzname</Name>
			<Properties />
			<LinkText>target.txt</LinkText>
		</SymLink>
		<Fifo>
			<Name Encoded="true">fifo%EF%BF%BF</Name>
			<Properties />
		</Fifo>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	// A malformed encoding does not fail the page; the name is returned as sent with Encoded still true.
	resp := listFirstPage(t, client, nil)
	symlinkName := resp.Segment.SymLinkItems[0].Name
	_require.Equal("bad%zzname", *symlinkName.Content)
	_require.True(*symlinkName.Encoded)
	_require.Equal("target.txt", *resp.Segment.SymLinkItems[0].LinkText)

	// Other items on the page are still decoded.
	fifoName := resp.Segment.FifoItems[0].Name
	_require.Equal("fifo￿", *fifoName.Content)
	_require.False(*fifoName.Encoded)
}

func TestListFilesAndDirectoriesSMBMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="smbshare" DirectoryPath="dir">
	<Entries>
		<File>
			<Name>file.txt</Name>
			<FileId>13835128424026341376</FileId>
			<Properties>
				<Content-Length>1024</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<Etag>"0x8DEAF15A711D34B"</Etag>
			</Properties>
			<Attributes>Archive</Attributes>
			<PermissionKey>4066528134148476695*1</PermissionKey>
		</File>
		<Directory>
			<Name>subdir</Name>
			<FileId>13835093239654252544</FileId>
			<Properties />
			<Attributes>Directory</Attributes>
			<PermissionKey>4066528134148476695*1</PermissionKey>
		</Directory>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment
	_require.Len(segment.Files, 1)
	_require.Len(segment.Directories, 1)
	_require.Empty(segment.SymLinkItems)
	_require.Empty(segment.BlockDeviceItems)
	_require.Empty(segment.CharDeviceItems)
	_require.Empty(segment.FifoItems)
	_require.Empty(segment.SocketItems)

	f := segment.Files[0]
	_require.Equal(directory.FileTypeRegular, f.FileType())
	_require.Equal("file.txt", *f.Name)
	_require.Equal("Archive", *f.Attributes)
	_require.Equal("4066528134148476695*1", *f.PermissionKey)
	_require.Equal(int64(1024), *f.Properties.ContentLength)
	// NFS-only fields stay nil on SMB
	_require.Nil(f.LinkCount)
	_require.Nil(f.Properties.Owner)
	_require.Nil(f.Properties.Group)
	_require.Nil(f.Properties.FileMode)

	dir := segment.Directories[0]
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.Equal("subdir", *dir.Name)
	_require.Equal("Directory", *dir.Attributes)
	_require.Nil(dir.LinkCount)
	_require.Nil(dir.Properties.Owner)
	_require.Nil(dir.Properties.Group)
	_require.Nil(dir.Properties.FileMode)
}

func TestListFilesAndDirectoriesIncludeQueryMocked(t *testing.T) {
	testCases := []struct {
		name     string
		include  directory.ListFilesInclude
		expected string
	}{
		{"None", directory.ListFilesInclude{}, ""},
		{"SMB", directory.ListFilesInclude{Timestamps: true, ETag: true, Attributes: true, PermissionKey: true}, "Timestamps,Etag,Attributes,PermissionKey"},
		{"NFS", directory.ListFilesInclude{Timestamps: true, ETag: true, Permissions: true, LinkCount: true, NfsAttributes: true}, "Timestamps,Etag,Permissions,LinkCount,NfsAttributes"},
		{"All", directory.ListFilesInclude{All: true}, "All"},
		// All takes precedence and is sent as the literal "All" rather than being expanded.
		{"AllWithOthers", directory.ListFilesInclude{All: true, Timestamps: true, LinkCount: true}, "All"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, transport := newMockListClient(t, emptyListFilesAndDirectoriesXML)
			_ = listFirstPage(t, client, &directory.ListFilesAndDirectoriesOptions{Include: tc.include})
			require.Len(t, transport.Requests, 1)
			query := transport.Requests[0].URL.Query()
			require.Equal(t, tc.expected, query.Get("include"))
			_, ok := query["include"]
			require.Equal(t, tc.expected != "", ok)
		})
	}
}

func TestFileTypeValues(t *testing.T) {
	require.Equal(t, []directory.FileType{
		directory.FileTypeRegular,
		directory.FileTypeDirectory,
		directory.FileTypeSymLink,
		directory.FileTypeBlockDevice,
		directory.FileTypeCharacterDevice,
		directory.FileTypeSocket,
		directory.FileTypeFifo,
	}, directory.PossibleFileTypeValues())
}

const emptyListFilesAndDirectoriesXML = `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="myshare" DirectoryPath="mydir">
	<Entries />
	<NextMarker />
</EnumerationResults>`

// nfsListFilesAndDirectoriesXML is a List Directories and Files response containing every NFS item type.
const nfsListFilesAndDirectoriesXML = `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="http://localhost:10004/testaccount1/" ShareName="nfsshare" DirectoryPath="dir2">
	<DirectoryId>9223442405599084545</DirectoryId>
	<Entries>
		<BlockDevice>
			<Name>block_device</Name>
			<FileId>10952824662509355033</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:33.1472879Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:33.1472879Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:33.1472879Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:33.1472879Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified>
				<Etag>"0x8DEAF1480FD31EF"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0640</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<DeviceMajor>8</DeviceMajor>
			<DeviceMinor>0</DeviceMinor>
		</BlockDevice>
		<CharDevice>
			<Name>char_device</Name>
			<FileId>16717432185543589911</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:33.1120018Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:33.1120018Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:33.1120018Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:33.1120018Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified>
				<Etag>"0x8DEAF1480F7CF92"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<DeviceMajor>1</DeviceMajor>
			<DeviceMinor>7</DeviceMinor>
		</CharDevice>
		<Fifo>
			<Name>fifo_pipe</Name>
			<FileId>14988049928633319435</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:21.1982170Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.1982170Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.1982170Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.1982170Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479DDE95A"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</Fifo>
		<File>
			<Name>hardlink.txt</Name>
			<FileId>13835128424026472451</FileId>
			<Properties>
				<Content-Length>80</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.0288880Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.5326411Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.5326411Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A711D34B"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</File>
		<File>
			<Name>regular.txt</Name>
			<FileId>13835128424026472451</FileId>
			<Properties>
				<Content-Length>80</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.0288880Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.5326411Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.5326411Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A711D34B"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</File>
		<File>
			<Name>regularClose.txt</Name>
			<FileId>9799903157902508049</FileId>
			<Properties>
				<Content-Length>14</Content-Length>
				<CreationTime>2026-05-11T04:19:21.2920068Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2920068Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.6790214Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.6790214Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A7282946"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</File>
		<File>
			<Name>regularOpen.txt</Name>
			<FileId>17293892937847013391</FileId>
			<Properties>
				<Content-Length>14</Content-Length>
				<CreationTime>2026-05-11T04:19:21.2612277Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2612277Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:51.5075180Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:51.5075180Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:51 GMT</Last-Modified>
				<Etag>"0x8DEAF15AA08EE6C"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</File>
		<Directory>
			<Name>subdir</Name>
			<FileId>12682206919419625485</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:21.2233863Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2233863Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.2233863Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.2233863Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479E1C087"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0755</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</Directory>
		<SymLink>
			<Name>symlink.txt</Name>
			<FileId>10376363910205931529</FileId>
			<Properties>
				<Content-Length>24</Content-Length>
				<CreationTime>2026-05-11T04:19:21.1661608Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.1661608Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.1661608Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.1661608Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479D90528"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0777</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<LinkText>/mnt/s2/dir2/regular.txt</LinkText>
		</SymLink>
		<Socket>
			<Name>unix_socket</Name>
			<FileId>16429201809391878183</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:27:46.7308143Z</CreationTime>
				<LastAccessTime>2026-05-11T04:27:46.7308143Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.7308143Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.7308143Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A730106F"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0755</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</Socket>
	</Entries>
	<NextMarker />
</EnumerationResults>`
