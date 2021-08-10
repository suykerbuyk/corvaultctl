package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
type CvtApiStatus []struct {
	ObjectName          string `json:"object-name"`
	Meta                string `json:"meta"`
	ResponseType        string `json:"response-type"`
	ResponseTypeNumeric int    `json:"response-type-numeric"`
	Response            string `json:"response"`
	ReturnCode          int    `json:"return-code"`
	ComponentID         string `json:"component-id"`
	TimeStamp           string `json:"time-stamp"`
	TimeStampNumeric    int    `json:"time-stamp-numeric"`
}

func (s *CvtResponseStatus) Json() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("ResponseStatus to JSON string error: " + err.Error()))
	}
	return string(prettyJSON)
}

func (c *CvtCertificates) Text() []string {
	var ret []string
	for _, cert := range c.Certificate {
		var b bytes.Buffer
		fmt.Fprintln(&b, "               object-name: ", cert.ObjectName)
		fmt.Fprintln(&b, "                      meta: ", cert.Meta)
		fmt.Fprintln(&b, "                controller: ", cert.Controller)
		fmt.Fprintln(&b, "        controller-numeric: ", cert.ControllerNumeric)
		fmt.Fprintln(&b, "        certificate-status: ", cert.CertificateStatus)
		fmt.Fprintln(&b, "certificant-status-numeric: ", cert.CertificateStatusNumeric)
		fmt.Fprintln(&b, "          certificate-time: ", cert.CertificateTime)
		fmt.Fprintln(&b, "     certificate-signature: ", cert.CertificateSignature)
		CertificateTextList := strings.Split(cert.CertificateText, "\\n")
		for _, v := range CertificateTextList {
			fmt.Fprintln(&b, v)
		}
		ret = append(ret, b.String())
	}
	return (ret)
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
	//Response CvtResponseStatus
	Status CvtApiStatus `json:"status"`
}

func (s *CvtCertificates) Json() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("CvtCertificates to JSON string error: " + err.Error()))
	}
	return string(prettyJSON)
}

type CvtDiskGroups struct {
	DiskGroups []struct {
		ObjectName                       string `json:"object-name"`
		Meta                             string `json:"meta"`
		Name                             string `json:"name"`
		URL                              string `json:"url"`
		Blocksize                        int    `json:"blocksize"`
		Size                             string `json:"size"`
		SizeNumeric                      int64  `json:"size-numeric"`
		Freespace                        string `json:"freespace"`
		FreespaceNumeric                 int    `json:"freespace-numeric"`
		RawSize                          string `json:"raw-size"`
		RawSizeNumeric                   int64  `json:"raw-size-numeric"`
		Overhead                         string `json:"overhead"`
		OverheadNumeric                  int64  `json:"overhead-numeric"`
		StorageType                      string `json:"storage-type"`
		StorageTypeNumeric               int    `json:"storage-type-numeric"`
		Pool                             string `json:"pool"`
		PoolsURL                         string `json:"pools-url"`
		PoolSerialNumber                 string `json:"pool-serial-number"`
		StorageTier                      string `json:"storage-tier"`
		StorageTierNumeric               int    `json:"storage-tier-numeric"`
		TotalPages                       int    `json:"total-pages"`
		AllocatedPages                   int    `json:"allocated-pages"`
		AvailablePages                   int    `json:"available-pages"`
		PoolPercentage                   int    `json:"pool-percentage"`
		PerformanceRank                  int    `json:"performance-rank"`
		Owner                            string `json:"owner"`
		OwnerNumeric                     int    `json:"owner-numeric"`
		PreferredOwner                   string `json:"preferred-owner"`
		PreferredOwnerNumeric            int    `json:"preferred-owner-numeric"`
		Raidtype                         string `json:"raidtype"`
		RaidtypeNumeric                  int    `json:"raidtype-numeric"`
		Diskcount                        int    `json:"diskcount"`
		InterleavedVolumeCount           int    `json:"interleaved-volume-count"`
		Spear                            string `json:"spear"`
		SpearNumeric                     int    `json:"spear-numeric"`
		Sparecount                       int    `json:"sparecount"`
		Chunksize                        string `json:"chunksize"`
		Status                           string `json:"status"`
		StatusNumeric                    int    `json:"status-numeric"`
		Lun                              int64  `json:"lun"`
		MinDriveSize                     string `json:"min-drive-size"`
		MinDriveSizeNumeric              int64  `json:"min-drive-size-numeric"`
		CreateDate                       string `json:"create-date"`
		CreateDateNumeric                int    `json:"create-date-numeric"`
		CacheReadAhead                   string `json:"cache-read-ahead"`
		CacheReadAheadNumeric            int    `json:"cache-read-ahead-numeric"`
		CacheFlushPeriod                 int    `json:"cache-flush-period"`
		ReadAheadEnabled                 string `json:"read-ahead-enabled"`
		ReadAheadEnabledNumeric          int    `json:"read-ahead-enabled-numeric"`
		WriteBackEnabled                 string `json:"write-back-enabled"`
		WriteBackEnabledNumeric          int    `json:"write-back-enabled-numeric"`
		JobRunning                       string `json:"job-running"`
		CurrentJob                       string `json:"current-job"`
		CurrentJobNumeric                int    `json:"current-job-numeric"`
		CurrentJobCompletion             string `json:"current-job-completion"`
		NumArrayPartitions               int    `json:"num-array-partitions"`
		LargestFreePartitionSpace        string `json:"largest-free-partition-space"`
		LargestFreePartitionSpaceNumeric int    `json:"largest-free-partition-space-numeric"`
		NumDrivesPerLowLevelArray        int    `json:"num-drives-per-low-level-array"`
		NumExpansionPartitions           int    `json:"num-expansion-partitions"`
		NumPartitionSegments             int    `json:"num-partition-segments"`
		NewPartitionLba                  string `json:"new-partition-lba"`
		NewPartitionLbaNumeric           int    `json:"new-partition-lba-numeric"`
		ArrayDriveType                   string `json:"array-drive-type"`
		ArrayDriveTypeNumeric            int    `json:"array-drive-type-numeric"`
		DiskDescription                  string `json:"disk-description"`
		DiskDescriptionNumeric           int    `json:"disk-description-numeric"`
		IsJobAutoAbortable               string `json:"is-job-auto-abortable"`
		IsJobAutoAbortableNumeric        int    `json:"is-job-auto-abortable-numeric"`
		SerialNumber                     string `json:"serial-number"`
		Blocks                           int64  `json:"blocks"`
		DiskDsdEnableVdisk               string `json:"disk-dsd-enable-vdisk"`
		DiskDsdEnableVdiskNumeric        int    `json:"disk-dsd-enable-vdisk-numeric"`
		DiskDsdDelayVdisk                int    `json:"disk-dsd-delay-vdisk"`
		ScrubDurationGoal                int    `json:"scrub-duration-goal"`
		PoolSectorFormat                 string `json:"pool-sector-format"`
		PoolSectorFormatNumeric          int    `json:"pool-sector-format-numeric"`
		StripeWidth                      string `json:"stripe-width"`
		StripeWidthNumeric               int    `json:"stripe-width-numeric"`
		TargetSpareCapacity              string `json:"target-spare-capacity"`
		TargetSpareCapacityNumeric       int64  `json:"target-spare-capacity-numeric"`
		ActualSpareCapacity              string `json:"actual-spare-capacity"`
		ActualSpareCapacityNumeric       int64  `json:"actual-spare-capacity-numeric"`
		CriticalCapacity                 string `json:"critical-capacity"`
		CriticalCapacityNumeric          int    `json:"critical-capacity-numeric"`
		DegradedCapacity                 string `json:"degraded-capacity"`
		DegradedCapacityNumeric          int    `json:"degraded-capacity-numeric"`
		LinearVolumeBoundary             int    `json:"linear-volume-boundary"`
		MetadataSize                     string `json:"metadata-size"`
		MetadataSizeNumeric              int    `json:"metadata-size-numeric"`
		ExtendedStatus                   int    `json:"extended-status"`
		Health                           string `json:"health"`
		HealthNumeric                    int    `json:"health-numeric"`
		HealthReason                     string `json:"health-reason"`
		HealthReasonNumeric              int    `json:"health-reason-numeric"`
		HealthRecommendation             string `json:"health-recommendation"`
		HealthRecommendationNumeric      int    `json:"health-recommendation-numeric"`
		HealthConditions                 []struct {
			ObjectName                  string `json:"object-name"`
			Meta                        string `json:"meta"`
			HealthReason                string `json:"health-reason"`
			HealthReasonNumeric         int    `json:"health-reason-numeric"`
			ReasonID                    int    `json:"reason-id"`
			HealthRecommendation        string `json:"health-recommendation"`
			HealthRecommendationNumeric int    `json:"health-recommendation-numeric"`
		} `json:"health-conditions"`
	} `json:"disk-groups"`
	Status CvtApiStatus `json:"status"`
}

