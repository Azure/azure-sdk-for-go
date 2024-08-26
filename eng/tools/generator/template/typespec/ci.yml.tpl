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
    - sdk/resourcemanager/{{.rpName}}/{{.packageName}}/

pr:
  branches:
    include:
      - main
      - feature/*
      - hotfix/*
      - release/*
  paths:
    include:
    - sdk/resourcemanager/{{.rpName}}/{{.packageName}}/

extends:
  template: /eng/pipelines/templates/jobs/archetype-sdk-client.yml
  parameters:
    IncludeRelease: true
    ServiceDirectory: 'resourcemanager/{{.rpName}}/{{.packageName}}'
