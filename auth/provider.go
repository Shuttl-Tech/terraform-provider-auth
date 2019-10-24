package auth

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/vault/api"
)

type Config struct {
	Vault *api.Client
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vault": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_ADDR", nil),
							Description: "URL of the root of the target Vault server",
						},
						"ca_cert_file": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_CACERT", ""),
							Description: "Path to a CA certificate file to validate the server's certificate.",
						},
						"ca_cert_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_CAPATH", ""),
							Description: "Path to directory containing CA certificate files to validate the server's certificate.",
						},
						"client_auth": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Client authentication credentials.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_file": {
										Type:        schema.TypeString,
										Required:    true,
										DefaultFunc: schema.EnvDefaultFunc("VAULT_CLIENT_CERT", ""),
										Description: "Path to a file containing the client certificate.",
									},
									"key_file": {
										Type:        schema.TypeString,
										Required:    true,
										DefaultFunc: schema.EnvDefaultFunc("VAULT_CLIENT_KEY", ""),
										Description: "Path to a file containing the private key that the certificate was issued for.",
									},
								},
							},
						},
						"skip_tls_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_SKIP_VERIFY", ""),
							Description: "Set this to true only if the target Vault server is an insecure development instance.",
						},
					},
				},
			},
		},

		ConfigureFunc: providerConfigure,

		DataSourcesMap: map[string]*schema.Resource{
			"auth_vault": vaultAuthDataSource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfg := &Config{}
	attrs := d.Get("vault").([]interface{})
	if len(attrs) > 1 {
		return nil, fmt.Errorf("'vault' block can appear only once")
	}

	if len(attrs) > 0 {
		vcfg, err := configureVaultClient(attrs[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		cfg.Vault = vcfg
	}

	return cfg, nil
}

func configureVaultClient(attrs map[string]interface{}) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = attrs["address"].(string)

	clientAuthI := attrs["client_auth"].([]interface{})
	if len(clientAuthI) > 1 {
		return nil, fmt.Errorf("client_auth block may appear only once")
	}

	clientAuthCert := ""
	clientAuthKey := ""
	if len(clientAuthI) == 1 {
		clientAuth := clientAuthI[0].(map[string]interface{})
		clientAuthCert = clientAuth["cert_file"].(string)
		clientAuthKey = clientAuth["key_file"].(string)
	}

	err := config.ConfigureTLS(&api.TLSConfig{
		CACert:   attrs["ca_cert_file"].(string),
		CAPath:   attrs["ca_cert_dir"].(string),
		Insecure: attrs["skip_tls_verify"].(bool),

		ClientCert: clientAuthCert,
		ClientKey:  clientAuthKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to configure TLS for Vault API: %s", err)
	}

	config.HttpClient.Transport = logging.NewTransport("Vault", config.HttpClient.Transport)

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Vault API: %s", err)
	}

	return client, nil
}
