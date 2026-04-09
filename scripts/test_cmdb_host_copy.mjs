import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const repoRoot = process.cwd();
const targetPath = path.join(repoRoot, 'web', 'src', 'views', 'cmdb', 'cmdbHost.vue');
const source = fs.readFileSync(targetPath, 'utf8');

const expectedText = [
  '搜索',
  '重置',
  '新建',
  '导入主机',
  'Excel导入',
  '批量同步',
  '部署 Agent',
  '卸载 Agent',
  '选择文件',
  '取消',
  '主机详情',
  '基本信息',
  '连接中心',
  '终端审计'
];

const mojibakeText = [
  '鎼滅储',
  '閲嶇疆',
  '鏂板缓',
  '瀵煎叆涓绘満',
  '鎵归噺鍚屾',
  '閮ㄧ讲 Agent',
  '鍗歌浇 Agent',
  '闅惧嫚鏂囦欢',
  '????????'
];

const missing = expectedText.filter((text) => !source.includes(text));
const unexpected = mojibakeText.filter((text) => source.includes(text));

if (missing.length || unexpected.length) {
  if (missing.length) {
    console.error('Missing expected copy:', missing.join(', '));
  }
  if (unexpected.length) {
    console.error('Found mojibake copy:', unexpected.join(', '));
  }
  process.exit(1);
}

console.log('cmdbHost copy regression check passed');
