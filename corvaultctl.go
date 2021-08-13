package main

import (
	"crypto/tls"
	"encoding/json"
	//	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	"github.com/alecthomas/kong"
	"strings"
	//	"os"
	"time"
)

type CorvaultCtx struct {
	Credential CorvaultCredential
	Client     *http.Client
}

func OpenSession(tgtCtx *CorvaultCtx) (err error) {
	url := tgtCtx.Credential.Host + "api/login"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	tgtCtx.Client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("OpenSession, http.NewRequest failed: %w", err)
	}
	req.Header.Add("Authorization", "Basic "+tgtCtx.Credential.Auth)
	req.Header.Add("dataType", "json")
	//dump, err := httputil.DumpRequestOut(req, false)
	//fmt.Printf("%s", dump)
	resp, err := tgtCtx.Client.Do(req)
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
	tgtCtx.Credential.Key = status.Status[0].Response
	//fmt.Printf("sessionKey=%s\n", tgtCtx.Credential.Key)
	return

}
func (tgtCtx CorvaultCtx) Show(aspect string) (buffer []byte, err error) {
	url := tgtCtx.Credential.Host + "api/show/" + aspect
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - http.NewRequest failed for %s: %w", url, err)
	}
	req.Header.Add("dataType", "json")
	req.Header.Add("sessionKey", tgtCtx.Credential.Key)
	resp, err := tgtCtx.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show - tgtCtx.Client.Do failed for %s: %w", url, err)
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
func (tgtCtx CorvaultCtx) GetCertificate() (certs *CvtCertificates, err error) {
	buffer, err := tgtCtx.Show("certificate")
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
func (tgtCtx CorvaultCtx) GetDiskGroups() (dgs *CvtDiskGroups, err error) {
	buffer, err := tgtCtx.Show("disk-groups")
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show disk-groups failed: %w", err)
	}
	dgs = new(CvtDiskGroups)
	//fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &dgs)
	if err != nil {
		return nil, fmt.Errorf("GetDiskGroups failed to unmarshal json data: %w", err)
	}
	if 1 > len(dgs.DiskGroups) {
		return nil, fmt.Errorf("No Disk Groups Present!")
	}
	return
}
func (tgtCtx CorvaultCtx) GetDiskGroupStatistics() (data *CvtDiskGroupStatistics, err error) {
	buffer, err := tgtCtx.Show("disk-group-statistics")
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show disk-group-statistics failed: %w", err)
	}
	data = new(CvtDiskGroupStatistics)
	//fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return nil, fmt.Errorf("GetDiskGroupStatistics failed to unmarshal json data: %w", err)
	}
	if 1 > len(data.Statistics) {
		return nil, fmt.Errorf("No Disk Group Statistics Present!")
	}
	return
}
func (tgtCtx CorvaultCtx) GetSystem() (data *CvtSystem, err error) {
	buffer, err := tgtCtx.Show("system")
	if err != nil {
		return nil, fmt.Errorf("CorvaultCtx.Show disk-group-statistics failed: %w", err)
	}
	data = new(CvtSystem)
	fmt.Println(string(buffer))
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		return nil, fmt.Errorf("GetSystem failed to unmarshal json data: %w", err)
	}
	return
}
func main() {
	cli := CLI{
		CliGlobals: CliGlobals{
			Version: VersionFlag("0.1.1"),
		},
	}

	cliCtx := kong.Parse(&cli,
		kong.Name("corvaultctl"),
		kong.Description("Simple CLI interface into Seagate Enclosures to facilitate DevOps."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			//Tree:    true,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	err := cliCtx.Run(&cli.CliGlobals)
	fmt.Println("Command: ", cliCtx.Command())
	fmt.Println("Command Args: ", cliCtx.Args)
	cliCtx.FatalIfErrorf(err)
	cmdStrings := strings.Split(cliCtx.Command(), " ")
	switch cmdStrings[0] {
	case "target":
		fmt.Println("Target command found")
	case "raw":
		fmt.Println("Raw command found", cli.Raw.Cmd)
	case "show":
		switch cmdStrings[1] {
		case "disks":
			fmt.Println("Show Disks command found")
		case "advanced-settings":
			fmt.Println("Show AdvancedSettings command found")
		case "alert-condition-history":
			fmt.Println("Show AlertConditionHistory command found")
		case "certificates":
			fmt.Println("Show Certificates command found")
		case "volumes":
			fmt.Println("Show Volumes command found")
		}
	}
	cfg, err := GetCvtConfig()
	if err != nil {
		err = fmt.Errorf("GetCvtConfig Failed!: %v", err.Error())
		log.Fatal(err)
	}
	tgtCtx := CorvaultCtx{}
	tgtCtx.Credential = cfg.Targets["corvault-2a"]

	err = OpenSession(&tgtCtx)
	if err != nil {
		err = fmt.Errorf("OpenSession Failed!: %v", err.Error())
		log.Fatal(err)
	}
	//certStatus, err := tgtCtx.GetCertificate()
	//if err != nil {
	//	err = fmt.Errorf("FetchCeritificate Failed!: %v", err.Error())
	//	log.Fatal(err)
	//}
	//fmt.Println(certStatus.Text())
	//diskGroups, err := tgtCtx.GetDiskGroups()
	//if err != nil {
	//	err = fmt.Errorf("GetDiskGroups Failed: %v", err.Error())
	//	log.Fatal(err)
	//}
	//fmt.Println(diskGroups.Json())
	//diskGroupStatistics, err := tgtCtx.GetDiskGroupStatistics()
	//if err != nil {
	//	err = fmt.Errorf("GetDiskGroupStatistics Failed: %v", err.Error())
	//	log.Fatal(err)
	//}
	//fmt.Println(diskGroupStatistics.Json())
	system, err := tgtCtx.GetSystem()
	if err != nil {
		err = fmt.Errorf("GetSystem Failed: %v", err.Error())
		log.Fatal(err)
	}
	fmt.Println(system.Json())
}
