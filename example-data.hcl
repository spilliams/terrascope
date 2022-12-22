org "acme-inc" {
  aws_account_id = 12345678

  platform "gold" {
    domain "commerce" {
      environment "dev" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
      environment "stage" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
      environment "prod" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
    }
    domain "data-science" {
      environment "dev" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
      environment "prod" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
    }
    domain "networking" {
      environment "global" {
        aws_account_id = 12345678
        region "us-west-1" {}
        region "us-west-2" {}
      }
    }
    domain "product" {
      environment "dev" {
        aws_account_id = 12345678
        region "us-west-1" {}
        region "us-west-2" {}
      }
      environment "stage" {
        aws_account_id = 12345678
        region "us-west-1" {}
        region "us-west-2" {}
      }
      environment "prod" {
        aws_account_id = 12345678
        region "us-west-1" {}
        region "us-west-2" {}
      }
    }
    domain "security-portal" {
      environment "global" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
    }
  }
  platform "silver" {
    domain "networking" {
      environment "global" {
        aws_account_id = 12345678
        region "us-west-1" {}
        region "us-west-2" {}
      }
    }
    domain "red" {
      environment "dev" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
      environment "stage" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
      environment "prod" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
    }
    domain "security-portal" {
      environment "global" {
        aws_account_id = 12345678
        region "us-west-2" {}
      }
    }
  }
}
