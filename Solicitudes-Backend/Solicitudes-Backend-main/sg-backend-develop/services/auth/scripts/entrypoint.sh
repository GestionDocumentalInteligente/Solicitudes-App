#!/bin/sh

# shellcheck disable=SC2154  # Desactivar aviso de shellcheck para APP_NAME
# shellcheck source=./config/.env disable=SC1091 # Desactivar aviso de shellcheck de archivo no especificado

# Function to log messages
log() {
  echo "ENTRYPOINT: $1"
}

# Load environment variables from a file without overwriting existing ones
loadEnvFile() {
  ENV_FILE="$1"
  if [ -f "$ENV_FILE" ]; then
    log "Loading environment variables from $ENV_FILE"
    # Read each line in the .env file
    while IFS= read -r line || [ -n "$line" ]; do
      # Ignore empty lines and comments
      if [ -n "$line" ] && [ "${line#\#}" = "$line" ]; then
        VAR_NAME=$(echo "$line" | cut -d '=' -f1)
        VAR_VALUE=$(echo "$line" | cut -d '=' -f2-)
        # Check if variable is already set
        if [ -z "$(printenv "$VAR_NAME")" ]; then
          export "$VAR_NAME=$VAR_VALUE"
        else
          log "Variable $VAR_NAME is already set to $(printenv "$VAR_NAME"), not overwriting"
        fi
      fi
    done < "$ENV_FILE"
  else
    log "WARNING: $ENV_FILE file not found."
  fi
}

# Validate essential environment variables
validateEnv() {
  if [ -z "${APP_NAME}" ]; then
    log "ERROR: APP_NAME is not set. Please check your .env files"
    exit 1
  fi

  if [ -z "${DEBUG}" ]; then
    log "ERROR: DEBUG is not set. Please check your .env files"
    exit 1
  fi

  log "Environment variables loaded successfully"
  log "App Name: ${APP_NAME}"
  log "Debug: ${DEBUG}"
}

# Function to initialize the file change logger
initializeFileChangeLogger() {
  echo "" > /tmp/filechanges.log
  tail -f /tmp/filechanges.log &
}

# Function to run the server
runServer() {
  log "Running service"

  # Kill any existing server processes
  log "Killing old processes"
  pkill -f dlv || true
  pkill -f "/app/tmp/${APP_NAME}" || true

  if [ "${DEBUG}" = "true" ]; then
    log "Running in debug mode with Air and Delve"
    air -c "$AIR_CONFIG"
  else
    log "Running in production mode"
    # Aquí puedes agregar el comando para ejecutar tu aplicación en modo producción
    # Por ejemplo: ./app
  fi
}

# Main function to orchestrate the process
main() {
  log "Starting script"
  log "Current directory: $(pwd)"

  # Load environment variables from .env and .env.local
  loadEnvFile "./config/.env"
  loadEnvFile "./config/.env.local"

  validateEnv
  initializeFileChangeLogger

  # Start server with Air and possibly Delve
  runServer
}

# Call the main function to start the process
main
