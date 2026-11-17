#!/bin/sh
cd "$(dirname "$0")"
fail=0
for f in test_cases/graph*.csv; do
  exp="${f%.csv}_output.txt"
  if ./graph_solution "$f" | diff -q - "$exp" >/dev/null; then echo "PASS $f"; else echo "FAIL $f"; fail=1; fi
done
exit $fail
