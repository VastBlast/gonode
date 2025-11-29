
Step 1: Generate c/c++ code of nodejs addon bridging golang
> ./gonode-darwin generate

Step 2: Generate golang library
> ./gonode-darwin build

Step 3: Install dependencies
> ./gonode-darwin install
> 
Step 4: Compile NodeJS Addon
Tip: Ensure that nodejs, npm and node-gyp are installed on the OS
> ./gonode-darwin make