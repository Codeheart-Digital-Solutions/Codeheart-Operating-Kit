// Focused contract tests for the managed relay forwarder. Run: node --test tests/coordination-forwarder.test.cjs
// Fixtures use synthetic identities. The fake exec_command runs each generated command string
// through /bin/sh, so quoting is exercised for real; the native send tool is always a fake.
'use strict';

const { test, before, after } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { spawnSync } = require('node:child_process');

const ROOT = path.resolve(__dirname, '..');
const SCRIPT = path.join(ROOT, 'components/agent-interface/managed/scripts/forward-prepared-message.js');
const MIRROR = path.join(ROOT, 'src/codeheart_operating_kit/resources/components/agent-interface/managed/scripts/forward-prepared-message.js');
const CONTRACT = path.join(ROOT, 'components/agent-interface/managed/reference/cross-tool-coordination-contract.md');
const forwarder = require(SCRIPT);

const THREAD = 'thread-synthetic-coordinator';
const HOST = 'local';
const DESTINATION = { threadId: THREAD, hostId: HOST };
const ASSIGNMENT = 'assign-a';
const ATTEMPT = 'attempt-1';
const SENSITIVE = "Reply with 'quotes', \"double\", $(touch PWNED), `touch PWNED`, $HOME; rm -rf x\nline two é中😀";

let helperPath = '';
let helperSkip = '';
const scratch = fs.realpathSync(fs.mkdtempSync(path.join(os.tmpdir(), 'forwarder-test-')));

before(() => {
  const built = path.join(scratch, 'bin', 'codeheart-operating-kit');
  const result = spawnSync('go', ['build', '-o', built, './cmd/codeheart-operating-kit'], { cwd: ROOT, encoding: 'utf8' });
  if (result.error || result.status !== 0) {
    helperSkip = 'Go helper unavailable: ' + (result.error ? result.error.message : result.stderr);
  } else {
    helperPath = built;
  }
});

after(() => {
  fs.rmSync(scratch, { recursive: true, force: true });
});

function envelope(threadId) {
  return { content: [{ type: 'text', text: JSON.stringify({ threadId }) }], isError: false };
}

function errorEnvelope(text) {
  return { content: [{ type: 'text', text }], isError: true };
}

let fixtureCount = 0;

// fixture creates a synthetic attempt below a state root whose path contains shell-sensitive text.
function fixture({ prompt = SENSITIVE, message, delivery } = {}) {
  fixtureCount += 1;
  const base = path.join(scratch, 'case-' + fixtureCount);
  const cwd = path.join(base, 'cwd');
  const stateRoot = path.join(base, "state 'root' $(touch PWNED) `touch PWNED`; x");
  const attemptDir = path.join(stateRoot, 'attempts', ASSIGNMENT, ATTEMPT);
  fs.mkdirSync(cwd, { recursive: true });
  fs.mkdirSync(attemptDir, { recursive: true, mode: 0o700 });
  const messageValue = message || { threadId: THREAD, hostId: HOST, prompt };
  fs.writeFileSync(path.join(attemptDir, 'message.json'), JSON.stringify(messageValue, null, 2) + '\n');
  const deliveryValue = delivery || {
    status: 'pending',
    prepared_recipient: { thread_id: THREAD, host_id: HOST },
    recipient_confirmed: false,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z'
  };
  fs.writeFileSync(path.join(attemptDir, 'delivery.json'), JSON.stringify(deliveryValue, null, 2) + '\n');
  fs.writeFileSync(path.join(attemptDir, 'attempt.json'), '{"status":"response_captured"}\n');
  const outcome = {
    status: 'response_captured',
    assignment_id: ASSIGNMENT,
    attempt_id: ATTEMPT,
    session_id: '00000000-0000-4000-8000-000000000000',
    attempt_dir: attemptDir,
    attempt_file: path.join(attemptDir, 'attempt.json'),
    message_file: path.join(attemptDir, 'message.json'),
    delivery_file: path.join(attemptDir, 'delivery.json'),
    permission_denials: 0
  };
  return { base, cwd, stateRoot, attemptDir, outcome, messageValue };
}

