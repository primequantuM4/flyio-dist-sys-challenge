# flyio-dist-sys-challenge
Solutions to the fly.io distributed systems challenge
## Introduction
This repository contains my solutions to the Fly.io Distributed Systems Challenge. The challenge is designed to help developers learn and implement key concepts in distributed systems through hands-on problem-solving. Each solution reflects my approach to tackling distributed computing.

## Repository Structure
- `challenge-1/`: Solution and code for Task 1 of the challenge.

- `challenge-2/`: Solution and code for Task 2 of the challenge.
  
- `challenge-3/`: Solutions for Task 3, which has multiple parts:
  - `part-a/`: Code and solution for Task 3, Part 1.
  - `part-b/`: Code and solution for Task 3, Part 2.
  - `part-c/`: Code and solution for Task 3, Part 3.
    
- `challenge-4/`: Solution and code for Task 4 of the challenge.

- `challenge-5/`: Solution and code for Task 5, which has multiple parts:
  - `part-a/`: Code and solution for Task 5, Part 1.
  - `part-b/`: Code and solution for Task 5, Part 2.
  - `part-c/`: Code and solution for Task 5, Part 3.
 
- `challenge-6/`: Solution and code for Task 6, which has multiple parts:
  - `part-a/`: Code and solution for Task 6, Part 1.
  - `part-b/`: Code and solution for Task 6, Part 2.
  - `part-c/`: Code and solution for Task 6, Part 3.
 
 ## Technologies Used
- Go (Golang)
- Maelstrom

## How to Run
1. Clone the repository:
   ```bash
   git clone https://github.com/primequantuM4/flyio-dist-sys-challenge.git
   ```
2. Navigate to the challenge of your choosing:
    ```bash
    cd flyio-dist-sys-challenge/challenge-<number>/
    ```
    - For challenges with multiple parts:
      ```bash
        cd flyio-dist-sys-challenge/challenge-<number>/part-<letter>/
      ```
3. Follow the evaluation command listed on Fly.io for the corresponding task
  - Make sure to run `go install .` for the specific task you want to test.
  - Visit fly.io/dist-sys/<task-number> for specific instructions.
  - Modify the command as needed, changing the directory path from ` ~/go/bin/maelstrom... ` to ` ~/go/bin/challenge-<number> ` or ` ~/go/bin/challenge-<number>/part-<letter> ` for correct execution.
  - Optionally you can run the `test.sh` command found in every directory using:
    ```bash
        ./test.sh
    ```
    
