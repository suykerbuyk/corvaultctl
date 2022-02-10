#!/bin/sh

set -e
export SSHPASS='Testit123!'


USER='manage'
TARGETS=("corvault-1a" "corvault-2a" "corvault-3a")
#TARGETS=("corvault-1a")

# provides a wee bit more verbosity to stderr
DBG=0
# prepatory command to the corvault
BASE_CMD='set cli-parameters json; '

# Make sure we have a support utils installed.
if [[ $(which jq 2>&1>/dev/null) ]]; then
       echo "Please intall jq"
       exit 1
elif [[ $(which sshpass 2>&1>/dev/null) ]]; then
       echo "Please intall sshpass"
       exit 1
fi

# An example of using the output of this utility to configure and map 
# Initiators to volumes:
# set initiator id 500062b206989400 nickname SAS9305-16e-SPA2500634-p0
# set initiator id 500062b206989401 nickname SAS9305-16e-SPA2500634-p1
# set initiator id 500062b206989408 nickname SAS9305-16e-SPA2500634-p2
# set initiator id 500062b206989649 nickname SAS9305-16e-SPA2500494-p3
# set initiator id 500062b206989640 nickname SAS9305-16e-SPA2500494-p0
# set initiator id 500062b206989641 nickname SAS9305-16e-SPA2500494-p1
# set initiator id 500062b206989648 nickname SAS9305-16e-SPA2500494-p2
# set initiator id 500062b206989409 nickname SAS9305-16e-SPA2500634-p3
# set initiator id 500605b00de2aa80 nickname HBA_9405W-16e-SPA2500494-p0
# set initiator id 500605b00de2aa81 nickname HBA_9405W-16e-SPA2500494-p1
# set initiator id 500605b00de2aa88 nickname HBA_9405W-16e-SPA2500494-p2
# set initiator id 500605b00de2aa89 nickname HBA_9405W-16e-SPA2500494-p3
# set initiator id 500605b00de2b640 nickname HBA_9405W-16e-SP82331324-p0
# set initiator id 500605b00de2b641 nickname HBA_9405W-16e-SP82331324-p1
# set initiator id 500605b00de2b648 nickname HBA_9405W-16e-SP82331324-p2
# set initiator id 500605b00de2b649 nickname HBA_9405W-16e-SP82331324-p3
# corvault-1 map volume Volume_0000 initiator SAS9305-16e-SPA2500494-p1 lun 2
# corvault-1 map volume Volume_0001 initiator SAS9305-16e-SPA2500634-p1 lun 2
# corvault-2 map volume VOL0000 initiator SAS9305-16e-SPA2500494-p0 lun 2
# corvault-2 map volume VOL0001 initiator SAS9305-16e-SPA2500634-p0 lun 2
# corvault-3 map volume VOL0000 initiator SAS9305-16e-SPA2500494-p2 lun 2
# corvault-3 map volume VOL0001 initiator SAS9305-16e-SPA2500634-p2 lun 2

# interesting sysfs paths for coorelating LUNs to host HBA ports and kdevs
#cat /sys/devices/*/*/*/host*/phy-*/sas_phy/*/sas_address | sort -u | cut -c 15-
#cat /sys/devices/pci*/*/*/host*/port*/end_device*/target*/*/sas_address | sort -u

# cat stor.show.all.json | jq -r '."Controllers"[]? | ."Response Data"."IT System Overview"[] | (.Ctl |tostring) + ",\t" + .Model + ",\t" + ."AdapterType"  + ",\t" + ."PCI Address" + ",\t" '
#  0,      HBA 9405W-16e,    SAS3616(B0),  00:19:00:00,
#  1,      HBA 9405W-16e,    SAS3616(B0),  00:b0:00:00,
#  2,      SAS9305-16e,      SAS3216(A1),  00:86:00:00,
#  3,      SAS9305-16e,      SAS3216(A1),  00:af:00:00,
# cat store.show.c0.all.json | jq -r '."Controllers"[]? | ."Response Data"."Basics"| ."SAS Address" + ",\t" + ."PCI Address"'
#  500605b00de2aa80,      00:19:00:00



