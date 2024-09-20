
# SSHTrust

# SSH Key Server with CA Signing and Docker SSH Server

This project demonstrates how to create a simple HTTP API that generates SSH Certificate Authorities (CA), signs SSH public keys using those CAs, and runs an SSH server in Docker that validates certificates signed by the CA. The SSH server is configured to trust the CA’s public key and allow certificate-based login using SSH keys.

## Project Overview

The project consists of:
1. An HTTP server that:
   - Creates in-memory CAs (Certificate Authorities) and stores them.
   - Signs SSH public keys using the in-memory CA.
   - Provides CA public keys for use in external systems (like the SSH server).
2. A Docker-based SSH server configured to trust the in-memory CA's public key, enabling certificate-based SSH authentication.

## HTTP API Endpoints

### 1. Create a CA
- **URL**: `/create-ca?name=<CA_NAME>`
- **Method**: `GET`
- **Description**: This endpoint generates a new SSH Certificate Authority (CA) and stores it in memory under the given name.
- **Example**:
   ```
   curl "http://localhost:8080/create-ca?name=my-ca"
   ```

### 2. Get the CA Public Key
- **URL**: `/get-ca-public-key?name=<CA_NAME>`
- **Method**: `GET`
- **Description**: Retrieves the public key of the CA stored under the specified name, which can be used by other systems to verify certificates signed by the CA.
- **Example**:
   ```
   curl "http://localhost:8080/get-ca-public-key?name=my-ca" -o ssh_ca.pub
   ```

### 3. Sign a Public Key
- **URL**: `/sign-public-key?ca=<CA_NAME>`
- **Method**: `POST`
- **Description**: Signs a public key with the specified CA. The public key should be provided in the body of the request as a JSON object in the format `{"public_key": "<public_key>"}`. The API responds with the signed certificate.
- **Example**:
   ```
   curl -X POST -d "{\"public_key\":\"$(cat ~/.ssh/id_rsa.pub)\"}" \
   "http://localhost:8080/sign-public-key?ca=my-ca"
   ```

- To store the signed certificate in a file:
   ```
   curl -X POST -d "{\"public_key\":\"$(cat ~/.ssh/id_ed25519.pub)\"}" \
   "http://localhost:8080/sign-public-key?ca=my-ca" \
   -o ~/.ssh/id_ed25519-cert.pub
   ```

## Docker SSH Server

This project also includes an OpenSSH server running in a Docker container, which is configured to trust certificates signed by the CA. The server is accessible via SSH and uses certificate-based authentication.

### Steps to Build and Run the Docker SSH Server

1. **Create the Dockerfile**: Create an SSH server with the following Dockerfile (`ssh-test-server.Dockerfile`):

   ```
   FROM ubuntu:20.04
   RUN apt-get update && apt-get install -y openssh-server sudo
   RUN mkdir /var/run/sshd
   RUN useradd -rm -d /home/testuser -s /bin/bash -g root -G sudo -u 1000 testuser
   RUN echo 'testuser:testpassword' | chpasswd
   RUN mkdir -p /home/testuser/.ssh && chown -R testuser:root /home/testuser/.ssh && chmod 700 /home/testuser/.ssh
   COPY ssh_ca.pub /etc/ssh/ca.pub
   RUN echo "TrustedUserCAKeys /etc/ssh/ca.pub" >> /etc/ssh/sshd_config
   EXPOSE 22
   CMD ["/usr/sbin/sshd", "-D"]
   ```

2. **Build the Docker Image**:
   ```
   docker build -t ssh-server -f ssh-test-server.Dockerfile .
   ```

3. **Run the Docker SSH Server**:
   ```
   docker run -d -p 2222:22 --name my-ssh-server ssh-server
   ```

4. **Verify the SSH Server is Running**:
   ```
   docker ps
   ```

## Using the SSH Server

Once the SSH server is running, you can SSH into it using a certificate signed by the CA:

1. **Sign your Public Key**: 
   ```
   curl -X POST -d "{\"public_key\":\"$(cat ~/.ssh/id_ed25519.pub)\"}" \
   "http://localhost:8080/sign-public-key?ca=my-ca" \
   -o ~/.ssh/id_ed25519-cert.pub
   ```

2. **SSH into the Server**:
   ```
   ssh -i ~/.ssh/id_ed25519 -o CertificateFile=~/.ssh/id_ed25519-cert.pub -p 2222 testuser@localhost
   ```

### SSH Server Setup Recap:
- **Public Key**: The CA’s public key (`ssh_ca.pub`) is copied to the SSH server and used to validate certificates.
- **Docker**: The SSH server runs inside a Docker container and listens on port 2222.

## Full Command History

Here is a recap of the commands used during the project:

1. Start the HTTP server:
   ```
   go run cmd/server/main.go
   ```

2. Create a CA:
   ```
   curl "http://localhost:8080/create-ca?name=my-ca"
   ```

3. Get the CA's public key:
   ```
   curl "http://localhost:8080/get-ca-public-key?name=my-ca" -o ssh_ca.pub
   ```

4. Sign your SSH public key:
   ```
   curl -X POST -d "{\"public_key\":\"$(cat ~/.ssh/id_ed25519.pub)\"}" \
   "http://localhost:8080/sign-public-key?ca=my-ca" \
   -o ~/.ssh/id_ed25519-cert.pub
   ```

5. Build and run the Docker SSH server:
   ```
   docker build -t ssh-server -f ssh-test-server.Dockerfile .
   docker run -d -p 2222:22 --name my-ssh-server ssh-server
   ```

6. SSH into the Docker SSH server:
   ```
   ssh -i ~/.ssh/id_ed25519 -o CertificateFile=~/.ssh/id_ed25519-cert.pub -p 2222 testuser@localhost
   ```

## Conclusion

This project provides a simple, functional setup for SSH certificate-based authentication using in-memory CAs. It includes a fully functioning HTTP API to create and manage CAs, sign SSH keys, and a Dockerized SSH server that verifies SSH certificates issued by the CA. This setup is ideal for environments where certificate-based SSH authentication is required for enhanced security and control.
