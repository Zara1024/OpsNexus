from __future__ import annotations

import csv
import sys
import tempfile
import unittest
from pathlib import Path
from zipfile import ZipFile

sys.path.insert(0, str(Path(__file__).resolve().parent))

from update_execution_tables import HEADERS, apply_updates_to_demo_root


def write_csv(path: Path, rows: list[dict[str, str]]) -> None:
    with path.open("w", encoding="utf-8-sig", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=HEADERS)
        writer.writeheader()
        writer.writerows(rows)


class UpdateExecutionTablesTests(unittest.TestCase):
    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.demo_root = Path(self._tmp.name)

        self.base_rows = [
            {
                "序号": "1",
                "套件类型": "基础版",
                "一级模块": "资产管理",
                "菜单位置": "资产管理 > 主机管理",
                "菜单页面": "主机管理",
                "子菜单": "主机管理",
                "页面地址": "/cmdb/ecs",
                "用例编号": "HOST-01",
                "用例名称": "页面能正常打开",
                "来源文件": "02-资产管理\\01-主机管理-测试用例.md",
                "文件行号": "10",
                "是否通过": "",
                "执行人": "",
                "执行日期": "",
                "Bug编号": "",
                "优先级": "中",
                "备注": "",
            }
        ]
        self.button_rows = [
            {
                "序号": "1",
                "套件类型": "按钮级",
                "一级模块": "仪表盘",
                "菜单位置": "仪表盘",
                "菜单页面": "仪表盘",
                "子菜单": "仪表盘",
                "页面地址": "/dashboard",
                "用例编号": "DASH-BTN-01",
                "用例名称": "刷新 按钮",
                "来源文件": "15-重点菜单按钮级测试用例\\07-仪表盘-按钮级测试用例.md",
                "文件行号": "20",
                "是否通过": "",
                "执行人": "",
                "执行日期": "",
                "Bug编号": "",
                "优先级": "高",
                "备注": "existing note",
            }
        ]
        self.all_rows = [self.base_rows[0].copy(), self.button_rows[0].copy()]

        write_csv(self.demo_root / "opsnexus-base-test-execution.csv", self.base_rows)
        write_csv(self.demo_root / "opsnexus-button-test-execution.csv", self.button_rows)
        write_csv(self.demo_root / "opsnexus-all-test-execution.csv", self.all_rows)

    def tearDown(self) -> None:
        self._tmp.cleanup()

    def read_rows(self, filename: str) -> list[dict[str, str]]:
        with (self.demo_root / filename).open("r", encoding="utf-8-sig", newline="") as handle:
            return list(csv.DictReader(handle))

    def test_updates_matching_rows_in_csv_and_xlsx(self) -> None:
        summary = apply_updates_to_demo_root(
            self.demo_root,
            [
                {
                    "用例编号": "HOST-01",
                    "是否通过": "通过",
                    "执行人": "Codex",
                    "执行日期": "2026-03-23",
                    "备注": "host page ok",
                },
                {
                    "用例编号": "DASH-BTN-01",
                    "是否通过": "失败",
                    "执行人": "Codex",
                    "执行日期": "2026-03-23",
                    "Bug编号": "BUG-1",
                    "备注": "button failed",
                },
            ],
        )

        self.assertEqual(summary["updated_case_ids"], ["HOST-01", "DASH-BTN-01"])

        base_rows = self.read_rows("opsnexus-base-test-execution.csv")
        self.assertEqual(base_rows[0]["是否通过"], "通过")
        self.assertEqual(base_rows[0]["执行人"], "Codex")
        self.assertEqual(base_rows[0]["执行日期"], "2026-03-23")
        self.assertEqual(base_rows[0]["备注"], "host page ok")

        button_rows = self.read_rows("opsnexus-button-test-execution.csv")
        self.assertEqual(button_rows[0]["是否通过"], "失败")
        self.assertEqual(button_rows[0]["Bug编号"], "BUG-1")
        self.assertEqual(button_rows[0]["备注"], "existing note | button failed")

        all_rows = self.read_rows("opsnexus-all-test-execution.csv")
        by_case = {row["用例编号"]: row for row in all_rows}
        self.assertEqual(by_case["HOST-01"]["是否通过"], "通过")
        self.assertEqual(by_case["DASH-BTN-01"]["是否通过"], "失败")

        workbook = self.demo_root / "opsnexus-test-execution.xlsx"
        self.assertTrue(workbook.exists())
        with ZipFile(workbook) as archive:
            sheet_one = archive.read("xl/worksheets/sheet1.xml").decode("utf-8")

        self.assertIn("HOST-01", sheet_one)
        self.assertIn("DASH-BTN-01", sheet_one)
        self.assertIn("host page ok", sheet_one)
        self.assertIn("existing note | button failed", sheet_one)

    def test_raises_for_unknown_case_id(self) -> None:
        with self.assertRaisesRegex(ValueError, "UNKNOWN-01"):
            apply_updates_to_demo_root(
                self.demo_root,
                [{"用例编号": "UNKNOWN-01", "是否通过": "通过"}],
            )


if __name__ == "__main__":
    unittest.main()