func (s CvtDiskGroups) Json() string {
	prettyJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("CvtDiskGroups to JSON string error: " + err.Error()))
	}
	return string(prettyJSON)
}

type CvtDiskGroupStatistics struct {
	Statistics []struct {
		ObjectName            string `json:"object-name"`
		Meta                  string `json:"meta"`
		SerialNumber          string `json:"serial-number"`
		Name                  string `json:"name"`
		TimeSinceReset        int    `json:"time-since-reset"`
		TimeSinceSample       int    `json:"time-since-sample"`
		NumberOfReads         int64  `json:"number-of-reads"`
		NumberOfWrites        int64  `json:"number-of-writes"`
		DataRead              string `json:"data-read"`
		DataReadNumeric       int64  `json:"data-read-numeric"`
		DataWritten           string `json:"data-written"`
		DataWrittenNumeric    int64  `json:"data-written-numeric"`
		BytesPerSecond        string `json:"bytes-per-second"`
		BytesPerSecondNumeric int64  `json:"bytes-per-second-numeric"`
		Iops                  int    `json:"iops"`
		AvgRspTime            int    `json:"avg-rsp-time"`
		AvgReadRspTime        int    `json:"avg-read-rsp-time"`
		AvgWriteRspTime       int    `json:"avg-write-rsp-time"`
	} `json:"disk-group-statistics"`
	Status CvtResponseStatus
}

type CvtDiskParameters struct {
	Parameters []struct {
		ObjectName                 string `json:"object-name"`
		Meta                       string `json:"meta"`
		Smart                      string `json:"smart"`
		SmartNumeric               int    `json:"smart-numeric"`
		DriveWriteBackCache        string `json:"drive-write-back-cache"`
		DriveWriteBackCacheNumeric int    `json:"drive-write-back-cache-numeric"`
		DriveTimeoutRetryMax       int    `json:"drive-timeout-retry-max"`
		DriveAttemptTimeout        int    `json:"drive-attempt-timeout"`
		DriveOverallTimeout        int    `json:"drive-overall-timeout"`
		DiskDsdEnable              string `json:"disk-dsd-enable"`
		DiskDsdEnableNumeric       int    `json:"disk-dsd-enable-numeric"`
		DiskDsdDelay               int    `json:"disk-dsd-delay"`
		Remanufacture              string `json:"remanufacture"`
		RemanufactureNumeric       int    `json:"remanufacture-numeric"`
	} `json:"drive-parameters"`
	Status CvtResponseStatus
}

