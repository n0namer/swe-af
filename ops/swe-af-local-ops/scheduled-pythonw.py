"""Console-free Windows Task Scheduler entrypoint; same interpreter, no child manager.
Capture normal Python stdout/stderr and preserve SystemExit/exception return codes.
"""
from __future__ import annotations
import runpy
import sys
from pathlib import Path

def main() -> None:
    if len(sys.argv) < 3:
        raise SystemExit(64)
    script = Path(sys.argv[1]).resolve(strict=True)
    log_prefix = Path(sys.argv[2])
    log_prefix.parent.mkdir(parents=True, exist_ok=True)
    passed_argv = [str(script), *sys.argv[3:]]
    original_argv, original_path = sys.argv, sys.path[:]
    original_out, original_err = sys.stdout, sys.stderr
    with Path(str(log_prefix) + ".stdout.log").open("a", encoding="utf-8", buffering=1) as out, Path(str(log_prefix) + ".stderr.log").open("a", encoding="utf-8", buffering=1) as err:
        try:
            sys.stdout, sys.stderr = out, err
            sys.argv = passed_argv
            sys.path.insert(0, str(script.parent))
            runpy.run_path(str(script), run_name="__main__")
        finally:
            sys.stdout, sys.stderr = original_out, original_err
            sys.argv, sys.path = original_argv, original_path

if __name__ == "__main__":
    main()
