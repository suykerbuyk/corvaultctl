#!/bin/bash
set -e
POOL_NAME=cvt
LUN_PATTERN="/dev/disk/by-id/scsi-SSEAGATE_6575_00c0f*"
LOGDIR=test

if [ ! -d ${LOGDIR} ]; then
	mkdir ${LOGDIR}
fi


echo 1000 >/sys/module/zfs/parameters/zfs_multihost_fail_intervals

wipe_zfs_pools() {
	for POOL in $(zpool list -H | grep $POOL_NAME | awk '{print $1}')
	do
		echo "zpool destroy $POOL"
		zpool destroy $POOL &
	done
	wait
	# clear out the primary ZFS partion
	for X in $(ls ${LUN_PATTERN} | grep part1)
	do
		echo "zpool labelclear -f ${X} && wipefs -a ${X}"
		zpool labelclear -f ${X} && wipefs -a ${X} &
	done
	wait

	# clear out the residual end of disk partition
	for X in $(ls ${LUN_PATTERN} | grep part9)
	do
		echo "wipefs -a ${X}"
		wipefs -a ${X} &
	done
	wait

	# Finally, wipe out the GPT partions.
	for X in $(ls ${LUN_PATTERN} | grep -v part)
	do
		echo "sgdisk -Z ${X}"
		sgdisk -Z ${X} &
	done
	wait
	sleep 1
}

rescan_scsi_bus() {
	# Rescan the SCSI bus
	for X in $(ls /sys/class/scsi_host/)
	do
		echo "- - -" > /sys/class/scsi_host/$X/scan &
	done
	wait
}



create_draid_zpool() {
	wipe_zfs_pool
	zpool create ${POOL_NAME}  -O recordsize=512K -O atime=off -O dnodesize=auto -o ashift=12 draid2:4d:6c:0s  ${LUN_PATTERN}
}
create_raidz2_zpool() {
	wipe_zfs_pool
	zpool create ${POOL_NAME}  -O recordsize=512K -O atime=off -O dnodesize=auto -o ashift=12 raidz2  ${LUN_PATTERN}
}
create_individual_pools() {
	wipe_zfs_pools
	IDX=0
	for X in $(ls ${LUN_PATTERN} | grep -v part)
	do
		zpool create ${POOL_NAME}_${IDX}  -O recordsize=512K -O atime=off -O dnodesize=auto -o ashift=12 $X
		IDX=$((IDX+1))
	done
}


wipe_zfs_pools
rescan_scsi_bus
create_individual_pools
#clear all dmesg history
dmesg -c
#for IOENGINE in libaio io_uring; do
for IOENGINE in libaio ; do
	for IODEPTH in 1 8 16 32; do
		for JOBS in 1 4 8 16 32; do
			#for PAT in 'write' 'read' 'randrw' 'randread' 'randwrite'; do
			for PAT in 'write'  'randrw' 'randwrite'; do
				#for BLK in 32k 128k 256k 512k 1024k; do
				for BLK in 1024k 8192k 32768k; do
					for POOL in $(zpool list -H | grep $POOL_NAME | awk '{print $1}')
					do
						TEST="${POOL}-${IOENGINE}-${IODEPTH}-${PAT}-${BLK}-${JOBS}.fio.json"
						echo "Running $TEST"
						zpool clear ${POOL}
						fio --directory=/${POOL} \
						    --name="${TEST}" \
						    --rw=$PAT \
						    --group_reporting=1 \
						    --bs=$BLK \
						    --direct=1 \
						    --numjobs=$JOBS \
						    --time_based=1 \
						    --runtime=30 \
						    --iodepth=$IODEPTH \
						    --ioengine=$IOENGINE \
						    --size=128G \
						    --output-format=json | tee "$PWD/${LOGDIR}/${TEST}" && \
						echo "Completed" &
					done
					wait
				done
			done
		done
	done
done
