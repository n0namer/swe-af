from pathlib import Path
p=Path(r"D:\TEMP\swe-af-go-dag-fence-20261020\go\internal\orch\resume.go")
s=p.read_text(encoding="utf-8")
s=s.replace('"completed_issues": checkpointValueOrEmpty(checkpoint, "completed_issues"),','"completed_issues": checkpointValueOrEmpty(checkpoint, "completed_issues"),\n\t\t"last_resume": checkpoint["last_resume"],')
s=s.replace('checkpoint["last_resume_result"] = result','if owner == "" && !hasEpoch {\n\t\treturn result, nil\n\t}\n\tcheckpoint["last_resume_result"] = result')
p.write_text(s,encoding="utf-8")
p=Path(r"D:\TEMP\swe-af-go-dag-fence-20261020\go\internal\dag\checkpoint.go")
s=p.read_text(encoding="utf-8")
s=s.replace('if err := os.WriteFile(path, b, 0o644); err != nil {\n\t\treturn\n\t}', '''tmp, err := os.CreateTemp(filepath.Dir(path), ".checkpoint-*.tmp")
    if err != nil { return }
    tmpPath := tmp.Name()
    defer os.Remove(tmpPath)
    if _, err := tmp.Write(b); err != nil { tmp.Close(); return }
    if err := tmp.Sync(); err != nil { tmp.Close(); return }
    if err := tmp.Close(); err != nil { return }
    if err := os.Rename(tmpPath, path); err != nil { return }''',1)
p.write_text(s,encoding="utf-8")
