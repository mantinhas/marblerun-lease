#!/bin/bash

if [[ "${DCAP_LIBRARY}" == "intel" ]]
then
    # rename the library installed by az-dcap-client
    mv /usr/lib/libdcap_quoteprov.so /usr/lib/libdcap_quoteprov.so.azure
    # create a link to the intel quote provider library
    ln -s /usr/lib/x86_64-linux-gnu/libdcap_quoteprov.so.intel /usr/lib/x86_64-linux-gnu/libdcap_quoteprov.so
fi

exec erthost coordinator-enclave.signed
