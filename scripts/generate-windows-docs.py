#!/usr/bin/env python3
"""
Generate Windows installation documentation with real version info.

Usage:
    python scripts/generate-windows-docs.py \
        --version 6.0.0 \
        --node-version 22 \
        --output docs/install/windows.md
"""

import argparse
import re
from datetime import UTC, datetime
from pathlib import Path

from jinja2 import Environment, FileSystemLoader

from naming import normalize_acronyms, to_human_readable, to_title_case


def main():
    parser = argparse.ArgumentParser(description="Generate Windows docs with real version info")
    parser.add_argument("--version", help="xcsh version")
    parser.add_argument("--node-version", help="Node.js version used for build")
    parser.add_argument("--install-output", help="Captured Windows installation output")
    parser.add_argument("--output", default="docs/install/windows.md", help="Output file path")
    parser.add_argument("--templates", default="scripts/templates", help="Templates directory")

    args = parser.parse_args()

    # Setup Jinja2
    env = Environment(
        loader=FileSystemLoader(args.templates),
        trim_blocks=True,
        lstrip_blocks=True,
    )

    # Register naming filters for consistency across all generators
    env.filters["to_human_readable"] = to_human_readable
    env.filters["normalize_acronyms"] = normalize_acronyms
    env.filters["title_case"] = to_title_case

    template = env.get_template("windows.md.j2")

    # Clean up install output (remove ANSI codes if present)
    install_output = args.install_output
    if install_output:
        # Remove ANSI escape codes
        install_output = re.sub(r"\x1b\[[0-9;]*m", "", install_output)
        # Remove Windows-specific escape sequences
        install_output = re.sub(r"\x1b\[[0-9;]*[A-Za-z]", "", install_output)
        install_output = install_output.strip()

    content = template.render(
        version=args.version,
        node_version=args.node_version,
        install_output=install_output,
        generation_date=datetime.now(UTC).strftime("%Y-%m-%d"),
    )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content)

    print(f"Generated: {output_path}")
    if args.version:
        print(f"  Version:  {args.version}")
        print(f"  Node.js:  {args.node_version}")
    else:
        print("  (No version info provided, using template defaults)")
    if install_output:
        print(f"  Install output: {len(install_output)} characters")


if __name__ == "__main__":
    main()