type CvtDisks struct {
	Drives []struct {
		ObjectName                  string `json:"object-name"`
		Meta                        string `json:"meta"`
		DurableID                   string `json:"durable-id"`
		EnclosureID                 int    `json:"enclosure-id"`
		DrawerID                    int    `json:"drawer-id"`
		Slot                        int    `json:"slot"`
		Location                    string `json:"location"`
		URL                         string `json:"url"`
		Port                        int    `json:"port"`
		ScsiID                      int    `json:"scsi-id"`
		Blocksize                   int    `json:"blocksize"`
		Blocks                      int64  `json:"blocks"`
		SerialNumber                string `json:"serial-number"`
		Vendor                      string `json:"vendor"`
		Model                       string `json:"model"`
		Revision                    string `json:"revision"`
		SecondaryChannel            int    `json:"secondary-channel"`
		ContainerIndex              int    `json:"container-index"`
		MemberIndex                 int    `json:"member-index"`
		Description                 string `json:"description"`
		DescriptionNumeric          int    `json:"description-numeric"`
		Architecture                string `json:"architecture"`
		ArchitectureNumeric         int    `json:"architecture-numeric"`
		Interface                   string `json:"interface"`
		InterfaceNumeric            int    `json:"interface-numeric"`
		SinglePorted                string `json:"single-ported"`
		SinglePortedNumeric         int    `json:"single-ported-numeric"`
		Type                        string `json:"type"`
		TypeNumeric                 int    `json:"type-numeric"`
		Usage                       string `json:"usage"`
		UsageNumeric                int    `json:"usage-numeric"`
		JobRunning                  string `json:"job-running"`
		JobRunningNumeric           int    `json:"job-running-numeric"`
		State                       string `json:"state"`
		CurrentJobCompletion        string `json:"current-job-completion"`
		Remanufacture               string `json:"remanufacture"`
		RemanufactureNumeric        int    `json:"remanufacture-numeric"`
		SupportsUnmap               string `json:"supports-unmap"`
		SupportsUnmapNumeric        int    `json:"supports-unmap-numeric"`
		Blink                       int    `json:"blink"`
		LocatorLed                  string `json:"locator-led"`
		LocatorLedNumeric           int    `json:"locator-led-numeric"`
		Speed                       int    `json:"speed"`
		Smart                       string `json:"smart"`
		SmartNumeric                int    `json:"smart-numeric"`
		DualPort                    int    `json:"dual-port"`
		Error                       int    `json:"error"`
		FcP1Channel                 int    `json:"fc-p1-channel"`
		FcP1DeviceID                int    `json:"fc-p1-device-id"`
		FcP1NodeWwn                 string `json:"fc-p1-node-wwn"`
		FcP1PortWwn                 string `json:"fc-p1-port-wwn"`
		FcP1UnitNumber              int    `json:"fc-p1-unit-number"`
		FcP2Channel                 int    `json:"fc-p2-channel"`
		FcP2DeviceID                int    `json:"fc-p2-device-id"`
		FcP2NodeWwn                 string `json:"fc-p2-node-wwn"`
		FcP2PortWwn                 string `json:"fc-p2-port-wwn"`
		FcP2UnitNumber              int    `json:"fc-p2-unit-number"`
		DriveDownCode               int    `json:"drive-down-code"`
		Owner                       string `json:"owner"`
		OwnerNumeric                int    `json:"owner-numeric"`
		Index                       int    `json:"index"`
		Rpm                         int    `json:"rpm"`
		Size                        string `json:"size"`
		SizeNumeric                 int64  `json:"size-numeric"`
		SectorFormat                string `json:"sector-format"`
		SectorFormatNumeric         int    `json:"sector-format-numeric"`
		TransferRate                string `json:"transfer-rate"`
		TransferRateNumeric         int    `json:"transfer-rate-numeric"`
		Attributes                  string `json:"attributes"`
		AttributesNumeric           int    `json:"attributes-numeric"`
		EnclosureWwn                string `json:"enclosure-wwn"`
		EnclosuresURL               string `json:"enclosures-url"`
		Status                      string `json:"status"`
		ReconState                  string `json:"recon-state"`
		ReconStateNumeric           int    `json:"recon-state-numeric"`
		CopybackState               string `json:"copyback-state"`
		CopybackStateNumeric        int    `json:"copyback-state-numeric"`
		VirtualDiskSerial           string `json:"virtual-disk-serial"`
		DiskGroup                   string `json:"disk-group"`
		StoragePoolName             string `json:"storage-pool-name"`
		StorageTier                 string `json:"storage-tier"`
		StorageTierNumeric          int    `json:"storage-tier-numeric"`
		SsdLifeLeft                 string `json:"ssd-life-left"`
		SsdLifeLeftNumeric          int    `json:"ssd-life-left-numeric"`
		LedStatus                   string `json:"led-status"`
		LedStatusNumeric            int    `json:"led-status-numeric"`
		DiskDsdCount                int    `json:"disk-dsd-count"`
		SpunDown                    int    `json:"spun-down"`
		NumberOfIos                 int    `json:"number-of-ios"`
		TotalDataTransferred        string `json:"total-data-transferred"`
		TotalDataTransferredNumeric int    `json:"total-data-transferred-numeric"`
		AvgRspTime                  int    `json:"avg-rsp-time"`
		FdeState                    string `json:"fde-state"`
		FdeStateNumeric             int    `json:"fde-state-numeric"`
		LockKeyID                   string `json:"lock-key-id"`
		ImportLockKeyID             string `json:"import-lock-key-id"`
		FdeConfigTime               string `json:"fde-config-time"`
		FdeConfigTimeNumeric        int    `json:"fde-config-time-numeric"`
		AssuranceLevel              string `json:"assurance-level"`
		AssuranceLevelNumeric       int    `json:"assurance-level-numeric"`
		Temperature                 string `json:"temperature"`
		TemperatureNumeric          int    `json:"temperature-numeric"`
		TemperatureStatus           string `json:"temperature-status"`
		TemperatureStatusNumeric    int    `json:"temperature-status-numeric"`
		PiFormatted                 string `json:"pi-formatted"`
		PiFormattedNumeric          int    `json:"pi-formatted-numeric"`
		PowerOnHours                int    `json:"power-on-hours"`
		ExtendedStatus              int    `json:"extended-status"`
		Health                      string `json:"health"`
		HealthNumeric               int    `json:"health-numeric"`
		HealthReason                string `json:"health-reason"`
		HealthReasonNumeric         int    `json:"health-reason-numeric"`
		HealthRecommendation        string `json:"health-recommendation"`
		HealthRecommendationNumeric int    `json:"health-recommendation-numeric"`
	} `json:"drives"`
	Status CvtResponseStatus
}

