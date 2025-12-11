## SDK Scale Testing
This directory is an independent Go module containing scale testing workloads for the SDK. The workloads are designed to test the performance and scalability of the SDK under various conditions.

### Project Structure
This is a standalone Go module with its own `go.mod` file. It depends on the parent `azcosmos` module via a local replace directive.

### Setup Scale Testing
1. Create a VM in Azure with the following configuration:
   - 8 vCPUs
   - 32 GB RAM
   - Ubuntu
   - Accelerated networking
1. Give the VM necessary [permissions](https://learn.microsoft.com/azure/cosmos-db/nosql/how-to-grant-data-plane-access?tabs=built-in-definition%2Ccsharp&pivots=azure-interface-cli) to access the Cosmos DB account if using AAD (Optional). 
1. Fork and clone this repository
1. Go to the workloads folder
   - `cd azure-sdk-for-go/sdk/data/azcosmos/workloads`
1. Checkout the branch with the changes to test. 
1. Run `./setup_env.sh`
1. Fill out relevant configs using environment variables:
   - `COSMOS_URI` - required
   - `COSMOS_KEY` - required (optional if using AAD)
   - `COSMOS_DATABASE` - optional (defaults to "scale_db")
   - `COSMOS_CONTAINER` - optional (defaults to "scale_cont")
   - `PARTITION_KEY` - optional (defaults to "pk")
   - `NUMBER_OF_LOGICAL_PARTITIONS` - optional (defaults to 10000)
   - `THROUGHPUT` - optional (defaults to 100000)
   - `PREFERRED_LOCATIONS` - optional (comma-separated list)
1. Set `AZURE_SDK_GO_LOGGING` env variable to "all" for detailed logs
1. Build and run the scale workload
    - `go run ./main/main.go`
    
    Or build the binary:
    - `go build -o workload ./main/main.go`
    - `./workload`

### Monitor Run
- `ps -eaf | grep "go"` to see the running processes

