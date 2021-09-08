// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccountSASPermissions(t *testing.T) {
	a := AccountSASPermissions{
		Read:    true,
		Write:   true,
		Delete:  true,
		List:    true,
		Add:     true,
		Create:  true,
		Update:  true,
		Process: true,
	}
	require.Equal(t, a.String(), "rwdlacup")

	err := a.Parse("rwdl")
	require.NoError(t, err)
	require.True(t, a.Read)
	require.True(t, a.Write)
	require.True(t, a.Delete)
	require.True(t, a.List)
	require.False(t, a.Add)
	require.False(t, a.Create)
	require.False(t, a.Update)
	require.False(t, a.Process)
	err = a.Parse("z")
	require.Error(t, err)
}

func TestAccountSASResourceTypes(t *testing.T) {
	a := AccountSASResourceTypes{
		Service:   true,
		Container: true,
		Object:    true,
	}
	require.Equal(t, a.String(), "sco")

	err := a.Parse("o")
	require.NoError(t, err)
	require.False(t, a.Service)
	require.False(t, a.Container)
	require.True(t, a.Object)
	err = a.Parse("z")
	require.Error(t, err)
}

func TestSASPermissions(t *testing.T) {
	s := SASPermissions{
		Read:   true,
		Add:    true,
		Update: true,
		Delete: true,
	}
	require.Equal(t, s.String(), "raud")

	err := s.Parse("a")
	require.NoError(t, err)
	require.True(t, s.Add)
	require.False(t, s.Read)
	require.False(t, s.Update)
	require.False(t, s.Delete)
}
