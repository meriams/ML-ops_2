import torch
from torchvision.transforms import RandomHorizontalFlip
from torch.utils.data import WeightedRandomSampler 
from sklearn.metrics import classification_report 
from torchvision.transforms import RandomCrop 
from torchvision.transforms import Grayscale 
from torchvision.transforms import ToTensor 
from torch.utils.data import random_split 
from torch.utils.data import DataLoader 
import config as cfg
from src.models.utils  import EarlyStopping
from src.models.utils  import LRScheduler 
from torchvision import transforms 
from src.models.predict_model import EmotionNet 
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
from torch.profiler import profile, record_function, ProfilerActivity
from src.models.train_model import train_model
import torchvision.models as models



#test model for profiling reference
model = models.resnet18()
inputs = torch.randn(5, 3, 224, 224)

# Testing case not requiered but shows an average model that is one by torch and should be optimized so we can see the difference as benchmark
with profile( activities=[ProfilerActivity.CPU, ProfilerActivity.CUDA],
        profile_memory=True,
        record_shapes=True , use_cuda=True
        ) as prof:
     # stucture is model(inputs)
     #model(inputs) # to see a defined model as a base for osbervational reference
     model(inputs)
     pass

# Print Profiler results
print('CPU profiling',prof.key_averages().table(sort_by="cpu_time_total", row_limit=10))
print('Cude profiling',prof.key_averages().table(sort_by="cuda_time_total", row_limit=10))


@hydra.main(config_name='config_profiling.yaml')
def our_prof(cfg):
    with profile(
        activities=[ProfilerActivity.CPU, ProfilerActivity.CUDA],
        profile_memory=True,
        record_shapes=True,
        use_cuda=True
    ) as prof:
        # train_model is a function that takes a configuration argument
        train_model(cfg)
    
    # Print Profiler results
    print('CPU profiling', prof.key_averages().table(sort_by="cpu_time_total", row_limit=10))
    print('CUDA profiling', prof.key_averages().table(sort_by="cuda_time_total", row_limit=10))

if __name__ == "__main__":
    our_prof()


prof.export_chrome_trace("profile_trace.json") # trace of profiling