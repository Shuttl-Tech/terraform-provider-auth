package auth

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/vault/api"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	"strings"
	"time"
)

func vaultAuthDataSource() *schema.Resource {
	return &schema.Resource{
		Read: vaultAuthLoginRead,

		Schema: map[string]*schema.Schema{
			"auth_backend": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of Vault Auth backend",
				ValidateFunc: validation.StringInSlice([]string{"aws"}, true),
				ForceNew:     true,
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the role against which the token will be created",
				ForceNew:    true,
			},
			"mount_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mount path of the auth backend",
				ForceNew:    true,
			},
			"aws": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_ec2_metadata": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Use EC2 instance metadata to authenticate with vault.",
							ForceNew:    true,
							ConflictsWith: []string{
								"aws.identity", "aws.signature", "aws.pkcs7", "aws.iam_http_request_method",
								"aws.iam_request_url", "aws.iam_request_body", "aws.iam_request_headers",
							},
						},
						"identity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64-encoded EC2 instance identity document to authenticate with.",
							ForceNew:    true,
						},
						"signature": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Base64-encoded SHA256 RSA signature of the instance identtiy document to authenticate with.",
							ForceNew:    true,
						},
						"pkcs7": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PKCS7 signature of the identity document to authenticate with, with all newline characters removed.",
							ForceNew:    true,
						},
						"nonce": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The nonce to be used for subsequent login requests.",
							Computed:    true,
							ForceNew:    true,
						},
						"iam_http_request_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HTTP method used in the signed request.",
							ForceNew:    true,
						},
						"iam_request_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Base64-encoded HTTP URL used in the signed request.",
							ForceNew:    true,
						},
						"iam_request_body": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Base64-encoded body of the signed request.",
							ForceNew:    true,
						},
						"iam_request_headers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Base64-encoded, JSON serialized representation of the sts:GetCallerIdentity HTTP request headers.",
							ForceNew:    true,
						},
					},
				},
			},

			"lease_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Lease duration in seconds relative to the time in lease_start_time.",
			},
			"lease_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which the lease was read, using the clock of the system where Terraform was running",
			},
			"renewable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the duration of this lease can be extended through renewal.",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The metadata reported by the Vault server.",
				Elem:        schema.TypeString,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The policies assigned to this token.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"accessor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The accessor returned from Vault for this token.",
			},
			"client_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The token returned by Vault.",
				Sensitive:   true,
			},
		},
	}
}

func vaultAuthLoginRead(d *schema.ResourceData, meta interface{}) error {
	rconf := meta.(*Config)

	if d.Get("auth_backend").(string) == "aws" {
		awsAuthConf := d.Get("aws").([]interface{})
		if len(awsAuthConf) > 1 {
			return fmt.Errorf("'aws' block must appear only once")
		}
		if len(awsAuthConf) != 1 {
			return fmt.Errorf("'aws' block must be provided with auth_backend 'aws'")
		}

		return vaultAuthLoginWithAWS(d, rconf.Vault, awsAuthConf[0].(map[string]interface{}))
	}

	return fmt.Errorf("unsupported auth_backend %q", d.Get("auth_backend"))
}

func vaultAuthLoginWithAWS(d *schema.ResourceData, client *api.Client, awsConfig map[string]interface{}) error {
	var (
		id     string
		secret *api.Secret
		err    error
	)

	backend := strings.Trim(d.Get("mount_path").(string), "/")

	if awsConfig["use_ec2_metadata"].(bool) {
		handler := credAws.CLIHandler{}
		meta := map[string]string{
			"mount":                 backend,
			"role":                  d.Get("role").(string),
			"aws_access_key_id":     "",
			"aws_secret_access_key": "",
			"aws_security_token":    "",
		}

		secret, err = handler.Auth(client, meta)
		if err != nil {
			return err
		}
	} else {
		path := "auth/" + backend + "/login"
		data := map[string]interface{}{}

		if v, ok := d.GetOk("role"); ok {
			data["role"] = v
		}

		if v, ok := awsConfig["identity"]; ok {
			data["identity"] = v
		}

		if v, ok := awsConfig["signature"]; ok {
			data["signature"] = v
		}

		if v, ok := awsConfig["pkcs7"]; ok {
			data["pkcs7"] = v
		}

		if v, ok := awsConfig["nonce"]; ok {
			data["nonce"] = v
		}

		if v, ok := awsConfig["iam_http_request_method"]; ok {
			data["iam_http_request_method"] = v
		}

		if v, ok := awsConfig["iam_request_url"]; ok {
			data["iam_request_url"] = v
		}

		if v, ok := awsConfig["iam_request_body"]; ok {
			data["iam_request_body"] = v
		}

		if v, ok := awsConfig["iam_request_headers"]; ok {
			data["iam_request_headers"] = v
		}

		secret, err = client.Logical().Write(path, data)
		if err != nil {
			return fmt.Errorf("error reading from Vault: %s", err)
		}
	}

	id = "accessor:" + secret.Auth.Accessor
	nonce, ok := secret.Auth.Metadata["nonce"]
	if ok {
		id = "nonce:" + nonce
	}

	d.SetId(id)
	_ = d.Set("lease_id", secret.LeaseID)
	_ = d.Set("lease_duration", secret.Auth.LeaseDuration)
	_ = d.Set("lease_start_time", time.Now().Format(time.RFC3339))
	_ = d.Set("renewable", secret.Auth.Renewable)
	_ = d.Set("metadata", secret.Auth.Metadata)
	_ = d.Set("policies", secret.Auth.Policies)
	_ = d.Set("accessor", secret.Auth.Accessor)
	_ = d.Set("client_token", secret.Auth.ClientToken)
	_ = d.Set("nonce", nonce)

	return nil
}
