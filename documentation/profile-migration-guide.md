# Guide for migrating from `profiles/**`

This document is intended for users that are familiar with the previous versions of the Azure Profile SDK packages for Go (`github.com/Azure/azure-sdk-for-go/profiles/**`), and target to migrate applications to the new releases of the modules.

## Table of contents

- [Prerequisites](#prerequisites)
- [Package mapping](#package-mapping)
- [API changes](#API-changes)

## Prerequisites

- [Supported](https://aka.ms/azsdk/go/supported-versions) version of Go

## Package mapping

There are several released packages of Azure Profile SDK for Go (`profiles/**`): `2017-03-09`, `2018-03-01`, `2019-03-01`, `2020-09-01`, `latest`, `preview`.

For those who are using `latest` or `preview`, you need to migrate to our latest modules under `github.com/Azure/azure-sdk-for-go/sdk`. The latest SDK provides one module for each resource provider. For each resource provider you used in the previous version, you need to migrate to a cooresponding module. You can search in [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go) to find those modules with stable and preview version.

For example, if you are using `github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice`, then you need to migrate to `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice@v2.4.0`; if you are using `github.com/Azure/azure-sdk-for-go/profiles/preview/containerservice/mgmt/containerservice`, then you need to migrate to `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice@2.5.0-beta.1`.

For those who are using `2017-03-09`, `2018-03-01` or `2019-03-01`, there is no direct mapping in the new releases. You have to upgrade to use the latest profile `2020-09-01` with our latest profile module `github.com/Azure/azure-sdk-for-go/profile/p20200901`.

For those who are using `2020-09-01`, you need to migrate to our latest profile module `github.com/Azure/azure-sdk-for-go/profile/p20200901`.

## API changes

For `management` modules, you could refer [this](https://aka.ms/azsdk/go/mgmt/migration) for migrations. For `client` modules, you could refer to "Client modules" section in the [README.md](https://github.com/Azure/azure-sdk-for-go/blob/main/README.md#client-modules) file.

## Need help?

If you have encountered an issue during migration, please file an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues).
