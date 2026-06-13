#!/usr/bin/env bun
/**
 * learning-search.ts — Search learnings by keyword or tag
 *
 * Usage:
 *   bun run learning-search.ts <query>
 *   bun run learning-search.ts --tag <tag>
 *   bun run learning-search.ts --all
 *
 * Searches both frontmatter (title, tags) and content body.
 */

import { readdirSync, readFileSync } from "fs";
import { join, resolve } from "path";

const LEARNINGS_DIR = resolve(import.meta.dir, "../learnings");

interface SearchResult {
  filename: string;
  title: string;
  date: string;
  tags: string[];
  matches: string[];
}

function parseArgs(
  args: string[]
): { query?: string; tag?: string; all?: boolean } {
  const options: { query?: string; tag?: string; all?: boolean } = {};
  let i = 2;

  while (i < args.length) {
    const arg = args[i];
    if (arg === "--tag" || arg === "-t") {
      i++;
      if (i < args.length) {
        options.tag = args[i].toLowerCase();
      }
    } else if (arg === "--all" || arg === "-a") {
      options.all = true;
    } else if (!arg.startsWith("-")) {
      options.query = arg;
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

      if (value.startsWith("[") && value.endsWith("]")) {
        value = value.slice(1, -1);
      }

      frontmatter[key] = value;
    }
  }

  return frontmatter;
}

function extractDateFromFilename(filename: string): string {
  const match = filename.match(/^(\d{4}-\d{2}-\d{2})/);
  return match ? match[1] : "0000-00-00";
}

function searchLearnings(
  query: string | undefined,
  tag: string | undefined
): SearchResult[] {
  try {
    const files = readdirSync(LEARNINGS_DIR).filter((f) => f.endsWith(".md"));
    const results: SearchResult[] = [];

    for (const filename of files) {
      const content = readFileSync(join(LEARNINGS_DIR, filename), "utf-8");
      const meta = extractFrontmatter(content);

      const title = meta.title || filename.replace(/\.md$/, "");
      const date = meta.date || extractDateFromFilename(filename);
      const tags = meta.tags
        ? meta.tags.split(",").map((t) => t.trim())
        : [];

      // Check tag filter
      if (tag && !tags.some((t) => t.includes(tag))) {
        continue;
      }

      // Check query filter
      if (query) {
        const lowerQuery = query.toLowerCase();
        const matches: string[] = [];

        // Search in title
        if (title.toLowerCase().includes(lowerQuery)) {
          matches.push(`title: ${title}`);
        }

        // Search in tags
        if (tags.some((t) => t.includes(lowerQuery))) {
          matches.push(`tags: ${tags.join(", ")}`);
        }

        // Search in content body (after frontmatter)
        const bodyContent = content.replace(/^---[\s\S]*?---\n?/, "");
        const lines = bodyContent.split("\n");
        for (const line of lines) {
          if (line.toLowerCase().includes(lowerQuery)) {
            matches.push(line.trim());
          }
        }

        if (matches.length === 0) continue;

        results.push({ filename, title, date, tags, matches: matches.slice(0, 3) });
      } else {
        results.push({ filename, title, date, tags, matches: [] });
      }
    }

    return results.sort((a, b) => b.date.localeCompare(a.date));
  } catch {
    return [];
  }
}

function highlight(text: string, query: string): string {
  const regex = new RegExp(`(${query})`, "gi");
  return text.replace(regex, "\x1b[33m$1\x1b[0m");
}

function main(): void {
  const options = parseArgs(process.argv);

  if (!options.query && !options.tag && !options.all) {
    console.log("Usage:");
    console.log("  bun run learning-search.ts <query>       Search by keyword");
    console.log("  bun run learning-search.ts --tag <tag>   Search by tag");
    console.log("  bun run learning-search.ts --all         List all learnings");
    process.exit(1);
  }

  if (options.all) {
    // Delegate to learning-list.ts
    const { readdirSync } = require("fs");
    const files = readdirSync(LEARNINGS_DIR).filter((f: string) =>
      f.endsWith(".md")
    );
    console.log(`📚 All learnings (${files.length} total)\n`);
    for (const f of files.sort().reverse()) {
      console.log(`  ${f}`);
    }
    return;
  }

  const results = searchLearnings(options.query, options.tag);

  if (results.length === 0) {
    console.log("🔍 No learnings found.");
    if (options.query) {
      console.log(`\nSearch query: "${options.query}"`);
    }
    if (options.tag) {
      console.log(`\nTag filter: "${options.tag}"`);
    }
    console.log("\nTry different search terms or create a new learning:");
    console.log("  bun run learning-create.ts");
    return;
  }

  console.log(`🔍 Found ${results.length} learning(s)\n`);

  for (const result of results) {
    const tags = result.tags.map((t) => `\x1b[36m${t}\x1b[0m`).join(" ");
    const title = options.query
      ? highlight(result.title, options.query)
      : result.title;

    console.log(`  \x1b[33m${result.date}\x1b[0m ${title}`);
    console.log(`    📄 ${result.filename}`);
    if (tags) console.log(`    🏷️  ${tags}`);

    if (result.matches.length > 0) {
      console.log("    📝 Matches:");
      for (const match of result.matches.slice(0, 3)) {
        const highlighted = options.query ? highlight(match, options.query) : match;
        console.log(`       ${highlighted.slice(0, 80)}`);
      }
    }
    console.log();
  }
}

main();
