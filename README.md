![Banner](https://s-christy.com/sbs/status-banner.svg?icon=editor/table_chart&hue=200&title=Spreadsheet%20Editor&description=A%20terminal-based%20TSV%20spreadsheet%20editor%20with%20formula%20support)

## Overview

<p align="center">
  <img src="./assets/screenshot.png" width=500 />
</p>

Spreadsheet Editor is a lightweight, terminal-based editor for TSV
(tab-separated values) files. It supports cell formulas, vim-style navigation,
and a range of keyboard-driven editing commands.

## Features

- Vim-style navigation with `hjkl` (single cell) and `HJKL` (five cells)
- Cell formulas with `=` prefix (e.g. `=A1*B2`, `=sum("B1:D5")`)
- Built-in functions: `sum()` for ranges and `strlen()` for string length
- Formula chaining: formulas can reference other formula cells
- Copy and paste cells with `yy` / `p`
- Increment and decrement numeric cell values with `Ctrl-A` / `Ctrl-X`
- Shift-arrow keys to increment a cell's value and advance the cursor
- Color-mark cells with `c`
- Toggle grid display with `g`
- Equalize or reset column widths with `=` / `+`
- Open the file in `vim` directly with `e`
- In-editor help menu with `F1`

## Setup

```
go build
```

## Usage

```
./spreadsheet-editor <file.tsv>
```

A sample file is provided to demonstrate the available features:

```
./spreadsheet-editor assets/budget.tsv
```

### Keybindings

| Key | Action |
|---|---|
| `h` / `j` / `k` / `l` | Move left / down / up / right |
| `H` / `J` / `K` / `L` | Move five cells in that direction |
| `0` | Jump to origin |
| `Enter` | Edit cell |
| `Escape` / `q` | Quit |
| `s` | Save |
| `yy` / `p` | Copy / paste cell |
| `Ctrl-A` / `Ctrl-X` | Increment / decrement cell value |
| `c` | Toggle color mark on cell |
| `g` | Toggle grid |
| `=` | Equalize column widths |
| `+` | Reset column widths |
| `e` | Edit file in vim |
| `F1` | Show help |

### Formulas

Cells starting with `=` are evaluated as expressions. Cell references like `A1`
or `C3` are substituted with their values before evaluation. The `sum()`
function accepts a quoted range string (e.g. `sum("B1:B5")`), and `strlen()`
returns the character count of a cell's value.

```
=A1*B2
=sum("B1:D1")
=E6/3
=strlen(A1)
```

## License

This work is licensed under the GNU General Public License version 3 (GPLv3).

[<img src="https://s-christy.com/status-banner-service/GPLv3_Logo.svg" width="150" />](https://www.gnu.org/licenses/gpl-3.0.en.html)
