#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source $SCRIPT_DIR/util.sh
cd $SCRIPT_DIR/..

database_password=`grep "grafana_metrics_db_password" etc/secrets.yaml | awk '{print $NF}'`
database_create_timeout=5

echo "Starting e2e test"
echo

echo "Make sure the 'grafana-metrics' database has been created in timescaledb"
test_database_created=`echo $(kubectl exec grafana-metrics-timescaledb-postgresql-0 -c postgresql -- psql postgresql://postgres:$database_password@localhost -l) | grep grafana-metrics | wc -l`
if [[ $test_database_created -eq 0 ]]; then
  echo "The 'grafana-metrics' database has not been created in timescaledb"
  exit 1
fi
echo "Success!!"

echo
echo "Waiting for the 'connect_fitbit_intraday_steps' table to be crated in the grafana database..."
timeout=0
test_table_created=0
while [[ $test_table_created -eq 0 ]]; do
  test_table_created=`echo $(kubectl exec grafana-metrics-timescaledb-postgresql-0 -c postgresql -- psql postgresql://postgres:$database_password@localhost/grafana-metrics -c '\dt') | grep connect_fitbit_intraday_steps | wc -l`
  timeout=$((timeout+1))
  if [[ $timeout -ge $database_create_timeout ]]; then
    echo "Failure: timeout reached after $database_create_timeout seconds"
    exit 1
  fi
  echo "Waiting for the 'connect_fitbit_intraday_steps' table to be crated in the grafana database..."
  sleep 1
done
echo "Success!!"

echo
echo "Waiting for events to be written to the 'connect_fitbit_intraday_steps' table..."
timeout=0
tets_table_populated=0
while [[ $test_table_populated -eq 0 ]]; do
  test_table_populated=`echo $(kubectl exec grafana-metrics-timescaledb-postgresql-0 -c postgresql -- psql postgresql://postgres:$database_password@localhost/grafana-metrics -c 'select COUNT(*) from connect_fitbit_intraday_steps;' -t)`
  timeout=$((timeout+1))
  if [[ $timeout -ge $database_create_timeout ]]; then
    echo "Failure: timeout reached after $database_create_timeout seconds"
    exit 1
  fi
  echo "Waiting for events to be written to the 'connect_fitbit_intraday_steps' table..."
  sleep 1
done
echo "Success!!"
