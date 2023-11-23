import wandb
import random
import matplotlib.pyplot as plt

def plot_training_history(history, visualization_path):
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

def run_wandb_simulation():
    # start a new wandb run to track this script
    wandb.init(
        project="test",
        config={
            "learning_rate": 0.02,
            "architecture": "CNN",
            "dataset": "CIFAR-100",
            "epochs": 10,
        }
    )

    # simulate training
    epochs = 10
    offset = random.random() / 5
    for epoch in range(2, epochs):
        acc = 1 - 2 ** -epoch - random.random() / epoch - offset
        loss = 2 ** -epoch + random.random() / epoch + offset
    
        # log metrics to wandb
        wandb.log({"acc": acc, "loss": loss})
    
    # [optional] finish the wandb run, necessary in notebooks
    wandb.finish()

# Call the wandb simulation function
run_wandb_simulation()