# Here is how you would make a CNAME for domain validation where foo.uoregon.edu resolves to bar.aws.com

# create a resource record for foo.uoregon.edu
resource "netdot_rr" "foo" {
  name = "foo"
  zone = "uoregon.edu"
}

# create a CNAME record attached to foo.uoregon.edu that points to bar.aws.com
resource "netdot_cname" "bar" {
  cname = "bar.aws.com"
  rr_id = netdot_rr.foo.id
}
