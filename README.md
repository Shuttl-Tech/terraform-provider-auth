# Terraform Provider Auth

This repository contains the source code for terraform auth provider. The provider provides a simple mechanism for
terraform to authenticate with various identity providers.

## Usage

### Provider Configuration

The *auth* provider expects nested configuration block for each supported identity provider. Following identity
providers are supported at the moment:

 - vault: Authenticate with Hashicorp Vault
 
*vault* configuration block expects following attributes:

 - address: (required) Address of an active vault server including the protocol and port.
 - ca_cert_file: (optional) Path to a CA certificate file to validate the server's certificate.
 - ca_cert_dir: (optional) Path to directory containing CA certificate files to validate the server's certificate.
 - skip_tls_verify: (optional) Set to true to skip TLS verification.
 - client_auth: (optional) Nested configuration block for client authentication.
   - cert_file: (required) Path to a file containing the client certificate.
   - key_file: (required) Path to a file containing the private key that the certificate was issued for.
   
### Data Provider

Following data sources are provided by Auth provider:

 - auth_vault: Authenticate with Hashicorp Vault

*auth_vault* data source expects following attributes:

 - auth_backend: (required) Name of Vault auth backend. Only 'aws' is supported at this time.
 - role: (required) Name of the vault role against which Terraform will try to authenticate.
 - mount_path: (required) Vault mount path of the auth backend.
 - aws: (required when `auth_backend` is `aws`) Nested configuration block:
   - use_ec2_metadata: Use EC2 instance metadata and IAM role to authenticate with Vault. No other attribute is supported when this is set.
   - identity: Base64-encoded EC2 instance identity document to authenticate with.
   - signature: Base64-encoded SHA256 RSA signature of the instance identtiy document to authenticate with.
   - pkcs7: PKCS7 signature of the identity document to authenticate with, with all newline characters removed.
   - nonce: The nonce to be used for subsequent login requests.
   - iam_http_request_method: The HTTP method used in signed request.
   - iam_request_url: The Base64-encoded HTTP URL used in the signed request.
   - iam_request_body: The Base64-encoded body of the signed request.
   - iam_request_headers: The Base64-encoded, JSON serialized representation of the sts:GetCallerIdentity HTTP request headers

In addition to the attributes specified above, following attributes are exported by *auth_vault*:

 - lease_duration: Lease duration in seconds relative to the time in lease_start_time.
 - lease_start_time: Time at which the lease was read, using the clock of the system where Terraform was running
 - renewable: True if the duration of this lease can be extended through renewal
 - metadata: The metadata reported by the Vault server
 - policies: The policies assigned to this token
 - accessor: The accessor returned from Vault for this token
 - client_token: The token returned by Vault
 
## Install

It is recommended to download a binary from the Releases tab on Github but you can also build the provider from source.
To build from source see [Building][].  
After obtaining the binary please follow the [official documentation][] to install the provider.

## Building

To build the provider from source you must have Go v1.11 or higher installed on your machine.  
Clone the repository and run build task:

```
git clone git@github.com:Shuttl-Tech/terraform-provider-auth.git $GOPATH/src/github.com/Shuttl-Tech/terraform-provider-auth
cd $GOPATH/src/github.com/Shuttl-Tech/terraform-provider-auth

# build for linux
make linux/amd64

# build for osx
make darwin/amd64
```

# Contributing
All reasonable pull requests are welcome. Before you start working on things it is a good idea to first search through existing pull requests and open issue to make sure your work don't clash with another contributor. If you are unsure of anything please feel free to open an issue and start the discussion.

# License
This code is licensed under the MPLv2 license.

[Building]: #building
[official documentation]: https://www.terraform.io/docs/plugins/basics.html#installing-plugins
