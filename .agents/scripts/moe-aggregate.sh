#!/bin/bash
# MoE Aggregate - Combines expert outputs and runs aggregator
#
# Usage:
#   source moe-aggregate.sh
#   # Set TASK_ID before calling
#   aggregate_moe_results "architect security performance maintainer"
#
# Requires: pi, tmux
# Environment: TASK_ID (timestamp)
# Outputs: Final result to stdout

aggregate_moe_results() {
  local experts="${1:-architect security performance maintainer}"
  
  if [ -z "$TASK_ID" ]; then
    echo "Error: TASK_ID variable must be set"
    return 1
  fi
  
  # Create aggregator prompt
  AGGREGATOR_INPUT=""
  for expert in $experts; do
    session="moe-${TASK_ID}-${expert}"
    output="/tmp/${session}.out"
    
    echo "=== $expert analysis ===" >> /tmp/moe-${TASK_ID}-combined.txt
    cat "$output" | grep -v "___EXPERT_DONE___" >> /tmp/moe-${TASK_ID}-combined.txt
    echo "" >> /tmp/moe-${TASK_ID}-combined.txt
  done
  
  # Spawn aggregator
  AGG_SESSION="moe-${TASK_ID}-aggregator"
  tmux new-session -d -s "$AGG_SESSION"
  tmux send-keys -t "$AGG_SESSION" \
    "cat /tmp/moe-${TASK_ID}-combined.txt | pi -p 'Synthesize these expert analyses into a unified recommendation. Identify conflicts, consensus areas, and priority actions. Structure as: 1) Summary, 2) Areas of Agreement, 3) Conflicts/Trade-offs, 4) Recommended Actions (prioritized)' > /tmp/moe-${TASK_ID}-final.out 2>&1 && echo '___AGGREGATOR_DONE___' >> /tmp/moe-${TASK_ID}-final.out" C-m
  
  # Wait for aggregator
  while ! grep -q "___AGGREGATOR_DONE___" /tmp/moe-${TASK_ID}-final.out 2>/dev/null; do
    sleep 2
  done
  
  # Display final result
  cat /tmp/moe-${TASK_ID}-final.out | grep -v "___AGGREGATOR_DONE___"
}

# If script is executed directly (not sourced), run with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  if [ -z "$TASK_ID" ]; then
    echo "Usage: TASK_ID=<timestamp> moe-aggregate.sh 'expert1 expert2'"
    exit 1
  fi
  aggregate_moe_results "$@"
fi
