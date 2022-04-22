# Unreleased

## Additive Changes

### New Funcs

1. AvailabilityZonePeers.MarshalJSON() ([]byte, error)
1. CheckZonePeersResult.MarshalJSON() ([]byte, error)
1. Client.CheckZonePeers(context.Context, string, CheckZonePeersRequest) (CheckZonePeersResult, error)
1. Client.CheckZonePeersPreparer(context.Context, string, CheckZonePeersRequest) (*http.Request, error)
1. Client.CheckZonePeersResponder(*http.Response) (CheckZonePeersResult, error)
1. Client.CheckZonePeersSender(*http.Request) (*http.Response, error)
1. Peers.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. AvailabilityZonePeers
1. CheckZonePeersRequest
1. CheckZonePeersResult
1. Peers
