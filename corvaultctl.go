package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	//	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	//"strings"
	//	"os"
	"time"
)

type CorvaultCtx struct {
	Credential CorvaultCredential
	Client     *http.Client
}

func OpenSession(ctx *CorvaultCtx) (err error) {
	auth_string := base64.StdEncoding.EncodeToString([]byte(ctx.Credential.User + ":" + ctx.Credential.Pass))
	//fmt.Println("Base64 auth_string = " + auth_string)
	url := ctx.Credential.Host + "api/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	ctx.Client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("OpenSession, http.NewRequest failed: %w", err)
	}
	req.Header.Add("Authorization", "Basic "+auth_string)
	req.Header.Add("dataType", "json")
	//dump, err := httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	resp, err := ctx.Client.Do(req)
	if err != nil {
		return fmt.Errorf("OpenSession, Client.Do failed: %w", err)
	}
	defer resp.Body.Close()
	//fmt.Println("response Header: ", resp.Header)
	//fmt.Println("response Status: ", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//dump, _ = httputil.DumpResponse(resp, true)
		//fmt.Printf("%s", dump)
		return fmt.Errorf("OpenSession, read http body failed: %w", err)
	}
	var status CvtResponseStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return fmt.Errorf("OpenSession, json.Unmarshal failed: %w", err)
	}
	if status.Status[0].ReturnCode != 1 {
		return fmt.Errorf("OpenSession : API return code was not \"1\" %d : ", status.Status[0].ReturnCode)
	}
	ctx.Credential.Key = status.Status[0].Response
	//fmt.Printf("sessionKey=%s\n", ctx.Credential.Key)
	return

}
func (ctx CorvaultCtx) Show(aspect string) (buffer []byte, err error) {
	url := ctx.Credential.Host + "api/show/" + aspect
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - http.NewRequest failed for %s: %w", url, err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", ctx.Credential.Key)
	resp, err := ctx.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - ctx.Client.Do failed for %s: %w", url, err)
	}
	defer resp.Body.Close()
	//fmt.Println("Dumping the request:")
	//dump, err := httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	//fmt.Println("Dumping the response:")
	//dump, err := httputil.DumpResponse(resp, true)
	//fmt.Printf("%s", dump)
	//fmt.Println("Dump Complete")
	//fmt.Println("response Status: ", resp.Status)
	//fmt.Println("response Header: ", resp.Header)
	buffer, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(string(buffer))
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - ioutil.Readall for http body failed for %s: %w", url, err)
	}
	return
}
func (ctx CorvaultCtx) GetCertificate() (certs *CvtCertificates, err error) {
	buffer, err := ctx.Show("certificate")
	if err != nil {
		return nil, fmt.Errorf("GetCertificate - CorvaultCtx.Show certificate failed: %w", err)
	}
	certs = new(CvtCertificates)
	err = json.Unmarshal(buffer, &certs)
	if err != nil {
		log.Fatal(err)
	}
	if 1 > len(certs.Certificate) {
		return nil, fmt.Errorf("GetCertificate - No Certificates present!")
	}
	return
}
func (ctx CorvaultCtx) GetDiskGroups() (dgs *CvtDiskGroups, err error) {
	buffer, err := ctx.Show("disk-groups")
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show disk-groups failed: %w", err)
	}
	dgs = new(CvtDiskGroups)
	fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &dgs)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json data: %w", err)
	}
	if 1 > len(dgs.DiskGroups) {
		return nil, fmt.Errorf("No Disk Groups Present!")
	}
	return
}
func (ctx CorvaultCtx) GetDiskGroupStatistics() (data *CvtDiskGroupStatistics, err error) {
	buffer, err := ctx.Show("disk-group-statistics")
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show disk-group-statistics failed: %w", err)
	}
	data = new(CvtDiskGroupStatistics)
	fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json data: %w", err)
	}
	if 1 > len(data.Statistics) {
		return nil, fmt.Errorf("No Disk Group Statistics Present!")
	}
	return
}
func main() {
	cfg, err := GetCvtConfig()
	if err != nil {
		err = fmt.Errorf("GetCvtConfig Failed!: %v", err.Error())
		log.Fatal(err)
	}
	ctx := CorvaultCtx{}
	ctx.Credential = cfg.Targets["corvault-1a"]

	err = OpenSession(&ctx)
	if err != nil {
		err = fmt.Errorf("OpenSession Failed!: %v", err.Error())
		log.Fatal(err)
	}
	//certStatus, err := ctx.GetCertificate()
	//if err != nil {
	//	err = fmt.Errorf("FetchCeritificate Failed!: %v", err.Error())
	//	log.Fatal(err)
	//}
	//fmt.Println(certStatus.Text())
	//diskGroups, err := ctx.GetDiskGroups()
	//if err != nil {
	//	err = fmt.Errorf("GetDiskGroups Failed: %v", err.Error())
	//	log.Fatal(err)
	//}
	//fmt.Println(diskGroups.Json())
	diskGroupStatistics, err := ctx.GetDiskGroupStatistics()
	if err != nil {
		err = fmt.Errorf("GetDiskGroupStatistics Failed: %v", err.Error())
		log.Fatal(err)
	}
	fmt.Println(diskGroupStatistics.Json())
}
