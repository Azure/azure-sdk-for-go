function ShouldGenerate-AutorestConfig([string]$autorestPath) {
    $metaPath = $autorestPath.Replace("autorest.md", "_meta.json")
    return !(Test-Path $autorestPath -PathType Leaf) -and (Test-Path $metaPath -PathType Leaf)
}

function Generate-AutorestConfig([string]$autorestPath) {
    $metaPath = $autorestPath.Replace("autorest.md", "_meta.json")
    $meta = Get-Content $metaPath | ConvertFrom-Json
    $readme = "https://github.com/Azure/azure-rest-api-specs/blob/" + $meta.commit + $meta.readme.Split("/azure-rest-api-specs")[1];

    $contents = @'
### AutoRest Configuration
> see https://aka.ms/autorest
``` yaml
tag: {0}
require:
- {1}
```
'@;
    $contents -f $meta.tag, $readme | Out-File -FilePath $autorestPath
}