type CvtVolumes struct {
	Volumes []struct {
		ObjectName                        string `json:"object-name"`
		Meta                              string `json:"meta"`
		DurableID                         string `json:"durable-id"`
		URL                               string `json:"url"`
		VirtualDiskName                   string `json:"virtual-disk-name"`
		StoragePoolName                   string `json:"storage-pool-name"`
		StoragePoolsURL                   string `json:"storage-pools-url"`
		VolumeName                        string `json:"volume-name"`
		Size                              string `json:"size"`
		SizeNumeric                       int64  `json:"size-numeric"`
		TotalSize                         string `json:"total-size"`
		TotalSizeNumeric                  int64  `json:"total-size-numeric"`
		AllocatedSize                     string `json:"allocated-size"`
		AllocatedSizeNumeric              int64  `json:"allocated-size-numeric"`
		StorageType                       string `json:"storage-type"`
		StorageTypeNumeric                int    `json:"storage-type-numeric"`
		PreferredOwner                    string `json:"preferred-owner"`
		PreferredOwnerNumeric             int    `json:"preferred-owner-numeric"`
		Owner                             string `json:"owner"`
		OwnerNumeric                      int    `json:"owner-numeric"`
		SerialNumber                      string `json:"serial-number"`
		WritePolicy                       string `json:"write-policy"`
		WritePolicyNumeric                int    `json:"write-policy-numeric"`
		CacheOptimization                 string `json:"cache-optimization"`
		CacheOptimizationNumeric          int    `json:"cache-optimization-numeric"`
		ReadAheadSize                     string `json:"read-ahead-size"`
		ReadAheadSizeNumeric              int    `json:"read-ahead-size-numeric"`
		VolumeType                        string `json:"volume-type"`
		VolumeTypeNumeric                 int    `json:"volume-type-numeric"`
		VolumeClass                       string `json:"volume-class"`
		VolumeClassNumeric                int    `json:"volume-class-numeric"`
		TierAffinity                      string `json:"tier-affinity"`
		TierAffinityNumeric               int    `json:"tier-affinity-numeric"`
		Snapshot                          string `json:"snapshot"`
		SnapshotRetentionPriority         string `json:"snapshot-retention-priority"`
		SnapshotRetentionPriorityNumeric  int    `json:"snapshot-retention-priority-numeric"`
		VolumeQualifier                   string `json:"volume-qualifier"`
		VolumeQualifierNumeric            int    `json:"volume-qualifier-numeric"`
		Blocksize                         int    `json:"blocksize"`
		Blocks                            int64  `json:"blocks"`
		Capabilities                      string `json:"capabilities"`
		VolumeParent                      string `json:"volume-parent"`
		SnapPool                          string `json:"snap-pool"`
		ReplicationSet                    string `json:"replication-set"`
		Attributes                        string `json:"attributes"`
		VirtualDiskSerial                 string `json:"virtual-disk-serial"`
		CreationDateTime                  string `json:"creation-date-time"`
		CreationDateTimeNumeric           int    `json:"creation-date-time-numeric"`
		VolumeDescription                 string `json:"volume-description"`
		Wwn                               string `json:"wwn"`
		Progress                          string `json:"progress"`
		ProgressNumeric                   int    `json:"progress-numeric"`
		ContainerName                     string `json:"container-name"`
		ContainerSerial                   string `json:"container-serial"`
		AllowedStorageTiers               string `json:"allowed-storage-tiers"`
		AllowedStorageTiersNumeric        int    `json:"allowed-storage-tiers-numeric"`
		ThresholdPercentOfPool            string `json:"threshold-percent-of-pool"`
		ReservedSizeInPages               int    `json:"reserved-size-in-pages"`
		AllocateReservedPagesFirst        string `json:"allocate-reserved-pages-first"`
		AllocateReservedPagesFirstNumeric int    `json:"allocate-reserved-pages-first-numeric"`
		ZeroInitPageOnAllocation          string `json:"zero-init-page-on-allocation"`
		ZeroInitPageOnAllocationNumeric   int    `json:"zero-init-page-on-allocation-numeric"`
		LargeVirtualExtents               string `json:"large-virtual-extents"`
		LargeVirtualExtentsNumeric        int    `json:"large-virtual-extents-numeric"`
		Raidtype                          string `json:"raidtype"`
		RaidtypeNumeric                   int    `json:"raidtype-numeric"`
		PiFormat                          string `json:"pi-format"`
		PiFormatNumeric                   int    `json:"pi-format-numeric"`
		CsReplicationRole                 string `json:"cs-replication-role"`
		CsCopyDest                        string `json:"cs-copy-dest"`
		CsCopyDestNumeric                 int    `json:"cs-copy-dest-numeric"`
		CsCopySrc                         string `json:"cs-copy-src"`
		CsCopySrcNumeric                  int    `json:"cs-copy-src-numeric"`
		CsPrimary                         string `json:"cs-primary"`
		CsPrimaryNumeric                  int    `json:"cs-primary-numeric"`
		CsSecondary                       string `json:"cs-secondary"`
		CsSecondaryNumeric                int    `json:"cs-secondary-numeric"`
		MetadataInUse                     string `json:"metadata-in-use"`
		MetadataInUseNumeric              int    `json:"metadata-in-use-numeric"`
		Health                            string `json:"health"`
		HealthNumeric                     int    `json:"health-numeric"`
		HealthReason                      string `json:"health-reason"`
		HealthRecommendation              string `json:"health-recommendation"`
		VolumeGroup                       string `json:"volume-group"`
		GroupKey                          string `json:"group-key"`
	} `json:"volumes"`
	Status CvtResponseStatus
}

type CvtVolumeNames struct {
	VolumeNames []struct {
		ObjectName   string `json:"object-name"`
		Meta         string `json:"meta"`
		VolumeName   string `json:"volume-name"`
		SerialNumber string `json:"serial-number"`
		Volume       string `json:"volume"`
	} `json:"volume-names"`
	Status CvtResponseStatus
}