# kind of like atop, but for Corvault
monitor_io() {
	while [ 1 ] ; do 
	date
	sshpass -e ssh manage@corvault-1a \
		'set cli-parameters json; show controller-statistics' \
		| grep bytes-per-second \
		| grep -v numeric
		sleep 10
	done
}
# Dispatches commands to the Corvault in a way that's easy to capture.
DoSSH() {
	sshpass -e $@
}
# The "meat & potatoes" - dispatches commands parses the fubar returned JSON for Corvault into something useful.
DoCmd() {
	TGT="${1}"
	shift
	REPLY_FILE="${TGT}.json"
	[[ $DBG != 0 ]] && echo "TGT: $TGT  CMD: $BASE_CMD $@" 1>&2
	SSHSOCKET=/tmp/$TGT.ssh.socket
	SSHOPTS="-o ControlPath=$SSHSOCKET -o ControlMaster=auto -o ControlPersist=10m -o StrictHostKeyChecking=accept-new"
	REPLY=$(DoSSH "ssh ${SSHOPTS} ${USER}@${TGT} ${BASE_CMD} $@")
	# Pull off the commented lines that contain the commands sent to the target
	[[ $DBG != 0 ]] && printf "REPLY: %s\n" "$REPLY" 1>&2
	REQ=$(echo "$REPLY" | egrep '^#.*' | sed -e 's/^#[ ]*//g' -e '/^$/d' | sed -e :a -e '$!N; s/\n/; /; ta')
	[[ $DBG != 0 ]] && printf "REQ: $REQ\n" 1>&2
	JSON=$(printf "%s\n" "$REPLY" | awk '/#  /,0' | egrep -v '^# .*' |  sed -e :a -e '$!N;  ta')
	[[ $DBG != 0 ]] && printf "JSON: %s\n" "$JSON" 1>&2
	RESP=$(echo ${JSON} | jq -r '.status[].response')
	STAT=$(echo ${JSON} | jq -r '.status[]."response-type"')
	[[ $DBG != 0 ]] && printf "RESP: %s\n" "$RESP" 1>&2
	[[ $DBG != 0 ]] && printf "STAT: %s\n" "$STAT" 1>&2
	if [ "${STAT}" != "Success" ] ; then
		echo "${REPLY}" >"${REPLY_FILE}"
		echo "Error: $BASE_CMD $@" 1>&2;
		echo "Status: ${STAT}" 1>&2;
		echo "Response: ${RESP}" 1>&2;
		echo "See ${REPLY} for full JSON return data" 1>&2;
		exit 1
	fi
	[[ $DBG != 0 ]] && echo "Status: ${STAT}" 1>&2
	printf "%s\n" "$JSON"
}
ShowInquiryJSON() {
	TGT=$1
	CMD="show inquiry"
	DoCmd ${TGT} "${CMD}"
}
ShowSensorStatusJSON() {
	TGT=$1
	CMD="show sensor-status"
	DoCmd ${TGT} "${CMD}"
}
ShowConfigurationJSON() {
	TGT=$1
	CMD="show configuration"
	DoCmd ${TGT} "${CMD}"
}
ShowHostPhyStatisticsJSON() {
	TGT=$1
	CMD="show host-phy-statistics"
	DoCmd ${TGT} "${CMD}"
}
ShowDiskGroupsJSON() {
	TGT="$1"
	CMD="show disk-groups"
	DoCmd ${TGT} "${CMD}"
}
ShowDisksJSON() {
	TGT=$1
	CMD="show disks"
	DoCmd ${TGT} "${CMD}"
}
ShowVolumesJSON() {
	TGT="$1"
	CMD="show volumes"
	DoCmd ${TGT} "${CMD}"
}
ShowInitiatorsJSON() {
	TGT="$1"
	CMD="show initiators"
	DoCmd ${TGT} "${CMD}"
}
ShowMapsJSON() {
	TGT="$1"
	CMD="show maps"
	DoCmd ${TGT} "${CMD}"
}
ShowMpt3SasHBAsJSON() {
	PREFIX="  {\n"
	printf "{\n\"mpt3hba\":[\n"
	for X in $(find /sys/class/scsi_host/host*/ | grep host_sas_address)
	do
		CTRLR_PATH=$(dirname $(realpath $X))
		PCI_ADDR=$(printf "$CTRLR_PATH" | awk -F '/' '{print $6}')
		PCI_HOST_PATH="$(dirname $(dirname $(dirname $CTRLR_PATH)))"
		PCI_VENDOR="$(cat $PCI_HOST_PATH/vendor)"
		PCI_SUBSYSTEM_VENDOR="$(cat $PCI_HOST_PATH/subsystem_vendor)"
		PCI_SUBSYSTEM_DEVICE="$(cat $PCI_HOST_PATH/subsystem_device)"
		UNIQUE_ID="$(cat $CTRLR_PATH/unique_id)"
		SAS_ADDR="$(cat $CTRLR_PATH/host_sas_address)"
		BOARD_NAME="$(cat $CTRLR_PATH/board_name)"
		BOARD_ASSEMBLY="$(cat $CTRLR_PATH/board_assembly)"
		VERSION_BIOS="$(cat $CTRLR_PATH/version_bios)"
		VERSION_FW="$(cat $CTRLR_PATH/version_fw)"
		VERSION_MPI="$(cat $CTRLR_PATH/version_mpi)"
		VERSION_NVDATA="$(cat $CTRLR_PATH/version_nvdata_persistent)"
		VERSION_PRODUCT="$(cat $CTRLR_PATH/version_product)"
		printf $PREFIX
		printf "  \"sysfs-path\": \"$CTRLR_PATH\",\n"
		printf "  \"unique-id\": \"$UNIQUE_ID\",\n"
		printf "  \"pci-vendor\": \"$PCI_VENDOR\",\n"
		printf "  \"pci-subsystem-vendor\": \"$PCI_SUBSYSTEM_VENDOR\",\n"
		printf "  \"pci-subsystem-device\": \"$PCI_SUBSYSTEM_DEVICE\",\n"
		printf "  \"board-name\": \"$BOARD_NAME\",\n"
		printf "  \"board-assembly\": \"$BOARD_ASSEMBLY\",\n"
		printf "  \"sas-address\": \"$SAS_ADDR\",\n"
		printf "  \"pci-address\": \"$PCI_ADDR\",\n"
		printf "  \"version-fw\": \"$VERSION_FW\",\n"
		printf "  \"version-bios\": \"$VERSION_BIOS\",\n"
		printf "  \"version-mpi\": \"$VERSION_MPI\",\n"
		printf "  \"version-nvdata\": \"$VERSION_NVDATA\",\n"
		printf "  \"version-product\": \"$VERSION_PRODUCT\"\n"
		PREFIX="  },\n{\n"
	done
	printf "  }\n"
	printf "]}\n"
}
ShowMpt3SasHBAs() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	   .mpt3hba[]? | "\t"
           + ."unique-id" + ",\t"
           + ."board-name" + ",\t"
           + ."version-product" + ",\t"
           + ."sas-address" + ",\t"
           + ."pci-address" + ",\t"
           + ."pci-vendor" + ",\t"
           + ."pci-subsystem-device" + ",\t"
           + ."version-fw"  + ",\t"
           + ."version-bios"  + ",\t"
           + ."version-mpi"  + ",\t"
           + ."version-nvdata"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	HDR01=" UniqueID,"
	HDR02="          BoardName,"
	HDR03="  ChipSet,"
	HDR04="               SAS_Address,"
	HDR05="      PCI_Address,"
	HDR06="   Vendor,"
	HDR07=" Device,"
	HDR08=" FirmwareVer,"
	HDR09="        BiosVer,"
	HDR10="    MpiVer,"
	HDR11=" NvDataVer"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}${HDR08}${HDR09}${HDR10}${HDR11}"
	printf "$HDR\n"
	RESULT=$(ShowMpt3SasHBAsJSON | jq  -r "${JQ}")
	printf "${RESULT}\n"
}
GetInquiryNoHdr() {
	TGT=$1
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	 $T + ",\t" +
	 (
           ."product-info"[] | ."product-id" + ",\t"
         )
	 + ( 
             .inquiry[] 
             | ."object-name" + ",\t"
             + ."serial-number" + ",\t"
             + ."mc-fw" + ",\t" + ."sc-fw"
             + ",\t" + ."mc-loader" + ",\t"
             + ."sc-loader" + ",\t\t"
             + ."mac-address" + ",\t"
             + ."ip-address"
           )
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowInquiryJSON $TGT | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetInquiry() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	HDR00="controller,\t"
	HDR01="product-id,\t"
	HDR02="controller,\t"
	HDR03="serial,\t\t\t"
	HDR04="mc-fw,\t\t"
	HDR05="sc-fw,\t\t"
	HDR06="mc-loader,\t"
	HDR07="sc-loader,\t"
	HDR08="mac-address,\t\t"
	HDR09="ip-address"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}${HDR08}${HDR09}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetInquiryNoHdr $TGT
	done
}
GetVolumesNoHdr() {
	TGT=$1
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.volumes[]?
	 | $T + ",\t"
	 + ."volume-name" +",\t"
	 + ."virtual-disk-name" + ",\t"
	 + ."size" + ",\t"
	 + ."serial-number" + ",\t" + (."wwn" | ascii_downcase) + ",\t"
	 + ."creation-date-time"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowVolumesJSON $TGT | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetVolumes() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR00="controller,\t"
	HDR01="volume-name,\t"
	HDR02="name,\t"
	HDR03="size,\t\t"
	HDR04="serial-number,\t\t\t\t"
	HDR05="wwn-number,\t\t\t\t"
	HDR06="creation-date-time"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetVolumesNoHdr $TGT
	done
}
GetInitiatorsNoHdr() {
	TGT=$1
	FILTERED=$2
	if [[ $FILTERED == 0 ]]; then
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.initiator[]
 	 | if ."host-id" == "NOHOST" then ."host-id"="NOHOST,                          " else ."host-id" = ."host-id" + "," end
	 | $T + ",\t"
	 + ."durable-id" + ",\t"
	 +  .discovered + ",\t"
	 +  .id + ",\t"
	 + ."host-id" + "\t"
	 + ."host-key" + ",\t\t"
	 + .nickname
EOF
)
	else
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.initiator[] | select(.discovered == "Yes") 
	 | if ."host-id" == "NOHOST" then ."host-id"="NOHOST,                          " else ."host-id" = ."host-id" + "," end
	 | $T + ",\t"
	 +  ."durable-id" + ",\t"
	 +  .discovered + ",\t"
	 +  .id + ",\t"
	 + ."host-id" + "\t"
	 + ."host-key" + ",\t\t"
	 + .nickname
