from __future__ import annotations

import csv
import re
from collections import Counter
from pathlib import Path
from xml.sax.saxutils import escape
from zipfile import ZIP_DEFLATED, ZipFile


ROOT = Path(__file__).resolve().parents[1].parent / "demo"
BUTTON_DIR = next(path for path in ROOT.iterdir() if path.is_dir() and path.name.startswith("15-"))

CASE_RE = re.compile(r"^###\s+(\S+)\s+(.+?)\s*$")

HEADERS = [
    "\u5e8f\u53f7",
    "\u5957\u4ef6\u7c7b\u578b",
    "\u4e00\u7ea7\u6a21\u5757",
    "\u83dc\u5355\u4f4d\u7f6e",
    "\u83dc\u5355\u9875\u9762",
    "\u5b50\u83dc\u5355",
    "\u9875\u9762\u5730\u5740",
    "\u7528\u4f8b\u7f16\u53f7",
    "\u7528\u4f8b\u540d\u79f0",
    "\u6765\u6e90\u6587\u4ef6",
    "\u6587\u4ef6\u884c\u53f7",
    "\u662f\u5426\u901a\u8fc7",
    "\u6267\u884c\u4eba",
    "\u6267\u884c\u65e5\u671f",
    "Bug\u7f16\u53f7",
    "\u4f18\u5148\u7ea7",
    "\u5907\u6ce8",
]


def parse_markdown(md_path: Path, suite_type: str) -> list[dict[str, object]]:
    text = md_path.read_text(encoding="utf-8")
    lines = text.splitlines()

    title = ""
    menu = ""
    route = ""

    for line in lines:
        if line.startswith("# "):
            title = (
                line[2:]
                .strip()
                .replace(" \u6309\u94ae\u7ea7\u6d4b\u8bd5\u7528\u4f8b", "")
                .replace(" \u6d4b\u8bd5\u7528\u4f8b", "")
            )
            break

    for line in lines:
        if line.startswith("- \u83dc\u5355\u4f4d\u7f6e\uff1a"):
            menu = line.split("\uff1a", 1)[1].strip()
        elif line.startswith("- \u9875\u9762\u5730\u5740\uff1a"):
            route = line.split("\uff1a", 1)[1].strip().strip("`")

    module = menu.split(" > ")[0].strip() if " > " in menu else menu.strip()
    submenu = menu.split(" > ")[-1].strip() if menu else title

    rows: list[dict[str, object]] = []
    for line_no, line in enumerate(lines, start=1):
        match = CASE_RE.match(line)
        if not match:
            continue
        case_id, case_name = match.groups()
        rows.append(
            {
                "\u5957\u4ef6\u7c7b\u578b": suite_type,
                "\u4e00\u7ea7\u6a21\u5757": module,
                "\u83dc\u5355\u4f4d\u7f6e": menu,
                "\u83dc\u5355\u9875\u9762": title,
                "\u5b50\u83dc\u5355": submenu,
                "\u9875\u9762\u5730\u5740": route,
                "\u7528\u4f8b\u7f16\u53f7": case_id,
                "\u7528\u4f8b\u540d\u79f0": case_name,
                "\u6765\u6e90\u6587\u4ef6": str(md_path.relative_to(ROOT)),
                "\u6587\u4ef6\u884c\u53f7": line_no,
                "\u662f\u5426\u901a\u8fc7": "",
                "\u6267\u884c\u4eba": "",
                "\u6267\u884c\u65e5\u671f": "",
                "Bug\u7f16\u53f7": "",
                "\u4f18\u5148\u7ea7": "\u9ad8" if suite_type == "\u6309\u94ae\u7ea7" else "\u4e2d",
                "\u5907\u6ce8": "",
            }
        )
    return rows


def list_files() -> tuple[list[Path], list[Path]]:
    base_files = sorted(
        file
        for file in ROOT.rglob("*.md")
        if not file.is_relative_to(BUTTON_DIR)
        and not file.name.startswith("00-")
        and not file.name.startswith("opsnexus-")
    )
    button_files = sorted(
        file for file in BUTTON_DIR.glob("*.md") if not file.name.startswith("00-")
    )
    return base_files, button_files


