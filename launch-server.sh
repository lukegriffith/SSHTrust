#!/bin/bash
set -xe; 


echo "NO AUTH TEST"
./sshtrust serve --no-auth & 
PID=$(echo $!); 
echo "Server started with PID: $$PID"  # This will correctly print the PID.
trap "echo Killing process $PID; kill -9 $PID" EXIT; 
./sshtrust ca new -n myca -p testuser -t ssh-ed25519; 
./sshtrust sign -n myca --ttl 4 -p testuser -k "$(cat ~/.ssh/id_ed25519.pub )"; 
./sshtrust ca list; 

# clean up
echo "killing process $PID"
kill -9 $PID

echo "AUTH TEST"
./sshtrust serve & 
PID=$(echo $!); 
echo "Server started with PID: $$PID"  # This will correctly print the PID.
trap "echo Killing process $PID; kill -9 $PID" EXIT; 

echo "1234" | ./sshtrust register --stdin -u test
echo "1234" | ./sshtrust login --stdin -u test
./sshtrust ca new -n myca -p testuser -t ssh-ed25519; 
./sshtrust sign -n myca --ttl 4 -p testuser -k "$(cat ~/.ssh/id_ed25519.pub )"; 
./sshtrust ca list; 
