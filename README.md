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

## Usage

### Keybindings

### Formulas

## License

This work is licensed under the GNU General Public License version 3 (GPLv3).

[<img src="https://s-christy.com/status-banner-service/GPLv3_Logo.svg" width="150" />](https://www.gnu.org/licenses/gpl-3.0.en.html)
