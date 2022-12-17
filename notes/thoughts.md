# thoughts

The CLI should make sure scopes don't have the same name as each other (or
special characters, etc).
Then we can use something like `{{scope.org}}` in a template file in a root.

How do I want the CLI to operate? What's its command structure?

```sh
# housekeeping
terraboots # prints short help
terraboots [-h|--help] # prints long help
terraboots version # prints version of terraboots

# running existing roots
terraboots plan <root> # using primary terraform commands like init, plan, apply and output will do what they say on the tin
terraboots plan <root> -- --var-file=foo.tfvars # plans the root with some extra tf cli arguments
terraboots terraform|tf <root> -- state list # runs `terraform state list` in the given root
terraboots plan <root> -v|--verbose # makes sure the user knows what's happening in the templating

# inspecting your platform
terraboots scope list # list all your scopes
terraboots root list # list every single permutation of every root
terraboots root list --affected # lists every single permutation of every root *which has been affected* (for some definition of "affected")

# building new stuff
terraboots new monorepo acme-infrastructure # builds a whole new monorepo named "acme-infrastructure"
terraboots new scope team # adds a new scope into the mix
terraboots new root my-awesome-stack # builds a new root module named 'my-awesome-stack'
```

The primary goal of `terraboots root list` is to make sure we can build a matrix
for GitHub Actions or any other CI runner to know exactly what to run. It might
want to take root dependency into account too.

Speaking of, I do need it to have some concept of dynamic dependency. For
example: "this root depends on another root". We'd have to do the math of
excluding some scopes if necessary. Root A depending on root B means that we'd
take all the B scopes that match A's, and make sure the values are the same.
That "if necessary" is important though: maybe root A depends on root B but with
some key scope differences (e.g. "account-networking" for scope:domain=Product
depends on "transit-gateway" for scope:domain=Networking)

I'm not sure how much to build the parameter pattern into it. I can get away
with tagging dependencies. Where the example above is "root named A depends on
root named B", this example is more like "root tagged consumer:foo depends on
root tagged producer:foo". Building this into terraboots looks like a way to set
in the main repo config that "all roots tagged `consumer:` will depend on their
corresponding `producer:` roots".

## The Core

Here's the core of terraboots:

- you don't need anything other than terraform to run it. This means its config
  files are dead simple.
- it will tell you when it is improperly configured. This is important because
  this configuration is confusing! And guardrails are just as important as docs.
  This means the configs will all be in HCL, and the tool will validate them at
  runtime.
- it can solve the `affected` problem quickly. This means it has to be
  performant.
- it can solve the dependencies quickly to generate a parallelizable task
  manifest. Again, it has to be performant.
