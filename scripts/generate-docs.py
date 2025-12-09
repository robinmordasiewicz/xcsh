#!/usr/bin/env python3
"""
vesctl Documentation Generator

Generates comprehensive, AI-friendly documentation for all vesctl CLI commands
by parsing the vesctl --spec JSON output and rendering Jinja2 templates.

Usage:
    python scripts/generate-docs.py [--vesctl PATH] [--output DIR] [--clean]
"""

import argparse
import json
import os
import shutil
import subprocess
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Optional

import yaml
from jinja2 import Environment, FileSystemLoader, select_autoescape


# Canonical action order for consistent display
ACTION_ORDER = [
    'list', 'get', 'create', 'delete', 'replace',
    'apply', 'patch', 'status', 'add-labels', 'remove-labels'
]


def sort_actions(actions: list) -> list:
    """Sort actions by canonical order."""
    action_priority = {action: i for i, action in enumerate(ACTION_ORDER)}
    return sorted(
        actions,
        key=lambda c: action_priority.get(c.path[1] if len(c.path) > 1 else '', 999)
    )


@dataclass
class Flag:
    """Represents a CLI flag."""
    name: str
    type: str
    description: str
    shorthand: str = ""
    default: str = ""

    @classmethod
    def from_dict(cls, d: dict) -> "Flag":
        return cls(
            name=d.get("name", ""),
            type=d.get("type", ""),
            description=d.get("description", ""),
            shorthand=d.get("shorthand", ""),
            default=d.get("default", ""),
        )


@dataclass
class Command:
    """Represents a CLI command."""
    path: list[str]
    use: str
    short: str
    long: str = ""
    example: str = ""
    aliases: list[str] = field(default_factory=list)
    flags: list[Flag] = field(default_factory=list)
    subcommands: list["Command"] = field(default_factory=list)

    @classmethod
    def from_dict(cls, d: dict) -> "Command":
        return cls(
            path=d.get("path", []),
            use=d.get("use", ""),
            short=d.get("short", ""),
            long=d.get("long", ""),
            example=d.get("example", ""),
            aliases=d.get("aliases", []),
            flags=[Flag.from_dict(f) for f in d.get("flags", [])],
            subcommands=[cls.from_dict(s) for s in d.get("subcommands", [])],
        )

    @property
    def name(self) -> str:
        """Get command name from path."""
        return self.path[-1] if self.path else ""

    @property
    def full_command(self) -> str:
        """Get full command string."""
        return "vesctl " + " ".join(self.path)

    @property
    def depth(self) -> int:
        """Get command depth in hierarchy."""
        return len(self.path)


@dataclass
class CommandTree:
    """Hierarchical tree of commands."""
    name: str
    command: Optional[Command] = None
    children: dict[str, "CommandTree"] = field(default_factory=dict)

    def add_command(self, cmd: Command) -> None:
        """Add a command to the tree."""
        node = self
        for part in cmd.path:
            if part not in node.children:
                node.children[part] = CommandTree(name=part)
            node = node.children[part]
        node.command = cmd

        # Recursively add subcommands
        for subcmd in cmd.subcommands:
            self.add_command(subcmd)