type CvtPwrSupplies struct {
	PowerSupplies []struct {
		ObjectName                string `json:"object-name"`
		Meta                      string `json:"meta"`
		DurableID                 string `json:"durable-id"`
		URL                       string `json:"url"`
		EnclosuresURL             string `json:"enclosures-url"`
		EnclosureID               int    `json:"enclosure-id"`
		DomID                     int    `json:"dom-id"`
		SerialNumber              string `json:"serial-number"`
		PartNumber                string `json:"part-number"`
		Description               string `json:"description"`
		Name                      string `json:"name"`
		FwRevision                string `json:"fw-revision"`
		Revision                  string `json:"revision"`
		Model                     string `json:"model"`
		Vendor                    string `json:"vendor"`
		Location                  string `json:"location"`
		Position                  string `json:"position"`
		PositionNumeric           int    `json:"position-numeric"`
		DashLevel                 string `json:"dash-level"`
		FruShortname              string `json:"fru-shortname"`
		MfgDate                   string `json:"mfg-date"`
		MfgDateNumeric            int    `json:"mfg-date-numeric"`
		MfgLocation               string `json:"mfg-location"`
		MfgVendorID               string `json:"mfg-vendor-id"`
		ConfigurationSerialnumber string `json:"configuration-serialnumber"`
		Dc12V                     int    `json:"dc12v"`
		Dc5V                      int    `json:"dc5v"`
		Dc33V                     int    `json:"dc33v"`
		Dc12I                     int    `json:"dc12i"`
		Dc5I                      int    `json:"dc5i"`
		Dctemp                    int    `json:"dctemp"`
		Health                    string `json:"health"`
		HealthNumeric             int    `json:"health-numeric"`
		HealthReason              string `json:"health-reason"`
		HealthRecommendation      string `json:"health-recommendation"`
		Status                    string `json:"status"`
		StatusNumeric             int    `json:"status-numeric"`
	} `json:"power-supplies"`
	Status CvtResponseStatus
}

type CvtSensors struct {
	Sensors []struct {
		ObjectName          string `json:"object-name"`
		Meta                string `json:"meta"`
		DurableID           string `json:"durable-id"`
		EnclosureID         int    `json:"enclosure-id"`
		DrawerID            string `json:"drawer-id"`
		DrawerIDNumeric     int    `json:"drawer-id-numeric"`
		ControllerID        string `json:"controller-id"`
		ControllerIDNumeric int    `json:"controller-id-numeric"`
		SensorName          string `json:"sensor-name"`
		Value               string `json:"value"`
		Status              string `json:"status"`
		StatusNumeric       int    `json:"status-numeric"`
		Container           string `json:"container"`
		ContainerNumeric    int    `json:"container-numeric"`
		SensorType          string `json:"sensor-type"`
		SensorTypeNumeric   int    `json:"sensor-type-numeric"`
	} `json:"sensors"`
	Status CvtResponseStatus
}

type CvtProvisioning struct {
	Provisioning []struct {
		ObjectName        string `json:"object-name"`
		Meta              string `json:"meta"`
		Volume            string `json:"volume"`
		VolumeSerial      string `json:"volume-serial"`
		Wwn               string `json:"wwn"`
		Controller        string `json:"controller"`
		ControllerNumeric int    `json:"controller-numeric"`
		DiskDisplay       string `json:"disk-display"`
		DiskDisplayFull   string `json:"disk-display-full"`
		VirtualDisk       string `json:"virtual-disk"`
		VirtualDiskSerial string `json:"virtual-disk-serial"`
		Health            string `json:"health"`
		HealthNumeric     int    `json:"health-numeric"`
		Mapped            string `json:"mapped"`
	} `json:"provisioning"`
	Status CvtResponseStatus
}

type CvtVersions struct {
	Versions []struct {
		ObjectName                    string    `json:"object-name"`
		Meta                          string    `json:"meta"`
		ScCPUType                     string    `json:"sc-cpu-type"`
		BundleVersion                 string    `json:"bundle-version"`
		BundleStatus                  string    `json:"bundle-status"`
		BundleStatusNumeric           int       `json:"bundle-status-numeric"`
		BundleVersionOnly             string    `json:"bundle-version-only"`
		BundleBaseVersion             string    `json:"bundle-base-version"`
		BuildDate                     time.Time `json:"build-date"`
		ScFw                          string    `json:"sc-fw"`
		ScBaselevel                   string    `json:"sc-baselevel"`
		ScMemory                      string    `json:"sc-memory"`
		ScFuVersion                   string    `json:"sc-fu-version"`
		ScLoader                      string    `json:"sc-loader"`
		CapiVersion                   string    `json:"capi-version"`
		McFw                          string    `json:"mc-fw"`
		McLoader                      string    `json:"mc-loader"`
		McBaseFw                      string    `json:"mc-base-fw"`
		FwDefaultPlatformBrand        string    `json:"fw-default-platform-brand"`
		FwDefaultPlatformBrandNumeric int       `json:"fw-default-platform-brand-numeric"`
		EcFw                          string    `json:"ec-fw"`
		PldRev                        string    `json:"pld-rev"`
		PmCpldVersion                 string    `json:"pm-cpld-version"`
		PrmVersion                    string    `json:"prm-version"`
		HwRev                         string    `json:"hw-rev"`
		HimRev                        string    `json:"him-rev"`
		HimModel                      string    `json:"him-model"`
		BackplaneType                 int       `json:"backplane-type"`
		HostChannelRevision           int       `json:"host-channel_revision"`
		DiskChannelRevision           int       `json:"disk-channel_revision"`
		MrcVersion                    string    `json:"mrc-version"`
		CtkVersion                    string    `json:"ctk-version"`
		McosVersion                   string    `json:"mcos-version"`
		GemVersion                    string    `json:"gem-version"`
		PubsVersion                   string    `json:"pubs-version"`
		TranslationVersion            string    `json:"translation-version"`
	} `json:"versions"`
	Status CvtResponseStatus
}

