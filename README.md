# terraboots

My attempt at a Terraform build orchestrator, for large platform projects with
hundreds of root modules.

This repository contains both the source code for the tool `terraboots`, as well
as a sample monorepo managed by that tool.

## Example Monorepo

The main configuration file for the monorepo is `terraboots.hcl`. This defines
where our "modules" and "roots" live (in `terraform/modules/` and
`terraform/roots/` respectively).

The main concept to terraboots is that it expects you to manage your platform
with many small terraform root configurations. This might make more sense with
a hypothetical secenario.

### Hypothetically...

Let's say your company ("Acme Inc") wants to operate in the cloud. They've
chosen AWS as a cloud provider, and have given you a credit card to go set it
up.

You already know that your engineering department has several teams, that each
want to work in their own sandbox environment, and they each want to deploy
their applications on their own (read: they use DevOps methods). You also know
that it'll probably be more secure if you can put bulkheads between each teams'
environments.

You decide to use AWS Organizations to power your platform. You also want to
keep some bulkheads in place to help contain certain high-security elements from
lower-security elements. With this in ind, you set up the following AWS
accounts:

1. an Organization account. This will be your Consolidated Billing account for
   the entirety of Acme Inc's engineering department.
2. a Security Portal account. This is where you'll manage the IAM users for your
   engineers, managers, third-party vendors, etc.
3. a Networking account. This is where you'll set up a VPN that will connect all
   the domain accounts' private subnets to each other, so that engineers can do
   their work without sending database credentials over the Internet.
4. a set of Domain accounts: one set per team, and within each set will be a
   Dev, Stage and Prod account.

Since you're just starting out, you set up one set of Domain accounts for your
primary engineering team: the Product team. This means you now have 6 accounts
to manage! Organization, Security Portal, Networking, Product Dev, Product
Stage, and Product Prod.

Because you also need a sandbox for your own team to develop changes to your
platform, you need a second Security Portal, Networking and set of Domain
accounts. So this is 5 more accounts to manage, and the names "Security Portal"
and "Networking" are already taken. Time to introduce the first core concept of
terraboots: scope.

A "scope" here refers to a layer of infrastructure. You can define your own set
of scopes for your own application, and `terraboots` will pick up on it, but for
this hypothetical situation we want to define the following:

- the "org" scope contains an Organization account and several platforms (e.g.
  "Acme Inc")
- the "platform" scope contains one Security Portal and one Networking
  account, as well as several domains (e.g. "Gold", "Silver", "Bronze")
- the "domain" scope contains three environments (e.g. "Product",
  "DataScience", etc)
- the "environment" scope contains the AWS accounts for our engineers ("Product
  Dev" etc), and any regions within those environments (e.g. "Dev", "Stage",
  "Prod")
- the "region" scope contains regional resources in each environment account
  (e.g. "us-west-1", "us-west-2")

At every level of this infrastructure, we could have a terraform root
configuration that says "build me these resources". At any level we could also
want our roots to depend on the levels before it: say, at the environment scope
we provision a single S3 bucket in every domain account, and at the region scope
we provision prefixes inside that s3 bucket (something this example monorepo
shows off with its aws-config modules and roots).

The goal of `terraboots` is to be able to manage this kind of infrastructure
without having 500 separate directories which are mostly copy-pasted from each
other. Yes, there are several other projects out there that seek to DRY up a
terraform system like this, but they each have their limitations. Terragrunt
still wants you to maintain a whole tree of directories of `terragrunt.hcl`
files, and only has one layer of inheritance. Terraspace has you writing a lot
of ruby, and I don't like maintaining a set of dependencies just for
development.

## Terraboots CLI

See `src/README.md`.
