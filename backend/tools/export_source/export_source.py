from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent.parent
OUTPUT_FILE = ROOT / "tools" / "export_source" / ".out"

IGNORE_DIRS = {
    ".git",
    ".idea",
    ".vscode",
    "vendor",
    "node_modules",
    "tmp",
    "tools",
}

IGNORE_FILES = {
    OUTPUT_FILE.name,
    "go.mod",
    "go.sum",
}

INCLUDE_EXTENSIONS = {
    ".go",
    ".sql",
    ".mod",
    ".sum",
    ".env.example",
    ".yaml",
    ".yml",
    ".json",
}

OUTPUT_FILE.write_text("", encoding="utf-8")

for path in sorted(ROOT.rglob("*")):
    if not path.is_file():
        continue

    if any(part in IGNORE_DIRS for part in path.parts):
        continue

    if path.name in IGNORE_FILES:
        continue

    if (
        path.suffix not in INCLUDE_EXTENSIONS
        and path.name not in {"Dockerfile", "Makefile"}
    ):
        continue

    relative_path = path.relative_to(ROOT)

    with OUTPUT_FILE.open("a", encoding="utf-8") as out:
        out.write(f"{relative_path.as_posix()}\n")

        try:
            content = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            continue

        out.write(content)

        if not content.endswith("\n"):
            out.write("\n")

        out.write("\n")
        out.write("=" * 100)
        out.write("\n\n")

print("Done!")