class VesctlDocsGenerator:
    """Main documentation generator class."""

    def __init__(
        self,
        vesctl_path: str = "./vesctl",
        output_dir: str = "docs/commands",
        template_dir: str = "scripts/templates",
    ):
        # Resolve to absolute path to avoid PATH lookup issues
        self.vesctl_path = Path(vesctl_path).resolve()
        self.output_dir = Path(output_dir)
        self.template_dir = Path(template_dir)
        self.spec: dict = {}
        self.global_flags: list[Flag] = []
        self.tree = CommandTree(name="vesctl")
        self.generated_files: list[Path] = []

        # Setup Jinja2 environment
        self.env = Environment(
            loader=FileSystemLoader(self.template_dir),
            autoescape=select_autoescape(["html", "xml"]),
            trim_blocks=True,
            lstrip_blocks=True,
        )

        # Add custom filters
        self.env.filters["underscore_to_space"] = lambda s: s.replace("_", " ")
        self.env.filters["title_case"] = lambda s: s.replace("_", " ").title()

        # API specs mapping
        self.api_specs_dir = Path("docs/specifications/api")
        self.resource_api_map: dict[str, dict] = {}

    def load_api_specs(self) -> None:
        """Load and index OpenAPI spec files for API documentation links."""
        if not self.api_specs_dir.exists():
            print(f"Warning: API specs directory not found: {self.api_specs_dir}")
            return

        print(f"Loading API specs from {self.api_specs_dir}...")
        spec_count = 0

        for spec_file in self.api_specs_dir.glob("*.json"):
            # Extract resource name from filename
            # Pattern: docs-cloud-f5-com.XXXX.public.ves.io.schema.[path].ves-swagger.json
            parts = spec_file.stem.split(".")
            if "schema" in parts:
                schema_idx = parts.index("schema")
                # Get the last part before "ves-swagger" as resource name
                resource_name = parts[-2] if len(parts) >= 2 else ""

                if resource_name and resource_name != "ves-swagger":
                    try:
                        with open(spec_file) as f:
                            spec_data = json.load(f)

                        # Store spec with resource name as key
                        # If resource already exists, keep the first one (they should be the same)
                        if resource_name not in self.resource_api_map:
                            self.resource_api_map[resource_name] = {
                                "spec": spec_data,
                                "file": spec_file,
                            }
                            spec_count += 1
                    except (json.JSONDecodeError, IOError) as e:
                        print(f"  Warning: Failed to load {spec_file}: {e}")

        print(f"  Loaded {spec_count} API specs, {len(self.resource_api_map)} unique resources")

    def get_api_docs_url(self, resource: str, action: str) -> Optional[str]:
        """Get API documentation URL for a resource+action combination."""
        if resource not in self.resource_api_map:
            return None

        spec = self.resource_api_map[resource]["spec"]

        # Map vesctl action to API operation name
        action_to_op = {
            "create": "Create",
            "list": "List",
            "get": "Get",
            "delete": "Delete",
            "replace": "Replace",
            "apply": "Replace",
            "status": "Get",
            "patch": "Replace",
            "add-labels": "Create",
            "remove-labels": "Delete",
        }

        op_name = action_to_op.get(action)
        if not op_name:
            return None

        # Search paths for matching operation with .API. service type
        for path, methods in spec.get("paths", {}).items():
            for method, details in methods.items():
                if isinstance(details, dict):
                    proto_rpc = details.get("x-ves-proto-rpc", "")
                    # Match operations that end with .API.<OpName>
                    if proto_rpc.endswith(f".API.{op_name}"):
                        external_docs = details.get("externalDocs", {})
                        url = external_docs.get("url")
                        if url:
                            return url

        return None

    def load_spec(self) -> None:
        """Load CLI specification from vesctl --spec."""
        print(f"Loading spec from {self.vesctl_path}...")

        try:
            result = subprocess.run(
                [str(self.vesctl_path), "--spec"],
                capture_output=True,
                text=True,
                check=True,
            )
            self.spec = json.loads(result.stdout)
        except subprocess.CalledProcessError as e:
            print(f"Error running vesctl --spec: {e.stderr}")
            sys.exit(1)
        except json.JSONDecodeError as e:
            print(f"Error parsing spec JSON: {e}")
            sys.exit(1)

        # Parse global flags
        self.global_flags = [
            Flag.from_dict(f) for f in self.spec.get("global_flags", [])
        ]

        # Build command tree
        for cmd_dict in self.spec.get("commands", []):
            cmd = Command.from_dict(cmd_dict)
            self.tree.add_command(cmd)

        print(f"Loaded {len(self.spec.get('commands', []))} top-level commands")

    def count_commands(self, node: CommandTree = None) -> int:
        """Count total commands in tree."""
        if node is None:
            node = self.tree
        count = 1 if node.command else 0
        for child in node.children.values():
            count += self.count_commands(child)
        return count

    def ensure_dir(self, path: Path) -> None:
        """Ensure directory exists."""
        path.mkdir(parents=True, exist_ok=True)

    def write_file(self, path: Path, content: str) -> None:
        """Write content to file and track it."""
        self.ensure_dir(path.parent)
        path.write_text(content)
        self.generated_files.append(path)
        print(f"  Generated: {path}")

    def generate_front_matter(
        self,
        title: str,
        description: str,
        command: Optional[Command] = None,
        **extra,
    ) -> dict:
        """Generate front matter for a documentation page."""
        fm = {
            "title": title,
            "description": description,
        }

        if command:
            # Build keywords from command path
            keywords = ["vesctl", "F5 XC", "F5 Distributed Cloud"]
            keywords.extend(command.path)
            keywords.extend([p.replace("_", " ") for p in command.path])
            fm["keywords"] = list(set(keywords))

            fm["command"] = command.full_command
            if len(command.path) >= 1:
                fm["command_group"] = command.path[0]
            if len(command.path) >= 2:
                fm["action"] = command.path[1]
            if len(command.path) >= 3:
                fm["resource_type"] = command.path[2]
            if command.aliases:
                fm["aliases"] = command.aliases

        fm.update(extra)
        return fm

    def generate_commands_index(self) -> None:
        """Generate the main commands index page."""
        template = self.env.get_template("commands_index.md.j2")

        # Get top-level commands
        top_level = []
        for name, child in sorted(self.tree.children.items()):
            if child.command:
                top_level.append(child.command)

        content = template.render(
            title="Command Reference",
            description="Complete reference for all vesctl CLI commands",
            commands=top_level,
            global_flags=self.global_flags,
            version=self.spec.get("version", "dev"),
        )

        self.write_file(self.output_dir / "index.md", content)

    def generate_command_group(self, name: str, node: CommandTree) -> None:
        """Generate documentation for a command group."""
        if not node.command:
            return

        cmd = node.command
        group_dir = self.output_dir / name

        # Generate group index
        template = self.env.get_template("command_group.md.j2")

        # Get subcommands - for resource-first groups, list resources instead of actions
        if name == "configuration":
            # For configuration, list unique resources as subcommands
            resources = self.collect_resources_across_actions(node)
            # Create pseudo-commands for the index page display
            subcommands = []
            for resource_name in sorted(resources.keys()):
                actions = resources[resource_name]
                if actions:
                    # Use first action's command as template but adjust for display
                    first_action = actions[0]
                    subcommands.append(Command(
                        path=[name, resource_name],
                        use=resource_name,
                        short=f"Manage {resource_name.replace('_', ' ')} resources",
                    ))
        else:
            # Standard action-first listing
            subcommands = []
            for child_name, child in sorted(node.children.items()):
                if child.command:
                    subcommands.append(child.command)

        fm = self.generate_front_matter(
            title=f"vesctl {name}",
            description=cmd.short,
            command=cmd,
        )

        # Use flat links for resource-first groups (configuration, api-endpoint)
        use_flat_links = name == "configuration"

        content = template.render(
            front_matter=fm,
            command=cmd,
            subcommands=subcommands,
            global_flags=self.global_flags,
            flat_links=use_flat_links,
        )

        self.write_file(group_dir / "index.md", content)

        # Generate .meta.yml
        self.generate_meta_yml(
            group_dir,
            description=cmd.short,
            tags=["vesctl", name],
        )

        # Generate subcommand documentation based on group type
        if name == "configuration":
            # Use resource-first layout for configuration
            self.generate_configuration_resource_first(name, node)
        else:
            # Use action-first layout for other groups
            for child_name, child in node.children.items():
                if child.command:
                    self.generate_action(name, child_name, child)

    def generate_action(
        self, group: str, action: str, node: CommandTree
    ) -> None:
        """Generate documentation for an action."""
        if not node.command:
            return

        cmd = node.command
        action_dir = self.output_dir / group / action

        # Special handling for RPC: use service-level grouping
        if group == "request" and action == "rpc":
            # Generate RPC index with service count instead of flat list
            template = self.env.get_template("action.md.j2")

            # Get services for display
            services = self.collect_rpc_services(node)

            # Create pseudo-commands for the index page display
            resources = []
            for service_name in sorted(services.keys()):
                procedures = services[service_name]
                if procedures:
                    resources.append(Command(
                        path=[group, action, service_name],
                        use=service_name,
                        short=f"{service_name.replace('_', ' ').title()} service ({len(procedures)} procedures)",
                    ))

            fm = self.generate_front_matter(
                title=f"vesctl {group} {action}",
                description=cmd.short,
                command=cmd,
            )

            content = template.render(
                front_matter=fm,
                command=cmd,
                resources=resources,
                global_flags=self.global_flags,
                group=group,
                action=action,
            )

            self.write_file(action_dir / "index.md", content)

            # Generate .meta.yml
            self.generate_meta_yml(
                action_dir,
                description=f"RPC commands for {group}",
                tags=["vesctl", group, action],
            )

            # Generate service-grouped RPC docs
            self.generate_rpc_service_grouped(group, action, node)
            return

        # Standard action processing
        # Generate action index
        template = self.env.get_template("action.md.j2")

        # Get resource types (subcommands)
        resources = []
        for child_name, child in sorted(node.children.items()):
            if child.command:
                resources.append(child.command)

        fm = self.generate_front_matter(
            title=f"vesctl {group} {action}",
            description=cmd.short,
            command=cmd,
        )

        content = template.render(
            front_matter=fm,
            command=cmd,
            resources=resources,
            global_flags=self.global_flags,
            group=group,
            action=action,
        )

        self.write_file(action_dir / "index.md", content)

        # Generate .meta.yml
        self.generate_meta_yml(
            action_dir,
            description=f"{action.replace('_', ' ').title()} commands for {group}",
            tags=["vesctl", group, action],
        )

        # Generate resource type pages
        for child_name, child in node.children.items():
            if child.command:
                self.generate_resource_page(group, action, child_name, child)

    def generate_resource_page(
        self,
        group: str,
        action: str,
        resource: str,
        node: CommandTree,
    ) -> None:
        """Generate documentation for a resource type."""
        if not node.command:
            return

        cmd = node.command
        template = self.env.get_template("resource_type.md.j2")

        # Get API documentation URL for this resource+action
        api_docs_url = self.get_api_docs_url(resource, action)

        fm = self.generate_front_matter(
            title=f"vesctl {group} {action} {resource}",
            description=cmd.short,
            command=cmd,
        )

        content = template.render(
            front_matter=fm,
            command=cmd,
            global_flags=self.global_flags,
            group=group,
            action=action,
            resource=resource,
            api_docs_url=api_docs_url,
        )

        resource_path = self.output_dir / group / action / f"{resource}.md"
        self.write_file(resource_path, content)

    def find_related_commands(
        self, group: str, resource: str
    ) -> list[Command]:
        """Find commands for the same resource in different actions."""
        related = []
        group_node = self.tree.children.get(group)
        if not group_node:
            return related

        for action_name, action_node in group_node.children.items():
            resource_node = action_node.children.get(resource)
            if resource_node and resource_node.command:
                related.append(resource_node.command)

        return sorted(related, key=lambda c: c.path[1] if len(c.path) > 1 else "")

    def find_actions_for_resource(
        self, group: str, resource: str
    ) -> list[Command]:
        """Find all actions available for a specific resource type."""
        return self.find_related_commands(group, resource)

    def generate_action_examples(self, group: str, action: str, resource: str) -> str:
        """Generate example bash commands for an action."""
        resource_display = resource.replace('_', ' ')
        resource_kebab = resource.replace('_', '-')

        examples = {
            'list': f'''```bash
# List all {resource_display} resources
vesctl {group} {action} {resource}

# List in specific namespace
vesctl {group} {action} {resource} -n example-namespace

# List with JSON output
vesctl {group} {action} {resource} --output-format json
```''',
            'get': f'''```bash
# Get {resource_display} details
vesctl {group} {action} {resource} example-{resource_kebab}

# Get with YAML output
vesctl {group} {action} {resource} example-{resource_kebab} --output-format yaml
```''',
            'create': f'''```bash
# Create {resource_display} from file
vesctl {group} {action} {resource} -i {resource}.yaml
```''',
            'delete': f'''```bash
# Delete {resource_display}
vesctl {group} {action} {resource} example-{resource_kebab}

# Delete with confirmation bypass
vesctl {group} {action} {resource} example-{resource_kebab} --yes
```''',
            'replace': f'''```bash
# Replace {resource_display} from file
vesctl {group} {action} {resource} -i {resource}.yaml
```''',
            'apply': f'''```bash
# Apply {resource_display} from file
vesctl {group} {action} {resource} -i {resource}.yaml
```''',
            'patch': f'''```bash
# Patch {resource_display}
vesctl {group} {action} {resource} example-{resource_kebab} -i patch.yaml
```''',
            'status': f'''```bash
# Get {resource_display} status
vesctl {group} {action} {resource} example-{resource_kebab}
```''',
            'add-labels': f'''```bash
# Add labels to {resource_display}
vesctl {group} {action} {resource} example-{resource_kebab} --label-key app --label-value web
```''',
            'remove-labels': f'''```bash
# Remove labels from {resource_display}
vesctl {group} {action} {resource} example-{resource_kebab} --label-key app
```''',
        }
        return examples.get(action, f'''```bash
vesctl {group} {action} {resource}
```''')

    def collect_resources_across_actions(
        self, group_node: CommandTree
    ) -> dict[str, list[Command]]:
        """Collect all resources and their available actions for a group."""
        resources: dict[str, list[Command]] = {}

        for action_name, action_node in group_node.children.items():
            for resource_name, resource_node in action_node.children.items():
                if resource_node.command:
                    if resource_name not in resources:
                        resources[resource_name] = []
                    resources[resource_name].append(resource_node.command)

        # Sort actions for each resource
        for resource_name in resources:
            resources[resource_name] = sorted(
                resources[resource_name],
                key=lambda c: c.path[1] if len(c.path) > 1 else ""
            )

        return resources

    def generate_configuration_resource_first(
        self, group: str, node: CommandTree
    ) -> None:
        """Generate configuration docs with resource-first layout."""
        # Collect all resources across all actions
        resources = self.collect_resources_across_actions(node)

        print(f"  Found {len(resources)} resource types")

        # Generate documentation for each resource type
        for resource_name, actions in sorted(resources.items()):
            self.generate_resource_group(group, resource_name, actions)

    def generate_resource_group(
        self, group: str, resource: str, actions: list[Command]
    ) -> None:
        """Generate unified documentation for a resource type with all actions."""
        template = self.env.get_template("resource_unified.md.j2")

        # Sort actions by canonical order
        sorted_actions = sort_actions(actions)

        # Build action data with API URLs and generated examples
        action_data = []
        api_docs_urls = {}

        for action_cmd in sorted_actions:
            action_name = action_cmd.path[1] if len(action_cmd.path) > 1 else ""
            api_url = self.get_api_docs_url(resource, action_name)

            if api_url:
                api_docs_urls[action_name] = api_url

            action_data.append({
                'action_name': action_name,
                'command': action_cmd,
                'api_docs_url': api_url,
                'generated_examples': self.generate_action_examples(group, action_name, resource),
            })

        fm = self.generate_front_matter(
            title=f"vesctl {group} {resource}",
            description=f"Manage {resource.replace('_', ' ')} resources",
            command=sorted_actions[0] if sorted_actions else None,
            resource_type=resource,
        )

        content = template.render(
            front_matter=fm,
            group=group,
            resource=resource,
            actions=action_data,
            api_docs_urls=api_docs_urls,
            global_flags=self.global_flags,
        )

        # Write single file (not in a directory)
        resource_path = self.output_dir / group / f"{resource}.md"
        self.write_file(resource_path, content)

    # ===== RPC Service Grouping Methods =====

    def extract_rpc_service(self, procedure_name: str) -> str:
        """Extract service prefix from RPC procedure name.

        Example: 'alert.CustomAPI.Alerts' -> 'alert'
        """
        parts = procedure_name.split('.')
        return parts[0] if parts else procedure_name

    def extract_rpc_procedure_name(self, full_name: str) -> str:
        """Extract procedure name from full RPC procedure name.

        Example: 'alert.CustomAPI.Alerts' -> 'Alerts'
        """
        parts = full_name.split('.')
        return parts[-1] if parts else full_name

    def collect_rpc_services(
        self, rpc_node: CommandTree
    ) -> dict[str, list[dict]]:
        """Collect all RPC procedures grouped by service.

        Returns dict mapping service name to list of procedure info dicts.
        """
        services: dict[str, list[dict]] = {}

        for proc_name, proc_node in rpc_node.children.items():
            if proc_node.command:
                service = self.extract_rpc_service(proc_name)
                procedure_name = self.extract_rpc_procedure_name(proc_name)

                proc_info = {
                    'full_name': proc_name,
                    'procedure_name': procedure_name,
                    'service': service,
                    'command': proc_node.command,
                }

                if service not in services:
                    services[service] = []
                services[service].append(proc_info)

        # Sort procedures within each service
        for service in services:
            services[service] = sorted(
                services[service],
                key=lambda p: p['procedure_name']
            )

        return services

    def generate_rpc_service_grouped(
        self, group: str, action: str, node: CommandTree
    ) -> None:
        """Generate RPC docs with unified service pages."""
        services = self.collect_rpc_services(node)

        print(f"  Found {len(services)} RPC services")

        # Generate unified page for each service
        for service_name, procedures in sorted(services.items()):
            self.generate_rpc_service_unified(group, action, service_name, procedures)

    def generate_rpc_service_unified(
        self,
        group: str,
        action: str,
        service: str,
        procedures: list[dict],
    ) -> None:
        """Generate unified service page with all procedures."""
        template = self.env.get_template("rpc_service_unified.md.j2")

        fm = self.generate_front_matter(
            title=f"vesctl request rpc {service}",
            description=f"{service.replace('_', ' ').title()} service RPC procedures",
            rpc_service=service,
        )

        content = template.render(
            front_matter=fm,
            service=service,
            procedures=procedures,
            global_flags=self.global_flags,
        )

        # Write single file per service
        service_path = self.output_dir / group / action / f"{service}.md"
        self.write_file(service_path, content)

    def build_rpc_service_nav(
        self, group: str, action: str, node: CommandTree
    ) -> list[dict]:
        """Build flat service navigation for RPC commands."""
        nav_items = []

        services = self.collect_rpc_services(node)

        # Build flat navigation - one entry per service
        for service_name in sorted(services.keys()):
            service_display = service_name.replace("_", " ").replace("-", " ").title()
            nav_items.append({
                service_display: f"commands/{group}/{action}/{service_name}.md"
            })

        return nav_items

    def generate_meta_yml(
        self, directory: Path, description: str, tags: list[str]
    ) -> None:
        """Generate .meta.yml for a directory."""
        meta = {
            "description": description,
            "tags": tags,
        }
        meta_path = directory / ".meta.yml"
        self.write_file(meta_path, yaml.dump(meta, default_flow_style=False))

    def generate_navigation(self) -> dict:
        """Generate navigation structure for mkdocs.yml."""
        nav = []

        # Commands index
        nav.append({"Commands": "commands/index.md"})

        # Top-level command groups
        for group_name in sorted(self.tree.children.keys()):
            group_node = self.tree.children[group_name]
            group_nav = self.build_nav_tree(group_name, group_node)
            if group_nav:
                nav.append(group_nav)

        return nav

    def build_nav_tree(self, name: str, node: CommandTree) -> dict:
        """Build navigation tree for a command node."""
        display_name = name.replace("_", " ").replace("-", " ").title()

        # Command groups always have index.md in a directory
        # Even if they have no children (like 'completion')
        if not node.children:
            # No children - just the index page
            return {display_name: f"commands/{name}/index.md"}

        # Has children - build nested structure
        children = []

        # Add index first
        index_path = "/".join([name])
        children.append({f"{display_name} Overview": f"commands/{index_path}/index.md"})

        # Special handling for configuration: use resource-first navigation
        if name == "configuration":
            children.extend(self.build_resource_first_nav(name, node))
        else:
            # Standard action-first navigation
            for child_name in sorted(node.children.keys()):
                child_node = node.children[child_name]
                child_nav = self.build_child_nav(name, child_name, child_node)
                if child_nav:
                    children.append(child_nav)

        return {display_name: children}

    def build_resource_first_nav(
        self, group: str, node: CommandTree
    ) -> list[dict]:
        """Build flat resource navigation for configuration command."""
        nav_items = []

        # Collect all resources across all actions
        resources = self.collect_resources_across_actions(node)

        # Build flat navigation - one entry per resource
        for resource_name in sorted(resources.keys()):
            resource_display = resource_name.replace("_", " ").replace("-", " ").title()
            nav_items.append({
                resource_display: f"commands/{group}/{resource_name}.md"
            })

        return nav_items

    def build_child_nav(
        self, parent_path: str, name: str, node: CommandTree
    ) -> dict:
        """Build navigation for child nodes."""
        display_name = name.replace("_", " ").replace("-", " ").title()
        path = f"{parent_path}/{name}"

        if not node.children:
            # Check depth to determine if this is an action (directory) or resource (file)
            # Actions have depth 2 (e.g., ["api-endpoint", "control"])
            # Resources have depth 3+ (e.g., ["configuration", "list", "namespace"])
            if node.command and len(node.command.path) <= 2:
                # This is an action without resource types - has index.md in directory
                return {display_name: f"commands/{path}/index.md"}
            else:
                # This is a resource type - standalone .md file
                return {display_name: f"commands/{path}.md"}

        # Special handling for RPC: use service-grouped navigation
        if parent_path == "request" and name == "rpc":
            children = []
            children.append({f"{display_name} Overview": f"commands/{path}/index.md"})
            children.extend(self.build_rpc_service_nav(parent_path, name, node))
            return {display_name: children}

        # Has children - build nested structure
        children = []

        # Add index first
        children.append({f"{display_name} Overview": f"commands/{path}/index.md"})

        # Add children
        for child_name in sorted(node.children.keys()):
            child_node = node.children[child_name]
            child_path = f"{path}/{child_name}"

            if child_node.children:
                # Recurse for nested children
                child_nav = self.build_child_nav(path, child_name, child_node)
                if child_nav:
                    children.append(child_nav)
            else:
                # Leaf node - check if it's an action or resource
                child_display = child_name.replace("_", " ").replace("-", " ").title()
                if child_node.command and len(child_node.command.path) <= 2:
                    # Action without resources
                    children.append({child_display: f"commands/{child_path}/index.md"})
                else:
                    # Resource type
                    children.append({child_display: f"commands/{child_path}.md"})

        return {display_name: children}

    def save_navigation(self, nav: dict, output_path: Path = None) -> None:
        """Save navigation to a YAML file for manual integration."""
        if output_path is None:
            output_path = Path("docs/commands/_nav.yml")

        nav_content = {"nav": nav}
        self.write_file(output_path, yaml.dump(nav_content, default_flow_style=False, sort_keys=False))

    def update_mkdocs_yml(self, nav: list, mkdocs_path: Path = None) -> None:
        """Update mkdocs.yml with generated Commands navigation.

        Uses text-based replacement to preserve Python tags in mkdocs.yml.
        """
        import re

        if mkdocs_path is None:
            mkdocs_path = Path("mkdocs.yml")

        if not mkdocs_path.exists():
            print(f"Warning: {mkdocs_path} not found, skipping mkdocs.yml update")
            return

        print(f"\nUpdating {mkdocs_path}...")

        # Read current mkdocs.yml as text
        content = mkdocs_path.read_text()

        # Build the new Commands section
        commands_nav = self.build_commands_nav_section(nav)

        # Generate YAML for the Commands section
        commands_yaml = yaml.dump(
            [{"Commands": commands_nav}],
            default_flow_style=False,
            sort_keys=False,
            allow_unicode=True,
            indent=2,
        )

        # Indent the commands section properly (2 spaces for nav items)
        indented_commands = "\n".join(
            "  " + line if line.strip() else line
            for line in commands_yaml.strip().split("\n")
        )

        # Find and replace the Commands section in nav
        # Pattern: find "  - Commands:" and everything until the next "  - " at same level or end of nav
        pattern = r'(  - Commands:.*?)(?=\n  - [A-Z]|\ntheme:|\nextra:|\n[a-z_]+:|\Z)'

        if re.search(pattern, content, re.DOTALL):
            new_content = re.sub(pattern, indented_commands, content, flags=re.DOTALL)
            mkdocs_path.write_text(new_content)
            print(f"  Updated: {mkdocs_path}")
        else:
            print("Warning: Could not find Commands section in mkdocs.yml nav")
            print("  Saving navigation to _nav.yml for manual integration")

    def build_commands_nav_section(self, nav: list) -> list:
        """Build the Commands nav section from generated navigation."""
        commands_nav = []

        for item in nav:
            if isinstance(item, dict):
                for key, value in item.items():
                    if key == "Commands":
                        # This is the top-level Commands index
                        commands_nav.append({"Overview": value})
                    else:
                        # This is a command group
                        commands_nav.append({key: value})

        return commands_nav

    def clean_output(self) -> None:
        """Clean the output directory."""
        if self.output_dir.exists():
            print(f"Cleaning {self.output_dir}...")
            shutil.rmtree(self.output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)

    def generate_all(self, update_mkdocs: bool = False) -> None:
        """Generate all documentation."""
        print("\nGenerating documentation...")

        # Load API specs for documentation links
        self.load_api_specs()

        # Generate main index
        self.generate_commands_index()

        # Generate command groups
        for group_name, group_node in self.tree.children.items():
            print(f"\nGenerating {group_name}...")
            self.generate_command_group(group_name, group_node)

        # Generate navigation
        print("\nGenerating navigation...")
        nav = self.generate_navigation()
        self.save_navigation(nav)

        # Update mkdocs.yml if requested
        if update_mkdocs:
            self.update_mkdocs_yml(nav)

        # Summary
        total = self.count_commands()
        print(f"\nGeneration complete!")
        print(f"  Total commands documented: {total}")
        print(f"  Files generated: {len(self.generated_files)}")
        print(f"  Navigation saved to: docs/commands/_nav.yml")
        if update_mkdocs:
            print(f"  mkdocs.yml updated with Commands navigation")


