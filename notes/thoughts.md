# thoughts

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

## affected

what makes it so that a root is worth running?

1. the scope data changed
2. the root itself changed
3. the modules the root relied on changed? not sure how to do this quickly
   without either parsing the entire AST or assuming semver sources

## scope data

Each root has scopes it follows. It will hydrate itself based on the values of
those scopes, but it can't take a naive approach. Sometimes a scope is different
depending on its ancestry. For instance: the "Gold" platform has different
domains than the "Silver" platform.

So the solution seems to be to spell out exactly how the scope tree goes. This
might be too verbose though! Oh well, it's still better than spelling out all
the roots separately.

So, do I spell the scope tree out with YAML or HCL? After a short experiment
(see `scopeTree1.hcl` and `scopeTree2.yml`), it seems the HCL isn't that much
more verbose than YAML, and I think it could pose some benefits for the parser.

This also raises a new topic for internal scope validation: the scope names
cannot be "attributes" or "roots".

If it's hcl there's some fun stuff we can do like set locals to reuse.
