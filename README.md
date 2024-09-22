
# SSHTrust

## SSH CA Key Signing Server

This project is a simple rest API that generates SSH Certificate Authorities (CA), signs SSH public keys using those CAs.


## Project Overview

The project consists of:
1. An HTTP server that:
   - Creates in-memory CAs (Certificate Authorities) and stores them.
   - Signs SSH public keys using the in-memory CA.
   - Provides CA public keys for use in external systems (like the SSH server).
2. Docker test suite to demonstrate the use case. 


## API Documentation
Swagger UI is enabled for this project. You can access it by navigating to the below link when the server is active locally:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### HTTP API Endpoints

#### 1. Create a CA
- **URL**: `/CA`
- **Method**: `POST`
- **Description**: This endpoint generates a new SSH Certificate Authority (CA) and stores it in memory under the given name.
- **Example**:
   ```bash
   curl localhost:8080/CA -X POST \
      -H "Content-Type: application/json" \
      -d '{"name": "MyCA", "bits": 2048, "type": "rsa"}'
   ```

#### 2. Get the CA Public Key
- **URL**: `/CA/:id`
- **Method**: `GET`
- **Description**: Retrieves the public key of the CA stored under the specified name, which can be used by other systems to verify certificates signed by the CA.
- **Example**:
   ```bash
   curl -s "http://localhost:8080/CA/MyCA" | jq -r .public_key > ssh_ca.pub
   ```

#### 3. Sign a Public Key
- **URL**: `/CA/:id/Sign`
- **Method**: `POST`
- **Description**: Signs a public key with the specified CA. The public key should be provided in the body of the request as a JSON object in the format `{"public_key": "<public_key>"}`. The API responds with the signed certificate.
- **Example**:
   ```bash
   curl -X POST http://localhost:8080/CA/MyCA/Sign \
    -H "Content-Type: application/json" \
    -d "{\"public_key\": \"$(cat ~/.ssh/id_ed25519.pub)\", \"principals\": [\"testuser\"], \"ttl_minutes\": 50}"
   ```

- To store the signed certificate in a file:
   ```bash
    curl -X POST http://localhost:8080/CA/MyCA/Sign \
    -H "Content-Type: application/json" \
    -d "{\"public_key\": \"$(cat ~/.ssh/id_ed25519.pub)\"}" \
   | jq -r .signed_key > ~/.ssh/id_ed25519-cert.pub
   ```

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
   ./sshtrust sign -n myca --ttl 30 -p testuser -k "$(cat ~/.ssh/id_ed25519)" > ~/.ssh/id_ed25519-cert.pub
   ```

2. **SSH into the Server**:
   ```
   ssh -i ~/.ssh/id_ed25519 -o CertificateFile=~/.ssh/id_ed25519-cert.pub -p 2222 testuser@localhost
   ```

### SSH Server Setup Recap:
- **Public Key**: The CA’s public key (`ssh_ca.pub`) is copied to the SSH server and used to validate certificates.
- **Docker**: The SSH server runs inside a Docker container and listens on port 2222.

## Conclusion

This project provides a simple, functional setup for SSH certificate-based authentication using in-memory CAs. It includes a fully functioning HTTP API to create and manage CAs, sign SSH keys, and a Dockerized SSH server that verifies SSH certificates issued by the CA. This setup is ideal for environments where certificate-based SSH authentication is required for enhanced security and control.