def main():
    parser = argparse.ArgumentParser(
        description="Generate vesctl CLI documentation"
    )
    parser.add_argument(
        "--vesctl",
        default="./vesctl",
        help="Path to vesctl binary (default: ./vesctl)",
    )
    parser.add_argument(
        "--output",
        default="docs/commands",
        help="Output directory (default: docs/commands)",
    )
    parser.add_argument(
        "--templates",
        default="scripts/templates",
        help="Templates directory (default: scripts/templates)",
    )
    parser.add_argument(
        "--clean",
        action="store_true",
        help="Clean output directory before generating",
    )
    parser.add_argument(
        "--nav-only",
        action="store_true",
        help="Only generate navigation, skip documentation",
    )
    parser.add_argument(
        "--update-mkdocs",
        action="store_true",
        help="Update mkdocs.yml with generated Commands navigation",
    )

    args = parser.parse_args()

    generator = VesctlDocsGenerator(
        vesctl_path=args.vesctl,
        output_dir=args.output,
        template_dir=args.templates,
    )

    # Load spec
    generator.load_spec()

    if args.nav_only:
        nav = generator.generate_navigation()
        generator.save_navigation(nav)
        if args.update_mkdocs:
            generator.update_mkdocs_yml(nav)
        print("Navigation generated successfully!")
        return

    if args.clean:
        generator.clean_output()

    # Generate all documentation
    generator.generate_all(update_mkdocs=args.update_mkdocs)


if __name__ == "__main__":
    main()