function input(f, overrides = {}) {
  return {
    outcome: JSON.stringify(f.outcome),
    assignmentId: ASSIGNMENT,
    attemptId: ATTEMPT,
    destination: { ...DESTINATION },
    stateRoot: f.stateRoot,
    helperPath,
    ...overrides
  };
}

// host emulates the Codex tool namespace: exec_command runs the command string in /bin/sh and
// resolves {output, exit_code}; the native send tool is a fake that records each call.
function host(f, { send, exec } = {}) {
  const calls = { sends: [], commands: [] };
  const tools = {
    exec_command: async (params) => {
      calls.commands.push(params.cmd);
      if (exec) {
        const replaced = await exec(params, calls);
        if (replaced) {
          return replaced;
        }
      }
      const run = spawnSync('/bin/sh', ['-c', params.cmd], { cwd: f.cwd, encoding: 'utf8' });
      return { output: run.stdout + run.stderr, exit_code: run.status, wall_time_seconds: 0 };
    },
    mcp__codex_app__send_message_to_thread: async (args) => {
      calls.sends.push(args);
      return send ? send(args, calls) : envelope(args.threadId);
    }
  };
  return { tools, calls };
}

function readJSON(file) {
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}

function assertNoInjection(f) {
  for (const dir of [f.cwd, f.base, f.stateRoot, ROOT]) {
    assert.equal(fs.existsSync(path.join(dir, 'PWNED')), false, 'shell-sensitive text executed in ' + dir);
  }
}

function needsHelper(t) {
  if (helperSkip && process.env.CI) {
    assert.fail(helperSkip);
  }
  if (helperSkip) {
    t.skip(helperSkip);
    return true;
  }
  return false;
}

test('the managed script is ASCII-only and its packaged mirror is byte-identical', () => {
  const source = fs.readFileSync(SCRIPT);
  assert.ok(source.every((byte) => byte < 0x80), 'non-ASCII byte in script');
  assert.deepEqual(fs.readFileSync(MIRROR), source);
});

test('success forwards the parsed file object unchanged and records sent with the receipt', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'delivered', JSON.stringify(result));
  assert.equal(result.delivery_status, 'sent');
  assert.equal(result.recorded, true);
  assert.equal(result.next_owner, 'coordinator');
  assert.equal(calls.sends.length, 1);
  assert.deepEqual(Object.keys(calls.sends[0]).sort(), ['hostId', 'prompt', 'threadId']);
  assert.deepEqual(calls.sends[0], f.messageValue);
  assert.equal(calls.sends[0].prompt, SENSITIVE);
  const delivery = readJSON(path.join(f.attemptDir, 'delivery.json'));
  assert.equal(delivery.status, 'sent');
  assert.equal(delivery.recipient_confirmed, true);
  assert.deepEqual(JSON.parse(delivery.receipt), envelope(THREAD));
  assert.equal(delivery.receipt_file, path.join(f.attemptDir, 'delivery-claim', 'receipt.json'));
  assert.deepEqual(readJSON(delivery.receipt_file), envelope(THREAD));
  assert.ok(!JSON.stringify(result).includes('Reply with'), 'the compact result never carries the message');
  assertNoInjection(f);
});

test('the known argument rejection is recorded and one compact failure notice follows', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const rejection = errorEnvelope('send_message_to_thread received invalid arguments.');
  const { tools, calls } = host(f, { send: (args, c) => (c.sends.length === 1 ? rejection : envelope(args.threadId)) });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'rejected', JSON.stringify(result));
  assert.equal(result.delivery_status, 'rejected');
  assert.equal(calls.sends.length, 2);
  const notice = calls.sends[1];
  assert.equal(notice.threadId, THREAD);
  assert.equal(notice.hostId, HOST);
  assert.ok(notice.prompt.includes('rejected by the native send tool before delivery'));
  assert.ok(!notice.prompt.includes('Reply with'), 'the original reply must not be resent');
  assert.equal(result.failure_notice.status, 'sent');
  assert.equal(result.next_owner, 'coordinator');
  const delivery = readJSON(path.join(f.attemptDir, 'delivery.json'));
  assert.equal(delivery.status, 'rejected');
  assert.equal(delivery.reported_recipient, undefined, 'a rejection reports no recipient');
  assert.deepEqual(JSON.parse(delivery.receipt), rejection);
  assert.deepEqual(readJSON(path.join(f.attemptDir, 'delivery-claim', 'failure-notice-receipt.json')), envelope(THREAD));
});

