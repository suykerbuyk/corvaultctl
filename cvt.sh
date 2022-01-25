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
elif [[ $(which jo 2>&1>/dev/null) ]]; then
       echo "Please intall jo"
       exit 1
elif [[ $(which sshpass 2>&1>/dev/null) ]]; then
       echo "Please intall sshpass"
       exit 1
fi

# interesting sysfs paths for coorelating LUNs to host HBA ports and kdevs
#cat /sys/devices/*/*/*/host*/phy-*/sas_phy/*/sas_address | sort -u | cut -c 15-
#cat /sys/devices/pci*/*/*/host*/port*/end_device*/target*/*/sas_address | sort -u

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
GetVolumes() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR01="volume-name,\t"
	HDR02="name,\t"
	HDR03="size,\t\t"
	HDR04="serial-number,\t\t\t\t"
	HDR05="wwn-number,\t\t\t\t"
	HDR06="creation-date-time"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.volumes[]?
	 | ."volume-name" +",\t"
	 + ."virtual-disk-name" + ",\t"
	 + ."size" + ",\t"
	 + ."serial-number" + ",\t" + (."wwn" | ascii_downcase) + ",\t"
	 + ."creation-date-time"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowVolumesJSON $TGT | jq -r "${JQ}")
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetInitiators() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR01="durable-id,\t"
	HDR02="id,\t\t\t"
	HDR03="host-id,\t\t\t\t"
	HDR04="nickname"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	.initiator[] | select(.mapped == "Yes")
	 | "\t" + ."durable-id" + ",\t"
	 +  .id + ",\t"
	 + ."host-id" + ",\t"
	 + .nickname
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowInitiatorsJSON $TGT | jq -r "${JQ}")
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetHostPhyStatistics() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]} (filtered for non-zero counters)\n"
	HDR01="port,\t"
	HDR02="disparities,\t"
	HDR03="lost-dws,\t"
	HDR04="invalid-dws,\t"
	HDR05="reset-errs"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."sas-host-phy-statistics"[]
	 |  select((((."disparity-errors" != "00000000")
	 or ."lost-dwords" != "00000000")
	 or ."invalid-dwords" != "00000000")
	 or ."reset-error-counter" != "00000000")
	 | .port + "-" + (.phy|tostring) + ",\t"
	 + ."disparity-errors" +",\t"
	 + ."lost-dwords" + ",\t"
	 + ."invalid-dwords" + ",\t"
	 + ."reset-error-counter"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowHostPhyStatisticsJSON $TGT | jq -r "${JQ}")
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetMaps(){
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR01="volume-serial,\t                        "
	HDR02="volume-identifier,                      "
	HDR03="volume-name,    "
	HDR04="access,         "
	HDR05="ports,          "
	HDR06="nickname,       "
	HDR07="lun"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."volume-view"[]?
	 | ."volume-serial" + ",\t"
	 + ."volume-view-mappings"[].identifier + ",\t"
	 + ."volume-name" + ",\t"
	 + ."volume-view-mappings"[].access + ",\t"
	 + ."volume-view-mappings"[].ports + ",\t"
	 + ."volume-view-mappings"[].nickname  + ",\t" + ."volume-view-mappings"[].lun
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowMapsJSON $TGT | jq -r "${JQ}")
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetDiskGroups() {
	TGT="$1"
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	HDR01="name,   "
	HDR02="size,           "
	HDR03="storage-type,\t"
	HDR04="raid-type,\t"
	HDR05="disk-count,\t"
	HDR06="owner,\t"
	HDR07="serial-number"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	JQ=$(cat <<"EOF" | tr -d '\n\r\t'
	."disk-groups"[]?
	 | .name +",\t" + .size + ",\t"
	 + ."storage-type" + ",\t\t"
	 + .raidtype + ",\t\t"
	 + (."diskcount"|tostring)
	 + ",\t\t" + .owner + ",\t"
	 + ."serial-number"
EOF
)
	[[ $DBG != 0 ]] && printf "JQ : %s\n" "${JQ}" 1>&2
	RESULT=$(ShowDiskGroupsJSON "${TGT}" | jq -r "${JQ}")
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetDiskByGroup() {
	TGT=$1
	DG=$2
	printf "\nRUN: $TGT ${FUNCNAME[0]} $DG\n"
	ShowDisksJSON $TGT | jq -r '.drives[]?  | ."disk-group" + " " +  ."location" ' \
		| grep $DG | awk -F ' ' '{print $2}' | tr  '\n' ','
}
GetAllDiskInAllGroups() {
	TGT=$1
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	for DG in $(ShowDiskGroupsJSON "${TGT}" | jq -r '."disk-groups"[]? | .name')
	do
		DISKS=$(ShowDisksJSON $TGT | jq -r '.drives[]?  | ."disk-group" + " " +  ."location" ' \
			| grep $DG | awk -F ' ' '{print $2}' | tr '\n' ',' | sed 's/,$//g')
		echo "$TGT $DG: $DISKS"
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
	printf "\nRUN: $TGT ${FUNCNAME[0]}\n"
	ShowSensorStatusJSON $TGT  | jq -r '."sensors"[]? | ."sensor-name" + " " + ."value" '  | grep "Input Rail" | grep -i 'volt\|current'
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

RunCmdOnAllTargets() {
	for TGT in "${TARGETS[@]}"
	do
		$1 $TGT
	done
}


for CMD in GetDiskGroups GetPowerReadings GetAllDiskInAllGroups GetMaps GetInitiators GetVolumes GetHostPhyStatistics
do
	RunCmdOnAllTargets $CMD
done

