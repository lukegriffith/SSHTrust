# Use an official Ubuntu base image
FROM rocky:9.3

# Install OpenSSH server and required tools
RUN apt-get update && apt-get install -y \
    openssh-server \
    sudo

# Create the SSH directory and setup permissions
RUN mkdir /var/run/sshd

# Create a user to login with SSH
RUN useradd -rm -d /home/testuser -s /bin/bash -g root -G sudo -u 1000 testuser
RUN echo 'testuser:testpassword' | chpasswd

# Allow password authentication (just for initial setup)
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

# Create the .ssh directory and set permissions
RUN mkdir -p /home/testuser/.ssh && \
    chown -R testuser:root /home/testuser/.ssh && \
    chmod 700 /home/testuser/.ssh

# Copy the CA public key into the container
COPY ssh_ca.pub /etc/ssh/ca.pub

# Configure SSH to trust the CA
RUN echo "TrustedUserCAKeys /etc/ssh/ca.pub" >> /etc/ssh/sshd_config

# Expose the SSH port
EXPOSE 22

# Start the SSH service
CMD ["/usr/sbin/sshd", "-D"]
