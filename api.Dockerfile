FROM python:3.9


# Install system dependencies
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
WORKDIR /
COPY ./requirements.txt /requirements.txt
COPY ./src /

# Install production dependencies.
RUN pip install --no-cache-dir --upgrade -r /requirements.txt

# Ensure the script is executable
RUN chmod +x gcsfuse_run.sh

# Use tini to manage zombie processes and signal forwarding
# https://github.com/krallin/tini
ENTRYPOINT ["/usr/bin/tini", "--"] 

# Pass the startup script as arguments to Tini
CMD ["/gcsfuse_run.sh"]


# CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]

