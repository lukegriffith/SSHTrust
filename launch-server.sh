#!/bin/bash
set -e; 

./sshtrust serve --no-auth & 
PID=$(echo $!); 
echo "Server started with PID: $$PID"  # This will correctly print the PID.
trap "echo Killing process $PID; kill -9 $PID" EXIT; 
./sshtrust ca new -n myca -p testuser -t ssh-ed25519; 
./sshtrust sign -n myca --ttl 4 -p testuser -k "$(cat ~/.ssh/id_ed25519.pub )"; 
./sshtrust ca list; 