BASE_FILES, BUTTON_FILES = list_files()
BASE_ROWS = [row for path in BASE_FILES for row in parse_markdown(path, "\u57fa\u7840\u7248")]
BUTTON_ROWS = [row for path in BUTTON_FILES for row in parse_markdown(path, "\u6309\u94ae\u7ea7")]
ALL_ROWS = sorted(
    BASE_ROWS + BUTTON_ROWS,
    key=lambda row: (
        row["\u5957\u4ef6\u7c7b\u578b"],
        row["\u4e00\u7ea7\u6a21\u5757"],
        row["\u83dc\u5355\u9875\u9762"],
        row["\u7528\u4f8b\u7f16\u53f7"],
    ),
)


def add_serial(rows: list[dict[str, object]]) -> None:
    for index, row in enumerate(rows, start=1):
        row["\u5e8f\u53f7"] = index


def write_csv(path: Path, rows: list[dict[str, object]]) -> None:
    with path.open("w", encoding="utf-8-sig", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=HEADERS)
        writer.writeheader()
        for row in rows:
            writer.writerow({header: row.get(header, "") for header in HEADERS})


def excel_col(index: int) -> str:
    result = ""
    while index:
        index, rem = divmod(index - 1, 26)
        result = chr(65 + rem) + result
    return result


def make_sheet_xml(rows: list[dict[str, object]]) -> str:
    widths: list[str] = []
    for col in range(1, len(HEADERS) + 1):
        width = 18
        if col in (4, 5, 10, 17):
            width = 28
        if col == 7:
            width = 20
        if col == 8:
            width = 14
        widths.append(f'<col min="{col}" max="{col}" width="{width}" customWidth="1"/>')

    def cell(ref: str, value: object) -> str:
        text = escape(str(value))
        return f'<c r="{ref}" t="inlineStr"><is><t xml:space="preserve">{text}</t></is></c>'

    data = [HEADERS] + [[row.get(header, "") for header in HEADERS] for row in rows]
    row_xml: list[str] = []
    for row_index, line in enumerate(data, start=1):
        cells = [
            cell(f"{excel_col(col_index)}{row_index}", value)
            for col_index, value in enumerate(line, start=1)
        ]
        row_xml.append(f'<row r="{row_index}">{"".join(cells)}</row>')

    end_ref = f"A1:{excel_col(len(HEADERS))}{len(data)}"
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'
        f'<dimension ref="{end_ref}"/>'
        '<sheetViews><sheetView workbookViewId="0"/></sheetViews>'
        '<sheetFormatPr defaultRowHeight="15"/>'
        f'<cols>{"".join(widths)}</cols>'
        f'<sheetData>{"".join(row_xml)}</sheetData>'
        "</worksheet>"
    )


def make_workbook_xml(sheet_names: list[str]) -> str:
    sheets = "".join(
        f'<sheet name="{escape(name)}" sheetId="{index}" r:id="rId{index}"/>'
        for index, name in enumerate(sheet_names, start=1)
    )
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" '
        'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">'
        f"<sheets>{sheets}</sheets>"
        "</workbook>"
    )


def make_workbook_rels(sheet_count: int) -> str:
    rels = "".join(
        f'<Relationship Id="rId{index}" '
        'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" '
        f'Target="worksheets/sheet{index}.xml"/>'
        for index in range(1, sheet_count + 1)
    )
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
        f"{rels}"
        "</Relationships>"
    )


ROOT_RELS = (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
    '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">'
    '<Relationship Id="rId1" '
    'Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" '
    'Target="xl/workbook.xml"/>'
    "</Relationships>"
)


def make_content_types(sheet_count: int) -> str:
    overrides = [
        '<Override PartName="/xl/workbook.xml" '
        'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>'
    ]
    overrides.extend(
        f'<Override PartName="/xl/worksheets/sheet{index}.xml" '
        'ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>'
        for index in range(1, sheet_count + 1)
    )
    return (
        '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>'
        '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">'
        '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>'
        '<Default Extension="xml" ContentType="application/xml"/>'
        f'{"".join(overrides)}'
        "</Types>"
    )


