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
	fmt.Printf("sessionKey=%s\n", ctx.Credential.Key)
	return

}
func FetchCertificates(ctx *CorvaultCtx) (certs *CvtCertificates, err error) {
	url := ctx.Credential.Host + "api/show/certificate"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("FetchCertificates - http.NewRequest failed: %w", err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", ctx.Credential.Key)
	resp, err := ctx.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FetchCertificates - Client.Do Failed: %w", err)
	}
	defer resp.Body.Close()
	//fmt.Println("response Status: ", resp.Status)
	//fmt.Println("response Header: ", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("FetchCertificates - ioutil.ReadAll from body failed: %w", err)
	}
	//fmt.Println(string(body))
	certs = new(CvtCertificates)
	err = json.Unmarshal(body, &certs)
	if err != nil {
		log.Fatal(err)
	}
	if 1 > len(certs.Certificate) {
		return nil, fmt.Errorf("FetchCertificates - No Certificates present!")
	}
	//fmt.Println("Dumping the request:")
	//dump, err := httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	//fmt.Println("Dumping the response:")
	//dump, err = httputil.DumpResponse(resp, true)
	//fmt.Println("Dump Complete")
	return
}
func (ctx *CorvaultCtx) Show(aspect string) (buffer []byte, err error) {
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
	//fmt.Println("response Status: ", resp.Status)
	//fmt.Println("response Header: ", resp.Header)
	buffer, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - ioutil.Readall for http body failed for %s: %w", url, err)
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
	certStatus, err := FetchCertificates(&ctx)
	if err != nil {
		err = fmt.Errorf("FetchCeritificate Failed!: %v", err.Error())
		log.Fatal(err)
	}
	fmt.Println(certStatus)
}