EOF
)
	fi
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowInitiatorsJSON $TGT | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetInitiators() {
	TGT=$1
	FILTERED=1
	if [[ $FILTERED == 0 ]]; then
		printf "\nRUN: $TGT ${FUNCNAME[0]} (unfiltered)\n"
	else
		printf "\nRUN: $TGT ${FUNCNAME[0]} (filtered for only discovered initiators)\n"
	fi
	HDR00="controller,\t"
	HDR01="d-id,\t"
	HDR02="dscvrd,\t\t"
	HDR03="id,\t\t"
	HDR04="host-id,\t\t\t  "
	HDR05="host-key,\t\t"
	HDR06="nickname"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR05}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetInitiatorsNoHdr $TGT $FILTERED
	done
}
GetHostPhyStatisticsNoHdr() {
	TGT=$1
	FILTERED=$2
	if [[ $FILTERED == 0 ]]; then
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."sas-host-phy-statistics"[]
	 | $T + ",\t"
	 + .port + "-" + (.phy|tostring) + ",\t"
	 + ."disparity-errors" +",\t"
	 + ."lost-dwords" + ",\t"
	 + ."invalid-dwords" + ",\t"
	 + ."reset-error-counter"
EOF
)
	else
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."sas-host-phy-statistics"[]
	 |  select((((."disparity-errors" != "00000000")
	 or ."lost-dwords" != "00000000")
	 or ."invalid-dwords" != "00000000")
	 or ."reset-error-counter" != "00000000")
	 | $T + ",\t"
	 + .port + "-" + (.phy|tostring) + ",\t"
	 + ."disparity-errors" +",\t"
	 + ."lost-dwords" + ",\t"
	 + ."invalid-dwords" + ",\t"
	 + ."reset-error-counter"