def write_xlsx(path: Path, sheets: list[tuple[str, list[dict[str, object]]]]) -> None:
    with ZipFile(path, "w", ZIP_DEFLATED) as workbook:
        workbook.writestr("[Content_Types].xml", make_content_types(len(sheets)))
        workbook.writestr("_rels/.rels", ROOT_RELS)
        workbook.writestr("xl/workbook.xml", make_workbook_xml([name for name, _ in sheets]))
        workbook.writestr("xl/_rels/workbook.xml.rels", make_workbook_rels(len(sheets)))
        for index, (_, rows) in enumerate(sheets, start=1):
            workbook.writestr(f"xl/worksheets/sheet{index}.xml", make_sheet_xml(rows))


def write_guide(base_count: int, button_count: int, all_count: int) -> None:
    guide_text = (
        "# \u6d4b\u8bd5\u6267\u884c\u8868\u8bf4\u660e\n\n"
        "## 1. \u672c\u6b21\u751f\u6210\u7684\u6267\u884c\u8868\u6587\u4ef6\n"
        "- `opsnexus-base-test-execution.csv`\n"
        "- `opsnexus-button-test-execution.csv`\n"
        "- `opsnexus-all-test-execution.csv`\n"
        "- `opsnexus-test-execution.xlsx`\n\n"
        "## 2. \u5efa\u8bae\u600e\u4e48\u7528\n"
        "1. \u53ea\u505a\u83dc\u5355\u7ea7\u6d4b\u8bd5\u65f6\uff0c\u7528 `opsnexus-base-test-execution.csv`\u3002\n"
        "2. \u505a\u56de\u5f52\u6216\u7ec6\u6d4b\u65f6\uff0c\u7528 `opsnexus-button-test-execution.csv`\u3002\n"
        "3. \u60f3\u7edf\u4e00\u6392\u671f\u3001\u6d3e\u4eba\u3001\u8ddf\u8fdb Bug \u65f6\uff0c\u7528 `opsnexus-all-test-execution.csv` \u6216 `opsnexus-test-execution.xlsx`\u3002\n\n"
        "## 3. \u5217\u8bf4\u660e\n"
        "- `\u5957\u4ef6\u7c7b\u578b`\uff1a\u57fa\u7840\u7248 / \u6309\u94ae\u7ea7\u3002\n"
        "- `\u4e00\u7ea7\u6a21\u5757`\uff1a\u4f8b\u5982\u8d44\u4ea7\u7ba1\u7406\u3001\u5bb9\u5668\u7ba1\u7406\u3001\u670d\u52a1\u7ba1\u7406\u3002\n"
        "- `\u83dc\u5355\u4f4d\u7f6e`\uff1a\u83dc\u5355\u5b8c\u6574\u8def\u5f84\u3002\n"
        "- `\u83dc\u5355\u9875\u9762`\uff1a\u9875\u9762\u540d\u79f0\u3002\n"
        "- `\u7528\u4f8b\u7f16\u53f7`\uff1a\u6267\u884c\u65f6\u76f4\u63a5\u5f15\u7528\u8fd9\u4e2a\u7f16\u53f7\u3002\n"
        "- `\u662f\u5426\u901a\u8fc7`\uff1a\u5efa\u8bae\u586b\u5199 \u901a\u8fc7 / \u5931\u8d25 / \u963b\u585e\u3002\n"
        "- `\u6267\u884c\u4eba`\uff1a\u586b\u5199\u5b9e\u9645\u6d4b\u8bd5\u4eba\u3002\n"
        "- `\u6267\u884c\u65e5\u671f`\uff1a\u5efa\u8bae\u683c\u5f0f `2026-03-23`\u3002\n"
        "- `Bug\u7f16\u53f7`\uff1a\u6ca1\u6709\u5c31\u7559\u7a7a\u3002\n"
        "- `\u5907\u6ce8`\uff1a\u8bb0\u5f55\u73b0\u8c61\u3001\u73af\u5883\u5dee\u5f02\u6216\u8865\u5145\u8bf4\u660e\u3002\n\n"
        "## 4. \u5f53\u524d\u7edf\u8ba1\n"
        f"- \u57fa\u7840\u7248\u83dc\u5355\u6587\u4ef6\uff1a{len(BASE_FILES)}\n"
        f"- \u6309\u94ae\u7ea7\u83dc\u5355\u6587\u4ef6\uff1a{len(BUTTON_FILES)}\n"
        f"- \u57fa\u7840\u7248\u7528\u4f8b\u603b\u6570\uff1a{base_count}\n"
        f"- \u6309\u94ae\u7ea7\u7528\u4f8b\u603b\u6570\uff1a{button_count}\n"
        f"- \u5168\u91cf\u7528\u4f8b\u603b\u6570\uff1a{all_count}\n"
    )
    (ROOT / "00-test-execution-guide.md").write_text(guide_text, encoding="utf-8-sig")


