package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type CvtResponseStatus struct {
	Status []struct {
		ObjectName          string `json:"object-name"`
		Meta                string `json:"meta"`
		ResponseType        string `json:"response-type"`
		ResponseTypeNumeric int    `json:"response-type-numeric"`
		Response            string `json:"response"`
		ReturnCode          int    `json:"return-code"`
		ComponentID         string `json:"component-id"`
		TimeStamp           string `json:"time-stamp"`
		TimeStampNumeric    int    `json:"time-stamp-numeric"`
	} `json:"status"`
}

func (s CvtResponseStatus) String() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("ResponseStatus to JSON string error: " + err.Error()))
	}
	return string(prettyJSON)
}

type CvtCertificates struct {
	Certificate []struct {
		ObjectName               string `json:"object-name"`
		Meta                     string `json:"meta"`
		Controller               string `json:"controller"`
		ControllerNumeric        int    `json:"controller-numeric"`
		CertificateStatus        string `json:"certificate-status"`
		CertificateStatusNumeric int    `json:"certificate-status-numeric"`
		CertificateTime          string `json:"certificate-time"`
		CertificateSignature     string `json:"certificate-signature"`
		CertificateText          string `json:"certificate-text"`
	} `json:"certificate-status"`
	Response CvtResponseStatus
}

