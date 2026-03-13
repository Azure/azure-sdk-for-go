# armneonpostgres — Service Retirement

## Service Retirement

The **Neon Postgres** service on Azure is being permanently retired on **13/03/2026**. There is no replacement service or SDK module. The `armneonpostgres` SDK module (`github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/neonpostgres/armneonpostgres`) will no longer be maintained after that date.

For more information, please refer to the [Azure SDK deprecation policy](https://aka.ms/azsdk/support-policies).

## What You Should Do

1. **Stop using the `armneonpostgres` module.** The service will no longer be available after the retirement date.
2. **Remove the dependency** from your Go modules:

```sh
go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/neonpostgres/armneonpostgres@none
```

Then run:

```sh
go mod tidy
```

## Questions and Support

If you have questions about this retirement, please refer to:
- [Azure SDK Support Policies](https://aka.ms/azsdk/support-policies)
- File an issue in the [azure-sdk-for-go repository](https://github.com/Azure/azure-sdk-for-go/issues)
