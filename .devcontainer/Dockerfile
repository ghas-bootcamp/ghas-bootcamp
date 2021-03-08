FROM debian:stable-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -yq update && apt-get -yq upgrade
RUN apt-get -yq install unzip git curl apt-transport-https build-essential sudo

RUN curl -sL https://golang.org/dl/go1.15.3.linux-amd64.tar.gz --output go1.15.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz
RUN curl -sL https://packages.microsoft.com/config/debian/10/packages-microsoft-prod.deb --output packages-microsoft-prod.deb
RUN dpkg -i packages-microsoft-prod.deb
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -

RUN apt-get -yq install nodejs python3 python3-pip python3-venv dotnet-sdk-3.1

RUN useradd -ms /bin/bash developer

RUN npm install -g @vue/cli
RUN dotnet new -i IdentityServer4.Templates
