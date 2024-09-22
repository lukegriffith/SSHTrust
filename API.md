# Developer API Docs

Although the project comes bundled with a client, the api is simple to consume.

### HTTP API Endpoints

#### 1. Create a CA
- **URL**: `/CA`
- **Method**: `POST`
- **Description**: This endpoint generates a new SSH Certificate Authority (CA) and stores it in memory under the given name.
- **Example**:
   ```bash
   curl localhost:8080/CA -X POST \
      -H "Content-Type: application/json" \
      -d '{"name": "MyCA", "bits": 2048, "type": "rsa", "valid_principals": ["testuser"], "max_ttl_minutes": 60}'
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


