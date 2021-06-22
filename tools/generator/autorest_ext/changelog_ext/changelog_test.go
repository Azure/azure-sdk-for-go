package changelog_ext_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest_ext/changelog_ext"
)

func TestGetLinesBetween(t *testing.T) {
	testdata := []struct {
		input    string
		title    []string
		previous []string
		rest     []string
	}{
		{
			input:    "# CHANGELOG\n\n## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`",
			title:    strings.Split("# CHANGELOG\n", "\n"),
			previous: strings.Split("## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`", "\n"),
			rest:     nil,
		},
		{
			input:    "# CHANGELOG\n\n## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`\n\n## `v50.0.0`\n\nNOTE: Due to the changes requested in [this issue](https://github.com/Azure/azure-sdk-for-go/issues/14010), we changed the properties and functions of all future types, which does not affect their functionality and usage, but leads to a very long list of breaking changes. This change requires the latest version of `github.com/Azure/go-autorest/autorest v0.11.15` to work properly.\n\n### Renamed Packages\n\n- `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute` renamed to `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute` to align other naming pattern of the profile packages\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/healthbot/mgmt/2020-12-08/healthbot`\n- `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources`\n\n### Removed Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/preview/avs/mgmt/2019-08-09-preview/avs`\n\n### Breaking Changes\n\n| Package Path | Changelog |\n| :--- | :---: |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2016-05-16/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2016-05-16/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-07-14/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-07-14/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-08-01/analysisservices/CHANGELOG.md) |\n\n## `v49.2.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-12-01/containerservice`\n- `github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration`\n- `github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2020-09-01/netapp`\n- `github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform`\n- `github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs`\n",
			title:    strings.Split("# CHANGELOG\n", "\n"),
			previous: strings.Split("## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`\n", "\n"),
			rest:     strings.Split("## `v50.0.0`\n\nNOTE: Due to the changes requested in [this issue](https://github.com/Azure/azure-sdk-for-go/issues/14010), we changed the properties and functions of all future types, which does not affect their functionality and usage, but leads to a very long list of breaking changes. This change requires the latest version of `github.com/Azure/go-autorest/autorest v0.11.15` to work properly.\n\n### Renamed Packages\n\n- `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute` renamed to `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute` to align other naming pattern of the profile packages\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/healthbot/mgmt/2020-12-08/healthbot`\n- `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources`\n\n### Removed Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/preview/avs/mgmt/2019-08-09-preview/avs`\n\n### Breaking Changes\n\n| Package Path | Changelog |\n| :--- | :---: |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2016-05-16/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2016-05-16/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-07-14/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-07-14/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-08-01/analysisservices/CHANGELOG.md) |\n\n## `v49.2.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-12-01/containerservice`\n- `github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration`\n- `github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2020-09-01/netapp`\n- `github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform`\n- `github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs`\n", "\n"),
		},
	}

	for _, c := range testdata {
		lines := strings.Split(c.input, "\n")
		versionRange := changelog_ext.FindVersionTitles(lines, 2)
		title, previous, rest := changelog_ext.GetLinesBetween(lines, versionRange)
		if !reflect.DeepEqual(title, c.title) {
			t.Fatalf("expect %+v, but got %+v", c.title, title)
		}
		if !reflect.DeepEqual(previous, c.previous) {
			t.Fatalf("expect %+v, but got %+v", c.previous, previous)
		}
		if !reflect.DeepEqual(rest, c.rest) {
			t.Fatalf("expect %+v, but got %+v", c.rest, rest)
		}
	}
}

func TestFindVersionTitles(t *testing.T) {
	testdata := []struct {
		input    string
		expected []changelog_ext.VersionTitleLine
	}{
		{
			input: "# CHANGELOG\n\n## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`",
			expected: []changelog_ext.VersionTitleLine{
				{
					LineNumber: 2,
					Version:    "v50.1.0",
				},
			},
		},
		{
			input: "# CHANGELOG\n\n## `v50.1.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/confluent/mgmt/2020-03-01/confluent`\n- `github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-08-01/network`\n- `github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization`\n\n## `v50.0.0`\n\nNOTE: Due to the changes requested in [this issue](https://github.com/Azure/azure-sdk-for-go/issues/14010), we changed the properties and functions of all future types, which does not affect their functionality and usage, but leads to a very long list of breaking changes. This change requires the latest version of `github.com/Azure/go-autorest/autorest v0.11.15` to work properly.\n\n### Renamed Packages\n\n- `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute` renamed to `github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute` to align other naming pattern of the profile packages\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/healthbot/mgmt/2020-12-08/healthbot`\n- `github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources`\n\n### Removed Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/preview/avs/mgmt/2019-08-09-preview/avs`\n\n### Breaking Changes\n\n| Package Path | Changelog |\n| :--- | :---: |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2016-05-16/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2016-05-16/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-07-14/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-07-14/analysisservices/CHANGELOG.md) |\n| `github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices` | [details](https://github.com/Azure/azure-sdk-for-go/blob/v50.0.0/services/analysisservices/mgmt/2017-08-01/analysisservices/CHANGELOG.md) |\n\n## `v49.2.0`\n\n### New Packages\n\n- `github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-12-01/containerservice`\n- `github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration`\n- `github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2020-09-01/netapp`\n- `github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2020-11-01-preview/appplatform`\n- `github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs`\n",
			expected: []changelog_ext.VersionTitleLine{
				{
					LineNumber: 2,
					Version:    "v50.1.0",
				},
				{
					LineNumber: 10,
					Version:    "v50.0.0",
				},
				{
					LineNumber: 35,
					Version:    "v49.2.0",
				},
			},
		},
	}

	for _, c := range testdata {
		versions := changelog_ext.FindVersionTitles(strings.Split(c.input, "\n"), -1)
		if !reflect.DeepEqual(versions, c.expected) {
			t.Fatalf("expect %+v, but got %+v", c.expected, versions)
		}
	}
}
