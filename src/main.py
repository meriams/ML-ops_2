'''
Create a Fastapi app with a single POST endpoint. It should be possible to upload an image file. The image file is used as input to a vision model. The model should be run in inference mode and the api call should return the model output. 
'''

from fastapi import FastAPI, File, UploadFile
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
import os
import shutil
import uvicorn
import numpy as np
from fastapi import UploadFile, File
from typing import Optional
from fastapi import status
# import sys
# sys.path.append('/Users/kamal/Documents/ML-ops_2/src/models/predict_model')
import torch
from models import EmotionNet
# breakpoint()
from PIL import Image
import io
from torchvision.transforms import Compose, Grayscale, ToTensor


app = FastAPI()

@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.post("/model/")
async def model(data: UploadFile = File(...)):

    mnt_dir = os.environ.get("MNT_DIR", "/mnt/nfs/filestore")
    model_pth = os.path.join(mnt_dir, "my_model.pth")
    print("this is the dir")
    print(model_pth)
    model = EmotionNet(num_of_channels=1, num_of_classes=7)
    model.load_state_dict(torch.load(f=model_pth, map_location=torch.device('cpu')))
    model = model.to("cpu")
    model.eval()
    
    contents = await data.read()
    image = Image.open(io.BytesIO(contents))

    # Define the transform
    test_transform = Compose([
        Grayscale(num_output_channels=1),
        ToTensor()
    ])

    # Apply the transform
    transformed_image = test_transform(image)
    transformed_image = transformed_image.unsqueeze(0)

    out = model(transformed_image)

    res_mapping = {{"0": 'angry', "1": 'disgust', "2": 'fear', "3": 'happy', "4": 'neutral', "5": 'sad', "6": 'surprise'}} #! fix this later
    index = torch.argmax(torch.nn.functional.softmax(out, dim=1)).detach().numpy().item()
    
    response = {
        "output": res_mapping[str(index)],
        "message": "success",
        "status-code": status.HTTP_200_OK,
    }
    return response

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=80)