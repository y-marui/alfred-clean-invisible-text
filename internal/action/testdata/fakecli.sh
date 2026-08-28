#!/usr/bin/env bash
# A stand-in for the pinned clean-invisible-text binary, used only by
# cliinvoke_test.go to exercise Run()'s exit-code/JSON handling without a
# network dependency or duplicating the real CLI's classification logic.
set -euo pipefail

path="${*: -1}"
scenario="${FAKECLI_SCENARIO:-clean}"

case "$scenario" in
  clean)
    printf '[{"path":"%s","findings":[],"changed":false,"error":null}]\n' "$path"
    exit 0
    ;;
  cleaned)
    printf '[{"path":"%s","findings":[{"line":1,"column":1,"offset":0,"rune":"U+200B","name":"ZERO WIDTH SPACE","category":"zwsp","action":"remove","replacement":""}],"changed":true,"error":null}]\n' "$path"
    exit 1
    ;;
  warning)
    printf '[{"path":"%s","findings":[{"line":1,"column":1,"offset":0,"rune":"U+2065","name":"UNASSIGNED","category":"unclassified","action":"warn","replacement":""}],"changed":true,"error":null}]\n' "$path"
    exit 1
    ;;
  file-error)
    printf '[{"path":"%s","findings":[],"changed":false,"error":"invalid UTF-8"}]\n' "$path"
    exit 2
    ;;
  process-error)
    echo "fatal: something went wrong" >&2
    exit 2
    ;;
  malformed-json)
    echo "not json" >&2
    exit 0
    ;;
  *)
    echo "fakecli.sh: unknown FAKECLI_SCENARIO=$scenario" >&2
    exit 2
    ;;
esac
