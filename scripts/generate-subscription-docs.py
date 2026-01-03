#!/usr/bin/env python3
"""
Subscription Documentation Generator

Generates comprehensive documentation for the xcsh subscription command group
by parsing xcsh --spec JSON output and rendering Jinja2 templates.

Usage:
    python scripts/generate-subscription-docs.py [--xcsh PATH] [--output DIR] [--clean]
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
        # Use --spec without --output-format json for compatibility with older binaries
        # (--spec already outputs JSON by default)
        result = subprocess.run(
            [cli_binary_path, "--spec"],
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


def load_subscription_spec(cli_binary_path: str) -> dict:
    """Run xcsh subscription --spec for extended subscription-specific data."""
    try:
        # Use --spec without --output-format json for compatibility with older binaries
        result = subprocess.run(
            [cli_binary_path, "subscription", "--spec"],
            capture_output=True,
            text=True,
            check=True,
        )
        return json.loads(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Note: xcsh subscription --spec not available: {e.stderr}", file=sys.stderr)
        return {}
    except json.JSONDecodeError as e:
        print(f"Error parsing xcsh subscription --spec output: {e}", file=sys.stderr)
        return {}


def find_subscription_command(spec: dict) -> dict | None:
    """Extract the subscription command from the main spec."""
    commands = spec.get("commands", [])
    for cmd in commands:
        path = cmd.get("path", [])
        if path == ["subscription"]:
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
        title = "xcsh subscription"
        description = cmd.get("short", "Subscription and billing management")
    elif command_type == "group":
        title = f"xcsh subscription {name}"
        description = cmd.get("short", f"Manage {to_human_readable(name)}")
    else:
        title = f"xcsh subscription {' '.join(path[1:])}"
        description = cmd.get("short", "")

    keywords = [
        "xcsh",
        "F5",
        "F5 XC",
        "F5 Distributed Cloud",
        "subscription",
        "billing",
        "quota",
        name,
    ]

    aliases = cmd.get("aliases", [])
    keywords.extend(aliases)

    return {
        "title": title,
        "description": normalize_acronyms(description),
        "keywords": list(set(keywords)),
        "command": full_command,
        "command_group": "subscription",
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

    env.filters["to_human_readable"] = to_human_readable
    env.filters["normalize_acronyms"] = normalize_acronyms
    env.filters["to_title_case"] = to_title_case
    env.filters["underscore_to_space"] = lambda s: s.replace("_", " ") if s else ""

    return env


def generate_overview(
    env: Environment, subscription_cmd: dict, extended_spec: dict, output_dir: Path
) -> None:
    """Generate the main subscription/index.md overview page."""
    template = env.get_template("subscription.md.j2")

    front_matter = create_front_matter(subscription_cmd, "overview")
    subcommands = get_subcommands(subscription_cmd)

    leaf_commands = [cmd for cmd in subcommands if is_leaf_command(cmd)]
    group_commands = [cmd for cmd in subcommands if has_subcommands(cmd)]

    content = template.render(
        front_matter=front_matter,
        command=subscription_cmd,
        subcommands=subcommands,
        leaf_commands=leaf_commands,
        group_commands=group_commands,
        workflows=extended_spec.get("workflows", []),
        exit_codes=extended_spec.get("exit_codes", []),
        ai_hints=extended_spec.get("ai_hints", {}),
    )

    output_file = output_dir / "index.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")


def generate_subcommand_group(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate index.md for a subcommand group."""
    template = env.get_template("subscription_subcommand.md.j2")

    name = get_command_name(cmd)
    front_matter = create_front_matter(cmd, "group")
    subcommands = get_subcommands(cmd)

    content = template.render(
        front_matter=front_matter,
        command=cmd,
        name=name,
        subcommands=subcommands,
    )

    subdir = output_dir / name
    subdir.mkdir(parents=True, exist_ok=True)

    output_file = subdir / "index.md"
    output_file.write_text(content)
    print(f"  Generated: {output_file}")

    for subcmd in subcommands:
        if has_subcommands(subcmd):
            generate_subcommand_group(env, subcmd, subdir)
        else:
            generate_leaf_command(env, subcmd, subdir)


