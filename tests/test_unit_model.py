
import torch
from sklearn.metrics import classification_report 
from torchvision.transforms import RandomCrop 
from torchvision.transforms import Grayscale 
from torchvision.transforms import ToTensor 
from torch.utils.data import DataLoader 
from torchvision import transforms 
from torchvision import datasets 
import matplotlib.pyplot as plt 
from collections import Counter 
from datetime import datetime
import torch.nn as nn 
import pandas as pd 
import os 
import hydra

import numpy as np
import cv2
import pytest
# from google.cloud import storage not used


_TEST_ROOT = os.path.dirname(__file__)  # root of test folder
_PROJECT_ROOT = os.path.dirname(_TEST_ROOT)  # root of project
_PATH_DATA = os.path.join(_PROJECT_ROOT, "dataset_local_git")  # root of data

print("test,project,data",_TEST_ROOT,_PROJECT_ROOT,_PATH_DATA) 
#local data path
#import tests.path_data
#from . import _PATH_DATA
#_PATH_DATA
#pathdata
print('Imports ok')

###################################################################
PROJECT_DIR = os.environ.get("GITHUB_WORKSPACE", None)
print("project dir",PROJECT_DIR) # None

# Checking Input Shape of Train Images
# def test_input_image_shape():

img = cv2.imread(os.path.join(_PATH_DATA, "train/angry/Training_3908.jpg"))
train_height, train_width, train_channels = img.shape
print(' -> Train image shape', train_height, train_width, train_channels)


# Checking Input Shape of Test Images
# def test_input_image_shape():
img = cv2.imread(os.path.join(_PATH_DATA, "test/angry/PrivateTest_88305.jpg"))
# img = cv2.imread("data/raw/test/angry/PrivateTest_88305.jpg")
test_height, test_width, test_channels = img.shape
print (' -> Test image shape', test_height, test_width, test_channels)


###################################################################
            #Transformer and loader

def train_transformer():
    # Initialize a list of preprocessing steps to apply on each image during training/validation and testing 
        train_transform = transforms.Compose([
        Grayscale(num_output_channels=1),
        RandomHorizontalFlip(),
        RandomCrop((48,48)),
        ToTensor()
        ])
        return train_transform

     # Load all the images whithin the specified folder and apply different augmentation 
     # Get the absolute path of the current file's directory

current_file_dir = os.path.dirname(os.path.abspath(__file__))
# Traverse two directories up from the current file's directory
#ROOT_DIR = os.path.abspath(os.path.join(current_file_dir, '..'))
#ET_FOLDER = os.path.join(_PATH_DATA)
trainDirectory = os.path.join(_PATH_DATA, "train")

train_data = datasets.ImageFolder(trainDirectory, transform=train_transformer)
print(type(train_data))

classes = train_data.classes
class_labels_list = classes

number_of_classes = len(classes) 
print(number_of_classes) # 7 test for number of classes

###################################################################
    # Test for class number and labels correctness

def test_classes(class_rules):
    
    assert all(class_rules)
    
    if all (class_rules):
        class_labels_result = " -> Classes : Numbers + Labels Correct"
    else:
        class_labels_result = " -> Classes : Numbers + Labels Not Correct"
    print(class_labels_result)
    return class_labels_result

true_labels_list = ['angry', 'disgust', 'fear', 'happy', 'neutral', 'sad', 'surprise']

@pytest.fixture
def class_rules():
    return [class_labels_list == true_labels_list,
                number_of_classes == 7]

###################################################################

    #Test for Input Dimensions
def test_input_shape(rules):
    assert all(rules)
    if all(rules):
        input_dimention__result = " -> Input : Train + Test Shape Correct"
    else:
        input_dimention__result = " -> Input : Train + Test Shape Not Correct Dimensions"
    print(input_dimention__result)
    return input_dimention__result

@pytest.fixture
def rules():
    return [train_height == 48,
    train_width == 48,
    train_channels == 3,
    
    test_height == 48,
    test_width == 48,
    test_channels == 3]

@pytest.fixture
def weight_rules():
    t1 = torch.tensor([7., 7., 7., 7., 7., 7., 7.])
    return t1.numpy() #to array

def test_class_weight(weight_rules):
    class_count = Counter(train_data.classes)
    print(f'[INFO] Total sample: {class_count}')
    # depending on the number of samples available 
    class_weight = torch.Tensor([len(train_data.classes) / c
                                for c in pd.Series(class_count).sort_index().values])
    assert class_weight.numpy().all() == weight_rules.all() #all elements are equal
