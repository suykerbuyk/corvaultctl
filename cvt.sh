#!/bin/sh

set -e
export SSHPASS='Testit123!'

USER='manage'
TARGETS=("corvault-1a" "corvault-2a" "corvault-3b")

BASE_CMD='set cli-parameters json; '
#BASE_CMD='set cli-parameters wbi pager off ; '
#REQ=""
export RESP=""
#JSON=""
#STAT=""
# cat /sys/devices/*/*/*/host*/phy-*/sas_phy/*/sas_address | sort -u | cut -c 15-
# cat /sys/devices/pci*/*/*/host*/port*/end_device*/target*/*/sas_address | sort -u
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
# sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show disk-groups' >corvault1.disk-groups.json
# sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show volumes' >corvault1.volumes.json
# sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show initiators' >corvault1.initiators.json
# sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show maps' >corvault1.maps.json
DoSSH() {
	sshpass -e $@
}
DoCmd() {
	TGT="${1}"
	shift
	REPLY_FILE="${TGT}.json"
	echo "TGT: $TGT  CMD: $BASE_CMD $@" 1>&2
	SSHSOCKET=/tmp/$TGT.ssh.socket
	REPLY=$(DoSSH "ssh -o ControlPath=$SSHSOCKET -o ControlMaster=auto -o ControlPersist=10m ${USER}@${TGT} ${BASE_CMD} $@")
	# Pull off the commented lines that contain the commands sent to the target
	#printf "REPLY: %s\n" "$REPLY" 1>&2
	REQ=$(echo "$REPLY" | egrep '^#.*' | sed -e 's/^#[ ]*//g' -e '/^$/d' | sed -e :a -e '$!N; s/\n/; /; ta')
	#printf "REQ: $REQ\n" 1>&2
	JSON=$(printf "%s\n" "$REPLY" | awk '/#  /,0' | egrep -v '^# .*' |  sed -e :a -e '$!N;  ta')
	#printf "JSON: %s\n" "$JSON" 1>&2
	RESP=$(echo ${JSON} | jq -r '.status[].response')
	STAT=$(echo ${JSON} | jq -r '.status[]."response-type"')
	#printf "RESP: %s\n" "$RESP" 1>&2
	#printf "STAT: %s\n" "$STAT" 1>&2
	if [ "${STAT}" != "Success" ] ; then
		echo "${REPLY}" >"${REPLY_FILE}"
		echo "Error: $BASE_CMD $@" 1>&2;
		echo "Status: ${STAT}" 1>&2;
		echo "Response: ${RESP}" 1>&2;
		echo "See ${REPLY} for full JSON return data" 1>&2;
		exit 1
	fi
	echo "Status: ${STAT}" 1>&2
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
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	ShowVolumesJSON $1 | jq -r '.volumes[] | ."volume-name" +",\t" + ."virtual-disk-name" + ",\t" + ."size" + ",\t" + ."serial-number" + ",\t" + ."wwn" + ",\t" + ."creation-date-time"'
}
GetInitiators() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	HDR01="durable-id,\t     "
	HDR02="id,               "
	HDR03="host-id,          "
	HDR04="nickname,         "
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}"
	RESULT=$(ShowInitiatorsJSON $TGT | jq -r '.initiator[] | select(.mapped == "Yes") | ."durable-id" + ",\t" +  .id + ",\t" + ."host-id" + ",\t" + .nickname')
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetMaps(){
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	HDR01="volume-serial,\t                        "
	HDR02="volume-identifier,                      "
	HDR03="volume-name,    "
	HDR04="access,         "
	HDR05="ports,          "
	HDR06="nickname,       "
	HDR07="lun"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	RESULT=$(ShowMapsJSON $TGT | jq -r '."volume-view"[]? | ."volume-serial" + ",\t" + ."volume-view-mappings"[].identifier + ",\t" + ."volume-name" + ",\t" + ."volume-view-mappings"[].access + ",\t" + ."volume-view-mappings"[].ports + ",\t" + ."volume-view-mappings"[].nickname  + ",\t" + ."volume-view-mappings"[].lun' )
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetDiskGroups() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT="$1"
	HDR01="name,   "
	HDR02="size,           "
	HDR03="storage-type,\t"
	HDR04="raid-type,\t"
	HDR05="disk-count,\t"
	HDR06="owner,\t"
	HDR07="serial-number"
	HDR="${HDR01}${HDR02}${HDR03}${HDR04}${HDR05}${HDR06}${HDR07}"
	RESULT=$(ShowDiskGroupsJSON "${TGT}" | jq -r '."disk-groups"[]? | .name +",\t" + .size + ",\t" + ."storage-type" + ",\t\t" + .raidtype + ",\t\t" + (."diskcount"|tostring) + ",\t\t" + .owner + ",\t" + ."serial-number" ')
	printf "${HDR}\n"
	printf "${RESULT}\n"
}
GetDiskByGroup() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	DG=$2
	printf "\nRUN: ${FUNCNAME[0]} $TGT $DG\n"
	ShowDisksJSON $TGT | jq -r '.drives[]?  | ."disk-group" + " " +  ."location" ' \
		| grep $DG | awk -F ' ' '{print $2}' | tr  '\n' ','
}
GetAllDiskInAllGroups() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	for DG in $(ShowDiskGroupsJSON "${TGT}" | jq -r '."disk-groups"[]? | .name')
	do
		DISKS=$(ShowDisksJSON $TGT | jq -r '.drives[]?  | ."disk-group" + " " +  ."location" ' \
			| grep $DG | awk -F ' ' '{print $2}' | tr '\n' ',' | sed 's/,$//g')
		echo "$TGT $DG: $DISKS"
	done
}
RemoveDiskGroup() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	DG=$2
	CMD="remove disk-groups $DG"
	DoCmd ${TGT} ${CMD} | jq -r '.status[]."response-type"'
}
RemoveAllDiskGroups() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	for DG in $(GetDiskGroups $TGT)
	do
		RemoveDiskGroup $TGT $DG
	done
}

CreateDiskGroups() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	CMD="add disk-group"
	CMD="${CMD} type linear level adapt stripe-width 16+2 spare-capacity 20.0TiB interleaved-volume-count 1"
	POOL1="assigned-to a disks 0.0-11,0.24-35,0.48-59,0.72-83,0.96-100 dg01"
	POOL2="assigned-to b disks 0.12-23,0.36-47,0.60-71,0.84-95,0.101-105 dg02"
	DoCmd ${TGT} ${CMD} ${POOL1} >/dev/null #don't care about the output
	DoCmd ${TGT} ${CMD} ${POOL2} >/dev/null #don't care about the output
}
CreateFourLun8plus2ADAPT() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
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
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	ShowSensorStatusJSON $TGT  | jq -r '."sensors"[]? | ."sensor-name" + " " + ."value" '  | grep "Input Rail" | grep -i 'volt\|current'
}
GetEcliKeyData() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
	ShowConfigurationJSON $TGT | \
	 jq -r '(.versions[]? | ."object-name" + "   SC_Version: " + ."sc-fw" + "   MC_Version: " +."mc-fw"),(.controllers[]? | ."durable-id" + "_internal_serial_number: " + ."internal-serial-number")'
}
ProvisionSystem() {
	printf "\nRUN: ${FUNCNAME[0]}\n"
	TGT=$1
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


for CMD in GetDiskGroups GetPowerReadings GetAllDiskInAllGroups GetMaps GetInitiators GetVolumes
do
	RunCmdOnAllTargets $CMD
done
