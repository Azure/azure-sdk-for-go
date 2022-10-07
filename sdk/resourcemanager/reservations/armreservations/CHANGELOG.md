# Release History

## 1.1.0 (2022-09-16)
### Features Added

- New const `ErrorResponseCodeRefundLimitExceeded`
- New const `ErrorResponseCodeSelfServiceRefundNotSupported`
- New function `*ReservationClient.Archive(context.Context, string, string, *ReservationClientArchiveOptions) (ReservationClientArchiveResponse, error)`
- New function `NewCalculateRefundClient(azcore.TokenCredential, *arm.ClientOptions) (*CalculateRefundClient, error)`
- New function `*CalculateRefundClient.Post(context.Context, string, CalculateRefundRequest, *CalculateRefundClientPostOptions) (CalculateRefundClientPostResponse, error)`
- New function `*ReservationClient.Unarchive(context.Context, string, string, *ReservationClientUnarchiveOptions) (ReservationClientUnarchiveResponse, error)`
- New function `NewReturnClient(azcore.TokenCredential, *arm.ClientOptions) (*ReturnClient, error)`
- New function `*ReturnClient.Post(context.Context, string, RefundRequest, *ReturnClientPostOptions) (ReturnClientPostResponse, error)`
- New struct `CalculateRefundClient`
- New struct `CalculateRefundClientPostOptions`
- New struct `CalculateRefundClientPostResponse`
- New struct `CalculateRefundRequest`
- New struct `CalculateRefundRequestProperties`
- New struct `CalculateRefundResponse`
- New struct `RefundBillingInformation`
- New struct `RefundPolicyError`
- New struct `RefundPolicyResult`
- New struct `RefundPolicyResultProperty`
- New struct `RefundRequest`
- New struct `RefundRequestProperties`
- New struct `RefundResponse`
- New struct `RefundResponseProperties`
- New struct `ReservationClientArchiveOptions`
- New struct `ReservationClientArchiveResponse`
- New struct `ReservationClientUnarchiveOptions`
- New struct `ReservationClientUnarchiveResponse`
- New struct `ReturnClient`
- New struct `ReturnClientPostOptions`
- New struct `ReturnClientPostResponse`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).