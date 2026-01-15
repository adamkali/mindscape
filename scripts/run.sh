#!/usr/bin/bash

# Colors for output
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# point the output of the build to the tmp directory
#

# start running the frontend
# and keep track of the process

echo -e "$GREEN Starting frontend... $NC"
cd ./web/
pnpm dev &
FRONTEND_PID=$!
echo -e "$GREEN  Frontend started successfully with PID: $FRONTEND_PID $NC"

# Return to root directory for air
cd ..

# Cleanup function to kill frontend on exit
cleanup() {
    echo -e "\n$RED Killing frontend process: $FRONTEND_PID $NC"
    kill $FRONTEND_PID 2>/dev/null
    echo -e "$GREEN Frontend killed successfully $NC"
    echo -e "$GREEN Ok... bye now $NC"
}

# Set trap to cleanup on script exit or interrupt
trap cleanup EXIT INT TERM

# start running the backend
# and now run the backend using air
air
