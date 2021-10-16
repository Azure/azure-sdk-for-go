package azfile

import (
	"strings"
)

// FileAttributeFlags are a subset of the WinNT.h file attribute bitflags.
// The uint32 returned from golang.org/x/sys/windows' GetFileAttributes func can be converted directly without worry, since both types are uint32, and the numeric values line up.
// However, crafting these manually is acceptable via the .Add function. Removing properties is done via .Remove. Checking if one contains another can be done with .Has
// The subset is listed at: https://docs.microsoft.com/en-us/rest/api/storageservices/set-file-properties#file-system-attributes
// and the values are listed at: https://www.magnumdb.com/search?q=filename%3Awinnt.h+AND+FILE_ATTRIBUTE_*
// This is intended for easy conversion from the following function: https://pkg.go.dev/golang.org/x/sys/windows?tab=doc#GetFileAttributes
type FileAttributeFlags uint32

const (
	FileAttributeNone              FileAttributeFlags = 0
	FileAttributeReadonly          FileAttributeFlags = 1
	FileAttributeHidden            FileAttributeFlags = 2
	FileAttributeSystem            FileAttributeFlags = 4
	FileAttributeArchive           FileAttributeFlags = 32
	FileAttributeTemporary         FileAttributeFlags = 256
	FileAttributeOffline           FileAttributeFlags = 4096
	FileAttributeNotContentIndexed FileAttributeFlags = 8192
	FileAttributeNoScrubData       FileAttributeFlags = 131072
)

func (f FileAttributeFlags) String() (out string) {
	// We choose not to do a map here, as indexing over a map doesn't inherently retain order.
	attrFlags := []FileAttributeFlags{
		FileAttributeReadonly,
		FileAttributeHidden,
		FileAttributeSystem,
		FileAttributeArchive,
		FileAttributeTemporary,
		FileAttributeOffline,
		FileAttributeNotContentIndexed,
		FileAttributeNoScrubData,
	}
	attrStrings := []string{
		"ReadOnly",
		"Hidden",
		"System",
		"Archive",
		"Temporary",
		"Offline",
		"NotContentIndexed",
		"NoScrubData",
	}

	for idx, flag := range attrFlags {
		if f.Has(flag) {
			out += attrStrings[idx] + "|"
		}
	}

	out = strings.TrimSuffix(out, "|")

	if out == "" {
		out = defaultFileAttributes
	}

	return
}

func (f FileAttributeFlags) Add(new FileAttributeFlags) FileAttributeFlags {
	return f | new
}

func (f FileAttributeFlags) Remove(old FileAttributeFlags) FileAttributeFlags {
	return f &^ old
}

func (f FileAttributeFlags) Has(item FileAttributeFlags) bool {
	return f&item == item
}

// ParseFileAttributeFlagsString parses the service-side file attribute strings that the above enum strongly types.
func ParseFileAttributeFlagsString(input string) (out FileAttributeFlags) {
	// We don't worry about the order here, since the resulting bitflags will automagically be in order.
	attrStrings := map[string]FileAttributeFlags{
		"none":              FileAttributeNone,
		"readonly":          FileAttributeReadonly,
		"hidden":            FileAttributeHidden,
		"system":            FileAttributeSystem,
		"archive":           FileAttributeArchive,
		"temporary":         FileAttributeTemporary,
		"offline":           FileAttributeOffline,
		"notcontentindexed": FileAttributeNotContentIndexed,
		"noscrubdata":       FileAttributeNoScrubData,
	}

	for _, v := range strings.Split(input, "|") {
		// We trim the space because the service returns the flags back with spaces in between the pipes
		// We also lowercase out of an abundance of caution to ensure we're getting what we think we're getting.
		key := strings.ToLower(strings.TrimSpace(v))
		if val, ok := attrStrings[key]; ok {
			out = out.Add(val)
		} else if key == "directory" {
			// just skip it. Users this SDK will already know whether it's a directory or a file, and will have made a getProperties call on the appropriate object type
		} else {
			panic("service sided attribute flags should never fail")
		}
	}

	return
}
