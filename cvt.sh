#!/bin/sh

set -e
export SSHPASS='Testit123!'

USER='manage'
TARGETS=("corvault-1a" "corvault-2a" "corvault-3a")

BASE_CMD='set cli-parameters json; '
#BASE_CMD='set cli-parameters wbi pager off ; '
#REQ=""
#RESP=""
JSON=""
#STAT=""

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
sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show disk-groups' >corvault1.disk-groups.json
sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show volumes' >corvault1.volumes.json
sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show initiators' >corvault1.initiators.json
sshpass -e ssh manage@corvault-1a 'set cli-parameters json; show maps' >corvault1.maps.json

DoCmd() {
	TGT=$1
	shift
	REPLY="${TGT}.json"
	echo "TGT: $TGT  CMD: $BASE_CMD $@" 1>&2
	sshpass -e ssh ${USER}@${TGT} "$BASE_CMD $@" >"${REPLY}"
	# Pull off the commented lines that contain the commands sent to the target
	REQ=$(cat "${REPLY}" | egrep '^#.*' |\
		sed -e 's/^#[ ]*//g' -e '/^$/d' | sed -e :a -e '$!N; s/\n/; /; ta')
	JSON=$(cat ${REPLY} | awk '/#   /,0' | egrep -v '^# .*' |  sed -e :a -e '$!N;  ta')
	RESP=$(echo ${JSON} | jq -r '.status[].response')
	STAT=$(echo ${JSON} | jq -r '.status[]."response-type"')
	if [ "${STAT}" != "Success" ] ; then
		echo "Error: $BASE_CMD $@" 1>&2;
		echo "Status: ${STAT}" 1>&2;
		echo "Response: ${RESP}" 1>&2;
		echo "See ${REPLY} for full JSON return data" 1>&2;
		exit 1
	fi
	echo "Status: ${STAT}" 1>&2
	echo $JSON
}

ShowDiskGroupsJSON() {
	TGT="$1"
	CMD="show disk-groups"
	DoCmd ${TGT} "${CMD}"
	#echo "JSON: ${JSON}"
	#echo "RESP: ${RESP}"
	#echo "STAT: ${STAT}"
}
GetDiskGroups() {
	TGT="$1"
	ShowDiskGroupsJSON "${TGT}" | jq -r '."disk-groups"[]? | .name' | tee /dev/stderr
}
ShowDisksJSON() {
	TGT=$1
	CMD="show disks"
	DoCmd ${TGT} "${CMD}"
}
ShowSensorStatusJSON() {
	TGT=$1
	CMD="show sensor-status"
	DoCmd ${TGT} "${CMD}"
}
GetDiskByGroup() {
	TGT=$1
	DG=$2
	ShowDisksJSON $TGT | jq -r '.drives[]?  | ."disk-group" + " " +  ."location" ' \
		| grep $DG | awk -F ' ' '{print $2}' | tr  '\n' ','
}
GetAllDiskInAllGroups() {
	TGT=$1
	for DG in $(GetDiskGroups "${TGT}")
	do
		DISKS=$(GetDiskByGroup $TGT $DG | sed 's/,$//g')
		echo "$DG: $DISKS"
	done
}
RemoveDiskGroup() {
	TGT=$1
	DG=$2
	CMD="remove disk-groups $DG"
	DoCmd ${TGT} ${CMD} | jq -r '.status[]."response-type"'
}
RemoveAllDiskGroups() {
	TGT=$1
	for DG in $(GetDiskGroups $TGT)
	do
		RemoveDiskGroup $TGT $DG
	done
}

CreateDiskGroups() {
	TGT=$1
	#CMD="${BASE_CMD} add disk-group"
	CMD="add disk-group"
	CMD="${CMD} type linear level adapt stripe-width 16+2 spare-capacity 20.0TiB interleaved-volume-count 1"
	POOL1="assigned-to a disks 0.0-11,0.24-35,0.48-59,0.72-83,0.96-100 dg01"
	POOL2="assigned-to b disks 0.12-23,0.36-47,0.60-71,0.84-95,0.101-105 dg02"
	DoCmd ${TGT} ${CMD} ${POOL1} >/dev/null #don't care about the output
	DoCmd ${TGT} ${CMD} ${POOL2} >/dev/null #don't care about the output
}

GetPowerReadings() {
	TGT=$1
	ShowSensorStatusJSON $TGT  | jq -r '."sensors"[]? | ."sensor-name" + " " + ."value" '  | grep "Input Rail" | grep -i 'volt\|current'
}
ProvisionSystem() {
	TGT=$1
	RemoveAllDiskGroups $TGT
	CreateDiskGroups $TGT
}


CreateEightPlus2Adapt() {
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
#GetDiskGroups "corvault-1a"
#ShowDisksJSON "corvault-1a" | jq -r .drives[].size | sort -u | wc -l
#ShowDiskGroupsJSON "corvault-1a"
#GetDiskByGroup "corvault-1a"
#GetAllDiskInAllGroups "corvault-1a"
#RemoveDiskGroup "corvault-1a" dg01
#CreateDiskGroups "corvault-1a"
#ProvisionSystem "corvault-1a"

GetPowerReadings "corvault-3b"
