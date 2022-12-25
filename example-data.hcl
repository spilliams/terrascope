scope "org" "acme-inc" {
  aws_account_id = 12345678

  scope "platform" "gold" {
    scope "domain" "commerce" {
      scope "environment" "dev" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
      scope "environment" "stage" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
      scope "environment" "prod" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "data-science" {
      scope "environment" "dev" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
      scope "environment" "prod" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "networking" {
      scope "environment" "global" {
        aws_account_id = 12345678
        scope "region" "us-west-1" {}
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "product" {
      scope "environment" "dev" {
        aws_account_id = 12345678
        scope "region" "us-west-1" {}
        scope "region" "us-west-2" {}
      }
      scope "environment" "stage" {
        aws_account_id = 12345678
        scope "region" "us-west-1" {}
        scope "region" "us-west-2" {}
      }
      scope "environment" "prod" {
        aws_account_id = 12345678
        scope "region" "us-west-1" {}
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "security-portal" {
      scope "environment" "global" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
    }
  }
  scope "platform" "silver" {
    scope "domain" "networking" {
      scope "environment" "global" {
        aws_account_id = 12345678
        scope "region" "us-west-1" {}
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "red" {
      scope "environment" "dev" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
      scope "environment" "stage" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
      scope "environment" "prod" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
    }
    scope "domain" "security-portal" {
      scope "environment" "global" {
        aws_account_id = 12345678
        scope "region" "us-west-2" {}
      }
    }
  }
}