test('a failed failure notice returns a blocked handoff to the relay without further sends', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const rejection = errorEnvelope('send_message_to_thread received invalid arguments.');
  const { tools, calls } = host(f, { send: () => rejection });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'rejected');
  assert.equal(calls.sends.length, 2);
  assert.equal(result.failure_notice.status, 'rejected');
  assert.equal(result.next_owner, 'relay');
});

for (const [name, receipt, reason] of [
  ['other error envelope (for example a permission refusal)', errorEnvelope('User denied the tool call.'), 'error_envelope'],
  ['unknown receipt', { ok: true }, 'unrecognized_receipt'],
  ['success envelope without recipient', { content: [{ type: 'text', text: '{}' }], isError: false }, 'unrecognized_receipt'],
  ['envelope with two text items', { content: [{ type: 'text', text: '{"threadId":"thread-synthetic-coordinator"}' }, { type: 'text', text: 'x' }], isError: false }, 'unrecognized_receipt'],
  ['undefined result', undefined, 'unrecognized_receipt']
]) {
  test('uncertain: ' + name + ' is never success and sends nothing further', async (t) => {
    if (needsHelper(t)) return;
    const f = fixture();
    const { tools, calls } = host(f, { send: () => receipt });
    const result = await forwarder.forwardWithCodexTools(tools, input(f));
    assert.equal(result.outcome, 'uncertain', JSON.stringify(result));
    assert.equal(result.reason, reason);
    assert.equal(result.failure_notice, null);
    assert.equal(calls.sends.length, 1);
    assert.equal(result.next_owner, 'relay');
    assert.equal(readJSON(path.join(f.attemptDir, 'delivery.json')).status, 'uncertain');
  });
}

test('a mismatched receipt recipient is recorded as uncertain with the mismatch', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const { tools, calls } = host(f, { send: () => envelope('thread-synthetic-other') });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'uncertain');
  assert.equal(result.reason, 'recipient_mismatch');
  assert.equal(calls.sends.length, 1);
  const delivery = readJSON(path.join(f.attemptDir, 'delivery.json'));
  assert.equal(delivery.status, 'uncertain');
  assert.equal(delivery.recipient_mismatch, true);
  assert.equal(delivery.reported_recipient.thread_id, 'thread-synthetic-other');
});

test('a throwing send (for example a permission refusal) is uncertain with no alternative send', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const { tools, calls } = host(f, { send: () => { throw new Error('permission denied for this tool'); } });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'uncertain');
  assert.equal(result.reason, 'send_threw');
  assert.equal(calls.sends.length, 1);
  assert.equal(result.failure_notice, null);
  const receipt = readJSON(path.join(f.attemptDir, 'delivery-claim', 'receipt.json'));
  assert.equal(receipt.forwarder_send_exception.message, 'permission denied for this tool');
});

test('the incident shape (thread instead of threadId) is refused before any send', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture({ message: { thread: THREAD, hostId: HOST, prompt: 'x' } });
  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'blocked');
  assert.equal(result.reason, 'invalid_message');
  assert.equal(result.send_attempted, false);
  assert.equal(calls.sends.length, 0);
  assert.equal(fs.existsSync(path.join(f.attemptDir, 'delivery-claim')), false);
});

