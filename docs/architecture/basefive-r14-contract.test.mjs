import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { createHash } from 'node:crypto';

const root = new URL('../../', import.meta.url);
const paths = ['AGENTS.md', 'docs/acceptance/c4-releasegate-reducer-frontier-20261013.md', '.archsteer/model.json', '.archsteer/architecture.md', 'swe_af/issue/build.py', 'go/internal/orch/build.go', 'graphify-out/graph.json'];
const bytes = Object.fromEntries(paths.map(path => [path, readFileSync(new URL(path, root))]));
const text = Object.fromEntries(paths.map(path => [path, bytes[path].toString('utf8')]));
const hash = value => createHash('sha256').update(value).digest('hex');
const arch = JSON.parse(text['.archsteer/model.json']);
const graph = JSON.parse(text['graphify-out/graph.json']);
const python = text['swe_af/issue/build.py'];
const go = text['go/internal/orch/build.go'];
const unique = values => [...new Set(values)].sort();
const graphPaths = unique(graph.nodes.map(node => node.source_file).filter(value => typeof value === 'string'));
const archPaths = Object.keys(arch.components).sort();
const targets = ['swe_af/issue/build.py', 'go/internal/orch/build.go'];
const coverage = {
  archsteerComponents: archPaths.length,
  graphifyNodes: graph.nodes.length,
  graphifyEdges: graph.edges.length,
  graphifyFiles: graphPaths.length,
  archsteerGoFiles: archPaths.filter(path => path.endsWith('.go')),
  graphifyGoFiles: graphPaths.filter(path => path.endsWith('.go')),
  archsteerTargets: Object.fromEntries(targets.map(path => [path, arch.components[path]?.exported_apis ?? []])),
  graphifyTargets: Object.fromEntries(targets.map(path => [path, graph.nodes.filter(node => node.source_file === path).length])),
  graphifyTargetLabels: Object.fromEntries(targets.map(path => [path, unique(graph.nodes.filter(node => node.source_file === path).map(node => node.label))])),
  graphifyAliasFiles: graphPaths.filter(path => path.startsWith('.graphify-code/')).length,
  archsteerAliasFiles: archPaths.filter(path => path.startsWith('.graphify-code/')).length,
};
const normalize = value => value.replace(/\s+/g, ' ').trim();
const observed = {
  pythonSignature: normalize(python.match(/async def implement_issue\(([\s\S]*?)\) -> dict:/)[0]),
  goSignature: normalize(go.match(/func Build\([^\n]+/)[0]),
  goRegistration: [...go.match(/func Handlers\(\)[\s\S]*?\n\}/)[0].matchAll(/"([^"]+)":\s+(\w+)/g)].map(match => [match[1], match[2]]),
  goInput: [...go.match(/type buildInput struct \{([\s\S]*?)\n\}/)[1].matchAll(/(\w+)\s+([^\n`]+)\s+`json:"([^"]+)"`/g)].map(match => ({ field: match[1], type: match[2].trim(), json: match[3] })),
  pythonResultFields: [...python.match(/return IssueBuildResult\(([\s\S]*?)\)\.model_dump\(\)/)[1].matchAll(/^\s*(\w+)=/gm)].map(match => match[1]),
  goResultFields: [...go.match(/buildResult := schemas.BuildResult\{([\s\S]*?)\n\t\}/)[1].matchAll(/^\s*(\w+):/gm)].map(match => match[1]),
  goCalls: unique([...go.matchAll(/deps\.Call(?:Raw)?\(ctx, "([^"]+)"/g), ...go.matchAll(/deps\.App\.Call\(ctx, deps\.target\("([^"]+)"\)/g)].map(match => match[1])),
};

const currentCoverage = {
  archsteerComponents: coverage.archsteerComponents,
  graphifyNodes: coverage.graphifyNodes,
  graphifyEdges: coverage.graphifyEdges,
  graphifyFiles: coverage.graphifyFiles,
  archsteerGoFileCount: coverage.archsteerGoFiles.length,
  graphifyGoFileCount: coverage.graphifyGoFiles.length,
  graphifyGoDirectories: unique(coverage.graphifyGoFiles.map(path => path.slice(0, path.lastIndexOf('/')))),
  archsteerTargets: coverage.archsteerTargets,
  graphifyTargets: coverage.graphifyTargets,
  pythonMissingGraphSymbols: unique([...python.matchAll(/^async def (\w+)\(/gm)].map(match => match[1]).filter(symbol => !coverage.graphifyTargetLabels['swe_af/issue/build.py'].includes(symbol + '()'))),
  graphifyAliasFiles: coverage.graphifyAliasFiles,
  archsteerAliasFiles: coverage.archsteerAliasFiles,
  repositoryCompleteness: 'UNVERIFIED',
  graphSourceFreshness: 'UNVERIFIED',
  packageRuntimeCoverage: 'UNVERIFIED',
};
const diagramFiles = ['docs/architecture/basefive-r14-C4-context.md', 'docs/architecture/basefive-r14-container.mmd', 'docs/architecture/basefive-r14-component.mmd'];
const diagrams = Object.fromEntries(diagramFiles.map(path => [path, readFileSync(new URL(path, root), 'utf8')]));
const contract = JSON.parse(readFileSync(new URL('basefive-r14-contracts.json', import.meta.url), 'utf8'));
const clone = value => structuredClone(value);

function parseDiagram(document, markdown = false) {
  let body = document;
  if (markdown) {
    const blocks = [...document.matchAll(/```mermaid\s*\n([\s\S]*?)```/g)];
    assert.equal(blocks.length, 1, 'exactly one context Mermaid block');
    body = blocks[0][1];
  }
  const lines = body.trim().split(/\r?\n/).map(line => line.trim()).filter(Boolean);
  assert.equal(lines.shift(), 'flowchart LR', 'restricted Mermaid header');
  const nodes = [];
  const edges = [];
  const identifiers = new Set();
  const groups = new Set();
  let depth = 0;
  for (const line of lines) {
    let match = line.match(/^subgraph ([A-Za-z][\w]*)\["([^"\n]+)"\]$/);
    if (match) {
      assert.ok(!groups.has(match[1]), 'unique subgraph');
      groups.add(match[1]);
      depth += 1;
      continue;
    }
    if (line === 'end') {
      assert.ok(depth > 0, 'balanced subgraph end');
      depth -= 1;
      continue;
    }
    match = line.match(/^([A-Za-z][\w]*)\["([^"\n]+)"\]$/);
    if (match) {
      assert.ok(!identifiers.has(match[1]), 'unique node');
      identifiers.add(match[1]);
      nodes.push({ id: match[1], label: match[2] });
      continue;
    }
    match = line.match(/^([A-Za-z][\w]*) -->\|([^|\n]+)\| ([A-Za-z][\w]*)$/);
    assert.ok(match, 'unsupported Mermaid line: ' + line);
    edges.push({ from: match[1], to: match[3], label: match[2] });
  }
  assert.equal(depth, 0, 'closed subgraphs');
  assert.ok(nodes.length > 0 && edges.length > 0, 'nonempty diagram');
  assert.equal(new Set(edges.map(edge => JSON.stringify(edge))).size, edges.length, 'unique edges');
  for (const identifier of groups) assert.ok(!identifiers.has(identifier), 'subgraph/node collision');
  for (const edge of edges) assert.ok(identifiers.has(edge.from) && identifiers.has(edge.to), 'no dangling edge');
  return { nodes, edges };
}

function validateContract(candidate, sourceBytes = bytes, documents = diagrams) {
  assert.equal(candidate.schema, 'swe-af.basefive-r14-contracts.v1', 'contract schema');
  assert.equal(candidate.project, 'swe-af', 'project');
  assert.equal(candidate.git_receipt, '064a7226a500365005b5b662044fc9014a328722', 'prelaunch identity, not live Git proof');
  assert.equal(candidate.acceptance, 'UNACCEPTED', 'no self-acceptance');
  assert.equal(candidate.proof_level, 'source-contract-only', 'no runtime proof');
  assert.deepEqual(Object.keys(candidate.sources).sort(), [...paths].sort(), 'exact source scope');
  for (const path of paths) {
    assert.match(candidate.sources[path], /^[a-f0-9]{64}$/, 'SHA-256 shape');
    assert.equal(hash(sourceBytes[path]), candidate.sources[path], 'source fingerprint: ' + path);
  }
  assert.deepEqual(candidate.api, observed, 'actual signatures, registration, input tags, result fields and calls');
  assert.deepEqual(candidate.coverage, currentCoverage, 'complete scoped index comparison');
  assert.deepEqual(candidate.dto, {
    pythonInput: 'IssueSpec.model_validate(issue)',
    pythonResult: 'IssueBuildResult(...).model_dump()',
    goInput: 'buildInput JSON tags; marshal/unmarshal errors ignored in Build',
    goResult: 'schemas.BuildResult constructor and conditional workspace_handle',
    wireSchemaVersion: null,
    wireCompatibility: 'UNVERIFIED',
    comparison: 'DIFFERENT_ENTRYPOINTS_NOT_PARITY',
    definitionsOutsideReadScope: ['swe_af/issue/schemas.py', 'go/internal/schemas'],
    runtimeExecution: 'UNVERIFIED',
  }, 'do not invent DTO versions or parity');
  const checkAnchor = item => {
    assert.ok(targets.includes(item.source), 'source-backed code anchor only');
    assert.ok(typeof item.anchor === 'string' && item.anchor.length > 0, 'nonempty anchor');
    assert.ok(sourceBytes[item.source].toString('utf8').includes(item.anchor), 'missing source anchor: ' + item.anchor);
  };
  assert.equal(candidate.invariants.length, 7, 'bounded behavioral invariants');
  candidate.invariants.forEach(checkAnchor);
  assert.deepEqual(candidate.diagrams.map(diagram => diagram.level), ['context', 'container', 'component'], 'three C4 levels');
  assert.deepEqual(candidate.diagrams.map(diagram => diagram.file), diagramFiles, 'exact diagram scope');
  for (const diagram of candidate.diagrams) {
    diagram.nodes.forEach(checkAnchor);
    diagram.edges.forEach(checkAnchor);
    const parsed = parseDiagram(documents[diagram.file], diagram.level === 'context');
    assert.deepEqual(parsed.nodes, diagram.nodes.map(({ id, label }) => ({ id, label })), 'Mermaid nodes equal contract');
    assert.deepEqual(parsed.edges, diagram.edges.map(({ from, to, label }) => ({ from, to, label })), 'Mermaid edges equal contract');
  }
}

test('source fingerprints, actual API/DTO observations, graph coverage and all C4 edges match', () => {
  validateContract(contract);
  console.log('SOURCEBOUND_RECEIPT ' + JSON.stringify({ fingerprints: contract.sources, coverage: { archsteerGoFiles: currentCoverage.archsteerGoFileCount, graphifyGoFiles: currentCoverage.graphifyGoFileCount, graphifyGoDirectories: currentCoverage.graphifyGoDirectories.length, missingPythonGraphSymbols: currentCoverage.pythonMissingGraphSymbols }, diagramCounts: contract.diagrams.map(diagram => ({ level: diagram.level, nodes: diagram.nodes.length, edges: diagram.edges.length })), proof: contract.proof_level }));
});

function rejectsMutation(name, mutate, pattern) {
  test(name, () => {
    validateContract(contract);
    const candidate = clone(contract);
    mutate(candidate);
    assert.throws(() => validateContract(candidate), pattern);
  });
}

rejectsMutation('rejects changed Python signature/default', candidate => { candidate.api.pythonSignature += ' changed'; }, /actual signatures/);
rejectsMutation('rejects renamed Go public signature', candidate => { candidate.api.goSignature = candidate.api.goSignature.replace('Build', 'BuildV2'); }, /actual signatures/);
rejectsMutation('rejects Go registration alias', candidate => { candidate.api.goRegistration[0][0] = 'implement_issue'; }, /actual signatures/);
rejectsMutation('rejects changed Go DTO JSON tag', candidate => { candidate.api.goInput[0].json = 'issue'; }, /actual signatures/);
rejectsMutation('rejects changed Go DTO field type', candidate => { candidate.api.goInput[7].type = 'string'; }, /actual signatures/);
rejectsMutation('rejects missing Python result constructor field', candidate => { candidate.api.pythonResultFields.pop(); }, /actual signatures/);
rejectsMutation('rejects fabricated Go result wire field', candidate => { candidate.api.goResultFields.push('DTO_VERSION'); }, /actual signatures/);
rejectsMutation('rejects omitted Go dispatch target', candidate => { candidate.api.goCalls.pop(); }, /actual signatures/);
rejectsMutation('rejects fabricated ArchSteer Go coverage', candidate => { candidate.coverage.archsteerGoFileCount = 200; }, /index comparison/);
rejectsMutation('rejects hidden Python checkpoint blind spot', candidate => { candidate.coverage.pythonMissingGraphSymbols = []; }, /index comparison/);
rejectsMutation('rejects repository completeness claim', candidate => { candidate.coverage.repositoryCompleteness = 'VERIFIED'; }, /index comparison/);
rejectsMutation('rejects DTO compatibility claim', candidate => { candidate.dto.wireCompatibility = 'VERIFIED'; }, /DTO versions or parity/);
rejectsMutation('rejects invented DTO schema version', candidate => { candidate.dto.wireSchemaVersion = '1.0.0'; }, /DTO versions or parity/);
rejectsMutation('rejects runtime proof claim', candidate => { candidate.proof_level = 'runtime'; }, /runtime proof/);
rejectsMutation('rejects self acceptance', candidate => { candidate.acceptance = 'ACCEPTED_GREEN'; }, /self-acceptance/);
rejectsMutation('rejects unsupported contract version', candidate => { candidate.schema += '.v2'; }, /contract schema/);
rejectsMutation('rejects traversal source reference without reading it', candidate => { candidate.sources['../../private'] = '0'.repeat(64); }, /exact source scope/);
rejectsMutation('rejects diagram path outside write scope', candidate => { candidate.diagrams[0].file = '../outside.mmd'; }, /exact diagram scope/);
rejectsMutation('rejects missing source anchor', candidate => { candidate.diagrams[0].edges[0].anchor = 'nonexistent_runtime_api'; }, /missing source anchor/);
rejectsMutation('rejects fabricated Mermaid edge in contract', candidate => { candidate.diagrams[0].edges[0].to = 'github'; }, /Mermaid edges equal/);
rejectsMutation('rejects source hash drift', candidate => { candidate.sources[targets[0]] = '0'.repeat(64); }, /source fingerprint/);

for (const path of paths) {
  test('rejects changed source bytes for ' + path, () => {
    const changed = { ...bytes, [path]: Buffer.concat([bytes[path], Buffer.from('\nchanged')]) };
    assert.throws(() => validateContract(contract, changed), /source fingerprint/);
  });
}

test('rejects diagram label drift against source-mapped contract', () => {
  const changed = { ...diagrams, [diagramFiles[1]]: diagrams[diagramFiles[1]].replace('Caller', 'Unverified caller') };
  assert.throws(() => validateContract(contract, bytes, changed), /Mermaid nodes equal/);
});

test('rejects dangling Mermaid node', () => {
  assert.throws(() => parseDiagram('flowchart LR\na["A"]\na -->|call| absent'), /dangling edge/);
});

test('rejects duplicate Mermaid node', () => {
  assert.throws(() => parseDiagram('flowchart LR\na["A"]\na["Again"]'), /unique node/);
});

test('rejects unclosed subgraph', () => {
  assert.throws(() => parseDiagram('flowchart LR\nsubgraph group["Group"]\na["A"]\na -->|call| a'), /closed subgraphs/);
});

test('rejects unsupported Mermaid grammar rather than ignoring it', () => {
  assert.throws(() => parseDiagram('flowchart LR\na["A"]\na -.-> b'), /unsupported Mermaid line/);
});

test('rejects extra context Mermaid block', () => {
  assert.throws(() => parseDiagram(diagrams[diagramFiles[0]] + '\n```mermaid\nflowchart LR\n```', true), /exactly one/);
});

test('indexed Go import directories are observed, not package execution proof', () => {
  const imports = [...go.matchAll(/"github\.com\/Agent-Field\/SWE-AF\/(go\/internal\/[^"\n]+)"/g)].map(match => match[1]);
  assert.deepEqual(imports, ['go/internal/config', 'go/internal/envelope', 'go/internal/hitl', 'go/internal/schemas', 'go/internal/workspace']);
  for (const directory of imports) {
    assert.ok(currentCoverage.graphifyGoDirectories.includes(directory));
    assert.ok(!archPaths.some(path => path.startsWith(directory + '/')));
  }
  assert.equal(currentCoverage.packageRuntimeCoverage, 'UNVERIFIED');
});

test('final authorized artifact fingerprints and line limits', () => {
  const files = [...diagramFiles, 'docs/architecture/basefive-r14-contracts.json', 'docs/architecture/basefive-r14-drift-analysis.md', 'docs/architecture/basefive-r14-contract.test.mjs'];
  const receipts = files.map(path => {
    const content = readFileSync(new URL(path, root));
    const lines = content.toString('utf8').trimEnd().split(/\r?\n/).length;
    assert.ok(lines < 500, 'project file line limit: ' + path);
    return { file: path, sha256: hash(content), lines };
  });
  console.log('FINAL_ARTIFACT_RECEIPT ' + JSON.stringify({ recorded_at: new Date().toISOString(), files: receipts }));
});
