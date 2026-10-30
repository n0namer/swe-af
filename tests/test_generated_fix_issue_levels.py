"""Regression tests for dependency ordering in generated verification fixes."""

import pytest

from swe_af.app import _generated_fix_issue_levels


def test_generated_fix_issues_follow_dependencies_when_input_is_reversed():
    issues = [
        {"name": "dependent-fix", "depends_on": ["prerequisite-fix"]},
        {"name": "prerequisite-fix", "depends_on": []},
    ]

    assert _generated_fix_issue_levels(issues) == [
        ["prerequisite-fix"],
        ["dependent-fix"],
    ]


def test_generated_fix_issues_keep_independent_fixes_parallel():
    issues = [
        {"name": "first-fix", "depends_on": []},
        {"name": "second-fix", "depends_on": []},
    ]

    assert _generated_fix_issue_levels(issues) == [["first-fix", "second-fix"]]


def test_generated_fix_issue_cycles_fail_closed():
    issues = [
        {"name": "first-fix", "depends_on": ["second-fix"]},
        {"name": "second-fix", "depends_on": ["first-fix"]},
    ]

    with pytest.raises(ValueError, match="Dependency cycle detected"):
        _generated_fix_issue_levels(issues)


def test_generated_fix_issues_with_unknown_dependencies_fail_closed():
    issues = [
        {"name": "dependent-fix", "depends_on": ["missing-fix"]},
        {"name": "independent-fix", "depends_on": []},
    ]

    with pytest.raises(ValueError, match="Unknown dependenc.*missing-fix"):
        _generated_fix_issue_levels(issues)