test('invalid inputs, unsafe paths and binding mismatches block before any command or send', async () => {
  const f = fixture();
  const cases = [
    [{ stateRoot: 'relative/state' }, 'invalid_input'],
    [{ stateRoot: f.stateRoot + '/../other' }, 'invalid_input'],
    [{ stateRoot: f.stateRoot + '/' }, 'invalid_input'],
    [{ helperPath: '/bin/helper\nrm' }, 'invalid_input'],
    [{ attemptId: '../attempt-1' }, 'invalid_input'],
    [{ destination: { threadId: THREAD, hostId: HOST, extra: 'x' } }, 'invalid_input'],
    [{ destination: { threadId: 'thread\u0007', hostId: HOST } }, 'invalid_input'],
    [{ extra: true }, 'invalid_input'],
    [{ outcome: 'helper failed to start: no such file' }, 'invalid_outcome'],
    [{ attemptId: 'attempt-2' }, 'binding_mismatch'],
    [{ outcome: JSON.stringify({ ...f.outcome, message_file: path.join(f.stateRoot, 'attempts', ASSIGNMENT, 'attempt-0', 'message.json') }) }, 'binding_mismatch'],
    [{ outcome: JSON.stringify({ ...f.outcome, message_file: undefined, delivery_file: undefined }) }, 'no_message'],
    [{ outcome: JSON.stringify({ ...f.outcome, fallback_message: { threadId: THREAD, hostId: HOST, prompt: 'x' } }) }, 'no_message']
  ];
  for (const [overrides, reason] of cases) {
    const { tools, calls } = host(f);
    const result = await forwarder.forwardWithCodexTools(tools, input(f, { helperPath: '/bin/helper', ...overrides }));
    assert.equal(result.outcome, 'blocked', JSON.stringify(overrides));
    assert.equal(result.reason, reason, JSON.stringify(overrides));
    assert.equal(result.next_owner, 'relay');
    assert.equal(calls.commands.length, 0, JSON.stringify(overrides));
    assert.equal(calls.sends.length, 0);
  }
});

test('a destination that differs from the prepared message or record blocks before sending', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f, { destination: { threadId: 'thread-synthetic-other', hostId: HOST } }));
  assert.equal(result.outcome, 'blocked');
  assert.equal(result.reason, 'recipient_mismatch');
  assert.equal(calls.sends.length, 0);

  const g = fixture({ message: { threadId: 'thread-synthetic-other', hostId: HOST, prompt: 'x' } });
  const second = host(g);
  const mismatch = await forwarder.forwardWithCodexTools(second.tools, input(g));
  assert.equal(mismatch.reason, 'invalid_message');
  assert.equal(second.calls.sends.length, 0);
});

test('missing host capabilities and an unusable helper block before sending', async () => {
  const f = fixture();
  const missing = await forwarder.forwardWithCodexTools({ exec_command: async () => ({}) }, input(f));
  assert.equal(missing.outcome, 'blocked');
  assert.equal(missing.reason, 'missing_capability');

  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f, { helperPath: path.join(f.base, 'missing-helper') }));
  assert.equal(result.reason, 'helper_unavailable');
  assert.equal(calls.sends.length, 0);
  assert.equal(fs.existsSync(path.join(f.attemptDir, 'delivery-claim')), false);
});

test('truncated or altered shell output is detected before sending', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture({ prompt: 'A'.repeat(3000) });
  const { tools, calls } = host(f, {
    exec: (params) => {
      if (!params.cmd.includes('message.json')) {
        return null;
      }
      const run = spawnSync('/bin/sh', ['-c', params.cmd], { cwd: f.cwd, encoding: 'utf8' });
      return { output: run.stdout.replace('A'.repeat(500), '...truncated...'), exit_code: run.status };
    }
  });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'blocked');
  assert.equal(result.reason, 'message_unreadable');
  assert.equal(calls.sends.length, 0);
});

test('a long-reply locator prompt and multi-byte text are forwarded unchanged', async (t) => {
  if (needsHelper(t)) return;
  const prompt = 'Claude CLI attempt attempt-1 for assignment assign-a: response_captured.\n' +
    'Original reply (5000 characters): /synthetic/state/attempts/assign-a/attempt-1/stdout.jsonl, line 7, JSON field "result".\n' +
    'éè中文😀'.repeat(400);
  const f = fixture({ prompt });
  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'delivered', JSON.stringify(result));
  assert.equal(calls.sends[0].prompt, prompt);
});

