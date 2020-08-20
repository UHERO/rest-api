# Examples

This directory holds examples to interact with the API from various languages/tools.

## Python/Jupyter Example

Jupyter Notebooks are a popular Python tool for interactively exploring data. The notebook shows the user how to:

1. Use the API key to pull data based on series ID, using the requests library
2. Transform the data into a Pandas dataframe
3. Plot the dataframe using matplotlib

To run the example, you will need the following:

* UHERO API Key
* Install (Docker)[https://docs.docker.com/get-docker/]

For convenience, Jupyter has created Dockerfiles containing common Python data science libraries like Pandas, NumPy, SciPy, etc.

Documentation for the images can be found (here)[https://github.com/jupyter/docker-stacks]

To run the example:
```
docker run --rm -p 8888:8888 -e JUPYTER_ENABLE_LAB=yes -v $PWD:/home/jovyan/work jupyter/scipy-notebook:latest
```