#!/usr/bin/env bun
/**
 * learning-show.ts — Display a specific learning entry
 *
 * Usage:
 *   bun run learning-show.ts <filename>
 *   bun run learning-show.ts --latest
 *
 * Shows the full content of a learning file with syntax highlighting.
 */

import { readFileSync, readdirSync } from "fs";
import { join, resolve } from "path";

const LEARNINGS_DIR = resolve(import.meta.dir, "../learnings");

function parseArgs(args: string[]): { filename?: string; latest?: boolean } {
  const options: { filename?: string; latest?: boolean } = {};
  let i = 2;

  while (i < args.length) {
    const arg = args[i];
    if (arg === "--latest" || arg === "-l") {
      options.latest = true;
    } else if (!arg.startsWith("-")) {
      options.filename = arg;
    }
    i++;
  }

  return options;
}

function getLatestLearning(): string | null {
  try {
    const files = readdirSync(LEARNINGS_DIR)
      .filter((f) => f.endsWith(".md"))
      .sort()
      .reverse();
    return files[0] || null;
  } catch {
    return null;
  }
}

function formatContent(content: string): string {
  // Basic Markdown formatting for terminal
  return content
    .replace(/^---\n[\s\S]*?\n---\n?/, "") // Remove frontmatter
    .replace(/^# (.+)$/gm, "\x1b[1;37m# $1\x1b[0m") // H1 bold white
    .replace(/^## (.+)$/gm, "\n\x1b[1;36m## $1\x1b[0m") // H2 cyan
    .replace(/^### (.+)$/gm, "\n\x1b[1;35m### $1\x1b[0m") // H3 magenta
    .replace(/`([^`]+)`/g, "\x1b[33m`$1`\x1b[0m") // Inline code yellow
    .replace(/```(\w+)?\n([\s\S]*?)```/g, (_, lang, code) => {
      // Code blocks
      const language = lang || "";
      return `\n\x1b[32m[${language}]\x1b[0m\n${code}`;
    })
    .replace(/\*\*(.+?)\*\*/g, "\x1b[1m$1\x1b[0m") // Bold
    .replace(/\*(.+?)\*/g, "\x1b[3m$1\x1b[0m") // Italic
    .replace(/- (.+)$/gm, "  • $1") // Bullet points
    .replace(/^\d+\. (.+)$/gm, (match) => `  ${match}`); // Numbered lists
}

function main(): void {
  const options = parseArgs(process.argv);

  let filename = options.filename;

  if (!filename && options.latest) {
    filename = getLatestLearning();
    if (!filename) {
      console.log("📚 No learnings found.");
      process.exit(1);
    }
    console.log(`📚 Latest learning: ${filename}\n`);
  }

  if (!filename) {
    console.log("Usage:");
    console.log("  bun run learning-show.ts <filename>    Show a specific learning");
    console.log("  bun run learning-show.ts --latest      Show the most recent learning");
    console.log("\nAvailable learnings:");

    try {
      const files = readdirSync(LEARNINGS_DIR)
        .filter((f) => f.endsWith(".md"))
        .sort()
        .reverse();

      if (files.length === 0) {
        console.log("  (none)");
      } else {
        for (const f of files.slice(0, 10)) {
          console.log(`  ${f}`);
        }
        if (files.length > 10) {
          console.log(`  ... and ${files.length - 10} more`);
        }
      }
    } catch {
      console.log("  (none)");
    }

    process.exit(1);
  }

  // Handle partial filenames - search for match
  let filepath = join(LEARNINGS_DIR, filename);
  const { readFileSync, existsSync } = require("fs");

  if (!existsSync(filepath)) {
    // Try to find a matching file
    try {
      const files = readdirSync(LEARNINGS_DIR).filter((f: string) =>
        f.includes(filename!)
      );
      if (files.length === 1) {
        filepath = join(LEARNINGS_DIR, files[0]);
        filename = files[0];
      } else if (files.length > 1) {
        console.log(`❌ Multiple matches found for "${filename}":`);
        for (const f of files) {
          console.log(`  ${f}`);
        }
        process.exit(1);
      } else {
        console.log(`❌ Learning not found: "${filename}"`);
        process.exit(1);
      }
    } catch {
      console.log(`❌ Learning not found: "${filename}"`);
      process.exit(1);
    }
  }

  const content = readFileSync(filepath, "utf-8");
  const formatted = formatContent(content);

  console.log(formatted);
}

main();