test('failed receipt persistence leaves the record pending, keeps the claim and refuses re-entry', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const failing = host(f, { exec: (params) => (params.cmd.includes('set -C') ? { output: 'disk full', exit_code: 1 } : null) });
  const result = await forwarder.forwardWithCodexTools(failing.tools, input(f));
  assert.equal(result.outcome, 'recording_failed', JSON.stringify(result));
  assert.equal(result.derived_status, 'sent');
  assert.equal(result.recorded, false);
  assert.equal(failing.calls.sends.length, 1);
  assert.equal(readJSON(path.join(f.attemptDir, 'delivery.json')).status, 'pending');
  assert.ok(fs.existsSync(path.join(f.attemptDir, 'delivery-claim')));

  const again = host(f);
  const retry = await forwarder.forwardWithCodexTools(again.tools, input(f));
  assert.equal(retry.outcome, 'refused');
  assert.equal(retry.reason, 'claim_unavailable');
  assert.equal(again.calls.sends.length, 0);
});

test('a failed record-delivery is reported as recording_failed, never as delivered', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const { tools, calls } = host(f, {
    exec: (params) => (params.cmd.includes("'--status'") ? { output: '{"status":"not_recorded"}', exit_code: 1 } : null)
  });
  const result = await forwarder.forwardWithCodexTools(tools, input(f));
  assert.equal(result.outcome, 'recording_failed');
  assert.equal(result.next_owner, 'relay');
  assert.equal(calls.sends.length, 1);
  assert.ok(fs.existsSync(path.join(f.attemptDir, 'delivery-claim', 'receipt.json')));
  assert.equal(readJSON(path.join(f.attemptDir, 'delivery.json')).status, 'pending');
});

test('an interruption after sending but before recording makes re-entry refuse', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const interrupted = host(f, {
    exec: (params, calls) => {
      if (calls.sends.length > 0) {
        throw new Error('host execution interrupted');
      }
      return null;
    }
  });
  const first = await forwarder.forwardWithCodexTools(interrupted.tools, input(f));
  assert.equal(first.outcome, 'recording_failed');
  assert.equal(first.send_attempted, true);
  assert.equal(interrupted.calls.sends.length, 1);

  for (let i = 0; i < 2; i++) {
    const again = host(f);
    const retry = await forwarder.forwardWithCodexTools(again.tools, input(f));
    assert.equal(retry.outcome, 'refused');
    assert.equal(retry.reason, 'claim_unavailable');
    assert.equal(again.calls.sends.length, 0);
  }
  assert.equal(readJSON(path.join(f.attemptDir, 'delivery.json')).status, 'pending');
});

test('concurrent and repeated entry produce exactly one send', async (t) => {
  if (needsHelper(t)) return;
  const f = fixture();
  const shared = host(f);
  const results = await Promise.all([1, 2, 3].map(() => forwarder.forwardWithCodexTools(shared.tools, input(f))));
  assert.equal(shared.calls.sends.length, 1);
  assert.equal(results.filter((r) => r.outcome === 'delivered').length, 1);
  for (const refused of results.filter((r) => r.outcome !== 'delivered')) {
    assert.equal(refused.outcome, 'refused');
    assert.equal(refused.send_attempted, false);
  }
  const later = host(f);
  const repeated = await forwarder.forwardWithCodexTools(later.tools, input(f));
  assert.equal(repeated.outcome, 'refused');
  assert.equal(repeated.reason, 'already_recorded');
  assert.equal(repeated.delivery_status, 'sent');
  assert.equal(later.calls.sends.length, 0);
});

function attemptExistsOutcome(f, fallback) {
  return {
    status: 'attempt_exists',
    detail: 'this attempt was already started',
    assignment_id: ASSIGNMENT,
    attempt_id: ATTEMPT,
    attempt_dir: f.attemptDir,
    attempt_file: path.join(f.attemptDir, 'attempt.json'),
    fallback_message: fallback,
    permission_denials: 0
  };
}

function snapshot(dir) {
  return fs.readdirSync(dir).sort().map((name) => [name, fs.readFileSync(path.join(dir, name))]);
}

