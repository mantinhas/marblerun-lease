#!/usr/bin/env bash
[[ $1 == "coord" ]] && {
	. /opt/edgelessrt/share/openenclave/openenclaverc
	bash signerid.sh c
	erthost build/coordinator-enclave.signed ;
}

[[ $1 == "lease" ]] && ./build/marblerun server lease localhost:4433 -c provider_certificate.crt -k provider_private.key

[[ $1 == "app" ]] && {
	cd test-api
	curl -k --data-binary @manifest.json https://localhost:4433/manifest
	. /opt/edgelessrt/share/openenclave/openenclaverc
	EDG_MARBLE_TYPE=hello EDG_MARBLE_UUID_FILE=server_uuid EDG_TEST_ADDR=localhost:8001 OE_SIMULATION=0 gramine-sgx python ;
}
