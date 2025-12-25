#!/usr/bin/env python3
"""
CloudStatus Documentation Generator

Generates comprehensive documentation for the xcsh cloudstatus command group
by parsing xcsh --spec JSON output and rendering Jinja2 templates.

Usage:
    python scripts/generate-cloudstatus-docs.py [--xcsh PATH] [--output DIR] [--clean]
"""

import argparse
import json
import shutil
import subprocess
import sys
from pathlib import Path

from jinja2 import Environment, FileSystemLoader, select_autoescape

from naming import normalize_acronyms, to_human_readable, to_title_case


def load_spec(cli_binary_path: str) -> dict:
    """Run xcsh --spec and return the full CLI spec."""
    try:
        result = subprocess.run(
            [cli_binary_path, "--spec", "--output-format", "json"],
            capture_output=True,
            text=True,
            check=True,
        )
        return json.loads(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error running xcsh --spec: {e.stderr}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"Error parsing xcsh --spec output: {e}", file=sys.stderr)
        sys.exit(1)


def load_cloudstatus_spec(cli_binary_path: str) -> dict:
    """Run xcsh cloudstatus --spec for extended cloudstatus-specific data."""
    try:
        result = subprocess.run(
            [cli_binary_path, "cloudstatus", "--spec", "--output-format", "json"],
            capture_output=True,
            text=True,
            check=True,
        )
        return json.loads(result.stdout)
    except subprocess.CalledProcessError as e:
        # cloudstatus --spec might not be available, return empty dict
        print(f"Note: xcsh cloudstatus --spec not available: {e.stderr}", file=sys.stderr)
        return {}
    except json.JSONDecodeError as e:
        print(f"Error parsing xcsh cloudstatus --spec output: {e}", file=sys.stderr)
        return {}


def find_cloudstatus_command(spec: dict) -> dict | None:
    """Extract the cloudstatus command from the main spec."""
    commands = spec.get("commands", [])
    for cmd in commands:
        path = cmd.get("path", [])
        if path == ["cloudstatus"]:
            return cmd
    return None


def get_command_name(cmd: dict) -> str:
    """Get the last element of the command path as the name."""
    path = cmd.get("path", [])
    return path[-1] if path else ""


def get_subcommands(cmd: dict) -> list:
    """Get direct subcommands of a command."""
    return cmd.get("subcommands", [])


def has_subcommands(cmd: dict) -> bool:
    """Check if command has subcommands."""
    return len(get_subcommands(cmd)) > 0


def is_leaf_command(cmd: dict) -> bool:
    """Check if command is a leaf (has no subcommands)."""
    return not has_subcommands(cmd)


def create_front_matter(cmd: dict, command_type: str = "command") -> dict:
    """Create YAML front matter for a documentation page."""
    path = cmd.get("path", [])
    name = get_command_name(cmd)
    full_command = " ".join(["xcsh"] + path)

    if command_type == "overview":
        title = f"Cloud Status - xcsh {name}"
        description = cmd.get("short", f"Manage {to_human_readable(name)} resources")
    elif command_type == "group":
        title = f"{to_human_readable(name)} - xcsh cloudstatus"
        description = cmd.get("short", f"Manage {to_human_readable(name)}")
    else:
        title = f"xcsh cloudstatus {' '.join(path[1:])}"
        description = cmd.get("short", "")

    keywords = [
        "xcsh",
        "F5",
        "F5 XC",
        "F5 Distributed Cloud",
        "cloud status",
        name,
    ]

    # Add aliases as keywords
    aliases = cmd.get("aliases", [])
    keywords.extend(aliases)

    return {
        "title": title,
        "description": normalize_acronyms(description),
        "keywords": list(set(keywords)),  # Remove duplicates
        "command": full_command,
        "command_group": "cloudstatus",
        "aliases": aliases,
    }


