#!/bin/bash
# MoE Spawn Experts - Spawns multiple expert agents in parallel
#
# Usage:
#   source moe-spawn.sh
#   # Set PROBLEM and TASK_ID before calling
#   spawn_moe_experts "architect security performance maintainer"
#
# Requires: tmux, pi
# Environment: PROBLEM (string), TASK_ID (timestamp)

spawn_moe_experts() {
  local experts="${1:-architect security performance maintainer}"
  
  if [ -z "$PROBLEM" ]; then
    echo "Error: PROBLEM variable must be set"
    return 1
  fi
  
  if [ -z "$TASK_ID" ]; then
    TASK_ID=$(date +%s)
  fi
  
  for expert in $experts; do
    session="moe-${TASK_ID}-${expert}"
    output="/tmp/${session}.out"
    
    # Create detached session
    tmux new-session -d -s "$session"
    
    # Build expert-specific prompt
    case $expert in
      architect)
        SYSTEM="You are a software architect. Analyze the design patterns, coupling, and long-term maintainability."
        ;;
      security)
        SYSTEM="You are a security engineer. Find vulnerabilities, injection risks, and auth flaws."
        ;;
      performance)
        SYSTEM="You are a performance engineer. Identify bottlenecks and optimization opportunities."
        ;;
      maintainer)
        SYSTEM="You are a senior maintainer. Evaluate readability, documentation, and testing."
        ;;
      minimalist)
        SYSTEM="You are a minimalist engineer. Find unnecessary complexity and YAGNI violations."
        ;;
      *)
        SYSTEM="You are an expert in your field. Analyze from your perspective."
        ;;
    esac
    
    # Spawn pi with expert system prompt
    tmux send-keys -t "$session" \
      "pi --system-prompt '$SYSTEM' -p '$PROBLEM' > $output 2>&1 && echo '___EXPERT_DONE___' >> $output" C-m
    
    echo "  → $expert spawned in session $session"
  done
}

# If script is executed directly (not sourced), run with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  if [ -z "$PROBLEM" ]; then
    echo "Usage: PROBLEM='your problem' moe-spawn.sh 'expert1 expert2'"
    exit 1
  fi
  spawn_moe_experts "$@"
fi
