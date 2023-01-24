# root dependencies

I do need it to have some concept of dynamic dependency. For example: "this root
depends on another root". We'd have to do the math of excluding some scopes
_if necessary_. Root A depending on root B means that we'd take all the B scopes
that match A's, and make sure the values are the same. That "if necessary" is
important though: maybe root A depends on root B but with some important scope
differences (e.g. root named "account-networking" with `scope:domain=Product`
depends on root named "transit-gateway" with scope `scope:domain=Networking`).

That probably looks like a `root` block of

```hcl
root "account-networking" {
  scopes = ["org", "platform", "domain", "environment", "region"]

  dependency {
    root = "transit-gateway"
    scopes = {
      # unset scopes here will retain the values of the dependant
      domain = "Networking"
    }
  }
}
```

I'm not sure how much to build the parameter pattern into it. I can get away
with tagging dependencies. Where the example above is "root named A depends on
root named B", this example is more like "root tagged consumer:foo depends on
root tagged producer:foo". Building this into terraboots looks like a way to set
in the main repo config that "all roots tagged `consumer:` will depend on their
corresponding `producer:` roots".

## implementation

1. list of roots
2. sort the list by dependency
3. make a map from root name to batch number
4. make a map from batch number to root name
5. for each root in the list
6. put the root in the maps based on its dependencies


7. for each root in the map[num]string 
8. generate the set of buildContexts [rename buildContext to rootContext?]
9. print it out
