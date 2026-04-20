# Skill: gga-shellcheck-standards

## Purpose
Maintain shell script quality and portability across all GGA bash scripts.

## When to Use
When writing or modifying bash scripts in `bin/` or `lib/`.

## Running ShellCheck

```bash
# Via make (preferred)
make lint

# Direct (matches CI configuration)
shellcheck -x -e SC1090,SC1091,SC2162,SC2129 bin/gga lib/*.sh
```

All new code must pass `make lint` before pushing.

## Accepted Exclusions

| Code | Rule | Why GGA Excludes It |
|------|------|---------------------|
| SC1090 | Can't follow non-constant source | GGA uses dynamic paths for lib loading |
| SC1091 | Not following sourced file | Same — dynamic lib sourcing |
| SC2162 | read without -r | GGA intentionally handles backslash input |
| SC2129 | Use { } >> file | Style preference — individual redirects are clearer here |

Do NOT add new exclusions without justification in a PR comment.

## Critical Rules

### Quoting
```bash
# GOOD — always quote variables
echo "$var"
"$cmd" "$arg"

# BAD — unquoted
echo $var
$cmd $arg
```

### Conditionals
```bash
# GOOD — use [[ ]] in bash
[[ "$var" == "value" ]]
[[ -f "$file" ]]

# BAD — use [ ] only for POSIX compatibility
[ "$var" = "value" ]
```

### Function Variables
```bash
# GOOD — always local
my_function() {
  local result
  result="something"
}

# BAD — leaks to global scope
my_function() {
  result="something"
}
```

### No eval
```bash
# BAD — never
eval "$user_input"

# If you think you need eval, you don't. Use arrays or parameter expansion.
```

### macOS vs Linux Portability
```bash
# sed -i behaves differently
# macOS requires an extension (even empty string)
sed -i '' 's/foo/bar/' file        # macOS
sed -i 's/foo/bar/' file           # Linux

# GGA pattern — detect and branch:
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' 's/foo/bar/' "$file"
else
  sed -i 's/foo/bar/' "$file"
fi
```

## Common Fixes

| ShellCheck warning | Fix |
|-------------------|-----|
| SC2086: double quote | Add `"$var"` around the variable |
| SC2181: check $? | Replace `if [ $? -eq 0 ]` with `if command; then` |
| SC2155: declare + assign | Split: `local var; var=$(cmd)` |
| SC2206: word splitting | Use `read -ra arr <<< "$str"` |

## Cookbook

| If... | Then... |
|-------|---------|
| Writing a new function | Add `local` for every variable, quote every expansion |
| Using `sed -i` | Add macOS/Linux compatibility check |
| Reading user input | Use `read -r` unless backslash handling is intentional |
| Writing a loop over files | Use `while IFS= read -r line` or glob, never `ls` |
