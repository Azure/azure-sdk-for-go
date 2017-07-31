package dns

import (
	 original "github.com/Azure/azure-sdk-for-go/service/dns/management/2016-04-01/dns"
)

type (
	 ZonesClient = original.ZonesClient
	 ManagementClient = original.ManagementClient
	 HTTPStatusCode = original.HTTPStatusCode
	 OperationStatus = original.OperationStatus
	 RecordType = original.RecordType
	 AaaaRecord = original.AaaaRecord
	 ARecord = original.ARecord
	 CloudError = original.CloudError
	 CloudErrorBody = original.CloudErrorBody
	 CnameRecord = original.CnameRecord
	 MxRecord = original.MxRecord
	 NsRecord = original.NsRecord
	 PtrRecord = original.PtrRecord
	 RecordSet = original.RecordSet
	 RecordSetListResult = original.RecordSetListResult
	 RecordSetProperties = original.RecordSetProperties
	 RecordSetUpdateParameters = original.RecordSetUpdateParameters
	 Resource = original.Resource
	 SoaRecord = original.SoaRecord
	 SrvRecord = original.SrvRecord
	 SubResource = original.SubResource
	 TxtRecord = original.TxtRecord
	 Zone = original.Zone
	 ZoneDeleteResult = original.ZoneDeleteResult
	 ZoneListResult = original.ZoneListResult
	 ZoneProperties = original.ZoneProperties
	 RecordSetsClient = original.RecordSetsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Accepted = original.Accepted
	 Ambiguous = original.Ambiguous
	 BadGateway = original.BadGateway
	 BadRequest = original.BadRequest
	 Conflict = original.Conflict
	 Continue = original.Continue
	 Created = original.Created
	 ExpectationFailed = original.ExpectationFailed
	 Forbidden = original.Forbidden
	 Found = original.Found
	 GatewayTimeout = original.GatewayTimeout
	 Gone = original.Gone
	 HTTPVersionNotSupported = original.HTTPVersionNotSupported
	 InternalServerError = original.InternalServerError
	 LengthRequired = original.LengthRequired
	 MethodNotAllowed = original.MethodNotAllowed
	 Moved = original.Moved
	 MovedPermanently = original.MovedPermanently
	 MultipleChoices = original.MultipleChoices
	 NoContent = original.NoContent
	 NonAuthoritativeInformation = original.NonAuthoritativeInformation
	 NotAcceptable = original.NotAcceptable
	 NotFound = original.NotFound
	 NotImplemented = original.NotImplemented
	 NotModified = original.NotModified
	 OK = original.OK
	 PartialContent = original.PartialContent
	 PaymentRequired = original.PaymentRequired
	 PreconditionFailed = original.PreconditionFailed
	 ProxyAuthenticationRequired = original.ProxyAuthenticationRequired
	 Redirect = original.Redirect
	 RedirectKeepVerb = original.RedirectKeepVerb
	 RedirectMethod = original.RedirectMethod
	 RequestedRangeNotSatisfiable = original.RequestedRangeNotSatisfiable
	 RequestEntityTooLarge = original.RequestEntityTooLarge
	 RequestTimeout = original.RequestTimeout
	 RequestURITooLong = original.RequestURITooLong
	 ResetContent = original.ResetContent
	 SeeOther = original.SeeOther
	 ServiceUnavailable = original.ServiceUnavailable
	 SwitchingProtocols = original.SwitchingProtocols
	 TemporaryRedirect = original.TemporaryRedirect
	 Unauthorized = original.Unauthorized
	 UnsupportedMediaType = original.UnsupportedMediaType
	 Unused = original.Unused
	 UpgradeRequired = original.UpgradeRequired
	 UseProxy = original.UseProxy
	 Failed = original.Failed
	 InProgress = original.InProgress
	 Succeeded = original.Succeeded
	 A = original.A
	 AAAA = original.AAAA
	 CNAME = original.CNAME
	 MX = original.MX
	 NS = original.NS
	 PTR = original.PTR
	 SOA = original.SOA
	 SRV = original.SRV
	 TXT = original.TXT
)
