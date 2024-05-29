# Script that copies the SignerID value from build/coordinator-config.json to the SignerID field in the manifest.json and era-config.json files.

rm -r marblerun-coordinator-data/*

awk -F'"' '/"SignerID"/ {print $4}' build/coordinator-config.json > build/SignerID.txt

sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" manifest.json

sed -i "s/\"SignerID\": \".*\"/\"SignerID\": \"$(cat build/SignerID.txt)\"/g" era-config.json