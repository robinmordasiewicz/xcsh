#!/usr/bin/env python3
"""
Generate install script documentation with real installation output.

Usage:
    python scripts/generate-script-docs.py \
        --version 4.6.0 \
        --install-output "$(cat install-output.txt)" \
        --output docs/install/script.md
"""

import argparse
from datetime import datetime, timezone
from pathlib import Path

from jinja2 import Environment, FileSystemLoader

from naming import to_human_readable, normalize_acronyms, to_title_case


def main():
    parser = argparse.ArgumentParser(
        description="Generate install script docs with real output"
    )
    parser.add_argument("--version", help="xcsh version for examples")
    parser.add_argument("--install-output", help="Captured installation output")
    parser.add_argument(
        "--output",
        default="docs/install/script.md",
        help="Output file path"
    )
    parser.add_argument(
        "--templates",
        default="scripts/templates",
        help="Templates directory"
    )

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

    template = env.get_template("script.md.j2")

    # Clean up install output (remove ANSI codes if present)
    install_output = args.install_output
    if install_output:
        # Remove common ANSI escape codes
        import re
        install_output = re.sub(r'\x1b\[[0-9;]*m', '', install_output)
        # Trim excessive whitespace
        install_output = install_output.strip()

    content = template.render(
        version=args.version,
        install_output=install_output,
        generation_date=datetime.now(timezone.utc).strftime("%Y-%m-%d"),
    )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content)

    print(f"Generated: {output_path}")
    if args.version:
        print(f"  Version: {args.version}")
    if install_output:
        print(f"  Install output: {len(install_output)} characters")


if __name__ == "__main__":
    main()