type CvtDnsMgmtHostnames struct {
	MgmtHostnames []struct {
		ObjectName             string `json:"object-name"`
		Meta                   string `json:"meta"`
		Controller             string `json:"controller"`
		ControllerNumeric      int    `json:"controller-numeric"`
		MgmtHostname           string `json:"mgmt-hostname"`
		DomainName             string `json:"domain-name"`
		DefaultHostname        string `json:"default-hostname"`
		DefaultHostnameNumeric int    `json:"default-hostname-numeric"`
	} `json:"mgmt-hostnames"`
	Status CvtResponseStatus
}

type CvtInquiry struct {
	ProductInfo []struct {
		ObjectName    string `json:"object-name"`
		Meta          string `json:"meta"`
		VendorName    string `json:"vendor-name"`
		ProductID     string `json:"product-id"`
		ScsiVendorID  string `json:"scsi-vendor-id"`
		ScsiProductID string `json:"scsi-product-id"`
	} `json:"product-info"`
	Inquiry []struct {
		ObjectName                  string `json:"object-name"`
		Meta                        string `json:"meta"`
		McFw                        string `json:"mc-fw"`
		McLoader                    string `json:"mc-loader"`
		ScFw                        string `json:"sc-fw"`
		ScLoader                    string `json:"sc-loader"`
		SerialNumber                string `json:"serial-number"`
		MacAddress                  string `json:"mac-address"`
		IPAddress                   string `json:"ip-address"`
		IP6LinkLocalAddress         string `json:"ip6-link-local-address"`
		IP6AutoAddress              string `json:"ip6-auto-address"`
		Dhcpv6                      string `json:"dhcpv6"`
		SlaacIP                     string `json:"slaac-ip"`
		IP6AutoAddressSource        string `json:"ip6-auto-address-source"`
		IP6AutoAddressSourceNumeric int    `json:"ip6-auto-address-source-numeric"`
		IP61Address                 string `json:"ip61-address"`
		IP62Address                 string `json:"ip62-address"`
		IP63Address                 string `json:"ip63-address"`
		IP64Address                 string `json:"ip64-address"`
		NvramDefaults               string `json:"nvram-defaults"`
	} `json:"inquiry"`
	Status CvtResponseStatus
}

type CvtVolumeStatistics struct {
	ProductInfo []struct {
		ObjectName    string `json:"object-name"`
		Meta          string `json:"meta"`
		VendorName    string `json:"vendor-name"`
		ProductID     string `json:"product-id"`
		ScsiVendorID  string `json:"scsi-vendor-id"`
		ScsiProductID string `json:"scsi-product-id"`
	} `json:"product-info"`
	Inquiry []struct {
		ObjectName                  string `json:"object-name"`
		Meta                        string `json:"meta"`
		McFw                        string `json:"mc-fw"`
		McLoader                    string `json:"mc-loader"`
		ScFw                        string `json:"sc-fw"`
		ScLoader                    string `json:"sc-loader"`
		SerialNumber                string `json:"serial-number"`
		MacAddress                  string `json:"mac-address"`
		IPAddress                   string `json:"ip-address"`
		IP6LinkLocalAddress         string `json:"ip6-link-local-address"`
		IP6AutoAddress              string `json:"ip6-auto-address"`
		Dhcpv6                      string `json:"dhcpv6"`
		SlaacIP                     string `json:"slaac-ip"`
		IP6AutoAddressSource        string `json:"ip6-auto-address-source"`
		IP6AutoAddressSourceNumeric int    `json:"ip6-auto-address-source-numeric"`
		IP61Address                 string `json:"ip61-address"`
		IP62Address                 string `json:"ip62-address"`
		IP63Address                 string `json:"ip63-address"`
		IP64Address                 string `json:"ip64-address"`
		NvramDefaults               string `json:"nvram-defaults"`
	} `json:"inquiry"`
	Status CvtResponseStatus
}

type CvtReservations struct {
	VolumeReservations []struct {
		ObjectName               string `json:"object-name"`
		Meta                     string `json:"meta"`
		VolumeName               string `json:"volume-name"`
		SerialNumber             string `json:"serial-number"`
		ReservationActive        string `json:"reservation-active"`
		ReservationActiveNumeric int    `json:"reservation-active-numeric"`
		PgrGeneration            string `json:"pgr-generation"`
		HostID                   string `json:"host-id"`
		Port                     string `json:"port"`
		ReserveKey               string `json:"reserve-key"`
		ReserveScope             string `json:"reserve-scope"`
		ReserveScopeNumeric      int    `json:"reserve-scope-numeric"`
		ReserveType              string `json:"reserve-type"`
		ReserveTypeNumeric       int    `json:"reserve-type-numeric"`
	} `json:"volume-reservations"`
	Status CvtResponseStatus
}

type CvtHostGroups struct {
	HostGroup []struct {
		ObjectName   string `json:"object-name"`
		Meta         string `json:"meta"`
		DurableID    string `json:"durable-id"`
		Name         string `json:"name"`
		SerialNumber string `json:"serial-number"`
		URL          string `json:"url"`
		MemberCount  int    `json:"member-count"`
		Host         []struct {
			ObjectName   string `json:"object-name"`
			Meta         string `json:"meta"`
			DurableID    string `json:"durable-id"`
			Name         string `json:"name"`
			SerialNumber string `json:"serial-number"`
			MemberCount  int    `json:"member-count"`
			HostGroup    string `json:"host-group"`
			GroupKey     string `json:"group-key"`
			Initiator    []struct {
				ObjectName         string `json:"object-name"`
				Meta               string `json:"meta"`
				DurableID          string `json:"durable-id"`
				Nickname           string `json:"nickname"`
				Discovered         string `json:"discovered"`
				Mapped             string `json:"mapped"`
				Profile            string `json:"profile"`
				ProfileNumeric     int    `json:"profile-numeric"`
				HostBusType        string `json:"host-bus-type"`
				HostBusTypeNumeric int    `json:"host-bus-type-numeric"`
				ID                 string `json:"id"`
				URL                string `json:"url"`
				HostID             string `json:"host-id"`
				HostKey            string `json:"host-key"`
				HostPortBitsA      int    `json:"host-port-bits-a"`
				HostPortBitsB      int    `json:"host-port-bits-b"`
			} `json:"initiator"`
		} `json:"host"`
	} `json:"host-group"`
	Host []struct {
		ObjectName   string `json:"object-name"`
		Meta         string `json:"meta"`
		DurableID    string `json:"durable-id"`
		Name         string `json:"name"`
		SerialNumber string `json:"serial-number"`
		MemberCount  int    `json:"member-count"`
		HostGroup    string `json:"host-group"`
		GroupKey     string `json:"group-key"`
		Initiator    []struct {
			ObjectName         string `json:"object-name"`
			Meta               string `json:"meta"`
			DurableID          string `json:"durable-id"`
			Nickname           string `json:"nickname"`
			Discovered         string `json:"discovered"`
			Mapped             string `json:"mapped"`
			Profile            string `json:"profile"`
			ProfileNumeric     int    `json:"profile-numeric"`
			HostBusType        string `json:"host-bus-type"`
			HostBusTypeNumeric int    `json:"host-bus-type-numeric"`
			ID                 string `json:"id"`
			URL                string `json:"url"`
			HostID             string `json:"host-id"`
			HostKey            string `json:"host-key"`
			HostPortBitsA      int    `json:"host-port-bits-a"`
			HostPortBitsB      int    `json:"host-port-bits-b"`
		} `json:"initiator"`
	} `json:"host"`
	Status CvtResponseStatus
}

