#!/bin/bash
# MoE Code Review - Complete code review workflow using Mixture of Experts
#
# Usage:
#   ./moe-code-review.sh <file-to-review>
#   ./moe-code-review.sh src/auth/login.ts
#
# Requires: pi, tmux, cat
# Environment: None required (sets own TASK_ID)

FILE="${1:-}"
TASK_ID=$(date +%s)
EXPERTS="architect security performance maintainer"

if [ -z "$FILE" ]; then
  echo "Usage: moe-code-review.sh <file-to-review>"
  echo "Example: moe-code-review.sh src/auth/login.ts"
  exit 1
fi

if [ ! -f "$FILE" ]; then
  echo "Error: File '$FILE' not found"
  exit 1
fi

echo "🔍 Starting Mixture of Experts review of $FILE..."

# Read file content
FILE_CONTENT=$(cat "$FILE")
PROBLEM="Review this code for issues and improvements:\n\n\`\`\`typescript\n$FILE_CONTENT\n\`\`\`"

# Spawn experts
for expert in $EXPERTS; do
  session="moe-${TASK_ID}-${expert}"
  tmux new-session -d -s "$session"
  
  case $expert in
    architect)
      PROMPT="As an architect, analyze: design patterns, separation of concerns, extensibility, and technical debt. $PROBLEM"
      ;;
    security)
      PROMPT="As a security engineer, find: vulnerabilities, injection risks, auth bypasses, and data leaks. $PROBLEM"
      ;;
    performance)
      PROMPT="As a performance engineer, identify: bottlenecks, unnecessary computations, memory issues, and N+1 queries. $PROBLEM"
      ;;
    maintainer)
      PROMPT="As a maintainer, evaluate: readability, comments, test coverage, error handling, and debuggability. $PROBLEM"
      ;;
  esac
  
  tmux send-keys -t "$session" "pi -p '$PROMPT' > /tmp/${session}.out 2>&1 && echo DONE >> /tmp/${session}.out" C-m
  echo "  → $expert spawned"
done

# Wait for completion
echo "⏳ Waiting for all experts..."
for expert in $EXPERTS; do
  while ! grep -q "DONE" "/tmp/moe-${TASK_ID}-${expert}.out" 2>/dev/null; do
    sleep 1
  done
  echo "  ✓ $expert complete"
done

# Aggregate
echo "🧠 Synthesizing insights..."
COMBINED=""
for expert in $EXPERTS; do
  OUT=$(cat "/tmp/moe-${TASK_ID}-${expert}.out" | grep -v "DONE")
  COMBINED="$COMBINED\n\n=== $expert ===\n$OUT"
done

echo -e "$COMBINED" | pi -p 'Synthesize these expert reviews into actionable recommendations. Format: Consensus (agreed by 3+ experts), Important (2 experts), Worth Considering (1 expert). Then give top 3 priority fixes.'

# Cleanup
for expert in $EXPERTS; do
  tmux kill-session -t "moe-${TASK_ID}-${expert}" 2>/dev/null
done
rm -f /tmp/moe-${TASK_ID}-*

echo "✅ MoE review complete"
