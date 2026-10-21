from __future__ import annotations

from swe_af.runtime.codex_harness_patch import _codex_windows_sandbox_args


def test_codex_windows_mxc_opt_in() -> None:
    assert _codex_windows_sandbox_args(
        {"SWE_CODEX_PREFER_MXC": "1"}, "nt"
    ) == ["-c", "features.prefer_mxc=true"]


def test_codex_windows_mxc_does_not_change_unrelated_runtimes() -> None:
    assert _codex_windows_sandbox_args({}, "nt") == []
    assert _codex_windows_sandbox_args({"SWE_CODEX_PREFER_MXC": "1"}, "posix") == []
    assert _codex_windows_sandbox_args({"SWE_CODEX_PREFER_MXC": "false"}, "nt") == []
