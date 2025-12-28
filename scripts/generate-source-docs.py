#!/usr/bin/env python3
"""
Generate source build documentation with real CLI output.

Usage:
    python scripts/generate-source-docs.py \
        --node-version "22" \
        --prereq-output "$(cat prereq.txt)" \
        --clone-output "$(cat clone.txt)" \
        --install-output "$(cat install.txt)" \
        --build-output "$(cat build.txt)" \
        --version-output "$(cat version.txt)" \
        --output docs/install/source.md
"""

import argparse
import re
from datetime import UTC, datetime
from pathlib import Path

from jinja2 import Environment, FileSystemLoader

from naming import normalize_acronyms, to_human_readable, to_title_case


def clean_output(text):
    """Remove ANSI codes and trim whitespace."""
    if not text:
        return text
    # Remove ANSI escape codes
    text = re.sub(r"\x1b\[[0-9;]*m", "", text)
    return text.strip()


def main():
    parser = argparse.ArgumentParser(description="Generate source build docs with real output")
    parser.add_argument("--node-version", help="Node.js version used for build")
    parser.add_argument("--prereq-output", help="Prerequisites check output")
    parser.add_argument("--clone-output", help="Git clone output")
    parser.add_argument("--install-output", help="npm install output")
    parser.add_argument("--build-output", help="npm build output")
    parser.add_argument("--version-output", help="xcsh version output")
    parser.add_argument("--output", default="docs/install/source.md", help="Output file path")
    parser.add_argument("--templates", default="scripts/templates", help="Templates directory")

    args = parser.parse_args()

    env = Environment(
        loader=FileSystemLoader(args.templates),
        trim_blocks=True,
        lstrip_blocks=True,
    )

    # Register naming filters for consistency across all generators
    env.filters["to_human_readable"] = to_human_readable
    env.filters["normalize_acronyms"] = normalize_acronyms
    env.filters["title_case"] = to_title_case

    template = env.get_template("source.md.j2")

    content = template.render(
        node_version=args.node_version,
        prereq_output=clean_output(args.prereq_output),
        clone_output=clean_output(args.clone_output),
        install_output=clean_output(args.install_output),
        build_output=clean_output(args.build_output),
        version_output=clean_output(args.version_output),
        generation_date=datetime.now(UTC).strftime("%Y-%m-%d"),
    )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content)

    print(f"Generated: {output_path}")
    if args.node_version:
        print(f"  Node.js version: {args.node_version}")
    if args.version_output:
        print(f"  Version output: {len(clean_output(args.version_output))} chars")


if __name__ == "__main__":
    main()
