package dbaas

import (
	"encoding/json"
	"fmt"
	//	"net/http"
	//	"bytes"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"
)

func resourcePgDb() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgCreate,
		ReadContext:   resourcePgRead,
		UpdateContext: resourcePgUpdate,
		DeleteContext: resourcePgDelete,
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tenant identifier",
			},
			"dbname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name with tenant prefix",
			},
			"dbsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Size of the database",
			},
			"dbconn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max connections",
			},
			"db_release": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the database",
			},
			"db_win": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Database window",
			},
			"role_map": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rol_type":    {Type: schema.TypeString, Required: true},
						"rol_name":    {Type: schema.TypeString, Required: true},
						"rol_group":   {Type: schema.TypeString, Optional: true},
						"rol_conn":    {Type: schema.TypeInt, Optional: true},
						"rol_timeout": {Type: schema.TypeInt, Optional: true},
					},
				},
			},
			"extension_map": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_name": {Type: schema.TypeString, Required: true},
						"ext_opt_map": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"hba_map": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hba_role":       {Type: schema.TypeString, Required: true},
						"hba_addr":       {Type: schema.TypeString, Required: true},
						"hba_src_tenant": {Type: schema.TypeString, Optional: true},
						"hba_auth":       {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"output": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"worker": {Type: schema.TypeString, Computed: true},
						"data":   {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Job name",
							ForceNew:    true,
						},
						"type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Default:  getMetadata("type", nil).(string),
						},
						"affinity": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
							Default:  getMetadata("affinity", nil).(string),
						},
						"retry": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  getMetadata("retry", nil).(bool),
						},
						"wait_retry": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      getMetadata("wait_retry", nil).(int),
							ValidateFunc: validation.IntAtLeast(15),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      getMetadata("timeout", nil).(int),
							ValidateFunc: validation.IntAtLeast(15),
						},
					},
				},
			},
		},
	}
}

// La fonction utilitaire pour préparer les données
func prepareDatabaseRequestData(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"tenant":        d.Get("tenant").(string),
		"dbname":        d.Get("dbname").(string),
		"dbsize":        d.Get("dbsize").(int),
		"dbconn":        d.Get("dbconn").(int),
		"db_release":    d.Get("db_release").(string),
		"db_win":        d.Get("db_win").(int),
		"role_map":      d.Get("role_map").([]interface{}),
		"extension_map": d.Get("extension_map").([]interface{}),
		"hba_map":       d.Get("hba_map").([]interface{}),
	}
}

func resourcePgCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client)
	name := d.Get("metadata.0.name").(string)
	jobType := d.Get("metadata.0.type").(string)
	if v, ok := d.GetOk("dbconn"); ok {
		maxDbConn := v.(int)

		if roles, ok := d.GetOk("role_map"); ok {
			roleList := roles.([]interface{})
			for _, roleItem := range roleList {
				roleData := roleItem.(map[string]interface{})
				if connVal, exists := roleData["rol_conn"]; exists && connVal != nil {
					rolConn := connVal.(int)
					if rolConn > maxDbConn {
						return diag.Errorf("La 'rol_conn' (%d) ne peut pas dépasser 'dbconn' (%d)", rolConn, maxDbConn)
					}
				}
			}
		}
	}
	payload := prepareDatabaseRequestData(d)
	dataBytes, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}
	data := fmt.Sprintf(`{"database": %s}`, string(dataBytes))

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + client.token,
	}

	path := fmt.Sprintf("/jobs/%s/%s", jobType, name)

	jobraw, err := client.send_request("POST", path, data, headers)
	baseMsg := fmt.Sprintf("Cannot create database '%s' -", name)
	fullurl := fmt.Sprintf("%s%s", client.uri, path)
	err = handleHTTPError(err, jobraw, fullurl, baseMsg)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitDatabaseJobState(ctx, "created", d, meta, false); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourcePgRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := PgDbRead(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func PgDbRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api_client)

	idArr := strings.Split(d.Id(), "/")
	var name, affinity, jobType string
	if len(idArr) >= 3 {
		jobType = idArr[0]
		affinity = idArr[1]
		name = idArr[2]
	} else {
		metadata := d.Get("metadata").([]interface{})[0].(map[string]interface{})
		name = metadata["name"].(string)
		jobType = metadata["type"].(string)
	}

	path := fmt.Sprintf("/jobs/%s/%s", jobType, name)
	headers := map[string]string{}

	jobraw, err := client.send_request("GET", path, "", headers)

	baseMsg := fmt.Sprintf("Cannot read database '%s' -", name)
	fullurl := fmt.Sprintf("%s%s", client.uri, path)
	err = handleHTTPError(err, jobraw, fullurl, baseMsg)
	if err != nil {
		if strings.Contains(err.Error(), "response code '404'") {
			d.SetId("")
			return nil
		}
		return err
	}

	job, err := getJob(client, jobType, name)
	if err != nil {
		if strings.Contains(err.Error(), "response code '404'") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Unable to retrieve database '%s': %s", name, err)
	}

	if err := d.Set("tenant", job.Tenant); err != nil {
		return err
	}

	metadata := []map[string]interface{}{{
		"name":       name,
		"affinity":   affinity,
		"type":       jobType,
		"retry":      getMetadata("retry", d).(bool),
		"timeout":    getMetadata("timeout", d).(int),
		"wait_retry": getMetadata("wait_retry", d).(int),
	}}
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}

	var jobOutput []JobOutput
	if len(job.Output) > 0 && job.Output[0].State == "finished" {
		jobOutput = append(jobOutput, JobOutput{
			Worker: job.Output[0].Worker,
			Data:   job.Output[0].Data,
		})
	}
	if err := d.Set("output", jobOutput); err != nil {
		return err
	}

	return nil
}

func resourcePgUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("dbsize", "dbconn", "db_release", "db_win", "role_map", "extension_map", "hba_map") {
		client := meta.(*api_client)

		name := d.Get("metadata.0.name").(string)
		jobType := d.Get("metadata.0.type").(string)
		affinity := d.Get("metadata.0.affinity").(string)
		tenant := d.Get("tenant").(string)

		payload := map[string]interface{}{
			"tenant": tenant,
			"dbname": d.Get("dbname").(string),
		}

		if v, ok := d.GetOk("dbconn"); ok {
			maxDbConn := v.(int)
			payload["dbconn"] = maxDbConn

			// Validation : chaque rol_conn <= dbconn
			if roles, ok := d.GetOk("role_map"); ok {
				roleList := roles.([]interface{})
				for _, roleItem := range roleList {
					roleData := roleItem.(map[string]interface{})
					if connVal, exists := roleData["rol_conn"]; exists && connVal != nil {
						rolConn := connVal.(int)
						if rolConn > maxDbConn {
							return diag.Errorf("La 'rol_conn' (%d) ne peut pas dépasser 'dbconn' (%d)", rolConn, maxDbConn)
						}
					}
				}
			}
		}

		if v, ok := d.GetOk("dbsize"); ok {
			payload["dbsize"] = v
		}
		if v, ok := d.GetOk("db_release"); ok {
			payload["db_release"] = v
		}
		if v, ok := d.GetOk("db_win"); ok {
			payload["db_win"] = v
		}
		if v, ok := d.GetOk("role_map"); ok {
			payload["role_map"] = v
		}
		if v, ok := d.GetOk("extension_map"); ok {
			payload["extension_map"] = v
		}
		if v, ok := d.GetOk("hba_map"); ok {
			payload["hba_map"] = v
		}

		dataBytes, err := json.Marshal(payload)
		if err != nil {
			return diag.FromErr(err)
		}
		data := fmt.Sprintf(`{"database": %s}`, string(dataBytes))

		headers := map[string]string{
			"Content-Type":     "application/json",
			"Moldapi-Affinity": affinity,
		}

		path := fmt.Sprintf("/jobs/%s/%s", jobType, name)
		jobraw, err := client.send_request("PUT", path, data, headers)
		fullurl := fmt.Sprintf("%s%s", client.uri, path)
		baseMsg := fmt.Sprintf("Cannot update database '%s' -", name)

		err = handleHTTPError(err, jobraw, fullurl, baseMsg)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitDatabaseJobState(ctx, "updated", d, meta, false)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diag.Diagnostics{}
}

func resourcePgDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api_client)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + client.token,
	}

	name := d.Get("metadata.0.name").(string)
	jobType := d.Get("metadata.0.type").(string)
	path := fmt.Sprintf("/jobs/%s/%s", jobType, name)

	_, err := client.send_request("DELETE", path, "", headers)
	if err != nil {
		// Ignore errors like HTTP 500 and 404, as well as specific task-related errors
		if strings.Contains(err.Error(), "500") ||
			strings.Contains(err.Error(), "response code '404'") ||
			strings.Contains(err.Error(), "Job not found") ||
			strings.Contains(err.Error(), "Code HTTP inattendu : 500") ||
			strings.Contains(err.Error(), "could not be deleted: Task error") ||
			strings.HasSuffix(err.Error(), "could not be deleted: Task error") {
			err = nil // Ignore error and continue
		} else {
			return diag.FromErr(fmt.Errorf(
				"Cannot delete repack '%s' at '%s': %v",
				name,
				fmt.Sprintf("%s%s", client.uri, path),
				err))
		}

	}

	err = waitDatabaseJobState(ctx, "deleted", d, meta, true)
	if err != nil {
		// If error 404, 500, or "Job not found", the resource is effectively gone, no issue
		if strings.Contains(err.Error(), "response code '404'") ||
			strings.Contains(err.Error(), "Job not found") ||
			strings.Contains(err.Error(), "Code HTTP inattendu : 500") ||
			strings.Contains(err.Error(), "could not be deleted: Task error") ||
			strings.HasSuffix(err.Error(), "could not be deleted: Task error") {
			err = nil
		} else {
			return diag.FromErr(fmt.Errorf("Failed to wait for deletion of '%s': %v", name, err))
		}
	}
	d.SetId("")

	return diag.Diagnostics{}
}
func isZero(job JobProperties) bool {
	return job.Name == "" && job.Type == "" && job.State == ""
}

func waitDatabaseJobState(
	ctx context.Context,
	state string,
	d *schema.ResourceData,
	meta interface{},
	isDeletion bool,
) error {
	client := meta.(*api_client)

	name := d.Get("metadata.0.name").(string)
	jobType := d.Get("metadata.0.type").(string)
	affinity := d.Get("metadata.0.affinity").(string)
	retry := d.Get("metadata.0.retry").(bool)
	timeout := d.Get("metadata.0.timeout").(int) // Maximum wait time for entire operation
	waitRetry := d.Get("metadata.0.wait_retry").(int)

	startTime := time.Now()

	return resource.Retry(time.Duration(timeout)*time.Second, func() *resource.RetryError {
		job, err := getJob(client, jobType, name)

		if err != nil {
			if isDeletion &&
				(strings.Contains(err.Error(), "response code '404'") ||
					strings.Contains(err.Error(), "Job not found") ||
					strings.Contains(err.Error(), "Code HTTP inattendu : 500") ||
					strings.HasSuffix(err.Error(), "could not be deleted: Task error")) {
				return nil // Assume deletion completed despite errors
			}
			return resource.NonRetryableError(fmt.Errorf("Error when getting exploit job '%s': %s", name, err))
		}

		if isZero(job) {
			if isDeletion {
				return nil // Deletion already done or job doesn't exist: fine
			}
			return resource.NonRetryableError(fmt.Errorf("Unexpected empty job returned for '%s'", name))
		}

		d.SetId(fmt.Sprintf("%s/%s/%s", jobType, affinity, name))

		if job.State == "finished" || (isDeletion && (job.State == "deleted" || job.State == "dead")) {
			if err := PgDbRead(ctx, d, meta); err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		}

		// Handle the "pending" state specifically for deletion
		if job.State == "pending" && isDeletion {
			if time.Since(startTime) >= time.Duration(timeout)*time.Second {
				return nil // Ignore pending state after maximum wait time
			}
			return resource.RetryableError(fmt.Errorf("Job '%s' is pending, waiting for state change", name))
		}
		errorMsg := fmt.Errorf("Expected exploit job '%s' to be %s but was in state %s", name, state, job.State)

		if job.State == "failed" && retry {
			time.Sleep(time.Duration(waitRetry) * time.Second)
			_, err := retryJob(client, jobType, name)
			if err != nil {
				errorMsg = fmt.Errorf("%v%v", errorMsg, friendlyYAMLError(job))
				return resource.NonRetryableError(errorMsg)
			}
		} else if job.State == "failed" || job.State == "dead" {
			errorMsg = fmt.Errorf("%v; details: %v", errorMsg, friendlyYAMLError(job))
			return resource.NonRetryableError(errorMsg)
		}

		errorMsg = fmt.Errorf("%v%v", errorMsg, friendlyYAMLError(job))
		return resource.RetryableError(errorMsg)
	})
}

