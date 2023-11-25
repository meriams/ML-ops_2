#############
# Use the official lightweight Python image.
# https://hub.docker.com/_/python
FROM nvidia/cuda:11.3.1-cudnn8-runtime-ubuntu20.04
FROM python:3.9


# Install system dependencies
# RUN set -e; \
#     apt-get update -y && \
#     apt-get install -y --no-install-recommends python3.9 python3.9-dev python3-pip && \
#     apt-get install -y lsb-release && \
#     update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.9 1 && \
#     update-alternatives --install /usr/local/bin/python python /usr/bin/python3.9 1 && \
#     tini \
#     lsb-release; \
#     gcsFuseRepo=gcsfuse-`lsb_release -c -s`; \
#     echo "deb https://packages.cloud.google.com/apt $gcsFuseRepo main" | \
#     tee /etc/apt/sources.list.d/gcsfuse.list; \
#     curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | \
#     apt-key add -; \
#     apt-get update; \
#     apt-get install -y gcsfuse \
#     && apt-get clean
RUN set -e; \
    apt-get update -y && apt-get install -y \
    tini \
    lsb-release; \
    gcsFuseRepo=gcsfuse-`lsb_release -c -s`; \
    echo "deb https://packages.cloud.google.com/apt $gcsFuseRepo main" | \
    tee /etc/apt/sources.list.d/gcsfuse.list; \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | \
    apt-key add -; \
    apt-get update; \
    apt-get install -y gcsfuse \
    && apt-get clean

# Set fallback mount directory
ENV MNT_DIR /mnt/gcs
ENV BUCKET fer2013_mlops

# Copy local code to the container image.
# Set the working directory in the container
WORKDIR /
# Copy the requirements file into the container at /app
COPY requirements.txt /
COPY . /

# Install production dependencies.
RUN pip3 install -r requirements.txt --no-cache-dir

# Ensure the script is executable
RUN chmod +x gcsfuse_run_train.sh

# Use tini to manage zombie processes and signal forwarding
# https://github.com/krallin/tini
ENTRYPOINT ["/usr/bin/tini", "--"] 

# Pass the startup script as arguments to Tini
CMD ["/gcsfuse_run_train.sh"]