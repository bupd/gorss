#!/bin/sh

# Start PostgreSQL Server
echo "Starting PostgreSQL server..."
sudo systemctl start postgresql 

# Activate Virtual Environment
echo "Activating virtual environment..."
# Replace 'path/to/your/venv' with the actual path to your virtual environment.
source /home/bupd/code/pg/venv/bin/activate

echo "spawning the pgadmin4..."
process= sudo pgadmin4 &

sleep 2

echo "process: $process"
# Open Link in Browser
echo "Opening link in the browser..."
# Replace 'https://example.com' with the actual link you want to open.
xdg-open http://127.0.0.1:5050  # Use 'open' on macOS or 'start' on Windows.

# Provide a message indicating the script execution is complete.
echo "pgadmin4 started and executed finally."
