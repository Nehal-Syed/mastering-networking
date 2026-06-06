#!/bin/bash

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN} Starting NetPlay Setup...${NC}\n"

# Check if MySQL is installed
if ! command -v mysql &> /dev/null; then
    echo -e "${RED} MySQL is not installed. Please install MySQL first.${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED} Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi

# Check if Nginx is installed
if ! command -v nginx &> /dev/null; then
    echo -e "${RED} Nginx is not installed. Please install Nginx first.${NC}"
    exit 1
fi

# Setup database
echo -e "${YELLOW} Setting up database...${NC}"
mysql -u root -p < sql/init.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN} Database setup complete${NC}"
else
    echo -e "${RED} Database setup failed${NC}"
    exit 1
fi

# Install Go dependencies
echo -e "${YELLOW} Installing Go dependencies...${NC}"
cd backend
go mod download
go mod tidy

if [ $? -eq 0 ]; then
    echo -e "${GREEN} Go dependencies installed${NC}"
else
    echo -e "${RED} Go dependencies installation failed${NC}"
    exit 1
fi

# Build backend
echo -e "${YELLOW} Building backend...${NC}"
go build -o netplay-backend cmd/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN} Backend build successful${NC}"
else
    echo -e "${RED} Backend build failed${NC}"
    exit 1
fi

# Start backend in background
echo -e "${YELLOW} Starting backend server...${NC}"
./netplay-backend &
BACKEND_PID=$!
echo -e "${GREEN} Backend started with PID: $BACKEND_PID${NC}"

# Wait for backend to start
sleep 3

# Configure Nginx
echo -e "${YELLOW}🔧 Configuring Nginx...${NC}"
# Update the root path in nginx.conf
sed -i "s|/path/to/netplay/frontend|$(pwd)/../frontend|g" ../nginx/nginx.conf

# Copy Nginx config
sudo cp ../nginx/nginx.conf /etc/nginx/sites-available/netplay
sudo ln -sf /etc/nginx/sites-available/netplay /etc/nginx/sites-enabled/
sudo nginx -t

if [ $? -eq 0 ]; then
    sudo systemctl reload nginx
    echo -e "${GREEN} Nginx configured successfully${NC}"
else
    echo -e "${RED} Nginx configuration failed${NC}"
    kill $BACKEND_PID
    exit 1
fi

echo -e "\n${GREEN}Setup Complete! ${NC}"
echo -e "${GREEN}Access the application: http://localhost${NC}"
echo -e "${YELLOW} API Endpoint: http://localhost/api/tasks${NC}"
echo -e "${YELLOW} Network Info: http://localhost/api/network-info${NC}"
echo -e "\n${RED}Press Ctrl+C to stop the backend server${NC}"

# Wait for user to press Ctrl+C
wait $BACKEND_PID