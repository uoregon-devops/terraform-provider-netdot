# How to Add New Types to Terraform Provider for Netdot
1. Using an example in internal/netdot/models/ create a new struct for the new type.
2. Using an example in internel/netdot/ create a new query builder for the new type.
3. Using an example named provider/<TYPE>.go create a schema for the new type.
4. Using an example named provider/<TYPE>_data_source.go create a data source for the new type.
5. Using an example named provider/<TYPE>_resource.go create a resource for the new type.
6. In provider/provider.go add the new resource and data source types to the Resources and DataSources returned by the netdot provider.
