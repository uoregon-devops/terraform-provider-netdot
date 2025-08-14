# Terraform Netdot Provider

## TODO
- [ ] implement import interface on resources

### Document
- [ ] adding a new table type from netdot to the provider

### Ipblock
- [x] don't reserve first n when changing between not-subnet and subnet.
- [ ] break early if creating an IP in a block that doesn't actually contain the IP
- [ ] allow specifying parent by cidr/address
- [ ] prevent ipblock deletion if children exist
- [ ] see if it makes sense to leverage IPBLOCK->add_range() to generate reverse DNS for blocks
- [ ] explore managing the gateway in ipblock subnet state, make it a flag
