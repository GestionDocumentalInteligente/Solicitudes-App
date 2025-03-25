#!/bin/sh

log() {
    echo "ENTRYPOINT: $1"
}

# Solo validar, no cargar
validateEnv() {
    if [ -z "${APP_NAME}" ]; then
        log "ERROR: APP_NAME not set"
        exit 1
    fi
}

cleanup() {
    log "Cleaning up processes"
    pkill -f dlv || true
    pkill -f "${APP_NAME}" || true
}

main() {
    validateEnv
    cleanup
    
    if [ "${STAGE}" = "dev" ]; then
        log "Development Mode: ${APP_NAME}"
        air -c "$AIR_CONFIG"
    else
        log "Production mode"
        ./${APP_NAME}
    fi
}

main