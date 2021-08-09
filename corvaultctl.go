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
	"net/http/httputil"
	"strings"
	//	"os"
	"time"
)

type CorvaultCreds struct {
	Host string `json:"ConnectionHost"`
	User string `json:"ConnectionUser"`
	Pass string `json:"ConnectionPass"`
	Key  string `json:"Key"`
}
type CorvaultCtx struct {
	Credentials CorvaultCreds
	Client      http.Client
}

func OpenSession(tgtCreds *CorvaultCreds) (err error, client *http.Client) {
	auth_string := base64.StdEncoding.EncodeToString([]byte(tgtCreds.User + ":" + tgtCreds.Pass))
	fmt.Println("Base64 auth_string = " + auth_string)
	url := tgtCreds.Host + "/api/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err, nil
	}
	req.Header.Add("Authorization", "Basic "+auth_string)
	req.Header.Add("dataType", "json")
	dump, err := httputil.DumpRequestOut(req, false)
	fmt.Printf("%s", dump)
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	fmt.Println("response Status: ", resp.Status)
	fmt.Println("response Header: ", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	var status CvtResponseStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatal("OpenSession Failed: ", err)
	}
	if status.Status[0].ReturnCode != 1 {
		log.Fatal("OpenSession : API return code was not \"1\" : ", status.Status[0].ReturnCode)
	}
	tgtCreds.Key = status.Status[0].Response
	fmt.Printf("sessionKey=%s\n", tgtCreds.Key)
	fmt.Println(status)
	return

}
func FetchCertificates(tgtCreds *CorvaultCreds, client *http.Client) (certs *CvtCertificates, err error) {
	url := tgtCreds.Host + "/api/show/certificate"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", tgtCreds.Key)
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
	fmt.Println(string(body))
	certs = new(CvtCertificates)
	err = json.Unmarshal(body, &certs)
	if err != nil {
		log.Fatal(err)
	}
	if 1 > len(certs.Certificate) {
		log.Fatal("Error, no certificate report present")
	}
	fmt.Println("Dumping the request:")
	dump, err := httputil.DumpRequestOut(req, false)
	fmt.Printf("%s", dump)
	fmt.Println("Dumping the response:")
	dump, err = httputil.DumpResponse(resp, true)
	fmt.Println("Dump Complete")
	return
}
func DumpCertificates(certs *CvtCertificates) {
	for _, cert := range certs.Certificate {
		fmt.Println("               object-name: ", cert.ObjectName)
		fmt.Println("                      meta: ", cert.Meta)
		fmt.Println("                controller: ", cert.Controller)
		fmt.Println("        controller-numeric: ", cert.ControllerNumeric)
		fmt.Println("        certificate-status: ", cert.CertificateStatus)
		fmt.Println("certificant-status-numeric: ", cert.CertificateStatusNumeric)
		fmt.Println("          certificate-time: ", cert.CertificateTime)
		fmt.Println("     certificate-signature: ", cert.CertificateSignature)
		//	fmt.Printf("          certificate-text: %s\n", cert.CertificateText)
		CertificateTextList := strings.Split(cert.CertificateText, "\\n")
		for _, v := range CertificateTextList {
			fmt.Println(v)
		}

	}
	//fmt.Println("GetCertificates Status: ", certs.Response)
}
func main() {
	var sessionCreds = CorvaultCreds{
		Host: "https://corvault-2a",
		User: "manage",
		Pass: "Testit123!",
		Key:  "",
	}

	err, client := OpenSession(&sessionCreds)
	if err != nil {
		err = fmt.Errorf("OpenSession Failed!: %v", err.Error())
		log.Fatal(err)
	}
	//DumpAuthStatusList(authStatus)
	certStatus, err := FetchCertificates(&sessionCreds, client)
	if err != nil {
		err = fmt.Errorf("FetchCeritificate Failed!: %v", err.Error())
		log.Fatal(err)
	}
	DumpCertificates(certStatus)
}
