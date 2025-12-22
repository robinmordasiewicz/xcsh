#!/usr/bin/env python3
"""
Generate homebrew installation documentation with real version info.

Usage:
    python scripts/generate-homebrew-docs.py \
        --version 1.2.0 \
        --commit abc1234 \
        --built 2024-12-09T10:00:00Z \
        --go-version go1.23.4 \
        --platform darwin/arm64 \
        --output docs/install/homebrew.md
"""

import argparse
from datetime import datetime, timezone
from pathlib import Path

from jinja2 import Environment, FileSystemLoader

from naming import to_human_readable, normalize_acronyms, to_title_case


def main():
    parser = argparse.ArgumentParser(
        description="Generate homebrew docs with real version info"
    )
    parser.add_argument("--version", help="xcsh version")
    parser.add_argument("--commit", help="Git commit hash")
    parser.add_argument("--built", help="Build timestamp")
    parser.add_argument("--go-version", help="Go version used for build")
    parser.add_argument("--platform", help="Target platform (os/arch)")
    parser.add_argument("--install-output", help="Captured homebrew installation output")
    parser.add_argument(
        "--output",
        default="docs/install/homebrew.md",
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

    template = env.get_template("homebrew.md.j2")

    # Clean up install output (remove ANSI codes if present)
    install_output = args.install_output
    if install_output:
        import re
        install_output = re.sub(r'\x1b\[[0-9;]*m', '', install_output)
        install_output = install_output.strip()

    content = template.render(
        version=args.version,
        commit=args.commit,
        built=args.built,
        go_version=args.go_version,
        platform=args.platform,
        install_output=install_output,
        generation_date=datetime.now(timezone.utc).strftime("%Y-%m-%d"),
    )

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(content)

    print(f"Generated: {output_path}")
    if args.version:
        print(f"  Version:  {args.version}")
        print(f"  Commit:   {args.commit}")
        print(f"  Built:    {args.built}")
        print(f"  Go:       {args.go_version}")
        print(f"  Platform: {args.platform}")
    else:
        print("  (No version info provided, using template defaults)")
    if install_output:
        print(f"  Install output: {len(install_output)} characters")


if __name__ == "__main__":
    main()
