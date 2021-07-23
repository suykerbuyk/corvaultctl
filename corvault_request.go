package main

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
