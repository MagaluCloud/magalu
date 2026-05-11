# Poetry
[Poetry](https://python-poetry.org/) is a python manager. With it, developers can control virtual environments, dependencies and builds

## Installation
To install Poetry, you can follow the instructions in the [oficial documentation](https://python-poetry.org/docs/#installation).

For quick installation, use the following command on Unix-like systems:
```bash
curl -sSL https://install.python-poetry.org | python3 -
```

Or this for Windows Systems:
```bash
(Invoke-WebRequest -Uri https://install.python-poetry.org -UseBasicParsing).Content | py -
```

## Usage
This project already has a [poetry configuration file](pyproject.toml).
Run these commands on the project root.

### Start the python environment
To install the correct python version and its dependencies in a virtual shell, use:
```bash
poetry shell
```
This will spawn a new shell based on your default system shell ready for python use.
Since it's a shell, simply type `exit` to return to the default environment

### Install dependencies
The first time you start a virtual environment (venv), or if you're installing it
globally, use `poetry install`. This will install all dependencies, including dev dependencies,
on the environment being used.