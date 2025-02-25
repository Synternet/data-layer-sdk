#!/bin/bash

NSC_HOME=./.nats

# Define exported subjects as associative arrays
# Format: [account]=("subject1" "subject2" ...)
declare -A EXPORTED_SERVICES=(
  [publisher]="synternet.rpc-example.users.service.*"
)

declare -A EXPORTED_STREAMS=(
  [publisher]="synternet.rpc-example.users.registrations"
)

# Define import relationships (importer -> exporter)
declare -A IMPORT_RELATIONSHIPS=(
  [subscriber]=publisher
)

# Define imported subjects explicitly per importer
declare -A IMPORTED_SERVICES=(
  [subscriber]="synternet.rpc-example.users.service.*"
)

declare -A IMPORTED_STREAMS=(
  [subscriber]="synternet.rpc-example.users.registrations"
)

# Function to configure exports (both services and streams)
configure_exports() {
  local account=$1
  shift
  local exports_services=($1)
  local exports_streams=($2)

  # Export services
  for subject in ${exports_services[@]}; do
    nsc add export --account "$account" --subject "$subject" --service --all-dirs "$NSC_HOME"
  done

  # Export streams
  for subject in ${exports_streams[@]}; do
    nsc add export --account "$account" --subject "$subject" --all-dirs "$NSC_HOME"
  done
}

# Function to configure imports (both services and streams)
configure_imports() {
  local account=$1
  local provider=$2

  # Import services (use specific subjects if defined, otherwise import all from provider)
  local import_services=(${IMPORTED_SERVICES[$account]:-${EXPORTED_SERVICES[$provider]}})
  for subject in ${import_services[@]}; do
    nsc add import --account "$account" --src-account "$provider" --remote-subject "$subject" --service --all-dirs "$NSC_HOME"
  done

  # Import streams (use specific subjects if defined, otherwise import all from provider)
  local import_streams=(${IMPORTED_STREAMS[$account]:-${EXPORTED_STREAMS[$provider]}})
  for subject in ${import_streams[@]}; do
    nsc add import --account "$account" --src-account "$provider" --remote-subject "$subject" --all-dirs "$NSC_HOME"
  done
}

# Function to configure a publisher
configure_publisher() {
  local account=$1
  local exports_services=(${EXPORTED_SERVICES[$account]})
  local exports_streams=(${EXPORTED_STREAMS[$account]})

  nsc add account --name "$account" --all-dirs "$NSC_HOME"
  configure_exports "$account" "${exports_services[@]}" "${exports_streams[@]}"
}

# Initialize NSC environment
rm -rf "$NSC_HOME"
mkdir -p "$NSC_HOME"
nsc init --all-dirs "$NSC_HOME" --dir "$NSC_HOME" --name test
nsc env -s "$NSC_HOME" --all-dirs "$NSC_HOME"

# Configure publisher accounts
for account in "${!EXPORTED_SERVICES[@]}"; do
  echo Configuring publisher: $account
  configure_publisher "$account"
  nsc add user --name "$account" --account "$account" --all-dirs "$NSC_HOME"
done

# Configure subscriber accounts
for account in "${!IMPORT_RELATIONSHIPS[@]}"; do
  provider=${IMPORT_RELATIONSHIPS[$account]}
  echo Configuring Subscriber: $account, $provider
  nsc add account --name "$account" --all-dirs "$NSC_HOME"
  configure_imports "$account" "$provider"
  nsc add user --name "$account" --account "$account" --all-dirs "$NSC_HOME"
done

# Generate NATS config file
nsc generate config --mem-resolver --config-file "$NSC_HOME/server.conf" --all-dirs "$NSC_HOME"
