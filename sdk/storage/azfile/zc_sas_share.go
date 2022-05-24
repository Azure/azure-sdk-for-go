package azfile

import (
	"bytes"
	"fmt"
)

// The ShareSASPermissions type simplifies creating the permissions string for an Azure Storage share SAS.
// Initialize an instance of this type and then call its String method to set FileSASSignatureValues's FilePermissions field.
type ShareSASPermissions struct {
	Read, Create, Write, Delete, List bool
}

// String produces the SAS permissions string for an Azure Storage share.
// Call this method to set FileSASSignatureValues's FilePermissions field.
func (p ShareSASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
	}
	if p.Create {
		b.WriteRune('c')
	}
	if p.Write {
		b.WriteRune('w')
	}
	if p.Delete {
		b.WriteRune('d')
	}
	if p.List {
		b.WriteRune('l')
	}
	return b.String()
}

// Parse initializes the ShareSASPermissions' fields from a string.
func (p *ShareSASPermissions) Parse(s string) error {
	*p = ShareSASPermissions{} // Clear the flags
	for _, r := range s {
		switch r {
		case 'r':
			p.Read = true
		case 'c':
			p.Create = true
		case 'w':
			p.Write = true
		case 'd':
			p.Delete = true
		case 'l':
			p.List = true
		default:
			return fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return nil
}
