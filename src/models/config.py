import os 

""" trying something different  
# Get the absolute path of the current file's directory
current_file_dir = os.path.dirname(os.path.abspath(__file__))

# Traverse two directories up from the current file's directory
ROOT_DIR = os.path.abspath(os.path.join(current_file_dir, '..', '..'))

DATASET_FOLDER = os.path.join(ROOT_DIR, "data")
trainDirectory = os.path.join(DATASET_FOLDER, "raw/train")
testDirectory = os.path.join(DATASET_FOLDER, "raw/test")
"""

# Get the absolute path of the current file's directory
current_file_dir = os.path.dirname(os.path.abspath(__file__))

# Traverse two directories up from the current file's directory
ROOT_DIR = os.path.abspath(os.path.join(current_file_dir, '..', '..'))

# DATASET_FOLDER = os.path.join(ROOT_DIR, "data")
# trainDirectory = os.path.join(DATASET_FOLDER, "raw/train")
# testDirectory = os.path.join(DATASET_FOLDER, "raw/test")
mnt_dir = os.environ.get("MNT_DIR", "/mnt/nfs/filestore")
trainDirectory = os.path.join(mnt_dir, "fer2013_mlops/data/train")
testDirectory = os.path.join(mnt_dir, "fer2013_mlops/data/test")

# Define paths for model output and visualization
MODELS_FOLDER = os.path.join(ROOT_DIR, 'models')
MODEL_PATH = os.path.join(MODELS_FOLDER, 'my_model.pth')
VISUALIZATION_FOLDER = os.path.join(ROOT_DIR, 'reports', 'figures')
VISUALIZATION_PATH = os.path.join(VISUALIZATION_FOLDER, 'training_plot.png')


