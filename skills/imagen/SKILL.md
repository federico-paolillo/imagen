---
name: imagen
description: Generate placeholder PNG images of an exact pixel size. Produces a solid magenta image with its own dimensions (e.g. "1024x768") centered in black. Use when you need a quick placeholder, test fixture, or dummy image of a specific width and height.
compatibility: Requires the bundled scripts/imagen executable to match the host OS and architecture. The current directory must be writable.
---

# imagen

`imagen` writes a PNG of an exact size: solid magenta (`#FF00FF`) with the text `WIDTHxHEIGHT` centered in black. The text scales with the image so it always fits. Output is deterministic.

## Running it

The executable is bundled with this skill at `scripts/imagen`, relative to this skill's directory:

```sh
scripts/imagen -w WIDTH -h HEIGHT
```

The image is written to your **current working directory**, not the skill's directory. To write it elsewhere, `cd` there first and call the executable by its absolute path:

```sh
cd assets/img && /path/to/skill/scripts/imagen -w 640 -h 480
```

## Flags

| Flag | Meaning |
| --- | --- |
| `-w` | Width in pixels, 1 to 8192 (required) |
| `-h` | Height in pixels, 1 to 8192 (required) |

`-h` is the **height**, not help. Use `-help` for usage text.

## Output

The file is named `WIDTHxHEIGHT.png` and written to the current working directory. There is no output-path option. An existing file with the same name is overwritten. The filename is printed to stdout, so it can be captured:

```sh
path=$(scripts/imagen -w 800 -h 600)   # path=800x600.png
```

## Examples

```sh
scripts/imagen -w 1024 -h 1024     # writes 1024x1024.png
scripts/imagen -w 1920 -h 1080     # writes 1920x1080.png
for s in 64 128 256; do scripts/imagen -w $s -h $s; done
```

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Image written (or `-help` shown) |
| 1 | Dimension negative or above 8192, or the file could not be written |
| 2 | Bad usage: missing or zero `-w`/`-h`, unknown flag, non-numeric value, extra arguments |

Errors go to stderr, and no file is created on failure.

## Verifying

`file 800x600.png` should report `PNG image data, 800 x 600`.
