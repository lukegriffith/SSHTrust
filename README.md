
# SSHTrust

## SSH CA Key Signing Server

This project is a simple rest API that generates SSH Certificate Authorities (CA), signs SSH public keys using those CAs.


## Project Overview

The project consists of:
1. An HTTP server that:
   - Creates in-memory CAs (Certificate Authorities) and stores them.
   - Signs SSH public keys using the in-memory CA.
   - Provides CA public keys for use in external systems (like the SSH server).
2. A client to interact with the server
3. Docker test suite to demonstrate the use case. 

## Roadmap

1. Support external storage backend
   - sql lite
   - postgres
   - etcd

2. Implement user authentication & acl

## Installation

Project is still in early development, and requires golang to be installed to install the system

```bash
go install github.com/lukegriffith/SSHTrust
```

## API Documentation
Swagger UI is enabled for this project. You can access it by navigating to the below link when the server is active locally:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The API developer docs are located here as well [api docs](API.md)

## Test SSH Server

This project also includes an OpenSSH server running in a Docker container, which is configured to trust certificates signed by the CA. The server is accessible via SSH and uses certificate-based authentication. This is an example container to demonstrate the configuration working.

### Steps to Build and Run the Test Docker SSH Server

1. **Make the server binary**: Using the golang compiler and the make file, create the sshtrust binary
   ```
   make
   ```

2. **Create a CA for Docker build**: Using the server, create a new CA under the project directory for the server to copy
   ```
   ./sshtrust serve &
   ./sshtrust ca new -n myca -p testuser
   ./sshtrust ca get myca | jq .public_key -r > ssh_ca.pub
   ```

3. **Build the Docker Image**:
   ```
   docker build -t ssh-server -f ssh-test-server.Dockerfile .
   ```

4. **Run the Docker SSH Server**:
   ```
   docker run -d -p 2222:22 --name my-ssh-server ssh-server
   ```

### Logging into the Test SSH Server

Once the SSH server is running, configured with a CA from the server. you can SSH into it using a certificate signed by the CA:

1. **Sign your Public Key**: 
   ```
   ./sshtrust sign -n myca --ttl 30 -p testuser -k "$(cat ~/.ssh/id_ed25519.pub)" > ~/.ssh/id_ed25519-cert.pub
   ```

2. **SSH into the Server**:
   ```
   ssh -i ~/.ssh/id_ed25519 -o CertificateFile=~/.ssh/id_ed25519-cert.pub -p 2222 testuser@localhost
   ```

### SSH Server Setup Recap:
- **Public Key**: The CA’s public key (`ssh_ca.pub`) is copied to the SSH server and used to validate certificates.
- **Docker**: The SSH server runs inside a Docker container and listens on port 2222.