function assertFallbackReceipt(f, result, expected) {
  const sidecarRoot = path.join(f.stateRoot, 'forwarder-receipts', ASSIGNMENT);
  assert.ok(result.receipt_file.startsWith(sidecarRoot + path.sep), result.receipt_file);
  assert.ok(!result.receipt_file.includes(path.sep + 'attempts' + path.sep));
  assert.deepEqual(readJSON(result.receipt_file), expected);
  assert.equal(fs.statSync(result.receipt_file).mode & 0o077, 0, 'the sidecar is private');
  assert.equal(result.receipt_problem, undefined);
}

test('an attempt_exists fallback is sent once, keeps its receipt in a sidecar and leaves the old attempt untouched', async () => {
  const f = fixture();
  const before = snapshot(f.attemptDir);
  const fallback = { threadId: THREAD, hostId: HOST, prompt: 'Claude CLI attempt attempt-1 for assignment assign-a: attempt_exists.\n' + SENSITIVE };
  const { tools, calls } = host(f);
  const result = await forwarder.forwardWithCodexTools(tools, input(f, { outcome: attemptExistsOutcome(f, fallback), helperPath: '/bin/unused-helper' }));
  assert.equal(result.outcome, 'delivered', JSON.stringify(result));
  assert.equal(result.once_only, 'invocation');
  assert.equal(result.recorded, false);
  assert.equal(calls.sends.length, 1);
  assert.equal(calls.sends[0], fallback, 'the parsed fallback object itself is forwarded');
  assertFallbackReceipt(f, result, envelope(THREAD));
  assert.ok(!JSON.stringify(result).includes('Reply with'), 'the compact result never carries the message');
  assert.deepEqual(snapshot(f.attemptDir), before, 'the earlier attempt is unchanged');
  assertNoInjection(f);

  // A second invocation is a separate send with its own sidecar; nothing is overwritten.
  const again = host(f);
  const second = await forwarder.forwardWithCodexTools(again.tools, input(f, { outcome: attemptExistsOutcome(f, fallback), helperPath: '/bin/unused-helper' }));
  assert.notEqual(second.receipt_file, result.receipt_file);
  assert.deepEqual(readJSON(result.receipt_file), envelope(THREAD));
  assert.deepEqual(snapshot(f.attemptDir), before);
});

test('fallback rejection, unknown and thrown receipts are preserved privately without a second notice', async () => {
  const f = fixture();
  const before = snapshot(f.attemptDir);
  const invalidRequest = { status: 'invalid_request', problems: ['brief_file must be an existing regular file'], fallback_message: { threadId: THREAD, hostId: HOST, prompt: 'rejected before launch' }, permission_denials: 0 };
  const rejection = errorEnvelope('send_message_to_thread received invalid arguments.');
  const second = host(f, { send: () => rejection });
  const rejected = await forwarder.forwardWithCodexTools(second.tools, input(f, { outcome: invalidRequest }));
  assert.equal(rejected.outcome, 'rejected');
  assert.equal(second.calls.sends.length, 1, 'a fallback is never followed by a second notice');
  assert.equal(rejected.next_owner, 'relay');
  assertFallbackReceipt(f, rejected, rejection);

  const unknown = { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD, status: 'failed' }) }], isError: false };
  const third = host(f, { send: () => unknown });
  const uncertain = await forwarder.forwardWithCodexTools(third.tools, input(f, { outcome: invalidRequest }));
  assert.equal(uncertain.outcome, 'uncertain');
  assert.equal(third.calls.sends.length, 1);
  assertFallbackReceipt(f, uncertain, unknown);

  const fourth = host(f, { send: () => { throw new Error('approval declined'); } });
  const thrown = await forwarder.forwardWithCodexTools(fourth.tools, input(f, { outcome: invalidRequest }));
  assert.equal(thrown.outcome, 'uncertain');
  assert.equal(thrown.reason, 'send_threw');
  assert.equal(fourth.calls.sends.length, 1);
  assertFallbackReceipt(f, thrown, { forwarder_send_exception: { name: 'Error', message: 'approval declined' } });
  assert.deepEqual(snapshot(f.attemptDir), before);
});

