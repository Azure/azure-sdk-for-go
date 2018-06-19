package servicebus

import (
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/mitchellh/mapstructure"
)

func (suite *serviceBusSuite) TestMapStructureEncode() {
	sp := new(SystemProperties)
	m, err := encodeStructureToMap(sp)
	if suite.NoError(err) {
		suite.Len(m, 0)
	}

	now := time.Now()
	pID := int16(1)
	sp.LockedUntil = &now
	m, err = encodeStructureToMap(sp)
	if suite.NoError(err) {
		suite.Equal(now, m["x-opt-locked-until"])
		suite.Len(m, 1)
	}

	sp.PartitionKey = to.StringPtr("foo")
	sp.PartitionID = &pID
	sp.SequenceNumber = to.Int64Ptr(42)
	sp.EnqueuedTime = &now
	sp.EnqueuedSequenceNumber = to.Int64Ptr(43)
	sp.DeadLetterSource = to.StringPtr("bar")
	sp.ScheduledEnqueuedTime = &now
	sp.ViaPartitionKey = to.StringPtr("via")

	m, err = encodeStructureToMap(sp)
	if suite.NoError(err) {
		var sp2 SystemProperties
		err = mapstructure.Decode(&m, &sp2)
		if suite.NoError(err) {
			suite.Equal(sp, &sp2)
		}
	}
}
