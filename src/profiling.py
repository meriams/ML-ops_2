import torch
from torch.profiler import profile, record_function, ProfilerActivity
import sys
import os
sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), "models"))
from models.train_model import train_model

def our_prof():
    with profile(
        activities=[ProfilerActivity.CPU, ProfilerActivity.CUDA],
        profile_memory=True,
        record_shapes=True,
        use_cuda=True
    ) as prof:
        train_model()
    
        # Print Profiler results
        print('CPU profiling', prof.key_averages().table(sort_by="cpu_time_total", row_limit=10))
        print('CUDA profiling', prof.key_averages().table(sort_by="cuda_time_total", row_limit=10))
        
        prof.export_chrome_trace("profile_trace.json") # trace of profiling

if __name__ == "__main__":
    our_prof()