test('a fallback receipt that cannot be saved is reported truthfully without another send', async () => {
  const f = fixture();
  const before = snapshot(f.attemptDir);
  const fallback = { threadId: THREAD, hostId: HOST, prompt: 'attempt_exists notice' };
  for (const failing of ['mkdir -p --', 'set -C;']) {
    const { tools, calls } = host(f, { exec: (params) => (params.cmd.includes(failing) ? { output: 'denied', exit_code: 1 } : null) });
    const result = await forwarder.forwardWithCodexTools(tools, input(f, { outcome: attemptExistsOutcome(f, fallback), helperPath: '/bin/unused-helper' }));
    assert.equal(result.outcome, 'recording_failed', failing);
    assert.equal(result.derived_status, 'sent');
    assert.equal(result.receipt_file, undefined);
    assert.match(result.receipt_problem, /^receipt not persisted: /);
    assert.equal(result.next_owner, 'relay');
    assert.equal(calls.sends.length, 1);
  }
  const unavailable = host(f, { exec: () => { throw new Error('shell refused'); } });
  const refused = await forwarder.forwardWithCodexTools(unavailable.tools, input(f, { outcome: attemptExistsOutcome(f, fallback), helperPath: '/bin/unused-helper' }));
  assert.equal(refused.outcome, 'recording_failed');
  assert.equal(refused.derived_status, 'sent');
  assert.equal(refused.receipt_problem, 'receipt not persisted: shell refused');
  assert.equal(refused.next_owner, 'relay');
  assert.equal(unavailable.calls.sends.length, 1);

  // A rejected fallback whose evidence is lost keeps its derived status and still gets no notice.
  const lost = host(f, {
    send: () => errorEnvelope('send_message_to_thread received invalid arguments.'),
    exec: (params) => (params.cmd.includes('set -C;') ? { output: 'denied', exit_code: 1 } : null)
  });
  const lostRejection = await forwarder.forwardWithCodexTools(lost.tools, input(f, { outcome: attemptExistsOutcome(f, fallback), helperPath: '/bin/unused-helper' }));
  assert.equal(lostRejection.outcome, 'recording_failed');
  assert.equal(lostRejection.derived_status, 'rejected');
  assert.equal(lostRejection.failure_notice, null);
  assert.equal(lost.calls.sends.length, 1);
  assert.deepEqual(snapshot(f.attemptDir), before);
});

test('receipt classification accepts only the qualified minimal success for the expected recipient', () => {
  const classify = forwarder.classifyReceipt;
  const wrapped = (payload, extra = {}) => ({ content: [{ type: 'text', text: JSON.stringify(payload) }], isError: false, ...extra });
  // Qualified success: the observed native envelope, or its exact payload.
  assert.deepEqual(classify(envelope(THREAD), THREAD), { status: 'sent', shape: 'mcp_envelope', reason: 'accepted', reportedThreadId: THREAD });
  assert.equal(classify(wrapped({ threadId: THREAD }), THREAD).status, 'sent');
  assert.equal(classify({ threadId: THREAD }, THREAD).status, 'sent');
  assert.equal(classify(JSON.stringify({ threadId: THREAD }), THREAD).status, 'sent');
  assert.equal(classify(errorEnvelope('send_message_to_thread received invalid arguments.'), THREAD).status, 'rejected');
  assert.equal(classify(errorEnvelope('  send_message_to_thread received invalid arguments  '), THREAD).status, 'rejected');
  const uncertain = [
    // Reviewer reproductions: conflicting inner or outer error/status fields beside the expected threadId.
    wrapped({ threadId: THREAD, isError: true }),
    wrapped({ threadId: THREAD, status: 'rejected' }),
    wrapped({ threadId: THREAD, status: 'failed' }),
    wrapped({ threadId: THREAD, error: 'x' }),
    wrapped({ threadId: THREAD }, { errors: ['late failure'] }),
    wrapped({ threadId: THREAD }, { error: 'late failure' }),
    wrapped({ threadId: THREAD }, { status: 'failed' }),
    wrapped({ threadId: THREAD }, { structuredContent: { threadId: THREAD } }),
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD }), annotations: {} }], isError: false },
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD }) }, { type: 'text', text: 'x' }], isError: false },
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD }) }], isError: 'false' },
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD }) }] },
    { ...envelope(THREAD), isError: true },
    { content: [], isError: false },
    { threadId: THREAD, status: 'failed' },
    { threadId: THREAD, isError: false },
    JSON.stringify({ threadId: THREAD, status: 'rejected' }),
    null,
    'sent',
    42
  ];
  for (const receipt of uncertain) {
    assert.equal(classify(receipt, THREAD).status, 'uncertain', JSON.stringify(receipt));
  }
  // A permission or security error never qualifies as the definite rejection that triggers a notice.
  for (const receipt of [
    errorEnvelope('Permission denied: send_message_to_thread received invalid arguments.'),
    errorEnvelope('send_message_to_thread received invalid arguments. Approval required.'),
    errorEnvelope('User denied the tool call.'),
    { ...errorEnvelope('send_message_to_thread received invalid arguments.'), errors: ['security policy'] },
    { content: [{ type: 'text', text: 'send_message_to_thread received invalid arguments.' }, { type: 'text', text: 'blocked by policy' }], isError: true }
  ]) {
    assert.equal(classify(receipt, THREAD).status, 'uncertain', JSON.stringify(receipt));
  }
});

