package dbaas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
	"time"
)

var (
	aclNameCompliantRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[\-_]+[a-zA-Z0-9]+)*$`)
	aclNameLengthRegexp    = regexp.MustCompile(`^.{1,40}$`)
	friendlyOutputRegexp   = regexp.MustCompile(`(?m)^.*?(Logs: .*)$`)
	escapeCRNL             = regexp.MustCompile(`(?m)(\s)*\\n`)
)

var defaultMetadataAttributes = map[string]interface{}{
	"type":       "database-postgres",
	"affinity":   "all",
	"retry":      true,
	"wait_retry": 5,
	"timeout":    150,
}

func getMetadata(key string, d *schema.ResourceData) interface{} {
	if d == nil {
		return defaultMetadataAttributes[key]
	}
	metadata := d.Get("metadata.0").(map[string]interface{})
	if len(metadata) != 0 {
		return metadata[key]
	} else {
		return defaultMetadataAttributes[key]
	}
}

func jsonPrettyPrint(input []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "  ")
	if err != nil {
		return string(input)
	}
	return out.String()
}

// Marshal is a UTF-8 friendly marshaler.  Go's json.Marshal is not UTF-8
// friendly because it replaces the valid UTF-8 and JSON characters "&". "<",
// ">" with the "slash u" unicode escaped forms (e.g. \u0026).  It preemptively
// escapes for HTML friendliness.
func JsonMarshalHTML(i interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(i)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

func handleHTTPError(err error, body string, url, baseMsg string) error {
	if err != nil {
		return fmt.Errorf("%s %v", baseMsg, err)
	}
	var result interface{}
	if json.Valid([]byte(body)) {
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			return fmt.Errorf("%s %v", baseMsg, err)
		}
	} else {
		return fmt.Errorf("%s Empty/Invalid json response from %s", baseMsg, url)
	}
	return nil
}

type JobProperties struct {
	Name        string                 `json:"name" example:"test"`
	State       string                 `json:"state" example:"finished"`
	TaskID      string                 `json:"taskid" example:"c1c6ef46-86c4-4443-a60a-94c6b7abadd5"`
	User        string                 `json:"user" example:"test"`
	Tenant      string                 `json:"tenant" example:"test"`
	Type        string                 `json:"type" example:"test"`
	CurrentData string                 `json:"current_data"`
	DesiredData string                 `json:"desired_data"`
	BeginTime   time.Time              `json:"begintime"`
	EndTime     time.Time              `json:"endtime"`
	Error       string                 `json:"error,omitempty"`
	Output      []TaskWorkerProperties `json:"output,omitempty"`
}

type TaskWorkerProperties struct {
	Worker  string                 `json:"worker" binding:"required"`
	State   string                 `json:"state" binding:"required"`
	Data    string                 `json:"data" example:"processed!"`
	Logs    string                 `json:"logs" example:"blabla"`
	EndTime time.Time              `json:"endtime"`
	Retries []TaskWorkerProperties `json:"retries,omitempty"`
}

type JobOutput struct {
	Worker string `json:"worker" mapstructure:"worker"`
	Data   string `json:"data" mapstructure:"data"`
}

func getJob(client *api_client, jobType, jobName string) (JobProperties, error) {
	var job JobProperties
	var headers map[string]string
	path := fmt.Sprintf("/jobs/%s/%s/details", jobType, jobName)
	jobraw, err := client.send_request("GET", path, "", headers)
	if err != nil {
		return job, err
	}
	json.Unmarshal([]byte(jobraw), &job)

	return job, nil
}

func retryJob(client *api_client, jobType, jobName string) (JobProperties, error) {
	var job JobProperties
	var headers map[string]string
	path := fmt.Sprintf("/jobs/%s/%s/retry", jobType, jobName)
	jobraw, err := client.send_request("POST", path, "", headers)
	if err != nil {
		return job, err
	}
	json.Unmarshal([]byte(jobraw), &job)

	return job, nil
}

// SliceFind takes a slice and looks for an element in it. If found it will
// return true otherwise false.
func SliceFind(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

type FriendlyErrorOutput struct {
	Name    string          `yaml:"name"`
	State   string          `yaml:"state"`
	TaskID  string          `yaml:"taskid"`
	Type    string          `yaml:"type"`
	Errors  []FriendlyError `yaml:"errors"`
	Retries []FriendlyError `yaml:"retries,omitempty"`
}

type FriendlyError struct {
	Workers []string `yaml:"workers"`
	Message string   `yaml:"message"`
}

func friendlyYAMLError(job JobProperties) error {

	var errorOutput FriendlyErrorOutput

	errorOutput.Name = job.Name
	errorOutput.State = job.State
	errorOutput.TaskID = job.TaskID
	errorOutput.Type = job.Type

	// build a user friendly error message
	var appendedWorkerOutputs []map[string]interface{}
	var appendedWorkerRetries []map[string]interface{}
	for _, workerOutput := range job.Output {
		var appended bool
		var appendedRetry bool
		// ignore the worker if it is successfull
		if workerOutput.State == "finished" {
			continue
		}

		successRetry := false
		// ignore the worker if it is successfull in retries
		for _, workerRetry := range workerOutput.Retries {
			if workerRetry.State == "finished" {
				successRetry = true
				break
			} else {
				for _, appendedWorkerRetry := range appendedWorkerRetries {
					if appendedWorkerRetry["logs"] == workerRetry.Logs {

						if !SliceFind(appendedWorkerRetry["workers"].([]string), workerRetry.Worker) {
							appendedWorkerRetry["workers"] = append(appendedWorkerRetry["workers"].([]string), workerRetry.Worker)
						}
						appendedRetry = true
					}
				}

				// next if the worker was appended
				if appendedRetry {
					continue
				}

				// add the worker and its output log if not already in appendedWorkerRetries
				retriesToAppend := make(map[string]interface{})
				retriesToAppend["workers"] = []string{workerRetry.Worker}

				retriesToAppend["logs"] = workerRetry.Logs
				appendedWorkerRetries = append(appendedWorkerRetries, retriesToAppend)
			}
		}

		if successRetry {
			continue
		}

		// append the worker if its output log is already in appendedWorkerOutputs
		for _, appendedWorkerOutput := range appendedWorkerOutputs {
			if appendedWorkerOutput["logs"] == workerOutput.Logs {
				appendedWorkerOutput["workers"] = append(appendedWorkerOutput["workers"].([]string), workerOutput.Worker)
				appended = true
			}
		}
		// next if the worker was appended
		if appended {
			continue
		}
		// add the worker and its output log if not already in appendedWorkerOutputs
		outputToAppend := make(map[string]interface{})
		outputToAppend["workers"] = []string{workerOutput.Worker}
		outputToAppend["logs"] = workerOutput.Logs
		appendedWorkerOutputs = append(appendedWorkerOutputs, outputToAppend)
	}

	// failed job outputs
	for _, appendedWorkerOutput := range appendedWorkerOutputs {
		var ferror FriendlyError
		ferror.Workers = appendedWorkerOutput["workers"].([]string)

		// remove useless first line containing command executed
		workerLog := friendlyOutputRegexp.ReplaceAllString(appendedWorkerOutput["logs"].(string), "$1")
		workerLog = strings.Replace(workerLog, " Error:", "\nError:", -1)
		workerLog = escapeCRNL.ReplaceAllString(workerLog, "\n")
		ferror.Message = workerLog
		errorOutput.Errors = append(errorOutput.Errors, ferror)
	}

	// failed retries job outputs
	for _, appendedWorkerRetry := range appendedWorkerRetries {
		var ferror FriendlyError
		ferror.Workers = appendedWorkerRetry["workers"].([]string)

		// remove useless first line containing command executed
		workerLog := friendlyOutputRegexp.ReplaceAllString(appendedWorkerRetry["logs"].(string), "$1")
		workerLog = strings.Replace(workerLog, " Error:", "\nError:", -1)
		workerLog = escapeCRNL.ReplaceAllString(workerLog, "\n")
		ferror.Message = workerLog
		errorOutput.Retries = append(errorOutput.Retries, ferror)
	}

	yamlStr, err := yaml.Marshal(errorOutput)
	if err != nil {
		return fmt.Errorf("Unable to yaml marshal: %v\n\nJob error is: %v", err, errorOutput)
	}
	return fmt.Errorf(". Task details:\n\n%s", yamlStr)

}

func buildSuccessOutput(job JobProperties) map[string]interface{} {
	return map[string]interface{}{
		"name":   job.Name,
		"state":  job.State,
		"taskid": job.TaskID,
		"type":   job.Type,
	}
}
