org "acme-inc" {
  attributes = {
    aws_account_id = 12345678
  }
  roots = ["org-account-setup", "sso-groups"]

  platform "gold" {
    domain "commerce" {
      roots = ["sso-permissions"]
      environment "dev" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = [
          "account-settings",
          "aws-config-global",
          "terraform-backend",
          "foo"
        ]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "stage" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "prod" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "data-science" {
      roots = ["sso-permissions"]
      environment "dev" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "prod" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "networking" {
      roots = ["sso-permissions"]
      environment "global" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global"]
        region "us-west-1" {
          roots = ["aws-config-regional"]
        }
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "product" {
      roots = ["sso-permissions"]
      environment "dev" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-1" {
          roots = ["aws-config-regional"]
        }
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "stage" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-1" {
          roots = ["aws-config-regional"]
        }
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "prod" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-1" {
          roots = ["aws-config-regional"]
        }
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "security-portal" {
      roots = ["sso-permissions"]
      environment "global" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
  }
  platform "silver" {
    domain "networking" {
      roots = ["sso-permissions"]
      environment "global" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global"]
        region "us-west-1" {
          roots = ["aws-config-regional"]
        }
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "red" {
      roots = ["sso-permissions"]
      environment "dev" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "stage" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
      environment "prod" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global", "terraform-backend"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
    domain "security-portal" {
      roots = ["sso-permissions"]
      environment "global" {
        attributes = {
          aws_account_id = 12345678
        }
        roots = ["account-settings", "aws-config-global"]
        region "us-west-2" {
          roots = ["aws-config-regional"]
        }
      }
    }
  }
}
