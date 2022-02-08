#!/bin/sh


GetMpt3SasHBAs() {
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
		printf "CTRLR_PATH: $CTRLR_PATH\n"
		printf "   PCI Vendor/Device: $PCI_VENDOR, $PCI_SUBSYSTEM_DEVICE\n"
		printf "   UNIQUE_ID: $UNIQUE_ID\n"
		printf "   SAS_ADDR: $SAS_ADDR\n"
		printf "   PCI_ADDR: $PCI_ADDR\n"
		printf "   BOARD_NAME: $BOARD_NAME\n"
		printf "   BOARD_ASSEMBLY: $BOARD_ASSEMBLY\n"
		printf "   BIOS VERSION: $VERSION_BIOS\n"
		printf "   FW VERSION: $VERSION_FW\n"
		printf "   MPI VERSION: $VERSION_MPI\n"
		printf "   NVDATA VERSION: $VERSION_NVDATA\n"
		printf "   PRODUCT VERSION: $VERSION_PRODUCT\n"
	done
}
GetMpt3SasHBAsJSON() {
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
		printf "  \"version-nvdata\": \"$VERSION_NVDATA\"\n"
		PREFIX="  },\n{\n"
	done
	printf "  }\n"
	printf "]}\n"
}
GetMpt3SasHBAsJSON
