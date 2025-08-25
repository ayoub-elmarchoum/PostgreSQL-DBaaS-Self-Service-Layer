package dbaas

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func dataSourcePgDb() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePgRead,

		Schema: map[string]*schema.Schema{
			"module": {
				Type:        schema.TypeString,
				Description: "Name of the module",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the database",
				Required:    true,
			},
			"data_map": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database description map",
			},
			"role_map": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Role configuration map",
			},
			"extension_map": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Extension configuration map",
			},
			"hba_map": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "HBA configuration map",
			},
		},
	}
}
func dataSourcePgRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api_client)
	name := d.Get("metadata.0.name").(string)
	jobType := d.Get("metadata.0.type").(string) // should be "database"

	job, err := getJob(client, jobType, name)
	if err != nil {
		if strings.Contains(err.Error(), "response code '404'") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to retrieve the database '%s': %s", name, err)
	}

	output := make([]JobOutput, 0)
	for _, item := range job.Output {
		if item.State == "finished" {
			output = append(output, JobOutput{Worker: item.Worker, Data: item.Data})
		} else {
			for _, retry := range item.Retries {
				if retry.State == "finished" {
					output = append(output, JobOutput{Worker: retry.Worker, Data: retry.Data})
				}
			}
		}
	}
	d.SetId(fmt.Sprintf("%s/%s", jobType, name))

	var data Database
	if json.Valid([]byte(job.CurrentData)) {
		err = json.Unmarshal([]byte(job.CurrentData), &data)
		if err != nil {
			return fmt.Errorf("unable to decode database '%s' data: %v", name, err)
		}
		if err := d.Set("database", flattenDatabase(data)); err != nil {
			return err
		}
	}

	d.Set("output", output)

	return nil
}
