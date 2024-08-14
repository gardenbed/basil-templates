FROM golang:1.23

# Install required tools
RUN apt-get update && \
    apt-get install -y git jq

# Install Basil CLI
RUN download_url=$(\
      wget -qO - https://api.github.com/repos/gardenbed/basil-cli/releases/latest | \
      jq -r '.assets[] | select(.name == "basil-linux-amd64") | .browser_download_url' \
    ) && \
    wget -q -O /usr/local/bin/basil ${download_url} && \
    chmod +x /usr/local/bin/basil && \
    basil -version

COPY entrypoint.sh /
ENTRYPOINT [ "/entrypoint.sh" ]
