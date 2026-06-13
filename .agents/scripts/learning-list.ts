#!/usr/bin/env bun
/**
 * learning-list.ts — List all learning entries
 *
 * Usage:
 *   bun run learning-list.ts [--recent N] [--tag tag]
 *
 * Lists all learnings in .agents/learnings/ sorted by date (newest first).
 */

import { readdirSync, readFileSync } from "fs";
import { join, resolve } from "path";

const LEARNINGS_DIR = resolve(import.meta.dir, "../learnings");

interface LearningMeta {
  filename: string;
  title: string;
  date: string;
  tags: string[];
}

function parseArgs(args: string[]): { recent?: number; tag?: string } {
  const options: { recent?: number; tag?: string } = {};
  let i = 2;

  while (i < args.length) {
    const arg = args[i];
    if (arg === "--recent" || arg === "-n") {
      i++;
      if (i < args.length) {
        options.recent = parseInt(args[i], 10);
      }
    } else if (arg === "--tag" || arg === "-t") {
      i++;
      if (i < args.length) {
        options.tag = args[i].toLowerCase();
      }
    }
    i++;
  }

  return options;
}

function extractFrontmatter(content: string): Record<string, string> {
  const match = content.match(/^---\n([\s\S]*?)\n---/);
  if (!match) return {};

  const frontmatter: Record<string, string> = {};
  const lines = match[1].split("\n");

  for (const line of lines) {
    const colonIndex = line.indexOf(":");
    if (colonIndex > 0) {
      const key = line.slice(0, colonIndex).trim();
      let value = line.slice(colonIndex + 1).trim();

      // Handle YAML array syntax
      if (value.startsWith("[") && value.endsWith("]")) {
        value = value.slice(1, -1);
      }

      frontmatter[key] = value;
    }
  }

  return frontmatter;
}

function getLearnings(): LearningMeta[] {
  try {
    const files = readdirSync(LEARNINGS_DIR).filter((f) => f.endsWith(".md"));

    return files
      .map((filename) => {
        const content = readFileSync(join(LEARNINGS_DIR, filename), "utf-8");
        const meta = extractFrontmatter(content);

        return {
          filename,
          title: meta.title || filename.replace(/\.md$/, ""),
          date: meta.date || extractDateFromFilename(filename),
          tags: meta.tags
            ? meta.tags.split(",").map((t) => t.trim())
            : [],
        };
      })
      .sort((a, b) => b.date.localeCompare(a.date));
  } catch {
    return [];
  }
}

function extractDateFromFilename(filename: string): string {
  const match = filename.match(/^(\d{4}-\d{2}-\d{2})/);
  return match ? match[1] : "0000-00-00";
}

function formatTags(tags: string[]): string {
  if (tags.length === 0) return "";
  return tags.map((t) => `\x1b[36m${t}\x1b[0m`).join(" ");
}

function main(): void {
  const options = parseArgs(process.argv);
  let learnings = getLearnings();

  // Filter by tag
  if (options.tag) {
    learnings = learnings.filter((l) =>
      l.tags.some((t) => t.includes(options.tag!))
    );
  }

  // Limit results
  if (options.recent) {
    learnings = learnings.slice(0, options.recent);
  }

  if (learnings.length === 0) {
    console.log("📚 No learnings found.");
    console.log("\nCreate one with: bun run learning-create.ts");
    return;
  }

  console.log(`📚 Learnings (${learnings.length} total)\n`);

  for (const learning of learnings) {
    const tags = formatTags(learning.tags);
    const tagsStr = tags ? ` ${tags}` : "";
    console.log(`  \x1b[33m${learning.date}\x1b[0m ${learning.title}${tagsStr}`);
    console.log(`    ${learning.filename}`);
  }

  console.log("\nCommands:");
  console.log("  bun run learning-show.ts <filename>    View a learning");
  console.log("  bun run learning-search.ts <query>    Search learnings");
}

main();
