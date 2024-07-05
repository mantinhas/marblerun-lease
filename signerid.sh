#All functions to be run in separate processes, starting from .../marblerun-modified

# Script that copies the SignerID value from build/coordinator-config.json to the SignerID field in the manifest.json and era-config.json files.
build() {
	echo "[BUILD] Starting build." 
	. /opt/edgelessrt/share/openenclave/openenclaverc
	cd build && make && cd ..
	cp build/premain-libos test-api/premain-libos
	cp build/premain-libos test-api/premain-libos.my
	cp build/premain-libos samples/gramine-hello/premain-libos
	cp build/premain-libos samples/gramine-hello/premain-libos.my
	echo "[BUILD] Finished build."
}

cleaningAndSignerID() {
	echo "[CLEAN] Starting clean." 
	rm -r marblerun-coordinator-data/*
	awk -F'"' '/"SignerID"/ {print $4}' build/coordinator-config.json > build/SignerID.txt
	sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" manifest.json
	sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" test-api/manifest.json
	sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" samples/gramine-hello/manifest.json
	sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" era-config.json
	echo "[CLEAN] Finished clean."
}

gramineMake(){
	echo "[GRAMINE] Starting gramine." 
	if [[ -n $1 ]]; then
		echo Detected path $1
		cd $1 || exit
	else
		echo Using test-api
		cd test-api
	fi
	make clean && make SGX=1 | tee /tmp/measurement.txt 
	if [[ $( grep "Measurement:" /tmp/measurement.txt) ]]; then 
		measure=$(grep -A1 "Measurement:" /tmp/measurement.txt | tail -1 | tr -d " ")
		sed -i "s/\"UniqueID\": \".*\"/\"UniqueID\": \"${measure}\"/g" manifest.json
	else 
		echo "Error in build" && cat /tmp/measurement.txt
	fi
	rm /tmp/measurement.txt
	cd ..
	echo "[GRAMINE] Finished gramine."
}

[[ $1 == "help" ]] && echo -e "b for build\nc for clean\ng for make gramine\n" && exit
[[ $( echo $1 | grep "b") ]] && build && cleaningAndSignerID && gramineMake $2
[[ $( echo $1 | grep "c") ]] && cleaningAndSignerID
[[ $( echo $1 | grep "g") ]] && gramineMake $2
