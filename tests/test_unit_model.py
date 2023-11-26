
import torch
from torchvision.transforms import RandomHorizontalFlip
from torch.utils.data import WeightedRandomSampler 
from sklearn.metrics import classification_report 
from torchvision.transforms import RandomCrop 
from torchvision.transforms import Grayscale 
from torchvision.transforms import ToTensor 
from torch.utils.data import random_split 
from torch.utils.data import DataLoader 
# import config as cfg 
from torchvision import transforms 
# from model import EmotionNet 
from torchvision import datasets 
import matplotlib.pyplot as plt 
from collections import Counter 
from datetime import datetime
from torch.optim import SGD 
import torch.nn as nn 
import pandas as pd 
import argparse 
import math 
import os 
import hydra

import numpy as np
import cv2
import pytest
from google.cloud import storage

print('Imports ok')

###################################################################

# Checking Input Shape of Train Images
#def test_input_image_shape():
# img = cv2.imread("/Users/kamal/Documents/ML-ops_2/data/raw/train/angry/Training_3908.jpg")
# train_height, train_width, train_channels = img.shape
# print(' -> Train image shape', train_height, train_width, train_channels)

storage_client = storage.Client(project="dtumlops-404710")
bucket = storage_client.bucket("fer2013_mlops")
blob = bucket.blob("my_model.pth")

# Download the model file to the specified path
blob.download_to_filename("my_model.pth")



# Checking Input Shape of Test Images
#def test_input_image_shape():
img = cv2.imread("/Users/kamal/Documents/ML-ops_2/data/raw/test/angry/PrivateTest_88305.jpg")
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

classes = 7 #! Should be fetched from the model def file
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

# test_input_shape(input_rules)



