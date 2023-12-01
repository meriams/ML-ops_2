import torch
from torch.profiler import profile, record_function, ProfilerActivity, tensorboard_trace_handler
import sys
import os
sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), "models"))
from models.train_model import train_model
import torchvision.models as models


def our_prof():
    with profile(
        activities=[ProfilerActivity.CPU, ProfilerActivity.CUDA],
        profile_memory=True,
        record_shapes=True,
        use_cuda=True,
        on_trace_ready=tensorboard_trace_handler("./log/emotionnet")
    ) as prof:

        #test model for profiling reference
        model = models.resnet18()
        inputs = torch.randn(5, 3, 224, 224)
        model(inputs)

        # train_model()
    
        # Print Profiler results
    print('CPU profiling', prof.key_averages().table(sort_by="cpu_time_total", row_limit=10))
    print('CUDA profiling', prof.key_averages().table(sort_by="cuda_time_total", row_limit=10))

    # prof.export_chrome_trace("profile_trace.json") # trace of profiling

if __name__ == "__main__":
    our_prof()
