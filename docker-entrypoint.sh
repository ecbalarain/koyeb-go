#!/bin/sh

# Docker entrypoint script for OXLOOK API

# Wait for database to be ready
wait_for_db() {
    echo "Waiting for database to be ready..."
    until nc -z -v -w30 $DATABASE_HOST $DATABASE_PORT
    do
        echo "Database is not ready yet..."
        sleep 2
    done
    echo "Database is ready!"
}

# Run migrations if needed
run_migrations() {
    echo "Running database migrations..."
    ./migrate
    if [ $? -ne 0 ]; then
        echo "Failed to run migrations!"
        exit 1
    fi
    echo "Migrations completed successfully."
}

# Start the application
start_app() {
    echo "Starting OXLOOK API..."
    ./main
}

# Main execution
main() {
    # Check if we need to run migrations
    if [ "$RUN_MIGRATIONS" = "true" ]; then
        wait_for_db
        run_migrations
    fi

    # Start the application
    start_app
}

# Run main function
main "$@"