#!/usr/bin/env bash
set -euo pipefail

# Usage:
# ./load-test.sh [URL] [TOTAL_REQUESTS] [CONCURRENCY]
# Example:
# ./load-test.sh http://localhost:8080/service-a/ 100 20

URL="${1:-http://localhost:8080/service-a/}"
TOTAL_REQUESTS="${2:-20}"
CONCURRENCY="${3:-5}"

# Optional settings
MAX_RETRIES=2          # number of retries on failure
SHOW_BODY=false        # set to true to print full response body

echo "Running load test:"
echo "→ URL: $URL"
echo "→ Total requests: $TOTAL_REQUESTS"
echo "→ Concurrency: $CONCURRENCY"
echo "→ Retries on failure: $MAX_RETRIES"
echo

sem() {
  local max_jobs=$1
  while [ "$(jobs -rp | wc -l)" -ge "$max_jobs" ]; do
    sleep 0.1
  done
}

for i in $(seq 1 "$TOTAL_REQUESTS"); do
  sem "$CONCURRENCY"
  {
    retries=0
    start_time=$(date +%s.%N)
    while true; do
      if $SHOW_BODY; then
        response=$(curl -sS -w "\n%{http_code}" "$URL")
        body=$(echo "$response" | head -n -1)
        http_code=$(echo "$response" | tail -n1)
      else
        http_code=$(curl -sS -o /dev/null -w "%{http_code}" "$URL" || echo "000")
      fi

      if [[ "$http_code" != "000" ]]; then
        break
      elif (( retries < MAX_RETRIES )); then
        ((retries++))
        sleep 0.2
      else
        break
      fi
    done

    end_time=$(date +%s.%N)
    elapsed=$(awk "BEGIN{printf \"%.3f\", $end_time - $start_time}")

    printf "req %3d: %s  %5ss (retries=%d)\n" "$i" "$http_code" "$elapsed" "$retries"
    if $SHOW_BODY; then
      echo "--- response body ---"
      echo "$body"
      echo "---------------------"
    fi
  } &
done

wait
echo
echo "✅ Completed $TOTAL_REQUESTS requests with concurrency $CONCURRENCY."