EOF
)
	fi
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowHostPhyStatisticsJSON $TGT | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetHostPhyStatistics() {
	TGT=$1
	FILTERED=0
	if [[ $FILTERED == 0 ]]; then
		printf "\nRUN: ${FUNCNAME[0]} (unfiltered)\n"
	else
		printf "\nRUN: ${FUNCNAME[0]} (filtered for non-zero counters)\n"
	fi
	HDR00="controller,\t"
	HDR01="port,\t"
	HDR02="disparities,\t"
	HDR03="lost-dws,\t"
	HDR04="invalid-dws,\t"
	HDR05="reset-errs"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetHostPhyStatisticsNoHdr "$TGT" $FILTERED
	done
}
GetMapsNoHdr(){
	TGT=$1
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."volume-view"[]?
	 | $T + ",\t"
	 + ."volume-serial" + ",\t"
	 + ."volume-view-mappings"[].identifier + ",\t"
	 + ."volume-name" + ",\t"
	 + ."volume-view-mappings"[].access + ",\t"
	 + ."volume-view-mappings"[].ports + ",\t"
	 + ."volume-view-mappings"[].nickname  + ",\t" + ."volume-view-mappings"[].lun
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowMapsJSON $TGT | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetMaps(){
	TGT=$1
	printf "\nRUN: ${FUNCNAME[0]}\n"
	HDR00="controller,\t"
	HDR01="volume-serial,\t                        "
	HDR02="volume-identifier,                      "
	HDR03="volume-name,    "
	HDR04="access,         "
	HDR05="ports,          "
	HDR06="nickname,       "
	HDR07="lun"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetMapsNoHdr "$TGT"
	done
}
GetDisksNoHdr() {
	TGT="$1"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.drives[]? 
	 | $T + ",\t"
	 + ."durable-id" + ",\t" 
	 + ."disk-group" + ",\t"
	 + ."vendor" + ",\t"
	 + ."model" + ",\t"
	 + ."serial-number" + ",\t"
	 + (."blocksize"|tostring) + ",\t"
	 + ."size" + ",\t"
	 + ."temperature" + ",\t"
	 + ."health"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowDisksJSON "${TGT}" | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetDisks() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	HDR00="controller,\t"
	HDR01="name,\t\t"
	HDR02="dgroup,\t"
	HDR03="vendor,\t\t"
	HDR04="model,\t\t"
	HDR05="serial,\t\t"
	HDR06="    blocksize,\t"
	HDR07="size,\t"
	HDR08="temperature,\t"
	HDR09="health"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}${HDR08}${HDR09}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetDisksNoHdr "$TGT"
	done
}
GetDiskGroupsNoHdr() {
	TGT="$1"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."disk-groups"[]?
	 | $T + ",\t"
	 + .name +",\t" + .size + ",\t"
	 + ."storage-type" + ",\t\t"
	 + .raidtype + ",\t\t"
	 + (."diskcount"|tostring)
	 + ",\t\t" + .owner + ",\t"
	 + ."serial-number"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowDiskGroupsJSON "${TGT}" | jq --arg T "$TGT" -r "${JQ}")
	printf "${RESULT}\n"
}
GetDiskGroups() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	HDR00="controller,\t"
	HDR01="name,   "
	HDR02="size,           "
	HDR03="storage-type,\t"
	HDR04="raid-type,\t"
	HDR05="disk-count,\t"
	HDR06="owner,\t"
	HDR07="serial-number"
	HDR="${HDR00}${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		GetDiskGroupsNoHdr "$TGT"
	done
}
GetDisksInDiskGroups() {
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR="controller,\tdisk-group,\tdisks"
	printf "${HDR}\n"
	for TGT in "${TARGETS[@]}"
	do
		SHOWDISK=$(ShowDisksJSON $TGT)
		for DG in $(echo $SHOWDISK | jq -r '.drives[]? | ."disk-group"' | sort -u)
		do
			printf "$TGT,\t$DG\t"
			printf "$SHOWDISK\n"  \
			| jq -r '.drives[]? | ."disk-group" + " " + ."location" ' \
			| grep $DG | awk -F ' ' '{print $2}' | tr '\n' ',' | sed 's/,$//g' ; printf "\n"
		done
	done
	
}
RemoveDiskGroup() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	DG=$2
	CMD="remove disk-groups $DG"
	DoCmd ${TGT} ${CMD} | jq -r '.status[]."response-type"'
}
RemoveAllDiskGroups() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	for DG in $(GetDiskGroups $TGT)
	do
		RemoveDiskGroup $TGT $DG
	done
}

