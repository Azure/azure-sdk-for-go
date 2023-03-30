# Guide for migrating from `profiles/**`

This document is intended for users that are familiar with the previous version of the Azure SDK For Go profiles package (`github.com/Azure/azure-sdk-for-go/profiles/**`) and wish to migrate their application to the next version of Azure SDK for Go.

## Table of contents

- [Prerequisites](#prerequisites)
- [Package mapping](#package-mapping)
- [API changes](#API-changes)

## Prerequisites

- Go 1.18

## Package mapping

There are several sub-packages under pervious version of the Azure SDK For Go profiles package (`profiles/**`): `2017-03-09`, `2018-03-01`, `2019-03-01`, `2020-09-01`, `latest`, `preview`.

For those who are using `latest` or `preview`, you need to migrate to our latest modules under `github.com/Azure/azure-sdk-for-go/sdk`. The latest SDK only provides one API version in one module for each RP. If you are using several RP's sub-packages, you need to migrate to several modules. You can search in [pkg.go.dev](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go) to find all the stable and preview version of all RPs.

For example, if you are using `github.com/Azure/azure-sdk-for-go/profiles/latest/containerservice/mgmt/containerservice`, then you need to migrate to `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice@v2.4.0`; if you are using `github.com/Azure/azure-sdk-for-go/profiles/preview/containerservice/mgmt/containerservice`, then you need to migrate to `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice@2.5.0-beta.1`.

For those who are using `2017-03-09`, `2018-03-01` or `2019-03-01`, unfortunately, we do not provide these old profile version in our latest Azure SDK for Go. Please upgrade to the latest profile `2020-09-01` first.

For those who are using `2020-09-01`, you need to migrate to our latest profile module `github.com/Azure/azure-sdk-for-go/profile/p20200901`.

## API changes

For `management` modules, you could refer [this](https://aka.ms/azsdk/go/mgmt/migration) for migrations. For `client` modules, you could refer README.md file of cooresponding module.

## Need help?

If you have encountered an issue during migration, please file an issue via [Github Issues](https://github.com/Azure/azure-sdk-for-go/issues) and make sure you add the "Preview" label to the issue
