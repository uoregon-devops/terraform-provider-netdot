# Here is how you would make an A record for foo.uoregon.edu a new server on 184.171.0.222
# A records in Netdot require that the target IP address exists in the Netdot database.

# You can reserve the IP address using the netdot_ip resource.
resource "netdot_ipblock" "foo" {
  address = "184.171.0.222"
}

# Then create a resource record for foo.uoregon.edu
resource "netdot_rr" "foo" {
  name = "foo"
  zone = "uoregon.edu"
}

# Finally attach an A record to the resource record pointing to the desired IP address
resource "netdot_arecord" "foo" {
  rr_id      = netdot_rr.foo.id
  ipblock_id = netdot_ipblock.foo.id
}

# If you wanted to implement DNS based load balancing, you could create multiple A records for the same resource record
# and point them to different IP addresses.
