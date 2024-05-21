# go-hzk

`go-hzk` 是一个用于 HZK12 和 HZK16 点阵字库的 Go 库，并提供配套命令行工具。它可以将支持的 GB2312 字符转换为点阵矩阵、终端文本或 PNG 图片。

## 功能特性

- 支持 HZK12：字模高度 12 行，每行按 16 点读取。
- 支持 HZK16：字模高度 16 行，每行按 16 点读取。
- 单字默认按 16 点宽度输出，可通过 `CellWidth` 或 `--cell-width` 设置单字输出宽度。
- 默认内嵌 HZK12 和 HZK16 字库数据。
- 支持加载外部 HZK 字库文件。
- 提供字模、矩阵、文本渲染和 PNG 编码 API。
- 命令行基于 Cobra 实现，提供 `text` 和 `image` 子命令。
- 支持配置前景、背景、字间距、行间距、图片缩放倍率和图片边距。

## 安装

```sh
go install github.com/secriy/go-hzk/cmd/hzk@latest
```

本地开发时可直接运行：

```sh
go run ./cmd/hzk text --size 16 中文
```

## 库用法

将文本转换为布尔矩阵并渲染为终端文本：

```go
package main

import (
 "fmt"

 "github.com/secriy/go-hzk"
)

func main() {
 font, err := hzk.New(hzk.HZK16)
 if err != nil {
  panic(err)
 }

 matrix, err := font.Matrix("中文")
 if err != nil {
  panic(err)
 }
 fmt.Println(len(matrix), len(matrix[0]))

 output, err := font.Render("中文", hzk.RenderOptions{
  Foreground:   "█",
  Background:   " ",
  GlyphSpacing: 1,
  CellWidth:    12,
 })
 if err != nil {
  panic(err)
 }
 fmt.Println(output)
}
```

生成 PNG 图片。以下示例假设已创建 `font`，并导入 `os`：

```go
file, err := os.Create("text.png")
if err != nil {
 panic(err)
}
defer file.Close()

err = font.EncodePNG(file, "中文", hzk.ImageOptions{
 Scale:        8,
 Padding:      8,
 GlyphSpacing: 1,
 CellWidth:    12,
})
if err != nil {
 panic(err)
}
```

加载外部字库：

```go
font, err := hzk.LoadFile(hzk.HZK16, "./fonts/HZK16")
```

## 命令行用法

渲染终端文本：

```sh
hzk text --size 12 --fg "#" --bg " " 中文
```

按 12 点单字宽度渲染：

```sh
hzk text --size 12 --cell-width 12 中文
```

从 stdin 读取文本：

```sh
printf "中文" | hzk text --size 16
```

生成 PNG 图片：

```sh
hzk image --size 16 --out text.png --scale 8 --padding 8 中文
```

将 PNG 数据写入 stdout：

```sh
hzk image --size 16 --out - 中文 > text.png
```

使用外部字库文件：

```sh
hzk text --size 16 --font ./fonts/HZK16 中文
```

## 字符支持

本库使用内嵌 GB2312 映射表定位 HZK 字模，支持 GB2312 汉字和常用全角符号。ASCII 空格会渲染为空白字模；其他 ASCII 字符不会自动转换为全角形式。

## 字库数据

默认构建会内嵌 HZK12 和 HZK16 字库数据。调用方可以通过 `LoadFile` 或 `NewFromBytes` 使用自定义字库数据。
