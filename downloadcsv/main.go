package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/docker/distribution/uuid"
)

var client http.Client

func main() {
	fmt.Println(uuid.Generate().String())
	/*
			query
			includeDeleted=false&filter=%7B%7D
			header
			GET /server-webservices/reportJob?includeDeleted=false&filter=%7B%7D HTTP/1.1
		Host: app.cantab.com
		Connection: keep-alive
		Accept: application/json
		Sec-Fetch-Dest: empty
		Authorization: AuthToken vAQUh6zZT0aQwqQywtjbzDgyhDw8-aDl_Rp6UuWK
		User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36
		Sec-Fetch-Site: same-origin
		Sec-Fetch-Mode: cors
		Referer: https://app.cantab.com/admin/index.html
		Accept-Encoding: gzip, deflate, br
		Accept-Language: zh-CN,zh;q=0.9
	*/
	tokenstr, err := gettocken()
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := GetReportDef(tokenstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result {
		if v["name"] == "Row per-measure – Normative data" {
			jobid, err := ReportJob(tokenstr, v)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("jobid: ", jobid)
			for !TryGetReportJob(tokenstr, jobid) {
				fmt.Println("file not exist .....")
				time.Sleep(2 * time.Second)
			}
			fmt.Println("csv can download .....")
			geterr := GetReportJob(tokenstr, jobid)
			if geterr != nil {
				fmt.Println(geterr)
			}
			return
		}
	}
	// fmt.Println(result)
}

/*

   {
       "clientId": "c-a751fed1-3860-4291-8e02-64cccd57ac45",
       "description": "A row per measure report with results from a single site",
       "details": {
           "allowXLSXLocking": true,
           "clientId": null,
           "fileName": "RowByMeasure_Site.zip",
           "type": "CLOVER_PROJECT_ZIP",
           "id": "57dbbe410cf24bd9ea65ce87"
       },
       "disabled": null,
       "level": "SITE",
       "name": "Row per-measure - Site",
       "orgTypeRestrictions": [
           "PHARMACEUTICAL"
       ],
       "organisation": null,
       "ownerOrganisation": null,
       "site": null,
       "study": null,
       "type": "FILE",
       "id": "57dbbe410cf24bd9ea65ce88",
       "version": 2
   }
*/
type StudyClientDetails struct {
	allowXLSXLocking bool
	clientId         string
	fileName         string
	Type             string
	id               string
}
type StudyClient struct {
	clientId            string
	description         string
	details             map[string]string
	disabled            string
	level               string
	name                string
	orgTypeRestrictions []string
	organisation        string
	ownerOrganisation   string
	site                string
	study               string
	id                  string
	version             int
}
type ReportDef struct {
	records []StudyClient
	total   int
	success bool
}

func GetReportDef(tok string) ([]map[string]interface{}, error) {
	fmt.Println(tok)
	req, err := http.NewRequest("GET", "https://app.cantab.com/server-webservices/reportDef", nil)
	req.URL.Query().Set("includeDeleted", "false")
	req.URL.Query().Set("filter", "{}")
	req.Header.Set("Authorization", fmt.Sprintf("AuthToken %s", tok))
	req.Header.Set("Host", "app.cantab.com")
	req.Header.Set("Origin", "https://app.cantab.com")
	req.Header.Set("Referer", "https://app.cantab.com/admin/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	var v interface{}
	json.Unmarshal(body, &v)
	data := (v.(map[string]interface{}))
	// fmt.Println(data["records"])
	studyclient := data["records"].([]interface{})
	result := make([]map[string]interface{}, len(studyclient))
	for index, v := range studyclient {
		tmp := v.(map[string]interface{})
		result[index] = tmp
	}
	return result, nil
}

//https://app.cantab.com/server-webservices/reportJob
/*
id: "c-c5bc90f7-2e03-49cb-898a-c64e0ea2b0b9"
clientId: null
timeSubmitted: null
status: null
reportUrl: null
reportDef: "5c642a090cf2b22c8bd8c236"
site: null
study: "5d5446760cf2ca8a450652c5"
subject: null
visit: null
organisation: "5cefe4f70cf21aa4f09ca39c"
*/
func ReportJob(tok string, data map[string]interface{}) (string, error) {
	jsonstr := make(map[string]interface{})
	jsonstr["id"] = fmt.Sprintf("c-%s", uuid.Generate().String())
	jsonstr["clientId"] = nil
	jsonstr["timeSubmitted"] = nil
	jsonstr["status"] = nil
	jsonstr["reportUrl"] = nil
	jsonstr["reportDef"] = data["id"].(string)
	jsonstr["site"] = nil
	jsonstr["study"] = "5d5446760cf2ca8a450652c5"
	jsonstr["subject"] = nil
	jsonstr["visit"] = nil
	jsonstr["organisation"] = "5cefe4f70cf21aa4f09ca39c"
	b, err := json.Marshal(jsonstr)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://app.cantab.com/server-webservices/reportJob", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("AuthToken %s", tok))
	req.Header.Set("Host", "app.cantab.com")
	req.Header.Set("Origin", "https://app.cantab.com")
	req.Header.Set("Referer", "https://app.cantab.com/admin/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	if m["success"].(bool) != true {
		return "", errors.New("ReportJob fail")
	}
	records := m["records"].([]interface{})
	if len(records) > 0 {
		record := records[0].(map[string]interface{})
		return record["id"].(string), nil
	}
	return "", errors.New("ReportJob fail")
	// fmt.Println(string(body))
}
func TryGetReportJob(tok string, jobid string) bool {
	//https://app.cantab.com/server-webservices/reportJob/download/5e71bd150cf25f5b8476ade9
	req, err := http.NewRequest("POST", fmt.Sprintf("https://app.cantab.com/server-webservices/reportJob/download/%s", jobid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("AuthToken %s", tok))
	req.Header.Set("Host", "app.cantab.com")
	req.Header.Set("Origin", "https://app.cantab.com")
	req.Header.Set("Referer", "https://app.cantab.com/admin/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.Header.Get("Content-Type") == "application/octet-stream" {
		return true
	}
	return false
}

func GetReportJob(tok string, jobid string) error {
	//https://app.cantab.com/server-webservices/reportJob/download/5e71bd150cf25f5b8476ade9
	req, err := http.NewRequest("POST", fmt.Sprintf("https://app.cantab.com/server-webservices/reportJob/download/%s", jobid), nil)
	req.Header.Set("Authorization", fmt.Sprintf("AuthToken %s", tok))
	req.Header.Set("Host", "app.cantab.com")
	req.Header.Set("Origin", "https://app.cantab.com")
	req.Header.Set("Referer", "https://app.cantab.com/admin/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	csvfile, err := os.Create(fmt.Sprintf("%s.csv", jobid))
	if err != nil {
		return err
	}
	io.Copy(csvfile, resp.Body)
	return nil
}
func gettocken() (string, error) {
	req, err := http.NewRequest("POST", "https://app.cantab.com/server-webservices/auth/token", nil)
	req.Header.Set("Authorization", "Basic NzA3NjE3NjNAcXEuY29tOkRpbmd5b25nbGkxMjM=")
	req.Header.Set("Host", "app.cantab.com")
	req.Header.Set("Origin", "https://app.cantab.com")
	req.Header.Set("Referer", "https://app.cantab.com/admin/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]string)

	json.Unmarshal(body, &m)

	return m["token"], nil
}
