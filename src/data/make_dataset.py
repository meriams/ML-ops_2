
# -*- coding: utf-8 -*-
"""
import click
import logging
from pathlib import Path
from dotenv import find_dotenv, load_dotenv


@click.command()
@click.argument('input_filepath', type=click.Path(exists=True))
@click.argument('output_filepath', type=click.Path())
def main(input_filepath, output_filepath):
    # Runs data processing scripts to turn raw data from (../raw) into
     #   cleaned data ready to be analyzed (saved in ../processed).
    #
    logger = logging.getLogger(__name__)
    logger.info('making final data set from raw data')


if __name__ == '__main__':
    log_fmt = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    logging.basicConfig(level=logging.INFO, format=log_fmt)

    # not used in this stub but often useful for finding various files
    project_dir = Path(__file__).resolve().parents[2]

    # find .env automagically by walking up directories until it's found, then
    # load up the .env entries as environment variables
    load_dotenv(find_dotenv())

    main()

"""


import click
import logging
from pathlib import Path
from dotenv import find_dotenv, load_dotenv
import subprocess
import zipfile
import os

@click.command()
@click.argument('output_dir', type=click.Path(), default='data/raw')
def main(output_dir):
    """Downloads data from Kaggle using the Kaggle API."""
    logger = logging.getLogger(__name__)
    logger.info('Downloading data from Kaggle')

    # Define dataset name
    dataset_name = 'msambare/fer2013'

    # Create output directory if it doesn't exist
    output_path = Path(output_dir)
    output_path.mkdir(parents=True, exist_ok=True)

    # Download dataset from Kaggle
    subprocess.run(['kaggle', 'datasets', 'download', '-d', dataset_name, '-p', str(output_path)])

    # Unzip the downloaded file
    zip_file_path = output_path / 'fer2013.zip'
    with zipfile.ZipFile(zip_file_path, 'r') as zip_ref:
        zip_ref.extractall(output_path)

    logger.info('Data downloaded and extracted successfully.')

if __name__ == '__main__':
    log_fmt = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    logging.basicConfig(level=logging.INFO, format=log_fmt)

    # not used in this stub but often useful for finding various files
    project_dir = Path(__file__).resolve().parents[2]

    # find .env automagically by walking up directories until it's found, then
    # load up the .env entries as environment variables
    load_dotenv(find_dotenv())

    main()

