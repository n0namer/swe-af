import './basefive-r14-contract.test.mjs';
import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { createHash } from 'node:crypto';

const root = new URL('../../', import.meta.url);
const pythonPath = 'swe_af/issue/build.py';
const goPath = 'go/internal/orch/build.go';
const contractPath = 'docs/architecture/basefive-r14-contracts.json';
const reviewPath = 'docs/architecture/basefive-r14-independent-c4-review.md';
const ownPath = 'docs/architecture/basefive-r14-independent-c4.test.mjs';
const diagramPaths = [
  'docs/architecture/basefive-r14-C4-context.md',
  'docs/architecture/basefive-r14-container.mmd',
  'docs/architecture/basefive-r14-component.mmd',
];
const pinned = {
  'AGENTS.md': '4a8876425a07de89597f3b56c6ba918b8dd89ddbce1abb72b4f9db9ec1782ae9',
  'docs/acceptance/c4-releasegate-reducer-frontier-20261013.md': 'feeb1fe6a88e53770831de890fe3674e6425dd40327deb913fbf73f3a01f1f17',
  '.archsteer/model.json': 'ed964b0a700a566bd9e42caa55551baee873706aa500fbd57c446a28b7971260',
  '.archsteer/architecture.md': '7d7ee7046812a9c3c5da1d3b7484015432c4cb3b3cb106a6d34e8bffeb65a3d2',
  [pythonPath]: '440909dfb60d257232658c18cd2f193f49d1c28446d7f630c9cebee4c80a1bcd',
  [goPath]: '05c39a2e3b41869b01a1482543c9875630a0777fbea24ea706c05727fe640cbc',
  'graphify-out/graph.json': 'e18cd58bd72e8975c07f305a988192227c529a39ab35752388ad0db706f5c277',
};
const read = path => readFileSync(new URL(path, root));
const hash = content => createHash('sha256').update(content).digest('hex');
const sourceBytes = Object.fromEntries(Object.keys(pinned).map(path => [path, read(path)]));
const source = Object.fromEntries(Object.entries(sourceBytes).map(([path, content]) => [path, content.toString('utf8').replace(/\r\n/g, '\n')]));
const contract = JSON.parse(read(contractPath));
const documents = Object.fromEntries(diagramPaths.map(path => [path, read(path).toString('utf8')]));
const graph = JSON.parse(source['graphify-out/graph.json']);
const arch = JSON.parse(source['.archsteer/model.json']);
const unique = values => [...new Set(values)].sort();
const normalize = value => value.replace(/\s+/g, ' ').trim();

function capture(text, pattern, name) {
  const matched = text.match(pattern);
  assert.ok(matched, 'source declaration exists: ' + name);
  return matched[1];
}

