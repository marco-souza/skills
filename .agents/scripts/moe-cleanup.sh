#!/bin/bash
# MoE Cleanup - Kills all expert sessions and removes temp files
#
# Usage:
#   source moe-cleanup.sh
#   # Set TASK_ID before calling
#   cleanup_moe "architect security performance maintainer"
#
# Requires: tmux
# Environment: TASK_ID (timestamp)

cleanup_moe() {
  local experts="${1:-architect security performance maintainer}"
  
  if [ -z "$TASK_ID" ]; then
    echo "Error: TASK_ID variable must be set"
    return 1
  fi
  
  echo "🧹 Cleaning up MoE sessions and temp files..."
  
  # Kill all expert sessions
  for expert in $experts; do
    tmux kill-session -t "moe-${TASK_ID}-${expert}" 2>/dev/null
  done
  
  # Kill aggregator session
  tmux kill-session -t "moe-${TASK_ID}-aggregator" 2>/dev/null
  
  # Remove temp files
  rm -f /tmp/moe-${TASK_ID}-*
  
  echo "Cleanup complete"
}

# If script is executed directly (not sourced), run with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  if [ -z "$TASK_ID" ]; then
    echo "Usage: TASK_ID=<timestamp> moe-cleanup.sh 'expert1 expert2'"
    exit 1
  fi
  cleanup_moe "$@"
fi