def setup_jinja_env(templates_dir: Path) -> Environment:
    """Set up Jinja2 environment with custom filters."""
    env = Environment(
        loader=FileSystemLoader(templates_dir),
        autoescape=select_autoescape(),
        trim_blocks=True,
        lstrip_blocks=True,
    )

    # Add custom filters
    env.filters["to_human_readable"] = to_human_readable
    env.filters["normalize_acronyms"] = normalize_acronyms
    env.filters["to_title_case"] = to_title_case
    env.filters["underscore_to_space"] = lambda s: s.replace("_", " ") if s else ""

    return env


def generate_overview(
    env: Environment, cloudstatus_cmd: dict, extended_spec: dict, output_dir: Path
) -> None:
    """Generate the main cloudstatus/index.md overview page."""
    template = env.get_template("cloudstatus.md.j2")

    front_matter = create_front_matter(cloudstatus_cmd, "overview")
    subcommands = get_subcommands(cloudstatus_cmd)

    # Separate leaf commands from group commands
    leaf_commands = [cmd for cmd in subcommands if is_leaf_command(cmd)]
    group_commands = [cmd for cmd in subcommands if has_subcommands(cmd)]

    content = template.render(
        front_matter=front_matter,
        command=cloudstatus_cmd,
        subcommands=subcommands,
        leaf_commands=leaf_commands,
        group_commands=group_commands,
        workflows=extended_spec.get("workflows", []),
        exit_codes=extended_spec.get("exit_codes", []),
        status_indicators=extended_spec.get("status_indicators", {}),
        ai_hints=extended_spec.get("ai_hints", {}),
    )

    output_file = output_dir / "index.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")


def generate_subcommand_group(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate index.md for a subcommand group (e.g., components, incidents)."""
    template = env.get_template("cloudstatus_subcommand.md.j2")

    name = get_command_name(cmd)
    front_matter = create_front_matter(cmd, "group")
    subcommands = get_subcommands(cmd)

    content = template.render(
        front_matter=front_matter,
        command=cmd,
        name=name,
        subcommands=subcommands,
    )

    # Create subdirectory
    subdir = output_dir / name
    subdir.mkdir(parents=True, exist_ok=True)

    output_file = subdir / "index.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")

    # Generate leaf command pages
    for subcmd in subcommands:
        generate_leaf_command(env, subcmd, subdir)


def generate_leaf_command(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate documentation for a leaf command."""
    template = env.get_template("cloudstatus_command.md.j2")

    path = cmd.get("path", [])
    name = get_command_name(cmd)
    front_matter = create_front_matter(cmd, "command")

    # Build the relative path from cloudstatus
    relative_path = path[1:] if len(path) > 1 else [name]

    content = template.render(
        front_matter=front_matter,
        command=cmd,
        name=name,
        path=relative_path,
        flags=cmd.get("flags", []),
    )

    output_file = output_dir / f"{name}.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")


