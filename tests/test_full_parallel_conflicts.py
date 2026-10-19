"""Native Full planning must preserve maximum *safe* independent width."""
from swe_af.reasoners.pipeline import (
    _compute_levels,
    _validate_file_conflicts,
    _serialize_file_conflicts,
)


def issue(name, *, files=(), depends=(), repo=""):
    return {
        "name": name,
        "depends_on": list(depends),
        "files_to_modify": list(files),
        "files_to_create": [],
        "target_repo": repo,
    }


def test_thirty_nonoverlapping_issues_remain_parallel():
    issues = [issue(f"job-{i:02}", files=[f"pkg/job{i:02}.go"]) for i in range(30)]
    changes = _serialize_file_conflicts(issues)
    assert changes == []
    assert len(_compute_levels(issues)[0]) == 30
    assert _validate_file_conflicts(issues, _compute_levels(issues)) == []


def test_file_overlap_is_serialized_without_serializing_independent_issues():
    issues = [
        issue("a", files=["src/Foo.go"]),
        issue("b", files=[r"src\foo.go"]),
        issue("c", files=["pkg/c.go"]),
        issue("d", files=["src/foo.go"]),
    ]
    changes = _serialize_file_conflicts(issues)
    assert changes
    levels = _compute_levels(issues)
    assert "a" in levels[0] and "c" in levels[0]
    assert "b" not in levels[0] and "d" not in levels[0]
    assert _validate_file_conflicts(issues, levels) == []


def test_later_level_collision_is_rechecked():
    issues = [
        issue("a", files=["x.go"]),
        issue("b", files=["x.go"]),
        issue("unblock", files=["y.go"]),
        issue("c", depends=["unblock"], files=["x.go"]),
    ]
    _serialize_file_conflicts(issues)
    assert _validate_file_conflicts(issues, _compute_levels(issues)) == []
    levels = _compute_levels(issues)
    assert not any(len(set(level) & {"a", "b", "c"}) > 1 for level in levels)


def test_same_relative_name_in_different_repositories_can_run_together():
    issues = [
        issue("one", repo="CodeAF", files=["internal/version.go"]),
        issue("two", repo="PR-AF", files=["internal/version.go"]),
    ]
    assert not _serialize_file_conflicts(issues)
    assert len(_compute_levels(issues)[0]) == 2
    assert _validate_file_conflicts(issues, _compute_levels(issues)) == []
