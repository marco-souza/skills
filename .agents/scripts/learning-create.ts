#!/usr/bin/env bun
/**
 * learning-create.ts — Create a new learning entry
 *
 * Usage:
 *   bun run learning-create.ts [title] [--tag tag1,tag2]
 *
 * If no title is provided, prompts interactively.
 * Creates a dated Markdown file in .agents/learnings/
 */

import { mkdirSync, writeFileSync, readdirSync } from "fs";
import { join, resolve } from "path";

const LEARNINGS_DIR = resolve(import.meta.dir, "../learnings");

interface LearningOptions {
  title?: string;
  tags?: string[];
  interactive?: boolean;
}

function parseArgs(args: string[]): LearningOptions {
  const options: LearningOptions = { tags: [] };
  let i = 2; // skip bun and script path

  while (i < args.length) {
    const arg = args[i];
    if (arg === "--tag" || arg === "-t") {
      i++;
      if (i < args.length) {
        options.tags = args[i].split(",").map((t) => t.trim().toLowerCase());
      }
    } else if (!arg.startsWith("-")) {
      options.title = arg;
    }
    i++;
  }

  options.interactive = !options.title;
  return options;
}

function formatDate(date: Date = new Date()): string {
  return date.toISOString().split("T")[0];
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/(^-|-$)/g, "");
}

function getExistingFiles(): string[] {
  try {
    return readdirSync(LEARNINGS_DIR).filter((f) => f.endsWith(".md"));
  } catch {
    return [];
  }
}

function generateFilename(title: string, date: string): string {
  const slug = slugify(title);
  const filename = `${date}-${slug}.md`;
  const existing = getExistingFiles();

  // Handle duplicates by appending a number
  if (existing.includes(filename)) {
    let counter = 2;
    while (existing.includes(`${date}-${slug}-${counter}.md`)) {
      counter++;
    }
    return `${date}-${slug}-${counter}.md`;
  }

  return filename;
}

function generateTemplate(
  title: string,
  tags: string[],
  date: string
): string {
  const tagsYaml = tags.length > 0 ? `[${tags.join(", ")}]` : "[add-tags]";

  return `---
title: ${title}
date: ${date}
tags: ${tagsYaml}
---

# ${title}

## What I Learned

[Clear, concise statement of the learning]

## Context

[When and why this was discovered — what problem were you solving?]

## Solution / Pattern

[The actual solution, pattern, or insight — make it actionable]

## Application

[How to apply this in the future — include code examples if helpful]

## References

- Related files: \`path/to/file.ts\`
- External docs: [link](url)
`;
}

async function prompt(question: string): Promise<string> {
  const { stdout, stdin } = process;
  return new Promise((resolve) => {
    stdout.write(question);
    let data = "";
    stdin.on("data", (chunk) => {
      data += chunk.toString();
      if (data.includes("\n")) {
        resolve(data.trim());
        stdin.removeAllListeners();
      }
    });
  });
}

async function interactiveCreate(): Promise<void> {
  console.log("🎓 Create New Learning\n");

  const title = await prompt("Title: ");
  if (!title) {
    console.error("❌ Title is required");
    process.exit(1);
  }

  const tagsInput = await prompt("Tags (comma-separated): ");
  const tags = tagsInput
    ? tagsInput.split(",").map((t) => t.trim().toLowerCase())
    : [];

  const date = formatDate();
  const filename = generateFilename(title, date);
  const filepath = join(LEARNINGS_DIR, filename);
  const content = generateTemplate(title, tags, date);

  mkdirSync(LEARNINGS_DIR, { recursive: true });
  writeFileSync(filepath, content);

  console.log(`\n✅ Created learning: ${filename}`);
  console.log(`📁 Location: ${filepath}`);
  console.log("\nNext steps:");
  console.log(`  1. Edit the file to fill in your learning`);
  console.log(`  2. Use \`bun run learning-search.ts\` to find it later`);
}

async function main(): Promise<void> {
  const options = parseArgs(process.argv);

  if (options.interactive) {
    await interactiveCreate();
    return;
  }

  const title = options.title!;
  const tags = options.tags || [];
  const date = formatDate();
  const filename = generateFilename(title, date);
  const filepath = join(LEARNINGS_DIR, filename);
  const content = generateTemplate(title, tags, date);

  mkdirSync(LEARNINGS_DIR, { recursive: true });
  writeFileSync(filepath, content);

  console.log(`✅ Created learning: ${filename}`);
  console.log(`📁 Location: ${filepath}`);
}

main().catch((err) => {
  console.error("❌ Error:", err.message);
  process.exit(1);
});