type CvtDiskGroup struct {
	ObjectName                      string `json:"object-name"`
	Meta                            string `json:"meta"`
	Name                            string `json:"name"`
	Url                             string `json:"url"`
	BlockSize                       string `json:"blocksize"`
	Size                            string `json:"size"`
	SizeNumeric                     int64  `json:"size-numeric"`
	FreeSpace                       string `json:"freespace"`
	FreeSpaceNumeric                int64  `json:"freespace-numeric"`
	RawSize                         string `json:"raw-size"`
	RawSizeNumeric                  int64  `json:"raw-size-numeric"`
	Overhead                        string `json:"overhead"`
	OverheadNumeric                 int64  `json:"overhead-numeric"`
	StorageType                     string `json:"storage-type"`
	StorageTypeNumeric              int64  `json:"storage-type-numeric"`
	Pool                            string `json:"pool"`
	PoolsUrl                        string `json:"pools-url"`
	PoolSerial                      string `json:"pool-serial-number"`
	StorageTier                     string `json:"storage-tier"`
	StorageTierNumeric              int64  `json:"storage-tier-numeric"`
	TotalPages                      string `json:"total-pages"`
	AllocatedPages                  string `json:"allocated-pages"`
	AvailablePages                  string `json:"available-pages"`
	PoolPercentage                  string `json:"pool-percentage"`
	PerformanceRank                 string `json:"performance-rank"`
	Owner                           string `json:"owner"`
	OwnerNumeric                    int64  `json:"owner-numeric"`
	PreferredOwner                  string `json:"preferred-owner"`
	PreferredOwnerNumeric           int64  `json:"preferred-owner-numeric"`
	RaidType                        string `json:"raidtype"`
	RaidTypeNumeric                 uint32 `json:"raidtype-numeric"`
	DiskCount                       uint16 `json:"diskcount"`
	InterleavedVolumeCount          uint16 `json:"interleaved-volume-count"`
	Spear                           string `json:"spear"`
	SpearNumeric                    uint32 `json:"spear-numeric"`
	SpareCount                      uint16 `json:"sparecount"`
	ChunkSize                       string `json:"chunksize"`
	Status                          string `json:"status"`
	StatusNumeric                   uint32 `json:"status-numeric"`
	Lun                             string `json:"lun"`
	MinDriveSize                    string `json:"min-drive-size"`
	MinDriveSizeNumeric             uint64 `json:"min-drive-size-numeric"`
	CreateDate                      string `json:"create-date"`
	CreateDateNumeric               uint32 `json:"create-date-numeric"`
	CacheReadAhead                  string `json:"cache-read-ahead"`
	CacheReadAheadNumeric           uint64 `json:"cache-read-ahead-numeric"`
	CacheFlushPeriod                uint32 `json:"cache-flush-period"`
	ReadAheadEnabled                string `json:"read-ahead-enabled"`
	ReadAheadEnabledNumeric         uint32 `json:"read-ahead-enabled-numeric"`
	WriteBackEnabled                string `json:"write-back-enabled"`
	WriteBackEnabledNumeric         uint32 `json:"write-back-enabled-numeric"`
	JobRunning                      string `json:"job-running"`
	CurrentJob                      string `json:"current-job"`
	CurrentJobNumeric               uint32 `json:"current-job-numeric"`
	CurrentJobCompletion            string `json:"current-job-completion"`
	NumArrayPartitions              uint32 `json:"num-array-partitions"`
	LargestFreePartitonSpace        string `json:"largest-free-partition-space"`
	LargestFreePartitonSpaceNumeric uint64 `json:"largest-free-partition-space-numeric"`
	NumDrivesPerLowLevelArray       uint8  `json:"num-drives-per-low-level-array"`
	NumExpansionPartitions          uint8  `json:"num-expansion-partitions"`
	NumPartitionSegments            uint8  `json:"num-partition-segments"`
	NewPartitionLba                 string `json:"new-partition-lba"`
	NewPartitionLbaNumeric          uint64 `json:"new-partition-lba-numeric"`
	ArrayDriveType                  string `json:"array-drive-type"`
	ArrayDriveTypeUint32            string `json:"array-drive-type-numeric"`
	DiskDescription                 string `json:"disk-description"`
	DiskDescriptionNumeric          uint32 `json:"disk-description-numeric"`
	IsJobAutoAbortable              string `json:"is-job-auto-abortable"`
	IsJobAutoAbortableNumeric       uint32 `json:"is-job-auto-abortable-numeric"`
	SerialNumber                    string `json:"serial-number"`
	Blocks                          uint64 `json:"blocks"`
	DiskDsdEnableVdisk              string `json:"disk-dsd-enable-vdisk"`
	DiskDsdEnableVdiskNumeric       uint32 `json:"disk-dsd-enable-vdisk-numeric"`
	DiskDsdEnableDelayVdisk         string `json:"disk-dsd-delay-vdisk"`
	ScrubDurationGoal               uint16 `json:"scrub-duration-goal"`
	PoolSectorFormat                string `json:"pool-sector-format"`
	PoolSectorFormatNumeric         uint32 `json:"pool-sector-format-numeric"`
	StripeWidth                     string `json:"stripe-width"`
	StripeWidthNumeric              uint32 `json:"stripe-width-numeric"`
	TargetSpareCapacity             string `json:"target-spare-capacity"`
	TargetSpareCapacityNumeric      uint64 `json:"target-spare-capacity-numeric"`
	ActualSpareCapacity             string `json:"actual-spare-capacity"`
	ActualSpareCapacityNumeric      uint64 `json:"actual-spare-capacity-numeric"`
	CriticalCapacity                string `json:"critical-capacity"`
	CriticalCapacityNumeric         uint64 `json:"critical-capacity-numeric"`
	DegradedCapacity                string `json:"degraded-capacity"`
	DegradedCapacityNumeric         uint64 `json:"degraded-capacity-numeric"`
	LinearVolumeBoundary            uint32 `json:"linear-volume-boundary"`
	MetaDataSize                    string `json:"metadata-size"`
	MetaDataSizeNumeric             uint64 `json:"metadata-size-numeric"`
	ExtendStatus                    uint64 `json:"extended-status"`
	Health                          string `json:"health"`
	HealthNumeric                   uint32 `json:"health-numeric"`
	HealthReason                    string `json:"health-reason"`
	HealthReasonNumeric             uint32 `json:"health-reason-numeric"`
	HealthRecommendation            string `json:"health-recommendation"`
	HealthRecommendationNumeric     uint32 `json:"health-recommendation-numeric"`
}

type CvtDiskGroups struct {
	Group  []CvtDiskGroup    `json:"disk-groups"`
	Status CvtResponseStatus `json:"status"`
}
