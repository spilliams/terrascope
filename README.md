# terrascope

A build orchestrator for terraform monorepos.

This repository contains both the source code for the tool `terrascope`, as well
as a sample monorepo managed by that tool.

## The Big Idea

Say you manage an engineering platform comprising a dozen or so product domains.

You want to give each product domain separate environments for their development,
staging and production efforts, so you make an AWS account for each of these
environments.

Now imagine that as a security-conscious engineer, you want to establish certain
settings and configuration in each of these AWS accounts. You want each account
to have: a VPC, AWS Config resource-recording, AWS Backup, a New Relic log
forwarder, and a couple of faceless IAM Users.

You could do all that in one terraform configuration, but every time you `init`,
`plan`, or `apply` that module, you'll be waiting for half an hour while it
completes.

Operating this system could be faster, and maintaining it could be less error-
prone, if each of these use-cases in each AWS account had their own modules. The
problem you have then is 180 or so individual terraform modules.

This last is the problem that `terrascope` seeks to simplify.

### Scopes

The main concept to understand in order to use `terrascope` effectively is that
of "scopes". Here we use that term to mean a sort of "build context". The goal
is to create and maintain a tree of nested scopes that matches your
organization's engineering domain architecture.

The root of the tree will likely be a scope type with a single value,
representing your entire `organization`.

The next scope type could be `department` or `team`, but tracking those might be
unnecessarily granular. To continue our example from above, I'll use `domain`
as the second scope.

As we mentioned earlier, each domain gets a set of three environments, so under
the `domain` scope we have `environment` scopes.

Here is how our scope tree looks at this point:

```text
organization
└── domain
    └── environment
```

With `terrascope` handling the scope values, it will be much easier to maintain
any number of terraform modules per environment.

## Installation & Usage

Come back later :P

### Example Monorepo

The main configuration file for the monorepo is `terrascope.hcl`. This defines
where our root configurations live, as well as what scope types our monorepo
deals with.

Each of the roots has its own `terrascope.hcl` which contains configuration
details about that root. This includes what scope types apply to the root, and
what dependencies the root might have.

Scope values are stored in another hcl file at the top directory of the
monorepo, by default this is `data.hcl` (see `example-data.hcl` for an example).

## Reference

### Terrascope HCL

For schema documentation, please see `docs/hcl-schema.md`

### Terrascope CLI

See `src/README.md`.
