from __future__ import annotations

import argparse
import csv
import json
from pathlib import Path
from xml.sax.saxutils import escape
from zipfile import ZIP_DEFLATED, ZipFile


HEADERS = [
    "序号",
    "套件类型",
    "一级模块",
    "菜单位置",
    "菜单页面",
    "子菜单",
    "页面地址",
    "用例编号",
    "用例名称",
    "来源文件",
    "文件行号",
    "是否通过",
    "执行人",
    "执行日期",
    "Bug编号",
    "优先级",
    "备注",
]

RESULT_FIELDS = ("是否通过", "执行人", "执行日期", "Bug编号", "备注")
VALID_STATUS = {"通过", "失败", "阻塞"}
SHEET_NAMES = (
    "全量执行表",
    "基础版执行表",
    "按钮级执行表",
)


def default_demo_root() -> Path:
    return Path(__file__).resolve().parents[2] / "demo"


def read_csv_rows(path: Path) -> list[dict[str, str]]:
    with path.open("r", encoding="utf-8-sig", newline="") as handle:
        return list(csv.DictReader(handle))


def write_csv_rows(path: Path, rows: list[dict[str, str]]) -> None:
    with path.open("w", encoding="utf-8-sig", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=HEADERS)
        writer.writeheader()
        for row in rows:
            writer.writerow({header: row.get(header, "") for header in HEADERS})


def merge_remark(existing: str, incoming: str) -> str:
    existing = (existing or "").strip()
    incoming = (incoming or "").strip()
    if not existing:
        return incoming
    if not incoming or incoming == existing:
        return existing
    return f"{existing} | {incoming}"


def apply_update_to_rows(rows: list[dict[str, str]], case_id: str, update: dict[str, str]) -> bool:
    matched = False
    for row in rows:
        if row.get("用例编号") != case_id:
            continue
        matched = True
        for field in RESULT_FIELDS:
            if field not in update:
                continue
            value = (update.get(field) or "").strip()
            if field == "备注":
                row[field] = merge_remark(row.get(field, ""), value)
            else:
                row[field] = value
    return matched


def validate_update(update: dict[str, str]) -> dict[str, str]:
    normalized = {}
    for key, value in update.items():
        normalized[str(key)] = "" if value is None else str(value)

    case_id = normalized.get("用例编号", "").strip()
    if not case_id:
        raise ValueError(f"Missing 用例编号 in update: {update}")

    status = normalized.get("是否通过", "").strip()
    if status and status not in VALID_STATUS:
        raise ValueError(f"Unsupported 是否通过={status} for 用例编号={case_id}")

    normalized["用例编号"] = case_id
    if status:
        normalized["是否通过"] = status
    return normalized


def excel_col(index: int) -> str:
    result = ""
    while index:
        index, rem = divmod(index - 1, 26)
        result = chr(65 + rem) + result
    return result


def make_sheet_xml(rows: list[dict[str, str]]) -> str:
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

    def cell(ref: str, value: str) -> str:
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


def make_workbook_xml(sheet_names: tuple[str, ...]) -> str:
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


def write_xlsx(path: Path, sheets: list[tuple[str, list[dict[str, str]]]]) -> None:
    with ZipFile(path, "w", ZIP_DEFLATED) as workbook:
        workbook.writestr("[Content_Types].xml", make_content_types(len(sheets)))
        workbook.writestr("_rels/.rels", ROOT_RELS)
        workbook.writestr("xl/workbook.xml", make_workbook_xml(tuple(name for name, _ in sheets)))
        workbook.writestr("xl/_rels/workbook.xml.rels", make_workbook_rels(len(sheets)))
        for index, (_, rows) in enumerate(sheets, start=1):
            workbook.writestr(f"xl/worksheets/sheet{index}.xml", make_sheet_xml(rows))


def write_workbook(demo_root: Path, base_rows: list[dict[str, str]], button_rows: list[dict[str, str]], all_rows: list[dict[str, str]]) -> None:
    write_xlsx(
        demo_root / "opsnexus-test-execution.xlsx",
        [
            (SHEET_NAMES[0], all_rows),
            (SHEET_NAMES[1], base_rows),
            (SHEET_NAMES[2], button_rows),
        ],
    )


def apply_updates_to_demo_root(demo_root: Path, updates: list[dict[str, str]]) -> dict[str, list[str] | int]:
    base_path = demo_root / "opsnexus-base-test-execution.csv"
    button_path = demo_root / "opsnexus-button-test-execution.csv"
    all_path = demo_root / "opsnexus-all-test-execution.csv"

    base_rows = read_csv_rows(base_path)
    button_rows = read_csv_rows(button_path)
    all_rows = read_csv_rows(all_path)

    updated_case_ids: list[str] = []
    for raw_update in updates:
        update = validate_update(raw_update)
        case_id = update["用例编号"]

        matched_any = False
        matched_any |= apply_update_to_rows(all_rows, case_id, update)
        matched_any |= apply_update_to_rows(base_rows, case_id, update)
        matched_any |= apply_update_to_rows(button_rows, case_id, update)

        if not matched_any:
            raise ValueError(f"Unknown 用例编号: {case_id}")

        updated_case_ids.append(case_id)

    write_csv_rows(base_path, base_rows)
    write_csv_rows(button_path, button_rows)
    write_csv_rows(all_path, all_rows)
    write_workbook(demo_root, base_rows, button_rows, all_rows)

    return {
        "updated_count": len(updated_case_ids),
        "updated_case_ids": updated_case_ids,
    }


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Update OpsNexus demo execution tables by case id.")
    parser.add_argument("--demo-root", default=str(default_demo_root()))
    parser.add_argument("--updates-file", help="JSON file containing a list of update objects.")
    parser.add_argument("--case-id", help="Single case id to update.")
    parser.add_argument("--status", choices=sorted(VALID_STATUS))
    parser.add_argument("--executor")
    parser.add_argument("--date")
    parser.add_argument("--bug-id")
    parser.add_argument("--remark")
    return parser.parse_args()


def load_updates_from_args(args: argparse.Namespace) -> list[dict[str, str]]:
    if args.updates_file:
        updates = json.loads(Path(args.updates_file).read_text(encoding="utf-8"))
        if not isinstance(updates, list):
            raise ValueError("--updates-file must contain a JSON list")
        return updates

    if not args.case_id:
        raise ValueError("Either --updates-file or --case-id is required")

    return [
        {
            "用例编号": args.case_id,
            "是否通过": args.status or "",
            "执行人": args.executor or "",
            "执行日期": args.date or "",
            "Bug编号": args.bug_id or "",
            "备注": args.remark or "",
        }
    ]


def main() -> None:
    args = parse_args()
    updates = load_updates_from_args(args)
    summary = apply_updates_to_demo_root(Path(args.demo_root), updates)
    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
