## SDK Scale Testing
This directory contains the scale testing workloads for the SDK. The workloads are designed to test the performance 
and scalability of the SDK under various conditions. 

### Setup Scale Testing
1. Create a VM in Azure with the following configuration:
   - 8 vCPUs
   - 32 GB RAM
   - Ubuntu
   - Accelerated networking
1. Give the VM necessary [permissions](https://learn.microsoft.com/azure/cosmos-db/nosql/how-to-grant-data-plane-access?tabs=built-in-definition%2Ccsharp&pivots=azure-interface-cli) to access the Cosmos DB account if using AAD (Optional). 
1. Fork and clone this repository
1. Go to azcosmos folder
   - `cd azure-sdk-for-go/sdk/data/azcosmos/workloads`
1. Checkout the branch with the changes to test. 
1. Run `./setup_env.sh`
1. Fill out relevant configs in `workload_configs.go`: key, host, etc using env variables
   - `COSMOS_URI` - required
   - `COSMOS_KEY` - required
   - `COSMOS_DATABASE` 
   - `COSMOS_CONTAINER` 
   - `PARTITION_KEY`
   - `NUMBER_OF_LOGICAL_PARTITIONS`
   - `THROUGHPUT`
   - `PREFERRED_LOCATIONS`
1. Run the scale workload
    - `go run ./main/main.go`

### Monitor Run
- `ps -eaf | grep "go"` to see the running processes

`
