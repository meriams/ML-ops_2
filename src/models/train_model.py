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
from utils  import EarlyStopping
from utils  import LRScheduler 
from torchvision import transforms 
from predict_model import EmotionNet 
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
import wandb
import cProfile
from google.cloud import storage
#from src.visualization.visualize import plot_training_history

# Initialize wandb with your project name and your entity (your username or organization name)
wandb.init(project='MLops', entity='ml_ops_dtu')
wandb.init(config=cfg) # for running a hyperparameter optimization sweep

# Get the absolute path to the directory of the current script (train_model.py)
current_dir = os.path.dirname(os.path.abspath(__file__))
root_dir = os.path.abspath(os.path.join(current_dir, '..', '..'))  # Move two levels up to the root directory

# Define paths for model output and visualization
models_folder = os.path.join(root_dir, 'models')  # Path to the models directory
model_output_path = os.path.join(models_folder, 'my_model.pth')  # Model output path
# Define path for visualization (report/figures)
visualization_folder = os.path.join(root_dir, 'reports', 'figures')  # Path to the visualization directory
visualization_path = os.path.join(visualization_folder, 'training_plot.png')

# Create directories if they don't exist
#os.makedirs(models_folder, exist_ok=True)
#os.makedirs(visualization_folder, exist_ok=True)


