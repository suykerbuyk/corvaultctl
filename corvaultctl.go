package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//	"net/http/httputil"
	"strings"
	//	"os"
	"time"
)

type CorvaultCtx struct {
	Host string `json:"ConnectionHost"`
	User string `json:"ConnectionUser"`
	Pass string `json:"ConnectionPass"`
	Key  string `json:"Key"`
}

func OpenSession(tgtCtx *CorvaultCtx) (client *http.Client, err error) {
	auth_string := base64.StdEncoding.EncodeToString([]byte(tgtCtx.User + ":" + tgtCtx.Pass))
	fmt.Println("Base64 auth_string = " + auth_string)
	url := tgtCtx.Host + "/api/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+auth_string)
	req.Header.Add("dataType", "json")
	//dump, err := httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("response Status: ", resp.Status)
	fmt.Println("response Header: ", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	authStatus := new(AuthStatusList)
	err = json.Unmarshal(body, &authStatus)
	if err != nil {
		return nil, err
	}
	//decodeStr, err := json.MarshalIndent(authStatus, "", "  ")
	//fmt.Println(string(decodeStr))
	if 1 > len(authStatus.List) {
		err = errors.New("Error, no status report present")
		return nil, err
	}
	for _, auth := range authStatus.List {
		fmt.Println("                 meta:", auth.Meta)
		fmt.Println("          object-name:", auth.ObjectName)
		fmt.Println("             response:", auth.Response)
		fmt.Println("        response-type:", auth.ResponseType)
		fmt.Println("response-type-numeric:", auth.ResponseTypeNumeric)
		fmt.Println("          return-code:", auth.ReturnCode)
		fmt.Println("         component-id:", auth.ComponentId)
		fmt.Println("           time-stamp:", auth.TimeStamp)
		fmt.Println("   time-stamp-numeric:", auth.TimeStampNumeric)
	}
	if authStatus.List[0].ReturnCode != 1 {
		log.Fatal("API return code was not \"1\" : ", authStatus.List[0].ReturnCode)
	}
	tgtCtx.Key = authStatus.List[0].Response
	fmt.Printf("sessionKey=%s\n", tgtCtx.Key)
	return

}

func main() {
	var SessionCtx = CorvaultCtx{
		Host: "https://corvault-2a",
		User: "manage",
		Pass: "Testit123!",
		Key:  "",
	}

	client, err := OpenSession(&SessionCtx)
	if err != nil {
		err = fmt.Errorf("OpenSession Failed!: %v", err.Error())
		log.Fatal(err)
	}

	url := SessionCtx.Host + "/api/show/certificate"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", SessionCtx.Key)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status: ", resp.Status)
	fmt.Println("response Header: ", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	certStatus := new(CertificateStatusList)
	err = json.Unmarshal(body, &certStatus)
	if err != nil {
		log.Fatal(err)
	}
	if 1 > len(certStatus.List) {
		log.Fatal("Error, no certificate report present")
	}
	for _, cert := range certStatus.List {
		fmt.Println("               object-name: ", cert.ObjectName)
		fmt.Println("                      meta: ", cert.Meta)
		fmt.Println("                controller: ", cert.Controller)
		fmt.Println("        controller-numeric: ", cert.ControllerNumeric)
		fmt.Println("        certificate-status: ", cert.CertificateStatus)
		fmt.Println("certificant-status-numeric: ", cert.CertificateStatusNumeric)
		fmt.Println("          certificate-time: ", cert.CertificateTime)
		fmt.Println("     certificate-signature: ", cert.CertificateSignature)
		//fmt.Printf("          certificate-text: %s\n", certStatus.List[idx].CertificateText)
		cert.CertificateTextList = strings.Split(cert.CertificateText, "\\n")
		for _, v := range cert.CertificateTextList {
			fmt.Println(v)
		}
	}
}