CreateDiskGroups() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	CMD="add disk-group"
	CMD="${CMD} type linear level adapt stripe-width 16+2 spare-capacity 20.0TiB interleaved-volume-count 1"
	POOL1="assigned-to a disks 0.0-11,0.24-35,0.48-59,0.72-83,0.96-100 dg01"
	POOL2="assigned-to b disks 0.12-23,0.36-47,0.60-71,0.84-95,0.101-105 dg02"
	DoCmd ${TGT} ${CMD} ${POOL1} >/dev/null #don't care about the output
	DoCmd ${TGT} ${CMD} ${POOL2} >/dev/null #don't care about the output
}
CreateFourLun8plus2ADAPT() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	CMD="add disk-group"
	CMD="${CMD} type linear level adapt stripe-width 8+2 spare-capacity 10.0TiB interleaved-volume-count 1 interleaved-basename V "
	DGS=("assigned-to a disks 0.0-11   dg01"\
	     "assigned-to a disks 0.12-23  dg02"\
	     "assigned-to a disks 0.24-35  dg03"\
	     "assigned-to a disks 0.36-47  dg04"\
	     "assigned-to b disks 0.53-64  dg05"\
	     "assigned-to b disks 0.65-76  dg06"\
	     "assigned-to b disks 0.77-88  dg07"\
	     "assigned-to b disks 0.89-100 dg08")
	for DG in "${DGS[@]}"
	do
		#echo ${TGT} ${CMD} ${DG}
		DoCmd ${TGT} ${CMD} ${DG} >/dev/null #don't care about the output
	done
}