@hydra.main(config_name='config.yaml')
def train_model(hcfg):
    # Access hyperparameters from the config
    #batch_size = hcfg.hyperparameters.batch_size
    #learning_rate = hcfg.hyperparameters.lr
    #num_epochs = hcfg.hyperparameters.num_epochs
    #train_split = hcfg.dataset.train_split
    #val_split = hcfg.dataset.val_split

    # configure the device to use for the training the mode
    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    print(f'[INFO] Current training device: {device}')

    # Initialize a list of preprocessing steps to apply on each image during training/validation and testing 
    train_transform = transforms.Compose([
        Grayscale(num_output_channels=1),
        RandomHorizontalFlip(),
        RandomCrop((48,48)),
        ToTensor()
    ])

    test_transform = transforms.Compose([
        Grayscale(num_output_channels=1),
        ToTensor()
    ])
    
    # Load all the images whithin the specified folder and apply different augmentation 
    train_data = datasets.ImageFolder(cfg.trainDirectory, transform=train_transform)
    test_data = datasets.ImageFolder(cfg.testDirectory, transform=test_transform)

    # extract class labels and the total number of classes
    classes = train_data.classes
    num_of_classes = len(classes)
    print(f'[INFO] Class labels: {classes}')

    # Using tain samples to generate train/validation set 
    num_train_samples = len(train_data)
    train_size = math.floor(num_train_samples *hcfg.dataset.train_split)
    val_size = math.ceil(num_train_samples*hcfg.dataset.val_split)
    print(f'[INFO] Train samples: {train_size} ... \t Validaiton samples: {val_size} ...')

    # random split the training dataset into train and validation set 
    train_data, val_data = random_split(train_data, [train_size, val_size])

    # Modify the data transform applied towards the validaiton set 
    val_data.dataset.transform = test_transform

    #get the labels within the traning set 
    train_classes = [label for _, label in train_data]

    # count each labels within each classes 
    class_count = Counter(train_classes)
    print(f'[INFO] Total sample: {class_count}')

    # depending on the number of samples available 
    class_weight = torch.Tensor([len(train_classes) / c
                                for c in pd.Series(class_count).sort_index().values])

    # corresponding class weight already computed 
    sample_weight = [0] * len(train_data)
    for idx, (image, label) in enumerate(train_data):
        weight = class_weight[label]
        sample_weight[idx] = weight 

    # Define a sampler which randomly sample labels fromthe train dataset 
    sampler = WeightedRandomSampler(weights=sample_weight, num_samples=len(train_data),
                                    replacement=True)

    # Load our own dataset and store each sample with their corresponding labels 
    train_dataloader = DataLoader(train_data, batch_size=hcfg.hyperparameters.batch_size, sampler=sampler)
    val_dataloader = DataLoader(val_data, batch_size=hcfg.hyperparameters.batch_size)
    test_dataloader = DataLoader(test_data, batch_size=hcfg.hyperparameters.batch_size)

    # initialize the mdoel and send it to device  
    model = EmotionNet(num_of_channels=1, num_of_classes=num_of_classes)
    model = model.to(device)

    # initialize our optimizer and loss function 
    optimizer = SGD(model.parameters(), hcfg.hyperparameters.lr)
    criterion = nn.CrossEntropyLoss()

    # Initialize the learning rate scheduler and early stopping mechanism 
    lr_scheduler = LRScheduler(optimizer)
    early_stopping = EarlyStopping()

    #Calculate the steps per epoch for training and validation 
    train_steps = len(train_dataloader.dataset) // hcfg.hyperparameters.batch_size
    val_steps = len(val_dataloader.dataset) // hcfg.hyperparameters.batch_size

    # initialize dictionary 
    history = {
        'train_acc': [], 
        'train_loss': [], 
        'val_acc': [], 
        'val_loss': []
    }

    #iterate through the epochs 
    print(f'[INFO] Training the model...')
    start_time = datetime.now()

    for epoch in range(0, hcfg.hyperparameters.num_epochs):
        print(f'[INFO] epoch: {epoch + 1}/{hcfg.hyperparameters.num_epochs}')

        '''
        Training 
        '''

        # set the model to training mode

        model.train()

        # initialize the total training and validation loss and the total 
        # number of correct predictions in both steps 
        total_train_loss = 0 
        total_val_loss = 0 
        train_correct = 0 
        val_correct = 0 

        # iterate through the training set
        for idx, (data, target) in enumerate(train_dataloader):
            print("printing data shape")
            print(data.shape)
            if idx==1:
                break
            # move the data into the device used for training,
            data, target = data.to(device), target.to(device)
    
            # perform a forward pass and calculate the training loss
            predictions = model(data)
            loss = criterion(predictions, target)
    
            # zero the gradients accumulated from the previous operation,
            # perform a backward pass, and then update the model parameters
            optimizer.zero_grad()
            loss.backward()
            optimizer.step()


            # add the training loss and keep track of the number of correct predictions
            total_train_loss += loss
            train_correct += (predictions.argmax(1) == target).type(torch.float).sum().item()

        '''
        Validation the model 
        '''

        model.eval()  # disable dropout and dropout layers
    
        # prevents pytorch from calculating the gradients, reducing
        # memory usage and speeding up the computation time (no back prop)
        with torch.set_grad_enabled(False):
    
            # iterate through the validation set
            for idx2, (data, target) in enumerate(val_dataloader):
                if idx2==1:
                    break
                # move the data into the device used for testing
                data, target = data.to(device), target.to(device)
    
                # perform a forward pass and calculate the training loss
                predictions = model(data)
                loss = criterion(predictions, target)
    
                # add the training loss and keep track of the number of correct predictions
                total_val_loss += loss
                val_correct += (predictions.argmax(1) == target).type(torch.float).sum().item()


        # calculate the average training and validation loss
        avg_train_loss = total_train_loss / train_steps
        avg_val_loss = total_val_loss / val_steps
    
        # calculate the train and validation accuracy
        train_correct = train_correct / len(train_dataloader.dataset)
        val_correct = val_correct / len(val_dataloader.dataset)
    
        # print model training and validation records
        print(f"train loss: {avg_train_loss:.3f}  .. train accuracy: {train_correct:.3f}")
        print(f"val loss: {avg_val_loss:.3f}  .. val accuracy: {val_correct:.3f}", end='\n\n')
    
        # update the training and validation results
        history['train_loss'].append(avg_train_loss.cpu().detach().numpy())
        history['train_acc'].append(train_correct)
        history['val_loss'].append(avg_val_loss.cpu().detach().numpy())
        history['val_acc'].append(val_correct)
    
        # execute the learning rate scheduler and early stopping
        validation_loss = avg_val_loss.cpu().detach().numpy()
        lr_scheduler(validation_loss) 
        early_stopping(validation_loss)
    
        # stop the training procedure due to no improvement while validating the model
        if early_stopping.early_stop_enabled:
            break
    
    print(f"[INFO] Total training time: {datetime.now() - start_time}...")

    # move model back to cpu and save the trained model to disk
    if device == "cpu":
        model = model.to("cpu")
    #model_output = os.path.join('models', 'my_model.pth')
    torch.save(model.state_dict(), model_output_path)

    # bucket_name = "your-bucket-name"
    # source_file_name = "local/path/to/file"
    # destination_blob_name = "storage-object-name"
    # project_number = os.environ["CLOUD_ML_PROJECT_ID"]
    storage_client = storage.Client(project="dtumlops-404710")
    bucket = storage_client.bucket("fer2013_mlops")
    blob = bucket.blob("my_model.pth")
    blob.upload_from_filename(model_output_path)

    # # Upload model artifact to Cloud Storage
    # model_directory = os.environ['AIP_MODEL_DIR']
    # storage_path = os.path.join(model_directory,'my_model.pth')
    # blob = storage.blob.Blob.from_string(storage_path, client=storage.Client())
    # blob.upload_from_filename(model_output_path)

    # save the visualization
    plt.style.use("ggplot")
    plt.figure()
    plt.plot(history['train_acc'], label='train_acc')
    plt.plot(history['val_acc'], label='val_acc')
    plt.plot(history['train_loss'], label='train_loss')
    plt.plot(history['val_loss'], label='val_loss')
    plt.ylabel('Loss/Accuracy')
    plt.xlabel("#No of Epochs")
    plt.title('Training Loss and Accuracy on FER2013')
    plt.legend(loc='upper right')
    plt.savefig(visualization_path) 
    #plot_training_history(history, visualization_path)

    # Log metrics to wandb
    wandb.log({
       'epoch': epoch,
       'train_loss': avg_train_loss,  # calculated training loss
       'train_accuracy': train_correct,  # calculated training accuracy
       'val_loss': avg_val_loss,  # calculated validation loss
       'val_accuracy': val_correct  # calculated validation accuracy
    })

    #wandb.save('my_model.pth')  # Replace 'path_to_your_model.pth' with your actual model path
    wandb.save('training_plot.png')

    # evaluate the model based on the test set
    model = model.to(device)
    with torch.set_grad_enabled(False):
        # set the evaluation mode
        model.eval()
    
        # initialize a list to keep track of our predictions
        predictions = []
    
        # iterate through the test set
        for (data, _) in test_dataloader:
            # move the data into the device used for testing
            data = data.to(device)
    
            # perform a forward pass and calculate the training loss
            output = model(data)
            output = output.argmax(axis=1).cpu().numpy()
            predictions.extend(output)
    
    # evaluate the network
    print("[INFO] evaluating network...")
    actual = [label for _, label in test_data]
    print(classification_report(actual, predictions, target_names=test_data.classes))

    if __name__ == '__main__':
        train_model()

    # Run the function with cProfile
        #cProfile.run('train_model()', filename='train_model_profile.txt')
        
        # Your existing code for hyperparameter optimization
        sweep_config = {
            'method': 'grid',
            'metric': {'goal': 'maximize', 'name': 'val_accuracy'},
            'parameters': {
                'hyperparameters.batch_size': {'values': [16, 32]},
                'hyperparameters.lr': {'values': [0.1, 0.01]}
            }
        }

        # Initialize the sweep
        sweep_id = wandb.sweep(sweep_config, project='MLops', entity='ml_ops_dtu')

        # Run the sweep
        wandb.agent(sweep_id, function=train_model)