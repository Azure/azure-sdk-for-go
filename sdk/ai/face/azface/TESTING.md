# Azure Face SDK Live Testing

This document describes how to run live tests for the Azure Face SDK.

## Prerequisites

1. An Azure subscription
2. An Azure Cognitive Services Face resource deployed in your subscription

## Environment Variables

The following environment variables are required for live testing:

- `FACE_ENDPOINT` - The endpoint URL of your Face service (e.g., `https://my-face-service.cognitiveservices.azure.com/`)
- `FACE_SUBSCRIPTION_KEY` - The subscription key for your Face service

For Azure AD authentication (optional, will fall back to subscription key):
- Standard Azure authentication environment variables (e.g., `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`)

## Running Tests

### Live Mode
```bash
export AZURE_RECORD_MODE=live
export FACE_ENDPOINT="https://your-face-service.cognitiveservices.azure.com/"
export FACE_SUBSCRIPTION_KEY="your-subscription-key"
go test -v
```

### Record Mode
```bash
export AZURE_RECORD_MODE=record
export FACE_ENDPOINT="https://your-face-service.cognitiveservices.azure.com/"
export FACE_SUBSCRIPTION_KEY="your-subscription-key"
go test -v
```

### Playback Mode (default)
```bash
go test -v
```

## Test Coverage

The live tests cover:

1. **Face Detection** - `TestClient_DetectFromURL` and `TestClient_Detect`
   - Detects faces in images from URL and image data
   - Returns face attributes like age, glasses, etc.

2. **Face Recognition** - `TestClient_FindSimilar`
   - Finds similar faces from a set of candidates

3. **Face Grouping** - `TestClient_Group`
   - Groups similar faces together

4. **Face Verification** - `TestClient_VerifyFaceToFace`
   - Verifies if two faces belong to the same person

5. **Administration** - `TestAdministrationLargeFaceList_*`
   - Large Face List creation, face addition, and management

## Notes

- Tests use a sample image from Azure's public samples for face detection
- Administration tests create temporary resources that are cleaned up automatically
- The test suite includes proper error handling and cleanup procedures