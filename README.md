# Branch Predictor

## Introduction

This project implements a branch predictor simulator that implements the following branch prediction algorithms:
- Always Taken
- Standard 2-bit predictor
- gshare
- Profiled

The main goal of this project is to compare the performance of the different branch prediction algorithms, identifying key strengths and weaknesses of each.

## Project Structure

The project is structured as follows:
- `main.go`: The main entry point of the program. This file reads the trace file and configuration file, and runs the branch predictor simulator.
- `branchpred/`: This directory contains the implementation of the branch prediction algorithms.
- `instruction/`: This directory contains the implementation of the reading of trace files in instruction structures.
- `utils/`: This directory contains utility functions used by the branch predictor simulator.
- `configs/`: This directory contains example configuration files for each algorithm. These examples can be modified to test different configurations of the algorithms.
- `outputs/`: This directory contains the output files generated by the branch predictor simulator. Within this directory are two subdirectories:
    - `results/`: This directory contains the results of the simulation in JSON format. Each algorithm has its won subdirectory within this directory displaying its results for each trace file.
    - `metadata/`: This directory contains the output files in JSON format. Each algorithm has its own subdirectory within this directory displaying the metadata for each trace file.

## Instructions

To run using the executable:
```bash
./branch-predictor /cs/studres/CS4202/Coursework/P2-BranchPredictor/trace-files/<trace-file> configs/<config-file>
```
To compile and run:
```bash
go run main.go /cs/studres/CS4202/Coursework/P1-CacheSim/trace-files/<trace-file>
```
To build the executable:
```bash
go build -o branch_predictor main.go
```

## Documentation

To view documentation:
```bash
godoc -http=:6060 -index
```
Then navigate to `http://localhost:6060/pkg/github.com/nsengupta5/Branch-Predictor/`

## Creating a Config File

The config file for each algorithm is a JSON file that contains the following fields:
- `algorithm`: The name of the algorithm to use. This can be one of `always-taken`, `two-bit`, `gshare`, or `profiled`.
- `max_lines`: The maximum number of lines to read from the trace file. Providing a value of -1 will read all lines from the trace file.
- `configs`: A list of configuration objects for the algorithm. Each configuration object represents a particular configuration settings for the various fields of the algorithm. The branch predictor simulator will run a simulation of the algorithm for each configuration oject specified. The fields in each configuration object depend on the algorithm. These fields are described in the following section.

### Always Taken

The `always-taken` algorithm does not require any configuration settings. 

An example configuration file for the `always-taken` algorithm is as follows:
```json
{
    "algorithm": "always-taken",
    "max_lines": 1000,
    "configs": []
}
```

### Two-Bit

The `two-bit` algorithm requires the following configuration settings:
- `table_size`: The size of the branch prediction table. This can be one of 512, 1024, 2048, or 4096.
- `initial_state`: The initial state of the branch prediction table. This can be one of `StronglyTaken`, `WeaklyTaken`, `WeaklyNotTaken`, or `StronglyNotTaken`.

An example configuration file for the `two-bit` algorithm is as follows:
```json
{
    "algorithm": "two-bit",
    "max_lines": -1,
    "configs": [
        {
            "table_size": 512,
            "initial_state": "StronglyNotTaken"
        },
        {
            "table_size": 2048,
            "initial_state": "WeaklyNotTaken"
        }
    ]
}
```

### Gshare

The `gshare` algorithm requires the following configuration settings:
- `table_size`: The size of the branch prediction table. This can be one of 512, 1024, 2048, or 4096. This will also be the size of the global history register.
- `initial_state`: The initial state of the branch prediction table. This can be one of `StronglyTaken`, `WeaklyTaken`, `WeaklyNotTaken`, or `StronglyNotTaken`.

An example configuration file for the `gshare` algorithm is as follows:
```json
{
    "algorithm": "gshare",
    "max_lines": 1000,
    "configs": [
        {
            "table_size": 1024,
            "initial_state": "StronglyTaken"
        },
        {
            "table_size": 2048,
            "initial_state": "WeaklyTaken"
        }
    ]
}
```

### Profiled

TODO

The `configs` directory contains example configuration files for each algorithm. These examples can be modified to test different configurations of the algorithms.
