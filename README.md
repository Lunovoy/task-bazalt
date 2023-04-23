# ALT Linux Package Comparison CLI

## Installation 

### __Build from binaries__

1. Clone repository
    ```
    git clone https://github.com/Lunovoy/task-bazalt.git
    cd task-bazalt
    ```
2. Build

    ```
    go build ./cmd/comparison.go
    ```
3. Run
    >Usage:
    > ./comparison "branch1" "branch2"

    Example:
    ```
    ./comparison p10 p9
    ```
## __Download from Release__

1. Download with wget
    ```
    wget https://github.com/Lunovoy/task-bazalt/releases/download/Latest/comparison
    ```

2. Give permissions 
    ```
    chmod +x ./comparison
    ```

3. Run
    >Usage:
    > ./comparison "branch1" "branch2"

    Example:
    ```
    ./comparison p10 p9
    ```