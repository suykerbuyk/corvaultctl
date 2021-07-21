package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	//	"os"
	"time"
)

const (
	CONN_HOST string = "https://corvault-2a"
	CONN_PORT        = "443"
)

type AuthStatus struct {
	ObjectName          string `json:"object-name"`
	Meta                string `json:"meta"`
	ResponseType        string `json:"response-type"`
	ResponseTypeNumeric int32  `json:"response-type-numeric"`
	Response            string `json:"response"`
	ReturnCode          int32  `json:"return-code"`
	ComponentId         string `json:"component-id"`
	TimeStamp           string `json:"time-stamp"`
	TimeStampNumeric    int64  `json:"time-stamp-numeric"`
}

type AuthStatusList struct {
	List []AuthStatus `json:"status"`
}

type CertificateStatus struct {
	ObjectName               string   `json:"object-name"`
	Meta                     string   `json:"meta"`
	Controller               string   `json:"controller"`
	ControllerNumeric        int64    `json:"controller-numeric"`
	CertificateStatus        string   `json:"certificate-status"`
	CertificateStatusNumeric int64    `json:"certificate-status-numeric"`
	CertificateTime          string   `json:"certificate-time"`
	CertificateSignature     string   `json:"certificate-signature"`
	CertificateText          string   `json:"certificate-text"`
	CertificateTextList      []string `json:"certificate-text-list,omitempty"`
}
type CertificateStatusList struct {
	List []CertificateStatus `json:"certificate-status"`
}

var user_name = "manage"
var user_pass = "Testit123!"

func main() {
	// Create the variables for the response and error
	//	var r *http.Response
	var err error
	auth_string := base64.StdEncoding.EncodeToString([]byte(user_name + ":" + user_pass))
	fmt.Println("Base64 auth_string = " + auth_string)
	url := CONN_HOST + "/api/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: time.Second * 5, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic "+auth_string)
	req.Header.Add("dataType", "json")
	dump, err := httputil.DumpRequestOut(req, false)
	fmt.Printf("%s", dump)
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
	//fmt.Printf("\nResponse Body:\n%s\n", body)

	authStatus := new(AuthStatusList)
	err = json.Unmarshal(body, &authStatus)
	if err != nil {
		log.Fatal(err)
	}
	//decodeStr, err := json.MarshalIndent(authStatus, "", "  ")
	//fmt.Println(string(decodeStr))
	statusCnt := len(authStatus.List)
	if statusCnt < 1 {
		log.Fatal("Error, no status report present")
	}
	fmt.Println("Status Report Count = ", statusCnt)
	for idx := 0; idx < statusCnt; idx++ {
		fmt.Println("                 meta: ", authStatus.List[idx].Meta)
		fmt.Println("          object-name: ", authStatus.List[idx].ObjectName)
		fmt.Println("             response: ", authStatus.List[idx].Response)
		fmt.Println("        response-type: ", authStatus.List[idx].ResponseType)
		fmt.Println("response-type-numeric: ", authStatus.List[idx].ResponseTypeNumeric)
		fmt.Println("          return-code: ", authStatus.List[idx].ReturnCode)
		fmt.Println("         component-id: ", authStatus.List[idx].ComponentId)
		fmt.Println("           time-stamp: ", authStatus.List[idx].TimeStamp)
		fmt.Println("   time-stamp-numeric: ", authStatus.List[idx].TimeStampNumeric)
	}
	if authStatus.List[0].ReturnCode != 1 {
		log.Fatal("API return code was not \"1\" : ", authStatus.List[0].ReturnCode)
	}

	sessionKey := authStatus.List[0].Response
	fmt.Printf("sessionKey=%s\n", sessionKey)

	url = CONN_HOST + "/api/show/certificate"
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", sessionKey)
	//dump, err = httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status: ", resp.Status)
	fmt.Println("response Header: ", resp.Header)
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("\nResponse Body:\n%s\n", body)
	certStatus := new(CertificateStatusList)
	err = json.Unmarshal(body, &certStatus)
	if err != nil {
		log.Fatal(err)
	}
	certCnt := len(certStatus.List)
	if certCnt < 1 {
		log.Fatal("Error, no certificate report present")
	}
	fmt.Println("Certificate Report Count = ", certCnt)
	for idx := 0; idx < certCnt; idx++ {
		fmt.Println("               object-name: ", certStatus.List[idx].ObjectName)
		fmt.Println("                      meta: ", certStatus.List[idx].Meta)
		fmt.Println("                controller: ", certStatus.List[idx].Controller)
		fmt.Println("        controller-numeric: ", certStatus.List[idx].ControllerNumeric)
		fmt.Println("        certificate-status: ", certStatus.List[idx].CertificateStatus)
		fmt.Println("certificant-status-numeric: ", certStatus.List[idx].CertificateStatusNumeric)
		fmt.Println("          certificate-time: ", certStatus.List[idx].CertificateTime)
		fmt.Println("     certificate-signature: ", certStatus.List[idx].CertificateSignature)
		//fmt.Printf("          certificate-text: %s\n", certStatus.List[idx].CertificateText)
		certStatus.List[idx].CertificateTextList = strings.Split(certStatus.List[idx].CertificateText, "\\n")
		for _, v := range certStatus.List[idx].CertificateTextList {
			fmt.Println(v)
		}
	}
}
