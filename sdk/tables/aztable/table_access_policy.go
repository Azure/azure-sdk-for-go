package aztable

import "time"

type TableAccessPolicy struct {
	Start      time.Time
	Expiry     time.Time
	Permission string
}

func castAccessPolicyToSignedIdentifier(accessPolicies []*AccessPolicy) []*SignedIdentifier {
	ret := make([]*SignedIdentifier, 0)

	for _, accessPolicy := range accessPolicies {
		ret = append(ret, &SignedIdentifier{AccessPolicy: accessPolicy})
	}
	return ret
}