test('conflicting success wrappers are recorded uncertain and never trigger a notice', async (t) => {
  if (needsHelper(t)) return;
  for (const receipt of [
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD, isError: true }) }], isError: false },
    { content: [{ type: 'text', text: JSON.stringify({ threadId: THREAD, status: 'rejected' }) }], isError: false },
    { ...envelope(THREAD), errors: ['late failure'] },
    { ...errorEnvelope('send_message_to_thread received invalid arguments.'), errors: ['security policy'] }
  ]) {
    const f = fixture();
    const { tools, calls } = host(f, { send: () => receipt });
    const result = await forwarder.forwardWithCodexTools(tools, input(f));
    assert.equal(result.outcome, 'uncertain', JSON.stringify(receipt));
    assert.equal(result.failure_notice, null);
    assert.equal(calls.sends.length, 1);
    const delivery = readJSON(path.join(f.attemptDir, 'delivery.json'));
    assert.equal(delivery.status, 'uncertain');
    assert.equal(delivery.recipient_confirmed, false);
    assert.deepEqual(JSON.parse(delivery.receipt), receipt);
  }
});

test('the documented loader reads, verifies and runs the installed script', async (t) => {
  if (needsHelper(t)) return;
  const contract = fs.readFileSync(CONTRACT, 'utf8');
  const match = contract.match(/<!-- forwarder-loader -->\n```js\n([\s\S]*?)```/);
  assert.ok(match, 'loader block missing from the contract');
  const f = fixture();
  const placeholder = '"<repository>/.codeheart/kit/docs/agent-interface/scripts/forward-prepared-message.js"';
  assert.ok(match[1].includes(placeholder));
  const value = input(f);
  value.outcome = f.outcome;
  let code = match[1].replace(placeholder, JSON.stringify(SCRIPT));
  code = code.replace(/const INPUT = \{[\s\S]*?\n\};/, 'const INPUT = ' + JSON.stringify(value) + ';');
  const { tools, calls } = host(f);
  const printed = [];
  const AsyncFunction = (async () => {}).constructor;
  await new AsyncFunction('tools', 'text', code)(tools, (value) => printed.push(value));
  assert.equal(printed.length, 1);
  const result = JSON.parse(printed[0]);
  assert.equal(result.outcome, 'delivered', printed[0]);
  assert.equal(calls.sends.length, 1);

  const truncated = host(f, { exec: (params) => (params.cmd.includes('wc -c') ? { output: '999\nvar x;', exit_code: 0 } : null) });
  await assert.rejects(new AsyncFunction('tools', 'text', code)(truncated.tools, () => {}), /unreadable or truncated/);
  assert.equal(truncated.calls.sends.length, 0);
});