def generate_leaf_command(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate documentation for a leaf command."""
    template = env.get_template("subscription_command.md.j2")

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


def generate_standalone_command(env: Environment, cmd: dict, output_dir: Path) -> None:
    """Generate documentation for a standalone leaf command."""
    template = env.get_template("subscription_command.md.j2")

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


def generate_nav_structure(subscription_cmd: dict) -> list:
    """Generate the navigation structure for mkdocs.yml."""
    nav = []

    nav.append({"Subscription Overview": "commands/subscription/index.md"})

    subcommands = get_subcommands(subscription_cmd)

    leaf_commands = sorted(
        [cmd for cmd in subcommands if is_leaf_command(cmd)], key=lambda c: get_command_name(c)
    )
    group_commands = sorted(
        [cmd for cmd in subcommands if has_subcommands(cmd)], key=lambda c: get_command_name(c)
    )

    for cmd in leaf_commands:
        name = get_command_name(cmd)
        human_name = to_title_case(to_human_readable(name))
        nav.append({human_name: f"commands/subscription/{name}.md"})

    for cmd in group_commands:
        name = get_command_name(cmd)
        human_name = to_title_case(to_human_readable(name))
        group_nav = [{f"{human_name} Overview": f"commands/subscription/{name}/index.md"}]

        for subcmd in sorted(get_subcommands(cmd), key=lambda c: get_command_name(c)):
            sub_name = get_command_name(subcmd)
            sub_human_name = to_title_case(to_human_readable(sub_name))

            if has_subcommands(subcmd):
                nested_nav = [
                    {
                        f"{sub_human_name} Overview": f"commands/subscription/{name}/{sub_name}/index.md"
                    }
                ]
                for nested_cmd in sorted(
                    get_subcommands(subcmd), key=lambda c: get_command_name(c)
                ):
                    nested_name = get_command_name(nested_cmd)
                    nested_human_name = to_title_case(to_human_readable(nested_name))
                    nested_nav.append(
                        {
                            nested_human_name: f"commands/subscription/{name}/{sub_name}/{nested_name}.md"
                        }
                    )
                group_nav.append({sub_human_name: nested_nav})
            else:
                group_nav.append({sub_human_name: f"commands/subscription/{name}/{sub_name}.md"})

        nav.append({human_name: group_nav})

    return nav


def main():
    parser = argparse.ArgumentParser(
        description="Generate subscription documentation from CLI --spec"
    )
    parser.add_argument(
        "--cli-binary", default="./xcsh", help="Path to CLI binary (default: ./xcsh)"
    )
    parser.add_argument(
        "--output",
        default="docs/commands/subscription",
        help="Output directory (default: docs/commands/subscription)",
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

    cli_binary_path = Path(args.cli_binary).resolve()
    output_dir = Path(args.output)
    templates_dir = Path(args.templates)

    if not cli_binary_path.exists():
        print(f"Error: CLI binary not found at {cli_binary_path}", file=sys.stderr)
        sys.exit(1)

    if not templates_dir.exists():
        print(f"Error: Templates directory not found at {templates_dir}", file=sys.stderr)
        sys.exit(1)

    print(f"Loading spec from {cli_binary_path}...")
    spec = load_spec(str(cli_binary_path))
    extended_spec = load_subscription_spec(str(cli_binary_path))

    subscription_cmd = find_subscription_command(spec)
    if not subscription_cmd:
        print("Error: subscription command not found in spec", file=sys.stderr)
        sys.exit(1)

    if args.print_nav:
        nav = generate_nav_structure(subscription_cmd)
        import yaml

        print("\n# Navigation structure for mkdocs.yml:")
        print("- Subscription:")
        for item in nav:
            print(f"  {yaml.dump([item], default_flow_style=False).strip()}")
        return

    if args.clean and output_dir.exists():
        print(f"Cleaning {output_dir}...")
        shutil.rmtree(output_dir)

    output_dir.mkdir(parents=True, exist_ok=True)

    env = setup_jinja_env(templates_dir)

    print(f"Generating subscription documentation to {output_dir}...")

    generate_overview(env, subscription_cmd, extended_spec, output_dir)

    subcommands = get_subcommands(subscription_cmd)

    for cmd in subcommands:
        if has_subcommands(cmd):
            generate_subcommand_group(env, cmd, output_dir)
        else:
            generate_standalone_command(env, cmd, output_dir)

    print("\nGenerated documentation for subscription command group")
    print("\nTo add to mkdocs.yml, run with --print-nav flag")


if __name__ == "__main__":
    main()
