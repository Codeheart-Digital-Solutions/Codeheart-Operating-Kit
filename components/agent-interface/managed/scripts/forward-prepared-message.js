// Codeheart Operating Kit managed workflow: forward one prepared coordination message.
//
// Owner: agent-interface cross-tool coordination. Contract:
// ../reference/cross-tool-coordination-contract.md ("Deterministic Forwarding").
// A relay loads this file through the documented loader inside the Codex tool-execution
// runtime, which has no filesystem or module imports. Host capabilities are injected: one
// native send callback and one shell callback. The workflow validates the helper outcome
// against the authorized assignment, attempt and destination, claims a recorded attempt
// before sending, passes the parsed message object unchanged to the native send tool,
// preserves the native receipt and records a receipt-derived status with the existing
// `coordination record-delivery` helper command. It never resumes Claude, rewrites or
// summarizes a message, retries, replays, or switches tool or recipient. It makes no
// exactly-once claim. Keep this file ASCII-only; the loader checks its byte length.
'use strict';

var CodeheartRelayForwarder = (function () {
  var IDENTIFIER = /^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$/;
  var REJECTED_ARGUMENTS = /^send_message_to_thread received invalid arguments\.?$/;
  var MESSAGE_KEYS = ['hostId', 'prompt', 'threadId'];
  var INPUT_KEYS = ['assignmentId', 'attemptId', 'destination', 'helperPath', 'outcome', 'stateRoot'];
  var CLAIM_DIR = 'delivery-claim';
  var RECEIPT_FILE = 'receipt.json';
  var NOTICE_RECEIPT_FILE = 'failure-notice-receipt.json';
  var FALLBACK_RECEIPT_DIR = 'forwarder-receipts';
  var SEND_TOOL = 'mcp__codex_app__send_message_to_thread';
  var SHELL_TOOL = 'exec_command';
  var TRANSPORT_ONLY = 'Transport status only; not task acceptance.';
  var RELAY_STEP = 'Report this result unchanged in your own reply. Do not resend, retry, edit ' +
    'records, or switch tool or recipient; the coordinator reconciles the known attempt.';
  // Reads one file between two markers with its byte count, so truncated or altered shell
  // output is detected instead of forwarded.
  var READ_SCRIPT = 'n=$(wc -c < "$1") || exit 3; printf \'%s %s\\n\' "$2" "$n"; ' +
    'cat -- "$1" || exit 4; printf \'\\n%s\\n\' "$2"';
  // Writes one new private file; noclobber refuses an existing file.
  var WRITE_SCRIPT = 'umask 077; set -C; printf \'%s\\n\' "$1" > "$2"';
  // Creates one private directory and its missing parents.
  var MKDIR_SCRIPT = 'umask 077; mkdir -p -- "$1"';

  function Stop(outcome, reason, detail) {
    this.outcome = outcome;
    this.reason = reason;
    this.detail = detail;
  }

  function blocked(reason, detail) {
    return new Stop('blocked', reason, detail);
  }

  function isPlainObject(value) {
    return value !== null && typeof value === 'object' && !Array.isArray(value);
  }

  function plainValue(value) {
    if (typeof value !== 'string' || value.trim() === '' || value.length > 512) {
      return false;
    }
    for (var i = 0; i < value.length; i++) {
      var code = value.charCodeAt(i);
      if (code < 0x20 || code === 0x7f) {
        return false;
      }
    }
    return true;
  }

  // cleanAbsolutePath accepts only a normalized absolute POSIX path, so derived attempt
  // paths compare exactly with the helper's own output.
  function cleanAbsolutePath(value) {
    if (typeof value !== 'string' || value.length < 2 || value.length > 4096 || value[0] !== '/') {
      return false;
    }
    for (var i = 0; i < value.length; i++) {
      var code = value.charCodeAt(i);
      if (code < 0x20 || code === 0x7f) {
        return false;
      }
    }
    var segments = value.slice(1).split('/');
    for (var j = 0; j < segments.length; j++) {
      if (segments[j] === '' || segments[j] === '.' || segments[j] === '..') {
        return false;
      }
    }
    return true;
  }

  function sortedKeys(value) {
    return Object.keys(value).sort();
  }

  function sameKeys(value, expected) {
    var keys = sortedKeys(value);
    return keys.length === expected.length && keys.every(function (key, index) {
      return key === expected[index];
    });
  }

  // shellQuote makes one argument literal shell data. Text such as $(...), backquotes,
  // quotes and newlines is never interpreted.
  function shellQuote(argument) {
    if (typeof argument !== 'string' || argument.indexOf('\u0000') !== -1) {
      throw new Error('shell arguments must be strings without NUL characters');
    }
    return "'" + argument.split("'").join("'\\''") + "'";
  }

  function shellCommand(argv) {
    if (!Array.isArray(argv) || argv.length === 0) {
      throw new Error('a command needs at least one argument');
    }
    return argv.map(shellQuote).join(' ');
  }

  function utf8Length(text) {
    var bytes = 0;
    for (var i = 0; i < text.length; i++) {
      var code = text.charCodeAt(i);
      if (code < 0x80) {
        bytes += 1;
      } else if (code < 0x800) {
        bytes += 2;
      } else if (code >= 0xd800 && code <= 0xdbff && i + 1 < text.length &&
        text.charCodeAt(i + 1) >= 0xdc00 && text.charCodeAt(i + 1) <= 0xdfff) {
        bytes += 4;
        i++;
      } else {
        bytes += 3;
      }
    }
    return bytes;
  }

  function frameMarker() {
    return 'CODEHEART-FORWARD-' + Math.random().toString(36).slice(2) + Date.now().toString(36);
  }

  async function readFile(run, path) {
    var marker = frameMarker();
    var result = await run(['/bin/sh', '-c', READ_SCRIPT, 'sh', path, marker]);
    if (result.exitCode !== 0) {
      throw new Error('read exited ' + result.exitCode);
    }
    var output = result.output;
    var head = output.indexOf(marker + ' ');
    var lineEnd = head < 0 ? -1 : output.indexOf('\n', head);
    var tail = output.lastIndexOf('\n' + marker);
    if (head < 0 || lineEnd < 0 || tail < lineEnd) {
      throw new Error('read output is not framed');
    }
    var size = Number(output.slice(head + marker.length + 1, lineEnd).trim());
    var content = output.slice(lineEnd + 1, tail);
    if (!Number.isInteger(size) || utf8Length(content) !== size) {
      throw new Error('read output is truncated or altered');
    }
    return content;
  }

  async function writeNewFile(run, path, data) {
    var result = await run(['/bin/sh', '-c', WRITE_SCRIPT, 'sh', data, path]);
    if (result.exitCode !== 0) {
      throw new Error('write exited ' + result.exitCode);
    }
  }

  function parseJSONObject(text) {
    var value;
    try {
      value = JSON.parse(text);
    } catch (error) {
      return null;
    }
    return isPlainObject(value) ? value : null;
  }

  // validateMessage accepts exactly the helper's three native send arguments for the
  // authorized destination.
  function validateMessage(message, destination) {
    if (!isPlainObject(message) || !sameKeys(message, MESSAGE_KEYS)) {
      return 'the message must contain exactly threadId, hostId and prompt';
    }
    if (!plainValue(message.threadId) || !plainValue(message.hostId) ||
      typeof message.prompt !== 'string' || message.prompt.trim() === '') {
      return 'threadId, hostId and prompt must be non-empty strings';
    }
    if (message.threadId !== destination.threadId || message.hostId !== destination.hostId) {
      return 'the message recipient differs from the authorized destination';
    }
    return '';
  }

  function validateInput(input) {
    if (!isPlainObject(input)) {
      throw blocked('invalid_input', 'the input must be an object');
    }
    var unknown = Object.keys(input).filter(function (key) {
      return INPUT_KEYS.indexOf(key) === -1;
    });
    if (unknown.length > 0) {
      throw blocked('invalid_input', 'unknown input fields: ' + unknown.sort().join(', '));
    }
    if (!IDENTIFIER.test(String(input.assignmentId)) || !IDENTIFIER.test(String(input.attemptId)) ||
      typeof input.assignmentId !== 'string' || typeof input.attemptId !== 'string') {
      throw blocked('invalid_input', 'assignmentId and attemptId must be helper identifiers');
    }
    var destination = input.destination;
    if (!isPlainObject(destination) || !sameKeys(destination, ['hostId', 'threadId']) ||
      !plainValue(destination.threadId) || !plainValue(destination.hostId)) {
      throw blocked('invalid_input', 'destination must contain exactly non-empty threadId and hostId');
    }
    if (!cleanAbsolutePath(input.stateRoot) || !cleanAbsolutePath(input.helperPath)) {
      throw blocked('invalid_input', 'stateRoot and helperPath must be normalized absolute paths');
    }
    var outcome = input.outcome;
    if (typeof outcome === 'string') {
      outcome = parseJSONObject(outcome);
    }
    if (!isPlainObject(outcome)) {
      throw blocked('invalid_outcome', 'the helper outcome must be one JSON object');
    }
    return outcome;
  }

  // classifyReceipt derives a delivery status from the native tool result. Only the observed
  // minimal shapes qualify: an envelope of exactly {content, isError} holding exactly one
  // {type: 'text', text} item, whose success text is exactly {"threadId": <expected>}, or that
  // payload itself. Only the exact known argument-rejection text in such an envelope is
  // rejected. Any extra, conflicting or unfamiliar field (error, errors, status, a nested
  // isError) and any other error, such as a permission refusal, is uncertain.
  function classifyReceipt(receipt, expectedThreadId) {
    function reported(payload, shape) {
      if (!isPlainObject(payload) || !sameKeys(payload, ['threadId'])) {
        return { status: 'uncertain', shape: shape, reason: 'unrecognized_receipt' };
      }
      var threadId = payload.threadId;
      if (!plainValue(threadId)) {
        return { status: 'uncertain', shape: shape, reason: 'no_recipient_in_receipt' };
      }
      if (threadId !== expectedThreadId) {
        return { status: 'uncertain', shape: shape, reason: 'recipient_mismatch', reportedThreadId: threadId };
      }
      return { status: 'sent', shape: shape, reason: 'accepted', reportedThreadId: threadId };
    }
    if (isPlainObject(receipt) && 'content' in receipt) {
      var item = Array.isArray(receipt.content) && receipt.content.length === 1 ? receipt.content[0] : null;
      var text = isPlainObject(item) && sameKeys(item, ['text', 'type']) && item.type === 'text' &&
        typeof item.text === 'string' ? item.text.trim() : null;
      if (!sameKeys(receipt, ['content', 'isError']) || typeof receipt.isError !== 'boolean' || text === null) {
        return { status: 'uncertain', shape: 'unrecognized', reason: 'unrecognized_receipt' };
      }
      if (receipt.isError) {
        if (REJECTED_ARGUMENTS.test(text)) {
          return { status: 'rejected', shape: 'mcp_envelope', reason: 'invalid_arguments' };
        }
        return { status: 'uncertain', shape: 'mcp_envelope', reason: 'error_envelope' };
      }
      return reported(parseJSONObject(text), 'mcp_envelope');
    }
    if (typeof receipt === 'string') {
      var parsed = parseJSONObject(receipt);
      if (parsed !== null) {
        return reported(parsed, 'thread_text');
      }
    } else if (isPlainObject(receipt) && 'threadId' in receipt) {
      return reported(receipt, 'thread_object');
    }
    return { status: 'uncertain', shape: 'unrecognized', reason: 'unrecognized_receipt' };
  }

  function receiptText(receipt) {
    var text;
    try {
      text = JSON.stringify(receipt);
    } catch (error) {
      text = undefined;
    }
    if (typeof text !== 'string') {
      text = JSON.stringify({ forwarder_unserializable_receipt: typeof receipt });
    }
    return text;
  }

  function thrownText(error) {
    var name = error && typeof error.name === 'string' ? error.name : typeof error;
    var message = error && typeof error.message === 'string' ? error.message : String(error);
    return JSON.stringify({ forwarder_send_exception: { name: name, message: message } });
  }

  // send calls the native tool once. A throw is kept as evidence and classified uncertain,
  // because the host may have accepted the message before failing.
  async function send(callbacks, message, expectedThreadId) {
    try {
      var receipt = await callbacks.sendMessage(message);
      var classification = classifyReceipt(receipt, expectedThreadId);
      classification.evidence = receiptText(receipt);
      return classification;
    } catch (error) {
      return { status: 'uncertain', shape: 'thrown', reason: 'send_threw', evidence: thrownText(error) };
    }
  }

  function failureNotice(assignmentId, attemptId, attemptDir, deliveryFile) {
    return 'Relay delivery for assignment ' + assignmentId + ', attempt ' + attemptId +
      ' was rejected by the native send tool before delivery.\n' +
      'The original result is retained and was not resent. Attempt record: ' + attemptDir +
      '/attempt.json. Delivery record: ' + deliveryFile + '.\n' + TRANSPORT_ONLY;
  }

  function finish(result) {
    result.transport_only = TRANSPORT_ONLY;
    var coordinatorNotified = result.outcome === 'delivered' ||
      (result.failure_notice && result.failure_notice.status === 'sent');
    result.next_owner = coordinatorNotified ? 'coordinator' : 'relay';
    if (result.outcome === 'delivered') {
      result.next_step = 'Delivered to the authorized coordinator; the relay has nothing further to do.';
    } else if (coordinatorNotified) {
      result.next_step = 'A failure notice reached the coordinator, who reads the known attempt once. ' +
        'Do not resend.';
    } else {
      result.next_step = RELAY_STEP;
    }
    return result;
  }

  // forwardFallback sends a fallback notice once. It owns no attempt record, so it keeps the
  // native receipt in a new uniquely named sidecar below <state_root>/forwarder-receipts/ and
  // never touches attempts/, including an earlier attempt named by attempt_exists.
  async function forwardFallback(input, outcome, callbacks, result) {
    result.once_only = 'invocation';
    result.detail = 'Fallback notice without a delivery record; it is sent at most once by this ' +
      'invocation and never replayed.';
    var message = outcome.fallback_message;
    var problem = validateMessage(message, input.destination);
    if (problem !== '') {
      throw blocked('invalid_message', problem);
    }
    result.send_attempted = true;
    var sent = await send(callbacks, message, input.destination.threadId);
    result.derived_status = sent.status;
    result.receipt_shape = sent.shape;
    result.reason = sent.reason;
    result.outcome = sent.status === 'sent' ? 'delivered' : sent.status;
    var receiptDir = input.stateRoot + '/' + FALLBACK_RECEIPT_DIR + '/' + input.assignmentId;
    var receiptFile = receiptDir + '/' + input.attemptId + '-fallback-' + Date.now().toString(36) + '-' +
      Math.random().toString(36).slice(2) + '.json';
    try {
      var made = await callbacks.runCommand(['/bin/sh', '-c', MKDIR_SCRIPT, 'sh', receiptDir]);
      if (made.exitCode !== 0) {
        throw new Error('directory creation exited ' + made.exitCode);
      }
      await writeNewFile(callbacks.runCommand, receiptFile, sent.evidence);
      result.receipt_file = receiptFile;
    } catch (error) {
      result.receipt_problem = 'receipt not persisted: ' + error.message;
    }
    return result;
  }

  async function forwardRecorded(input, outcome, callbacks, result) {
    var run = callbacks.runCommand;
    var attemptDir = result.attempt_dir;
    var deliveryFile = attemptDir + '/delivery.json';
    var claimDir = attemptDir + '/' + CLAIM_DIR;
    result.once_only = 'attempt_claim';
    if (outcome.message_file !== attemptDir + '/message.json' || outcome.delivery_file !== deliveryFile) {
      throw blocked('binding_mismatch', 'message_file and delivery_file must be this attempt\'s ' +
        'message.json and delivery.json');
    }
    result.delivery_file = deliveryFile;

    var help;
    try {
      help = await run([input.helperPath, 'coordination', 'record-delivery', '--help']);
    } catch (error) {
      throw blocked('helper_unavailable', 'the helper could not be checked: ' + error.message);
    }
    if (help.exitCode !== 0 || help.output.indexOf('--attempt-dir') === -1 ||
      help.output.indexOf('--receipt-file') === -1) {
      throw blocked('helper_unavailable', 'the verified helper does not provide record-delivery');
    }

    var delivery;
    try {
      delivery = parseJSONObject(await readFile(run, deliveryFile));
    } catch (error) {
      throw blocked('delivery_record_unreadable', error.message);
    }
    var prepared = delivery && delivery.prepared_recipient;
    if (!delivery || typeof delivery.status !== 'string' || !isPlainObject(prepared)) {
      throw blocked('delivery_record_unreadable', 'the delivery record is not a recognized record');
    }
    if (prepared.thread_id !== input.destination.threadId || prepared.host_id !== input.destination.hostId) {
      throw blocked('recipient_mismatch', 'the delivery record names a different prepared recipient');
    }
    if (delivery.status !== 'pending') {
      result.delivery_status = delivery.status;
      throw new Stop('refused', 'already_recorded', 'the delivery was already recorded as ' +
        delivery.status + '; nothing was sent and the record is unchanged');
    }
    result.delivery_status = 'pending';

    var message;
    try {
      message = parseJSONObject(await readFile(run, outcome.message_file));
    } catch (error) {
      throw blocked('message_unreadable', error.message);
    }
    var problem = validateMessage(message, input.destination);
    if (problem !== '') {
      throw blocked('invalid_message', problem);
    }

    var claim;
    try {
      claim = await run(['/bin/mkdir', '-m', '700', '--', claimDir]);
    } catch (error) {
      claim = { exitCode: -1 };
    }
    if (claim.exitCode !== 0) {
      throw new Stop('refused', 'claim_unavailable', 'a delivery claim already exists for this attempt ' +
        'or could not be created, so a send may already have started. Nothing was sent. The owner ' +
        'checks destination history before any manual action; age or a dead process never ' +
        'authorizes replay.');
    }
    result.claim_dir = claimDir;

    result.send_attempted = true;
    var sent = await send(callbacks, message, input.destination.threadId);
    result.derived_status = sent.status;
    result.receipt_shape = sent.shape;
    result.reason = sent.reason;

    var receiptFile = claimDir + '/' + RECEIPT_FILE;
    var problems = [];
    try {
      await writeNewFile(run, receiptFile, sent.evidence);
      result.receipt_file = receiptFile;
    } catch (error) {
      problems.push('receipt not persisted: ' + error.message);
    }
    if (problems.length === 0) {
      var argv = [input.helperPath, 'coordination', 'record-delivery', '--attempt-dir', attemptDir,
        '--status', sent.status, '--receipt-file', receiptFile];
      if (sent.reportedThreadId) {
        argv.push('--reported-thread-id', sent.reportedThreadId);
      }
      argv.push('--json');
      try {
        var recorded = await run(argv);
        var record = parseJSONObject(recorded.output.trim());
        var confirmed = sent.status !== 'sent' || (record && record.recipient_confirmed === true);
        if (recorded.exitCode !== 0 || !record || record.status !== 'recorded' ||
          record.delivery_status !== sent.status || !confirmed) {
          problems.push('record-delivery did not confirm the ' + sent.status + ' record');
        } else {
          result.recorded = true;
          result.delivery_status = record.delivery_status;
        }
      } catch (error) {
        problems.push('record-delivery failed: ' + error.message);
      }
    }
    if (problems.length > 0) {
      result.outcome = 'recording_failed';
      result.detail = problems.join('; ') + '. The send ' +
        (sent.status === 'rejected' ? 'was rejected' : 'may have been delivered') +
        '; the claim is retained and the delivery record was not completed by this workflow.';
    } else {
      result.outcome = sent.status === 'sent' ? 'delivered' : sent.status;
    }

    if (sent.status === 'rejected') {
      result.failure_notice = await sendFailureNotice(input, callbacks, attemptDir, deliveryFile, claimDir);
    }
    return result;
  }

  // sendFailureNotice sends one compact notice after a definite argument rejection only. It
  // never resends the original reply and is never attempted twice.
  async function sendFailureNotice(input, callbacks, attemptDir, deliveryFile, claimDir) {
    var notice = {
      threadId: input.destination.threadId,
      hostId: input.destination.hostId,
      prompt: failureNotice(input.assignmentId, input.attemptId, attemptDir, deliveryFile)
    };
    var sent = await send(callbacks, notice, input.destination.threadId);
    var report = { status: sent.status, reason: sent.reason, receipt_shape: sent.shape };
    var receiptFile = claimDir + '/' + NOTICE_RECEIPT_FILE;
    try {
      await writeNewFile(callbacks.runCommand, receiptFile, sent.evidence);
      report.receipt_file = receiptFile;
    } catch (error) {
      report.receipt_problem = 'notice receipt not persisted: ' + error.message;
    }
    return report;
  }

  // forwardPreparedMessage is the single entry point. callbacks.sendMessage(args) calls the
  // native send tool; callbacks.runCommand(argv) runs one command and resolves to
  // {exitCode, output}.
  async function forwardPreparedMessage(input, callbacks) {
    var result = {
      outcome: 'blocked',
      reason: '',
      detail: '',
      assignment_id: isPlainObject(input) && typeof input.assignmentId === 'string' ? input.assignmentId : '',
      attempt_id: isPlainObject(input) && typeof input.attemptId === 'string' ? input.attemptId : '',
      send_attempted: false,
      recorded: false,
      delivery_status: null,
      failure_notice: null
    };
    try {
      if (!callbacks || typeof callbacks.sendMessage !== 'function' ||
        typeof callbacks.runCommand !== 'function') {
        throw blocked('missing_capability', 'the native send and shell capabilities are required');
      }
      var outcome = validateInput(input);
      var attemptDir = input.stateRoot + '/attempts/' + input.assignmentId + '/' + input.attemptId;
      result.attempt_dir = attemptDir;
      if (outcome.assignment_id !== undefined && outcome.assignment_id !== input.assignmentId ||
        outcome.attempt_id !== undefined && outcome.attempt_id !== input.attemptId ||
        outcome.attempt_dir !== undefined && outcome.attempt_dir !== attemptDir) {
        throw blocked('binding_mismatch', 'the helper outcome belongs to a different assignment, ' +
          'attempt or state root');
      }
      result.helper_status = typeof outcome.status === 'string' ? outcome.status : '';
      var hasFile = outcome.message_file !== undefined;
      var hasFallback = outcome.fallback_message !== undefined;
      if (hasFile && !hasFallback) {
        return finish(await forwardRecorded(input, outcome, callbacks, result));
      }
      if (hasFallback && !hasFile) {
        return finish(await forwardFallback(input, outcome, callbacks, result));
      }
      throw blocked('no_message', hasFile ? 'the outcome names both a message file and a fallback ' +
        'message' : 'the helper outcome names no message file or fallback message');
    } catch (stop) {
      if (!(stop instanceof Stop)) {
        // An unexpected failure after a send started must never read as a pre-send blocker.
        stop = new Stop(result.send_attempted ? 'recording_failed' : 'blocked', 'internal_error',
          stop && stop.message ? stop.message : String(stop));
      }
      result.outcome = stop.outcome;
      result.reason = stop.reason;
      result.detail = stop.detail;
      return finish(result);
    }
  }

  // codexCallbacks adapts the Codex tool-execution namespace. exec_command resolves to
  // {output, exit_code}; exit_code is absent while the command is still running.
  function codexCallbacks(tools) {
    if (!tools || typeof tools[SEND_TOOL] !== 'function' || typeof tools[SHELL_TOOL] !== 'function') {
      return null;
    }
    return {
      sendMessage: function (args) {
        return tools[SEND_TOOL](args);
      },
      runCommand: async function (argv) {
        var result = await tools[SHELL_TOOL]({
          cmd: shellCommand(argv),
          login: false,
          tty: false,
          yield_time_ms: 30000,
          max_output_tokens: 60000
        });
        if (!isPlainObject(result) || typeof result.output !== 'string') {
          throw new Error('unrecognized shell result');
        }
        if (typeof result.exit_code !== 'number') {
          throw new Error('shell command did not finish');
        }
        return { exitCode: result.exit_code, output: result.output };
      }
    };
  }

  function forwardWithCodexTools(tools, input) {
    return forwardPreparedMessage(input, codexCallbacks(tools));
  }

  return {
    forwardWithCodexTools: forwardWithCodexTools,
    forwardPreparedMessage: forwardPreparedMessage,
    codexCallbacks: codexCallbacks,
    classifyReceipt: classifyReceipt,
    shellCommand: shellCommand
  };
})();

if (typeof module === 'object' && module && module.exports) {
  module.exports = CodeheartRelayForwarder;
}
