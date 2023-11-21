import os 

# Get the absolute path of the current file's directory
current_file_dir = os.path.dirname(os.path.abspath(__file__))

# Traverse two directories up from the current file's directory
ROOT_DIR = os.path.abspath(os.path.join(current_file_dir, '..', '..'))

DATASET_FOLDER = os.path.join(ROOT_DIR, "data")
trainDirectory = os.path.join(DATASET_FOLDER, "raw/train")
testDirectory = os.path.join(DATASET_FOLDER, "raw/test")

