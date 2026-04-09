import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';

const repoRoot = process.cwd();
const mainPath = path.join(repoRoot, 'web', 'src', 'main.js');
const mainSource = fs.readFileSync(mainPath, 'utf8');

const vueFiles = [
  path.join(repoRoot, 'web', 'src', 'components', 'Tags.vue'),
  path.join(repoRoot, 'web', 'src', 'views', 'cmdb', 'Host', 'HostSsh.vue'),
  path.join(repoRoot, 'web', 'src', 'views', 'task', 'Job', 'AnsibleFlowTemp.vue')
];

const failures = [];

if (mainSource.includes('app.use(handleTree)')) {
  failures.push('main.js should not register handleTree via app.use(...)');
}

for (const file of vueFiles) {
  const source = fs.readFileSync(file, 'utf8');
  if (source.includes('size="medium"')) {
    failures.push(`${path.relative(repoRoot, file)} still uses size="medium"`);
  }
}

if (failures.length) {
  console.error(failures.join('\n'));
  process.exit(1);
}

console.log('frontend warning regression check passed');
