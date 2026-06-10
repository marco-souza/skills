#!/bin/bash
# MoE Wait - Waits for all expert sessions to complete
#
# Usage:
#   source moe-wait.sh
#   # Set TASK_ID before calling
#   wait_moe_experts "architect security performance maintainer"
#
# Requires: tmux
# Environment: TASK_ID (timestamp)

wait_moe_experts() {
  local experts="${1:-architect security performance maintainer}"
  
  if [ -z "$TASK_ID" ]; then
    echo "Error: TASK_ID variable must be set"
    return 1
  fi
  
  echo "⏳ Waiting for all experts..."
  for expert in $experts; do
    session="moe-${TASK_ID}-${expert}"
    output="/tmp/${session}.out"
    
    echo "Waiting for $expert..."
    while ! grep -q "___EXPERT_DONE___" "$output" 2>/dev/null; do
      sleep 2
    done
    echo "  ✓ $expert complete"
  done
  echo "All experts complete"
}

# If script is executed directly (not sourced), run with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  if [ -z "$TASK_ID" ]; then
    echo "Usage: TASK_ID=<timestamp> moe-wait.sh 'expert1 expert2'"
    exit 1
  fi
  wait_moe_experts "$@"
fi
