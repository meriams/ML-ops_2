# Use an official Python runtime as a parent image
# linux/amd to run via emulation on apple silicon, removing builds image against arm64 but crashes on run
# caveat is to give flag when running linux/amd64 image: docker run --platform linux/amd64 mlops
# FROM --platform=linux/amd64 python:3.9-slim-buster 
FROM nvidia/cuda:11.3.1-cudnn8-runtime-ubuntu20.04

# Set the working directory in the container
WORKDIR /

# Copy the requirements file into the container at /app
COPY requirements.txt /

# Install any needed packages specified in requirements.txt
RUN apt-get update && \
    apt-get install -y --no-install-recommends python3-pip && \
    pip3 install -r requirements.txt --no-cache-dir && \
    rm -rf /var/lib/apt/lists/*


# Copy the rest of the application code into the container
COPY . /

# Expose a port if your application needs it (e.g., for a web server)
# EXPOSE 80

# Define the command to run your application
# unbuffered output otherwise can't see print statements
CMD [ "python", "-u", "src/models/train_model.py" ]
