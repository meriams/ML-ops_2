import os 

# path to folder with data 
ROOT_DIR = os.path.dirname(os.path.abspath(__file__))
DATASET_FOLDER = os.path.join(ROOT_DIR, "data")
trainDirectory = os.path.join(DATASET_FOLDER, "raw/train")
testDirectory = os.path.join(DATASET_FOLDER, "raw/test")