GetPowerReadings() {
	TGT=$1
	printf "\nRUN: ${FUNCNAME[0]}\n"
	printf "controller,\tL1_VOLT,\tL1_AMP,\tL1_WATT,\tL2_VOLT,\tL2_AMP,\tL2_WATT,\tTotalWatts\n"
	for TGT in ${TARGETS[*]}; do
		RESULT=$(ShowSensorStatusJSON $TGT)
		LVOLT1=$(echo ${RESULT} | jq -r '.sensors[]? | select(."durable-id" == "sensor_volt_psu_0.0.1").value')
		LVOLT2=$(echo ${RESULT} | jq -r '.sensors[]? | select(."durable-id" == "sensor_volt_psu_0.1.1").value')
		LCURR1=$(echo ${RESULT} | jq -r '.sensors[]? | select(."durable-id" == "sensor_curr_psu_0.0.1").value')
		LCURR2=$(echo ${RESULT} | jq -r '.sensors[]? | select(."durable-id" == "sensor_curr_psu_0.1.1").value')
		LWATT1=$(echo "scale=2; $LVOLT1 * $LCURR1" | bc -l)
		LWATT2=$(echo "scale=2; $LVOLT2 * $LCURR2" | bc -l)
		LWATT_TOTAL=$(echo "scale=2; $LWATT1 + $LWATT2" | bc -l)
		printf "$TGT,\t$LVOLT1,\t\t$LCURR1,\t$LWATT1,\t\t$LVOLT2,\t\t$LCURR2,\t$LWATT2,\t\t$LWATT_TOTAL\n"
	done
}
GetEcliKeyData() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	ShowConfigurationJSON $TGT | \
	 jq -r '(.versions[]? | ."object-name" + "   SC_Version: " + ."sc-fw" + "   MC_Version: " +."mc-fw"),(.controllers[]? | ."durable-id" + "_internal_serial_number: " + ."internal-serial-number")'
}
ProvisionSystem() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	RemoveAllDiskGroups $TGT
	CreateDiskGroups $TGT
}
CreateEightPlus2Adapt() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	for TGT in ${TARGETS[*]}; do
		RemoveDiskGroups "${TGT}" &
	done
	wait
	for TGT in ${TARGETS[*]}; do
		CreateDiskGroups "${TGT}" &
	done
	wait
	for TGT in ${TARGETS[*]}; do
		echo "DiskGroups on $TGT"
		GetDiskGroups "${TGT}"
	done
}
Provision8plus24lun() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	for TGT in "${TARGETS[@]}"
	do
		RemoveAllDiskGroups "$TGT"
	done
	wait
	for TGT in "${TARGETS[@]}"
	do
		CreateFourLun8plus2ADAPT "$TGT"
	done
	wait
}
for CMD in ShowMpt3SasHBAs GetInquiry GetPowerReadings GetVolumes GetInitiators GetMaps GetDiskGroups GetDisksInDiskGroups GetHostPhyStatistics GetDisks
do
	$CMD
done