def generate_standalone_command(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate documentation for a standalone leaf command at cloudstatus level."""
    template = env.get_template("cloudstatus_command.md.j2")

    path = cmd.get("path", [])
    name = get_command_name(cmd)
    front_matter = create_front_matter(cmd, "command")

    relative_path = path[1:] if len(path) > 1 else [name]

    content = template.render(
        front_matter=front_matter,
        command=cmd,
        name=name,
        path=relative_path,
        flags=cmd.get("flags", []),
    )

    output_file = output_dir / f"{name}.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")


def generate_nav_structure(cloudstatus_cmd: dict) -> list:
    """Generate the navigation structure for mkdocs.yml."""
    nav = []

    # Overview
    nav.append({"Overview": "commands/cloudstatus/index.md"})

    subcommands = get_subcommands(cloudstatus_cmd)

    # Sort subcommands: leaf commands first (alphabetically), then groups (alphabetically)
    leaf_commands = sorted(
        [cmd for cmd in subcommands if is_leaf_command(cmd)], key=lambda c: get_command_name(c)
    )
    group_commands = sorted(
        [cmd for cmd in subcommands if has_subcommands(cmd)], key=lambda c: get_command_name(c)
    )

    # Add leaf commands
    for cmd in leaf_commands:
        name = get_command_name(cmd)
        human_name = to_human_readable(name)
        nav.append({human_name: f"commands/cloudstatus/{name}.md"})

    # Add group commands with their children
    for cmd in group_commands:
        name = get_command_name(cmd)
        human_name = to_human_readable(name)
        group_nav = [{"Overview": f"commands/cloudstatus/{name}/index.md"}]

        for subcmd in sorted(get_subcommands(cmd), key=lambda c: get_command_name(c)):
            sub_name = get_command_name(subcmd)
            sub_human_name = to_human_readable(sub_name)
            group_nav.append({sub_human_name: f"commands/cloudstatus/{name}/{sub_name}.md"})

        nav.append({human_name: group_nav})

    return nav


def main():
    parser = argparse.ArgumentParser(
        description="Generate cloudstatus documentation from CLI --spec"
    )
    parser.add_argument(
        "--cli-binary", default="./xcsh", help="Path to CLI binary (default: ./xcsh)"
    )
    parser.add_argument(
        "--output",
        default="docs/commands/cloudstatus",
        help="Output directory (default: docs/commands/cloudstatus)",
    )
    parser.add_argument(
        "--templates",
        default="scripts/templates",
        help="Templates directory (default: scripts/templates)",
    )
    parser.add_argument(
        "--clean", action="store_true", help="Clean output directory before generating"
    )
    parser.add_argument(
        "--print-nav", action="store_true", help="Print navigation structure for mkdocs.yml"
    )

    args = parser.parse_args()

    # Resolve paths
    cli_binary_path = Path(args.cli_binary).resolve()
    output_dir = Path(args.output)
    templates_dir = Path(args.templates)

    # Verify CLI binary exists
    if not cli_binary_path.exists():
        print(f"Error: CLI binary not found at {cli_binary_path}", file=sys.stderr)
        sys.exit(1)

    # Verify templates exist
    if not templates_dir.exists():
        print(f"Error: Templates directory not found at {templates_dir}", file=sys.stderr)
        sys.exit(1)

    print(f"Loading spec from {cli_binary_path}...")
    spec = load_spec(str(cli_binary_path))
    extended_spec = load_cloudstatus_spec(str(cli_binary_path))

    # Find cloudstatus command
    cloudstatus_cmd = find_cloudstatus_command(spec)
    if not cloudstatus_cmd:
        print("Error: cloudstatus command not found in spec", file=sys.stderr)
        sys.exit(1)

    # Print nav structure if requested
    if args.print_nav:
        nav = generate_nav_structure(cloudstatus_cmd)
        import yaml

        print("\n# Navigation structure for mkdocs.yml:")
        print("    - Cloud Status:")
        for item in nav:
            print(f"      {yaml.dump([item], default_flow_style=False).strip()}")
        return

    # Clean output directory if requested
    if args.clean and output_dir.exists():
        print(f"Cleaning {output_dir}...")
        shutil.rmtree(output_dir)

    # Create output directory
    output_dir.mkdir(parents=True, exist_ok=True)

    # Set up Jinja2 environment
    env = setup_jinja_env(templates_dir)

    print(f"Generating cloudstatus documentation to {output_dir}...")

    # Generate overview page
    generate_overview(env, cloudstatus_cmd, extended_spec, output_dir)

    # Generate documentation for each subcommand
    subcommands = get_subcommands(cloudstatus_cmd)

    for cmd in subcommands:
        if has_subcommands(cmd):
            # Generate group index and children
            generate_subcommand_group(env, cmd, output_dir)
        else:
            # Generate standalone leaf command
            generate_standalone_command(env, cmd, output_dir)

    print("\nGenerated documentation for cloudstatus command group")

    # Print nav structure hint
    print("\nTo add to mkdocs.yml, run with --print-nav flag")


if __name__ == "__main__":
    main()
