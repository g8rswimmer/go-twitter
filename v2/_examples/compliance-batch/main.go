package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/**
	In order to run, the user will need to provide the bearer token and the list of tweet ids.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	jobType := flag.String("type", "", "job type")
	upload := flag.String("upload", "", "upload file")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Compliance Batch Job Example")

	opts := twitter.CreateComplianceBatchJobOpts{
		Name: "go twitter job example",
	}

	// 1. Create a compliance batch job
	fmt.Println("1. Create a compliance batch job")
	complianceResponse, err := client.CreateComplianceBatchJob(context.Background(), twitter.ComplianceBatchJobType(*jobType), opts)
	if err != nil {
		log.Panicf("create compliance job error: %v", err)
	}

	enc, err := json.MarshalIndent(complianceResponse, "", "    ")
	if err != nil {
		log.Panicf("create compliance job error: %v", err)
	}
	fmt.Println(string(enc))

	job := complianceResponse.Raw.Job

	// 2. Upload ids from file
	fmt.Println("2. Upload ids from file")
	f, err := os.Open(*upload)
	if err != nil {
		log.Panicf("open upload file error: %v", err)
	}
	defer f.Close()

	err = job.Upload(context.Background(), f)
	if err != nil {
		log.Panicf("upload ids error: %v", err)
	}

	// 3. Check the job status
	fmt.Println("3. Check the job status")
	for {
		time.Sleep(time.Second)

		resp, err := client.ComplianceBatchJob(context.Background(), job.ID)
		if err != nil {
			log.Panicf("check status error: %v", err)
		}
		jobStatus := resp.Raw.Job
		fmt.Println("Status: " + jobStatus.Status)
		if jobStatus.Status != twitter.ComplianceBatchJobStatusInProgress {
			break
		}
	}

	// 4. Download results
	fmt.Println("4. Download results")
	downloadResponse, err := job.Download(context.Background())
	if err != nil {
		log.Panicf("download results error: %v", err)
	}

	enc, err = json.MarshalIndent(downloadResponse, "", "    ")
	if err != nil {
		log.Panicf("download results error: %v", err)
	}
	fmt.Println(string(enc))

}