type Role struct {
	RolType    string `json:"rol_type"`
	RolName    string `json:"rol_name"`
	RolGroup   string `json:"rol_group,omitempty"`
	RolConn    int    `json:"rol_conn,omitempty"`
	RolTimeout int    `json:"rol_timeout,omitempty"`
}

type Extension struct {
	ExtName   string            `json:"ext_name"`
	ExtOptMap map[string]string `json:"ext_opt_map,omitempty"`
}

type HBA struct {
	HbaRole      string `json:"hba_role"`
	HbaAddr      string `json:"hba_addr"`
	HbaSrcTenant string `json:"hba_src_tenant,omitempty"`
	HbaAuth      string `json:"hba_auth,omitempty"`
}

type Database struct {
	Tenant       string      `json:"tenant"`
	DBName       string      `json:"dbname"`
	DBSize       int         `json:"dbsize"`
	DBConn       int         `json:"dbconn"`
	DBRelease    string      `json:"db_release"`
	DBWin        int         `json:"db_win"`
	RoleMap      []Role      `json:"role_map,omitempty"`
	ExtensionMap []Extension `json:"extension_map,omitempty"`
	HbaMap       []HBA       `json:"hba_map,omitempty"`
}

func flattenDatabase(db Database) map[string]interface{} {
	roleMaps := make([]map[string]interface{}, 0, len(db.RoleMap))
	for _, r := range db.RoleMap {
		roleMaps = append(roleMaps, map[string]interface{}{
			"rol_type":    r.RolType,
			"rol_name":    r.RolName,
			"rol_group":   r.RolGroup,
			"rol_conn":    r.RolConn,
			"rol_timeout": r.RolTimeout,
		})
	}

	extensionMaps := make([]map[string]interface{}, 0, len(db.ExtensionMap))
	for _, ext := range db.ExtensionMap {
		extMap := map[string]interface{}{
			"ext_name": ext.ExtName,
		}
		if ext.ExtOptMap != nil {
			extMap["ext_opt_map"] = ext.ExtOptMap
		}
		extensionMaps = append(extensionMaps, extMap)
	}

	hbaMaps := make([]map[string]interface{}, 0, len(db.HbaMap))
	for _, h := range db.HbaMap {
		hbaMaps = append(hbaMaps, map[string]interface{}{
			"hba_role":       h.HbaRole,
			"hba_addr":       h.HbaAddr,
			"hba_src_tenant": h.HbaSrcTenant,
			"hba_auth":       h.HbaAuth,
		})
	}

	return map[string]interface{}{
		"tenant":        db.Tenant,
		"dbname":        db.DBName,
		"dbsize":        db.DBSize,
		"dbconn":        db.DBConn,
		"db_release":    db.DBRelease,
		"db_win":        db.DBWin,
		"role_map":      roleMaps,
		"extension_map": extensionMaps,
		"hba_map":       hbaMaps,
	}
}