def main() -> None:
    add_serial(BASE_ROWS)
    add_serial(BUTTON_ROWS)
    add_serial(ALL_ROWS)

    write_csv(ROOT / "opsnexus-base-test-execution.csv", BASE_ROWS)
    write_csv(ROOT / "opsnexus-button-test-execution.csv", BUTTON_ROWS)
    write_csv(ROOT / "opsnexus-all-test-execution.csv", ALL_ROWS)

    base_counter = Counter(
        (row["\u4e00\u7ea7\u6a21\u5757"], row["\u83dc\u5355\u9875\u9762"]) for row in BASE_ROWS
    )
    button_counter = Counter(
        (row["\u4e00\u7ea7\u6a21\u5757"], row["\u83dc\u5355\u9875\u9762"]) for row in BUTTON_ROWS
    )

    stat_rows: list[dict[str, object]] = []
    for index, key in enumerate(sorted(set(base_counter) | set(button_counter)), start=1):
        module, page = key
        stat_rows.append(
            {
                "\u5e8f\u53f7": index,
                "\u5957\u4ef6\u7c7b\u578b": "\u7edf\u8ba1",
                "\u4e00\u7ea7\u6a21\u5757": module,
                "\u83dc\u5355\u4f4d\u7f6e": page,
                "\u83dc\u5355\u9875\u9762": "",
                "\u5b50\u83dc\u5355": "",
                "\u9875\u9762\u5730\u5740": "",
                "\u7528\u4f8b\u7f16\u53f7": f"\u57fa\u7840\u7248:{base_counter.get(key, 0)}",
                "\u7528\u4f8b\u540d\u79f0": f"\u6309\u94ae\u7ea7:{button_counter.get(key, 0)}",
                "\u6765\u6e90\u6587\u4ef6": "",
                "\u6587\u4ef6\u884c\u53f7": "",
                "\u662f\u5426\u901a\u8fc7": "",
                "\u6267\u884c\u4eba": "",
                "\u6267\u884c\u65e5\u671f": "",
                "Bug\u7f16\u53f7": "",
                "\u4f18\u5148\u7ea7": "",
                "\u5907\u6ce8": "",
            }
        )

    write_xlsx(
        ROOT / "opsnexus-test-execution.xlsx",
        [
            ("\u5168\u91cf\u6267\u884c\u8868", ALL_ROWS),
            ("\u57fa\u7840\u7248\u6267\u884c\u8868", BASE_ROWS),
            ("\u6309\u94ae\u7ea7\u6267\u884c\u8868", BUTTON_ROWS),
            ("\u8986\u76d6\u7edf\u8ba1", stat_rows),
        ],
    )

    write_guide(len(BASE_ROWS), len(BUTTON_ROWS), len(ALL_ROWS))

    print(f"ROOT={ROOT}")
    print(f"BASE_CASES={len(BASE_ROWS)}")
    print(f"BUTTON_CASES={len(BUTTON_ROWS)}")
    print(f"ALL_CASES={len(ALL_ROWS)}")


if __name__ == "__main__":
    main()
