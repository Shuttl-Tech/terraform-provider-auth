module github.com/Shuttl-Tech/terraform-provider-auth

go 1.13

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

require (
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.1.0
	github.com/hashicorp/vault v1.2.3
	github.com/hashicorp/vault/api v1.0.5-0.20190909201928-35325e2c3262
	github.com/mattn/go-isatty v0.0.6 // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
)