type CvtHostPhyStatistics struct {
	SasHostPhyStatistics []struct {
		ObjectName        string `json:"object-name"`
		Meta              string `json:"meta"`
		Port              string `json:"port"`
		Phy               int    `json:"phy"`
		DisparityErrors   string `json:"disparity-errors"`
		LostDwords        string `json:"lost-dwords"`
		InvalidDwords     string `json:"invalid-dwords"`
		ResetErrorCounter string `json:"reset-error-counter"`
	} `json:"sas-host-phy-statistics"`
	Status CvtResponseStatus
}

type CvtHostPortStatistics struct {
	HostPortStatistics []struct {
		ObjectName             string `json:"object-name"`
		Meta                   string `json:"meta"`
		DurableID              string `json:"durable-id"`
		BytesPerSecond         string `json:"bytes-per-second"`
		BytesPerSecondNumeric  int    `json:"bytes-per-second-numeric"`
		Iops                   int    `json:"iops"`
		NumberOfReads          int    `json:"number-of-reads"`
		NumberOfWrites         int    `json:"number-of-writes"`
		DataRead               string `json:"data-read"`
		DataReadNumeric        int    `json:"data-read-numeric"`
		DataWritten            string `json:"data-written"`
		DataWrittenNumeric     int    `json:"data-written-numeric"`
		QueueDepth             int    `json:"queue-depth"`
		AvgRspTime             int    `json:"avg-rsp-time"`
		AvgReadRspTime         int    `json:"avg-read-rsp-time"`
		AvgWriteRspTime        int    `json:"avg-write-rsp-time"`
		ResetTime              string `json:"reset-time"`
		ResetTimeNumeric       int    `json:"reset-time-numeric"`
		StartSampleTime        string `json:"start-sample-time"`
		StartSampleTimeNumeric int    `json:"start-sample-time-numeric"`
		StopSampleTime         string `json:"stop-sample-time"`
		StopSampleTimeNumeric  int    `json:"stop-sample-time-numeric"`
	} `json:"host-port-statistics"`
	Status CvtResponseStatus
}

type CvtAdvancedSettings struct {
	AdvancedSettingsTable []struct {
		ObjectName                             string `json:"object-name"`
		Meta                                   string `json:"meta"`
		BackgroundScrub                        string `json:"background-scrub"`
		BackgroundScrubNumeric                 int    `json:"background-scrub-numeric"`
		BackgroundScrubInterval                int    `json:"background-scrub-interval"`
		PartnerFirmwareUpgrade                 string `json:"partner-firmware-upgrade"`
		PartnerFirmwareUpgradeNumeric          int    `json:"partner-firmware-upgrade-numeric"`
		UtilityPriority                        string `json:"utility-priority"`
		UtilityPriorityNumeric                 int    `json:"utility-priority-numeric"`
		Smart                                  string `json:"smart"`
		SmartNumeric                           int    `json:"smart-numeric"`
		DynamicSpares                          string `json:"dynamic-spares"`
		EmpPollRate                            string `json:"emp-poll-rate"`
		HostCacheControl                       string `json:"host-cache-control"`
		HostCacheControlNumeric                int    `json:"host-cache-control-numeric"`
		SyncCacheMode                          string `json:"sync-cache-mode"`
		SyncCacheModeNumeric                   int    `json:"sync-cache-mode-numeric"`
		IndependentCache                       string `json:"independent-cache"`
		IndependentCacheNumeric                int    `json:"independent-cache-numeric"`
		MissingLunResponse                     string `json:"missing-lun-response"`
		MissingLunResponseNumeric              int    `json:"missing-lun-response-numeric"`
		ControllerFailure                      string `json:"controller-failure"`
		ControllerFailureNumeric               int    `json:"controller-failure-numeric"`
		SuperCapFailure                        string `json:"super-cap-failure"`
		SuperCapFailureNumeric                 int    `json:"super-cap-failure-numeric"`
		MemoryCardFailure                      string `json:"memory-card-failure"`
		MemoryCardFailureNumeric               int    `json:"memory-card-failure-numeric"`
		PowerSupplyFailure                     string `json:"power-supply-failure"`
		PowerSupplyFailureNumeric              int    `json:"power-supply-failure-numeric"`
		FanFailure                             string `json:"fan-failure"`
		FanFailureNumeric                      int    `json:"fan-failure-numeric"`
		TemperatureExceeded                    string `json:"temperature-exceeded"`
		TemperatureExceededNumeric             int    `json:"temperature-exceeded-numeric"`
		PartnerNotify                          string `json:"partner-notify"`
		PartnerNotifyNumeric                   int    `json:"partner-notify-numeric"`
		AutoWriteBack                          string `json:"auto-write-back"`
		AutoWriteBackNumeric                   int    `json:"auto-write-back-numeric"`
		DiskDsdEnable                          string `json:"disk-dsd-enable"`
		DiskDsdEnableNumeric                   int    `json:"disk-dsd-enable-numeric"`
		DiskDsdDelay                           int    `json:"disk-dsd-delay"`
		BackgroundDiskScrub                    string `json:"background-disk-scrub"`
		BackgroundDiskScrubNumeric             int    `json:"background-disk-scrub-numeric"`
		ManagedLogs                            string `json:"managed-logs"`
		ManagedLogsNumeric                     int    `json:"managed-logs-numeric"`
		SingleController                       string `json:"single-controller"`
		SingleControllerNumeric                int    `json:"single-controller-numeric"`
		AutoStallRecovery                      string `json:"auto-stall-recovery"`
		AutoStallRecoveryNumeric               int    `json:"auto-stall-recovery-numeric"`
		DeleteOverride                         string `json:"delete-override"`
		DeleteOverrideNumeric                  int    `json:"delete-override-numeric"`
		RestartOnCapiFail                      string `json:"restart-on-capi-fail"`
		RestartOnCapiFailNumeric               int    `json:"restart-on-capi-fail-numeric"`
		LargePools                             string `json:"large-pools"`
		LargePoolsNumeric                      int    `json:"large-pools-numeric"`
		SsdConcurrentAccess                    string `json:"ssd-concurrent-access"`
		SsdConcurrentAccessNumeric             int    `json:"ssd-concurrent-access-numeric"`
		SlotAffinity                           string `json:"slot-affinity"`
		SlotAffinityNumeric                    int    `json:"slot-affinity-numeric"`
		RandomIoPerformanceOptimization        string `json:"random-io-performance-optimization"`
		RandomIoPerformanceOptimizationNumeric int    `json:"random-io-performance-optimization-numeric"`
		CacheFlushTimeout                      string `json:"cache-flush-timeout"`
		CacheFlushTimeoutNumeric               int    `json:"cache-flush-timeout-numeric"`
		Remanufacture                          string `json:"remanufacture"`
		RemanufactureNumeric                   int    `json:"remanufacture-numeric"`
		HedgedReadsTimeout                     string `json:"hedged-reads-timeout"`
		HedgedReadsTimeoutNumeric              int    `json:"hedged-reads-timeout-numeric"`
	} `json:"advanced-settings-table"`
	Status CvtResponseStatus
}

