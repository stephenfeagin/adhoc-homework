# NOC List

## Installing Dependencies

This project uses [Pipenv](https://pipenv.readthedocs.io) to manage dependencies. Simply
run `pipenv install` from the command line to install the required packages. 

## Running the Server

As per the instructions in the original [README](./README-original.md), run the "BADSEC"
server with:

```bash
$ docker run --rm -p 8888:8888 adhocteam/noclist
```

## Executing the Program

Once the server is running, you can run the program with:

```bash
$ pipenv run python noclist.py
```

