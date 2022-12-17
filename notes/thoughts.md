# thoughts

- The CLI should make sure scopes don't have the same name as each other (or
  special characters, etc).
  Then we can use something like `{{scope.org}}` in a template file in a root.
- People will probably want some kind of output wiring (not unlike Terragrunt's
  `dependency` blocks). This will likely have to be a new feature after v1
  though.

## The Core

Here's the core of terraboots:

- you don't need anything other than terraform to run it. This means its config
  files are dead simple.
- it will tell you when it is improperly configured. This is important because
  configuration is confusing! And guardrails are just as important as docs.
  This means the configs will all be in HCL, and the tool will validate them at
  runtime.
- it can solve the `affected` problem quickly. This means it has to be
  performant. And maybe it has to be extensible, because people have different
  ideas of what it means to be "affected".
- it can solve root dependencies quickly to generate a parallelizable task
  manifest. Again, it has to be performant.

The primary goal of `terraboots root list` is to make sure we can build a matrix
for GitHub Actions or any other CI runner to know exactly what to run. It might
want to take root dependency into account too.

## up next

1. I need to parse a single root, and build it.
2. I need a schema for the scope data.
3. scope validation is a future feature.
4. root graph is an interesting idea, to show all dependencies

## affected

what makes it so that a root is worth running?

1. the scope data changed
2. the root itself changed
3. the modules the root relied on changed? not sure how to do this quickly
   without either parsing the entire AST or assuming semver sources

## root dependencies

Speaking of, I do need it to have some concept of dynamic dependency. For
example: "this root depends on another root". We'd have to do the math of
excluding some scopes _if necessary_. Root A depending on root B means that we'd
take all the B scopes that match A's, and make sure the values are the same.
That "if necessary" is important though: maybe root A depends on root B but with
some important scope differences (e.g. root named "account-networking" with
`scope:domain=Product` depends on root named "transit-gateway" with scope
`scope:domain=Networking`).

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