type CvtSystem struct {
	System []struct {
		ObjectName               string `json:"object-name"`
		Meta                     string `json:"meta"`
		SystemName               string `json:"system-name"`
		SystemContact            string `json:"system-contact"`
		SystemLocation           string `json:"system-location"`
		SystemInformation        string `json:"system-information"`
		MidplaneSerialNumber     string `json:"midplane-serial-number"`
		URL                      string `json:"url"`
		VendorName               string `json:"vendor-name"`
		ProductID                string `json:"product-id"`
		ProductBrand             string `json:"product-brand"`
		ScsiVendorID             string `json:"scsi-vendor-id"`
		ScsiProductID            string `json:"scsi-product-id"`
		EnclosureCount           int    `json:"enclosure-count"`
		Health                   string `json:"health"`
		HealthNumeric            int    `json:"health-numeric"`
		HealthReason             string `json:"health-reason"`
		OtherMCStatus            string `json:"other-MC-status"`
		OtherMCStatusNumeric     int    `json:"other-MC-status-numeric"`
		PfuStatus                string `json:"pfuStatus"`
		PfuStatusNumeric         int    `json:"pfuStatus-numeric"`
		SupportedLocales         string `json:"supported-locales"`
		CurrentNodeWwn           string `json:"current-node-wwn"`
		FdeSecurityStatus        string `json:"fde-security-status"`
		FdeSecurityStatusNumeric int    `json:"fde-security-status-numeric"`
		PlatformType             string `json:"platform-type"`
		PlatformTypeNumeric      int    `json:"platform-type-numeric"`
		PlatformBrand            string `json:"platform-brand"`
		PlatformBrandNumeric     int    `json:"platform-brand-numeric"`
		Redundancy               []struct {
			ObjectName               string `json:"object-name"`
			Meta                     string `json:"meta"`
			RedundancyMode           string `json:"redundancy-mode"`
			RedundancyModeNumeric    int    `json:"redundancy-mode-numeric"`
			RedundancyStatus         string `json:"redundancy-status"`
			RedundancyStatusNumeric  int    `json:"redundancy-status-numeric"`
			ControllerAStatus        string `json:"controller-a-status"`
			ControllerAStatusNumeric int    `json:"controller-a-status-numeric"`
			ControllerASerialNumber  string `json:"controller-a-serial-number"`
			ControllerBStatus        string `json:"controller-b-status"`
			ControllerBStatusNumeric int    `json:"controller-b-status-numeric"`
			ControllerBSerialNumber  string `json:"controller-b-serial-number"`
			OtherMCStatus            string `json:"other-MC-status"`
			OtherMCStatusNumeric     int    `json:"other-MC-status-numeric"`
			SystemReady              string `json:"system-ready"`
			SystemReadyNumeric       int    `json:"system-ready-numeric"`
			LocalReady               string `json:"local-ready"`
			LocalReadyNumeric        int    `json:"local-ready-numeric"`
			LocalReason              string `json:"local-reason"`
			OtherReady               string `json:"other-ready"`
			OtherReadyNumeric        int    `json:"other-ready-numeric"`
			OtherReason              string `json:"other-reason"`
		} `json:"redundancy"`
		UnhealthyComponent []struct {
			ObjectName           string `json:"object-name"`
			Meta                 string `json:"meta"`
			ComponentType        string `json:"component-type"`
			ComponentTypeNumeric int    `json:"component-type-numeric"`
			ComponentID          string `json:"component-id"`
			Basetype             string `json:"basetype"`
			PrimaryKey           string `json:"primary-key"`
			Health               string `json:"health"`
			HealthNumeric        int    `json:"health-numeric"`
			HealthReason         string `json:"health-reason"`
			HealthRecommendation string `json:"health-recommendation"`
		} `json:"unhealthy-component"`
	} `json:"system"`
	Status CvtResponseStatus
}
