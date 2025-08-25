package dbaas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DBIAAS_URI", nil),
				Description: "DBaaS base URL",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DBIAAS_TOKEN", nil),
				Description: "Token for Bearer auth to the API.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DBIAAS_USERNAME", nil),
				Description: "Username for BASIC auth to the API.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DBIAAS_PASSWORD", nil),
				Description: "Password for BASIC auth to the API.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DBIAAS_INSECURE", nil),
				Description: "Disables TLS verification if using HTTPS.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Timeout in seconds for requests.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable debug mode to trace requests.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dbaas-postgres_pg_db": resourcePgDb(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"datasource_pg_db": dataSourcePgDb(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	headers := make(map[string]string)

	opt := &apiClientOpt{
		token:    d.Get("token").(string),
		username: d.Get("username").(string),
		password: d.Get("password").(string),
		insecure: d.Get("insecure").(bool),
		uri:      d.Get("uri").(string),
		headers:  headers,
		timeout:  d.Get("timeout").(int),
		debug:    d.Get("debug").(bool),
	}

	client, err := NewAPIClient(opt)
	return client, err
}