function checkDTO(candidate, inputs = source) {
  const python = inputs[pythonPath];
  const go = inputs[goPath];
  const inputBody = capture(go, /type buildInput struct \{\n([\s\S]*?)\n\}/, 'buildInput');
  const inputFields = inputBody.trim().split('\n').map(line => {
    const matched = line.trim().match(/^(\w+)\s+(.+?)\s+`json:"([^"]+)"`$/);
    assert.ok(matched, 'every Go input field parsed');
    return { field: matched[1], type: matched[2], json: matched[3] };
  });
  assert.deepEqual(candidate.api.goInput, inputFields, 'independent Go DTO fields/types/tags');
  const pythonResult = capture(python, /return IssueBuildResult\(\n([\s\S]*?)\n    \)\.model_dump\(\)/, 'IssueBuildResult');
  const goResult = capture(go, /buildResult := schemas.BuildResult\{\n([\s\S]*?)\n\t\}/, 'BuildResult');
  assert.deepEqual(candidate.api.pythonResultFields, [...pythonResult.matchAll(/^        (\w+)=/gm)].map(match => match[1]), 'independent Python result fields');
  assert.deepEqual(candidate.api.goResultFields, [...goResult.matchAll(/^\t\t(\w+):/gm)].map(match => match[1]), 'independent Go result fields');
  assert.equal(candidate.api.pythonSignature, normalize(capture(python, /^(async def implement_issue\([\s\S]*?\) -> dict:)/m, 'implement_issue')), 'independent Python entry');
  assert.equal(candidate.api.goSignature, capture(go, /^(func Build\([^\n]+)/m, 'Build'), 'independent Go entry');
  assert.match(python, /spec = IssueSpec\.model_validate\(issue\)/);
  assert.match(go, /if raw, err := json\.Marshal\(input\); err == nil \{\n\t\t_ = json\.Unmarshal\(raw, &in\)/);
  assert.match(go, /if workspaceHandle != nil \{\n\t\tbuildResultMap\["workspace_handle"\] = dumpToMap\(workspaceHandle\)/);
  assert.equal(candidate.dto.wireSchemaVersion, null, 'no imported DTO version proof');
  assert.equal(candidate.dto.wireCompatibility, 'UNVERIFIED', 'no imported DTO compatibility proof');
  assert.equal(candidate.dto.comparison, 'DIFFERENT_ENTRYPOINTS_NOT_PARITY');
  assert.equal(candidate.dto.runtimeExecution, 'UNVERIFIED');
}

const expectedEdges = [
  [
    ['caller', 'system', pythonPath, 'async def implement_issue('],
    ['system', 'agents', pythonPath, 'f"{node_id}.run_verifier"'],
    ['system', 'git', pythonPath, 'git_ops.add_worktree'],
    ['system', 'github', pythonPath, 'if cfg.enable_github_pr and coding_ok and commits:'],
  ],
  [
    ['caller', 'python', pythonPath, 'issue: dict,'],
    ['caller', 'go', goPath, 'input map[string]any'],
    ['python', 'agents', pythonPath, 'agent.call(target, **kwargs)'],
    ['go', 'agents', goPath, 'deps.App.Call(ctx, deps.target("plan"), planKwargs)'],
    ['python', 'local', pythonPath, 'git_ops.remove_worktree'],
    ['go', 'local', goPath, 'prepareSingleRepo(ctx, deps, cfg, repoPath)'],
    ['python', 'github', pythonPath, 'f"{node_id}.run_github_pr"'],
    ['go', 'github', goPath, 'singleRepoPR(ctx, deps, cfg'],
  ],
  [
    ['entry', 'impl', pythonPath, 'return await _implement_issue_impl('],
    ['impl', 'loop', pythonPath, 'loop_result = await run_coding_loop('],
    ['impl', 'checkpoint', pythonPath, 'await _checkpoint_interrupted_issue('],
    ['impl', 'verify', pythonPath, 'verification = await _run_verification('],
    ['impl', 'result', pythonPath, 'return IssueBuildResult('],
    ['handlers', 'build', goPath, '"build": Build'],
    ['build', 'dispatch', goPath, 'deps.Call(ctx, "generate_fix_issues"'],
    ['build', 'goresult', goPath, 'return buildResultMap, rf'],
  ],
];

function checkEdges(candidate, docs = documents, inputs = source) {
  assert.deepEqual(candidate.diagrams.map(diagram => diagram.file), diagramPaths);
  candidate.diagrams.forEach((diagram, index) => {
    const actual = diagram.edges.map(({ from, to, source: path, anchor }) => [from, to, path, anchor]);
    assert.deepEqual(actual, expectedEdges[index], 'independent edge ownership and topology');
    const renderedEdges = [...docs[diagram.file].matchAll(/^\s*(\w+) -->\|([^|\n]+)\| (\w+)\s*$/gm)].map(match => [match[1], match[3], match[2]]);
    assert.deepEqual(renderedEdges, diagram.edges.map(edge => [edge.from, edge.to, edge.label]), 'independent diagram edges');
    for (const item of [...diagram.nodes, ...diagram.edges]) {
      assert.ok([pythonPath, goPath].includes(item.source), 'only inspected code sources');
      assert.ok(item.anchor.length > 0 && inputs[item.source].includes(item.anchor), 'independent source witness');
    }
  });
  const python = inputs[pythonPath];
  const go = inputs[goPath];
  assert.match(python, /if cfg\.verify and coding_ok and commits:\n            verification = await _run_verification\(/, 'conditional verification edge');
  assert.match(python, /if cfg\.enable_github_pr and coding_ok and commits:\n            pr_url = await _maybe_create_pr\(/, 'PR is not verifier-gated');
  assert.match(python, /success = coding_ok and bool\(commits\) and verify_ok/);
  assert.match(python, /except asyncio\.CancelledError:\n        await _checkpoint_interrupted_issue\([\s\S]*?\n        raise\n    except Exception/);
  assert.match(python, /finally:\n        if not cfg\.keep_worktree:\n            await asyncio\.to_thread\(git_ops\.remove_worktree/);
  assert.match(go, /success := asBool\(verification\["passed"\]\)/);
  assert.match(go, /rf := &agent\.ReasonerFailed\{Message: msg\}[\s\S]*?rf\.Result = buildResultMap[\s\S]*?return buildResultMap, rf/);
}

function checkCoverage(candidate, graphData = graph, archData = arch) {
  const graphPaths = unique(graphData.nodes.map(node => node.source_file).filter(path => typeof path === 'string'));
  const archPaths = Object.keys(archData.components);
  const goFiles = graphPaths.filter(path => path.endsWith('.go'));
  const directories = unique(goFiles.map(path => path.slice(0, path.lastIndexOf('/'))));
  assert.equal(candidate.coverage.archsteerGoFileCount, archPaths.filter(path => path.endsWith('.go')).length, 'ArchSteer Go blind spot');
  assert.equal(candidate.coverage.graphifyGoFileCount, goFiles.length, 'Graphify Go inventory');
  assert.deepEqual(candidate.coverage.graphifyGoDirectories, directories);
  assert.equal(goFiles.length, 200);
  assert.equal(directories.length, 36);
  assert.equal(candidate.coverage.archsteerGoFileCount, 0);
  const labels = Object.fromEntries([pythonPath, goPath].map(path => [path, graphData.nodes.filter(node => node.source_file === path).map(node => node.label)]));
  assert.equal(labels[pythonPath].length, 10);
  assert.equal(labels[goPath].length, 14);
  for (const symbol of ['Build', 'Handlers', 'buildInput']) assert.ok(labels[goPath].some(label => label === symbol || label === symbol + '()'), 'indexed Go symbol: ' + symbol);
  for (const symbol of ['_checkpoint_if_dirty', '_checkpoint_interrupted_issue']) {
    assert.ok(archData.components[pythonPath].exported_apis.includes(symbol));
    assert.ok(!labels[pythonPath].includes(symbol + '()'));
    assert.ok(source[pythonPath].includes('async def ' + symbol + '('));
  }
  assert.deepEqual(candidate.coverage.pythonMissingGraphSymbols, ['_checkpoint_if_dirty', '_checkpoint_interrupted_issue']);
  const imports = [...source[goPath].matchAll(/"github\.com\/Agent-Field\/SWE-AF\/(go\/internal\/[^"\n]+)"/g)].map(match => match[1]);
  for (const directory of imports) assert.ok(directories.includes(directory), 'indexed imported directory: ' + directory);
  for (const field of ['repositoryCompleteness', 'graphSourceFreshness', 'packageRuntimeCoverage']) assert.equal(candidate.coverage[field], 'UNVERIFIED');
  return { archsteerGoFiles: 0, graphifyGoFiles: goFiles.length, goDirectories: directories.length, imports, targetLabels: labels, graphNodes: graphData.nodes.length, graphEdges: graphData.edges.length };
}

test('independent exact seven source fingerprints and non-promoting identity', () => {
  assert.deepEqual(contract.sources, pinned);
  for (const [path, digest] of Object.entries(pinned)) assert.equal(hash(sourceBytes[path]), digest, path);
  assert.equal(contract.acceptance, 'UNACCEPTED');
  assert.equal(contract.proof_level, 'source-contract-only');
  console.log('INDEPENDENT_SOURCE_RECEIPT ' + JSON.stringify(pinned));
});

test('independent source DTO declarations and all twenty C4 edges', () => {
  checkDTO(contract);
  checkEdges(contract);
});

test('independent complete snapshot Go inventory and bounded checkpoint gaps', () => {
  console.log('INDEPENDENT_COVERAGE_RECEIPT ' + JSON.stringify(checkCoverage(contract)));
});

function negative(name, change, checker, pattern) {
  test(name, () => {
    checker(contract);
    const changed = structuredClone(contract);
    change(changed);
    assert.throws(() => checker(changed), pattern);
  });
}

negative('negative DTO drift: Go JSON tag renamed', candidate => { candidate.api.goInput[0].json = 'issue'; }, checkDTO, /independent Go DTO/);
negative('negative edge drift: fabricated Python-to-Go call', candidate => { candidate.diagrams[1].edges[2].to = 'go'; }, checkEdges, /edge ownership and topology/);
negative('negative witness drift: existing but unrelated source anchor', candidate => { candidate.diagrams[2].edges[3].anchor = 'async def implement_issue('; }, checkEdges, /edge ownership and topology/);
negative('negative coverage drift: ArchSteer Go coverage invented', candidate => { candidate.coverage.archsteerGoFileCount = 200; }, checkCoverage, /Go blind spot/);
negative('negative wire-version drift: tool schema is not DTO version', candidate => { candidate.dto.wireSchemaVersion = arch.schema_version; }, checkDTO, /DTO version proof/);

test('negative source DTO drift is detected without relying on fingerprint mismatch', () => {
  checkDTO(contract);
  const changed = { ...source, [goPath]: source[goPath].replace('`json:"goal"`', '`json:"issue"`') };
  assert.notEqual(changed[goPath], source[goPath]);
  assert.throws(() => checkDTO(contract, changed), /independent Go DTO/);
});

test('negative joint diagram and contract topology drift cannot validate each other', () => {
  checkEdges(contract);
  const changed = structuredClone(contract);
  changed.diagrams[1].edges[2].to = 'go';
  const docs = { ...documents, [diagramPaths[1]]: documents[diagramPaths[1]].replace('python -->|coding loop and verifier calls| agents', 'python -->|coding loop and verifier calls| go') };
  assert.notEqual(docs[diagramPaths[1]], documents[diagramPaths[1]]);
  assert.throws(() => checkEdges(changed, docs), /edge ownership and topology/);
});

test('negative source control-flow drift cannot hide behind an existing call anchor', () => {
  checkEdges(contract);
  const changed = { ...source, [pythonPath]: source[pythonPath].replace('if cfg.verify and coding_ok and commits:', 'if coding_ok:') };
  assert.notEqual(changed[pythonPath], source[pythonPath]);
  assert.throws(() => checkEdges(contract, documents, changed), /conditional verification edge/);
});

test('QA artifacts and producer inputs have reproducible byte receipts', () => {
  const paths = [...Object.keys(pinned), contractPath, ...diagramPaths, 'docs/architecture/basefive-r14-drift-analysis.md', 'docs/architecture/basefive-r14-contract.test.mjs', ownPath, reviewPath];
  const files = paths.map(path => {
    const content = read(path);
    const lines = content.toString('utf8').trimEnd().split(/\r?\n/).length;
    if ([ownPath, reviewPath].includes(path)) assert.ok(lines < 500, 'QA file line limit');
    return { file: path, sha256: hash(content), lines };
  });
  const review = read(reviewPath).toString('utf8');
  assert.match(review, /UNACCEPTED/);
  assert.match(review, /PYTHON_VENV_UNAVAILABLE/);
  assert.match(review, /UNVERIFIED/);
  console.log('INDEPENDENT_ARTIFACT_RECEIPT ' + JSON.stringify({ recorded_at: new Date().toISOString(), files }));
});
