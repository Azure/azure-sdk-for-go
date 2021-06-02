# Generator

This is a command line tool for generating new releases for `github.com/Azure/azure-sdk-for-go`.

## Commands

This CLI tool provides 2 commands now: `automation`, `issue`.

### The `issue` command and the configuration file

The `issue` command fetches the release request issues from `github.com/Azure/sdk-release-request/issues` and parses them into the configuration that other commands consume. The configuration will output to stdout.

The configuration is a JSON string, which has the following pattern:
```json
{
  "track1Requests": {
    "specification/network/resource-manager/readme.md": {
      "package-2020-12-01": [
        {
          "targetDate": "2021-02-11T00:00:00Z",
          "requestLink": "https://github.com/Azure/sdk-release-request/issues/1212"
        }
      ]
    }
  },
  "track2Requests": {},
  "refresh": {}
}
```
The keys of this JSON is the relative path of the `readme.md` file in `azure-rest-api-specs`.

To authenticate this command, you need to either
1. Populate a personal access token by assigning the `-t` flag.
1. Populate the username, password (and OTP if necessary) by assigning the `-u`, `-p` and `--otp` flags.

**Important notice:**
1. A release request by design can only have one RP in them, therefore if a release request is referencing a PR that contains changes of multiple RPs, the tool will just give an error.
1. A release request by design can only have one tag in them, therefore if a release request is requesting release on multiple tags, the tool will not give an error but output the plain value of the multiple tags without splitting them.
1. This command will try to output everything that it is able to parse, even some errors occur.

Example usage:
```shell
generator issue -t $YOUR_PERSONAL_ACCESS_TOKEN > sdk-release.json
```