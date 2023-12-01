import os
_TEST_ROOT = os.path.dirname(__file__)  # root of test folder
print(_TEST_ROOT)
_PROJECT_ROOT = os.path.dirname(_TEST_ROOT)  # root of project
print(_PROJECT_ROOT)
_PATH_DATA = os.path.join(_PROJECT_ROOT, "dataset")  # root of data
print(_PATH_DATA)

