#!/bin/bash

# Function to log messages
log_message() {
  echo -e "\n$1"
}

# Function to handle errors
handle_error() {
  if [ $? -ne 0 ]; then
    log_message "Error: $1 failed."
    exit 1
  fi
}

# Step 1: Add the first organization
log_message "Adding Organization One..."
response=$(curl -s -w "%{http_code}" -o /dev/null -X POST http://localhost:8080/orgs \
  -H "Content-Type: application/json" \
  -d '{
  "id": "org1",
  "name": "Organization One",
  "phone": "123-456-7890",
  "location" : {
    "latitude": 40.7128,
    "longitude": -74.0060
  }
}')

if [ "$response" -eq 201 ]; then
  log_message "Successfully added Organization One."
else
  handle_error "Adding Organization One"
fi

# Step 2: Add the second organization
log_message "Adding Organization Two..."
response=$(curl -s -w "%{http_code}" -o /dev/null -X POST http://localhost:8080/orgs \
  -H "Content-Type: application/json" \
  -d '{
  "id": "org2",
  "name": "Organization Two",
  "phone": "098-765-4321",
  "location":{
    "latitude": 34.0522,
    "longitude": -118.2437
  }
}')

if [ "$response" -eq 201 ]; then
  log_message "Successfully added Organization Two."
else
  handle_error "Adding Organization Two"
fi

# Step 3: Add service "1" to the first organization
log_message "Adding Service 1 to Organization One..."
response=$(curl -s -w "%{http_code}" -o /dev/null -X POST http://localhost:8080/orgs/org1/services \
  -H "Content-Type: application/json" \
  -d '["1"]')

if [ "$response" -eq 200 ]; then
  log_message "Successfully added Service 1 to Organization One."
else
  handle_error "Adding Service 1 to Organization One"
fi

# Step 4: Add service "2" to the second organization
log_message "Adding Service 2 to Organization Two..."
response=$(curl -s -w "%{http_code}" -o /dev/null -X POST http://localhost:8080/orgs/org2/services \
  -H "Content-Type: application/json" \
  -d '["2"]')

if [ "$response" -eq 200 ]; then
  log_message "Successfully added Service 2 to Organization Two."
else
  handle_error "Adding Service 2 to Organization Two"
fi

# Step 5: Test finding the closest organization offering both services
log_message "Finding closest organization offering both services (1 and 2)..."
response=$(curl -s -X GET http://localhost:8080/services/nearest \
  -H "Content-Type: application/json" \
  -d '{
  "services": ["1", "2"],
  "latitude": 37.7749,
  "longitude": -122.4194
}')

log_message "Response for both services (1 and 2): $response"

# Step 6: Test finding the closest organization offering only service 1
log_message "Finding closest organization offering only Service 1..."
response=$(curl -s -X GET http://localhost:8080/services/nearest \
  -H "Content-Type: application/json" \
  -d '{
  "services": ["1"],
  "latitude": 37.7749,
  "longitude": -122.4194
}')

log_message "Response for only Service 1: $response"

# Step 7: Test finding the closest organization offering only service 2
log_message "Finding closest organization offering only Service 2..."
response=$(curl -s -X GET http://localhost:8080/services/nearest \
  -H "Content-Type: application/json" \
  -d '{
  "services": ["2"],
  "latitude": 37.7749,
  "longitude": -122.4194
}')

log_message "Response for only Service 2: $response"

# Step 8: Test finding an organization that doesn't have both services
log_message "Testing organization that doesn't have both services (expecting no matches)..."
response=$(curl -s -X GET http://localhost:8080/services/nearest \
  -H "Content-Type: application/json" \
  -d '{
  "services": ["1", "3"], 
  "latitude": 37.7749,
  "longitude": -122.4194
}')

log_message "Response for non-existing service combination (1 and 3): $response"
