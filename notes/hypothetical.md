# Hypothetically...

Let's say your company ("Acme Inc") wants to operate in the cloud. They've
chosen AWS as a cloud provider, and have given you a credit card to go set it
up.

You already know that your engineering department has several teams, that each
want to work in their own sandbox environment, and they each want to deploy
their applications on their own (read: they use DevOps methods). You also know
that it'll probably be more secure if you can put bulkheads between each teams'
environments.

You decide to use AWS Organizations to power your infrastructure. You also want
to keep some bulkheads in place to help contain certain organization- or
networking-related elements from application- or product-related elements. With
this in ind, you set up the following AWS accounts:

1. an Organization account. This will be your Consolidated Billing account for
   the entirety of Acme Inc's engineering department. All other accounts will be
   made inside of this account.
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
infrastructure, you need a second Security Portal, Networking and set of Domain
accounts. So this is 5 more accounts to manage, and the names "Security Portal"
and "Networking" are already taken. Time to introduce the first core concept of
terrascope: scope.

A "scope" here refers to a layer of infrastructure. You can define your own set
of scopes for your infrastructure, and `terrascope` will pick up on it, but for
this hypothetical situation we want to define the following:

- an `org` (e.g. "Acme Inc") contains an Organization account and several
  platforms.
- a `platform` (e.g. "Gold", "Silver", or "Bronze") contains one Security Portal
  and one Networking account, as well as several domains.
- a `domain` (e.g. "Product", or "Red") contains three environments.
- an `environment` (e.g. "Dev", "Stage", or "Prod") contains the AWS accounts
  for our engineers ("Product Dev" etc), as well as any primary deployment
  regions within those environments.
- a `region` (e.g. "us-west-2") contains regional resources in an environment
  account.

The scope types for this hypothetical situation are: `org`, `platform`,
`domain`, `environment`, and `region`.

At every level of this infrastructure, we could have a terraform root
configuration that says "build me these resources". At any level we could also
want our roots to depend on the levels before it. For instance: at the
`environment` scope we provision a single S3 bucket in every domain account, and
at the `region` scope below it we provision prefixes inside that s3 bucket.

The goal of `terrascope` is to be able to manage this kind of infrastructure
without having 500 separate terraform root configurations which are mostly
copy-pasted from each other.
