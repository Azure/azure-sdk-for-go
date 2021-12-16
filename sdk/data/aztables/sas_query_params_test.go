// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatTimes(t *testing.T) {
	start := time.Date(2021, time.September, 8, 13, 6, 0, 0, time.UTC)
	expiry := start.AddDate(1, 0, 0)
	startString, expiryString := FormatTimesForSASSigning(start, expiry)
	require.Equal(t, "2021-09-08T13:06:00Z", startString)
	require.Equal(t, "2022-09-08T13:06:00Z", expiryString)
}

func TestFormatIPRange(t *testing.T) {
	i := IPRange{
		Start: net.IPv4(224, 0, 0, 250),
	}
	require.Equal(t, i.String(), "224.0.0.250")

	i2 := IPRange{
		End: net.IPv4(192, 0, 0, 168),
	}
	require.Equal(t, i2.String(), "")

	i3 := IPRange{
		Start: net.IPv4(192, 0, 0, 168),
		End:   net.IPv4(224, 0, 0, 250),
	}
	require.Equal(t, i3.String(), "192.0.0.168-224.0.0.250")
}

func TestSASQueryParameters(t *testing.T) {
	start := time.Date(2021, time.September, 8, 13, 45, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	i := IPRange{
		Start: net.IPv4(192, 0, 0, 168),
		End:   net.IPv4(224, 0, 0, 250),
	}
	s := SASQueryParameters{
		version:       "2020-08-04",
		services:      "t",
		resourceTypes: "sco",
		protocol:      SASProtocolHTTPS,
		startTime:     start,
		expiryTime:    end,
		ipRange:       i,
		identifier:    "i",
		resource:      "t",
		permissions:   "raud",
		signature:     "fakesignature",
		signedVersion: "signedVersion",
		tableName:     "tableName",
		startPk:       "startPK",
		startRk:       "startRK",
		endPk:         "endPK",
		endRk:         "endRk",
	}

	require.Equal(t, s.SignedVersion(), "signedVersion")
	require.Equal(t, s.Version(), "2020-08-04")
	require.Equal(t, s.Services(), "t")
	require.Equal(t, s.ResourceTypes(), "sco")
	require.Equal(t, s.Protocol(), SASProtocolHTTPS)
	require.Equal(t, s.StartTime(), start)
	require.Equal(t, s.ExpiryTime(), end)
	require.Equal(t, s.IPRange(), i)
	require.Equal(t, s.Identifier(), "i")
	require.Equal(t, s.Resource(), "t")
	require.Equal(t, s.Permissions(), "raud")
	require.Equal(t, s.Signature(), "fakesignature")
	require.Equal(t, s.StartPartitionKey(), "startPK")
	require.Equal(t, s.StartRowKey(), "startRK")
	require.Equal(t, s.EndPartitionKey(), "endPK")
	require.Equal(t, s.EndRowKey(), "endRk")

	encoded := s.Encode()
	require.Equal(t, "epk=endPK&erk=endRk&se=2022-09-08T13%3A45%3A00Z&si=i&sig=fakesignature&sip=192.0.0.168-224.0.0.250&sp=raud&spk=startPK&spr=https&sr=t&srk=startRK&srt=sco&ss=t&st=2021-09-08T13%3A45%3A00Z&sv=2020-08-04&tn=tableName", encoded)

	v := url.Values{}
	v.Add("qp1", "value1")
	v.Add("qp2", "value2")

	result := s.addToValues(v)
	require.Equal(t, result.Get("epk"), s.EndPartitionKey())
	require.Equal(t, result.Get("erk"), s.EndRowKey())
	require.Equal(t, result.Get("spk"), s.StartPartitionKey())
	require.Equal(t, result.Get("srk"), s.StartRowKey())
	require.Equal(t, result.Get("qp1"), "value1")
	require.Equal(t, result.Get("qp2"), "value2")
	require.Equal(t, result.Get("se"), "2022-09-08T13:45:00Z")
	require.Equal(t, result.Get("si"), s.Identifier())
	require.Equal(t, result.Get("sig"), s.Signature())
	require.Equal(t, result.Get("sip"), "192.0.0.168-224.0.0.250")
	require.Equal(t, result.Get("sp"), s.Permissions())
	require.Equal(t, result.Get("srt"), s.ResourceTypes())
	require.Equal(t, result.Get("ss"), s.Services())
	require.Equal(t, result.Get("st"), "2021-09-08T13:45:00Z")
	require.Equal(t, result.Get("sv"), "2020-08-04")
	require.Equal(t, result.Get("tn"), "tableName")
}
