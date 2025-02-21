# NOTE: Please refer to https://aka.ms/azsdk/engsys/ci-yaml before editing this file.
trigger:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

pr:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - {{.moduleRelativePath}}/

extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    IncludeRelease: true
    ServiceDirectory: '{{.serviceDir}}'
