# ML-ops

## Overall goal of the project
The overall goal of the project is to classify facial expressions into different emotions such as happiness, sadness, anger, surprise, fear, disgust, or neutral using deep learning techniques. This involves training a model to recognize patterns in facial images and associate them with specific emotions.
x
## What framework are you going to use (PyTorch Image Models, Transformer, Pytorch-Geometrics)

In this project, we intend to use PyTorch, a popular deep learning framework, for building and training the models. Specifically, we'll likely utilize standard PyTorch functionalities and potentially leverage pre-trained models and components from the PyTorch ecosystem to enhance model performance.

## How to you intend to include the framework into your project
We will integrate PyTorch into the project by setting up a PyTorch project structure, utilizing PyTorch for model construction, defining data loaders, loss functions, optimizers, and training loops. Additionally, we may employ specific PyTorch components like PyTorch Image Models or Transformers if they align with the project requirements.

## What data are you going to run on (initially, may change)
The initial dataset for training and testing is the [FER2013](https://www.kaggle.com/datasets/msambare/fer2013) dataset, which contains facial expression images labeled with corresponding emotions. The data consists of 48x48 pixel grayscale images of faces. The faces have been automatically registered so that the face is more or less centred and occupies about the same amount of space in each image. There is seven categories of emotions (0=Angry, 1=Disgust, 2=Fear, 3=Happy, 4=Sad, 5=Surprise, 6=Neutral). The training set consists of 28,709 examples and the public test set consists of 3,589 examples.

## What deep learning models do you expect to use
We expect to use a deep learning models for this project to find the best one, including:

- Convolutional Neural Networks (CNNs): A common choice for image classification tasks.
- Transfer Learning: Utilizing pre-trained CNNs (e.g., ResNet, VGG, etc.) and fine-tuning them for our specific emotion classification task.




==============================

## Running notes and reflections
### GCP
`Dockerfile` can be run on VertexAI CPU only. Getting cuda to work in the docker container was troublesome so we did not pursue it further on VertexAI. 

Instead we decided to train the model on GCP Engine on a P4 Tesla without docker. A virtual env was created using conda and reqs installed from `requirements.txt`. The `train_model.py` is used, the model checkpoint is uploaded to the GCP bucket which also contains the data files. `dvc pull` was used to successfully pull data from GCP bucket to the VM in GCP Engine.  

We use Cloud Run to deploy a FastAPI application for model inference. `api.Dockerfile` is used. The app fetches the model checkpoint from the cloud Bucket to build the model. Currently, the deployment is done manually, so whenever a new image is built and uploaded to the GCP container registry one has to manually update the cloud run. This could be automated in the future.

### Wandb
Simple scalar value profiling is done within the training script. Sweep did not work due to a networking issue in the GCP VM but the code is there and can be run locally.

### GCP Monitoring
Alerts are created in Cloud Run that get triggered on high CPU usage.

### OpenTelemetry
We did not deploy this as it required Kubernetes, but the `main.py` file supports integration with OpenTelemtetry. It can be run locally by first doing a docker-compose up inside the Signoz dir: 

```
docker-compose -f docker/clickhouse-setup/docker-compose.yaml up -d
```

Then the relevant code can be uncommented in the `train_model.py` file and then run the app. Data should now be logged to the running SigNoz app.

### GitHub CI
- pytest is run against the `test/` directory and a test coverage report is uploaded as an artefact to GitHub.
- `Dockerfile` and `api.Dockerfile` are built as images and uploaded to GCP container registry

Project Organization
------------

    ├── LICENSE
    ├── Makefile           <- Makefile with commands like `make data` or `make train`
    ├── README.md          <- The top-level README for developers using this project.
    ├── data
    │   ├── external       <- Data from third party sources.
    │   ├── interim        <- Intermediate data that has been transformed.
    │   ├── processed      <- The final, canonical data sets for modeling.
    │   └── raw            <- The original, immutable data dump.
    │
    ├── docs               <- A default Sphinx project; see sphinx-doc.org for details
    │
    ├── models             <- Trained and serialized models, model predictions, or model summaries
    │
    ├── notebooks          <- Jupyter notebooks. Naming convention is a number (for ordering),
    │                         the creator's initials, and a short `-` delimited description, e.g.
    │                         `1.0-jqp-initial-data-exploration`.
    │
    ├── references         <- Data dictionaries, manuals, and all other explanatory materials.
    │
    ├── reports            <- Generated analysis as HTML, PDF, LaTeX, etc.
    │   └── figures        <- Generated graphics and figures to be used in reporting
    │
    ├── requirements.txt   <- The requirements file for reproducing the analysis environment, e.g.
    │                         generated with `pip freeze > requirements.txt`
    │
    ├── setup.py           <- makes project pip installable (pip install -e .) so src can be imported
    ├── src                <- Source code for use in this project.
    │   ├── __init__.py    <- Makes src a Python module
    │   │
    │   ├── data           <- Scripts to download or generate data
    │   │   └── make_dataset.py
    │   │
    │   ├── features       <- Scripts to turn raw data into features for modeling
    │   │   └── build_features.py
    │   │
    │   ├── models         <- Scripts to train models and then use trained models to make
    │   │   │                 predictions
    │   │   ├── predict_model.py
    │   │   └── train_model.py
    │   │
    │   └── visualization  <- Scripts to create exploratory and results oriented visualizations
    │       └── visualize.py
    │
    └── tox.ini            <- tox file with settings for running tox; see tox.readthedocs.io


--------

<p><small>Project based on the <a target="_blank" href="https://drivendata.github.io/cookiecutter-data-science/">cookiecutter data science project template</a>. #cookiecutterdatascience</small></p>
