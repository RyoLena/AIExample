# Go dreaming of adventure sample

一个用 Go 编写的开发人员示例，展示了 Gemini 的创意写作能力。根据用户输入，Gemini 一次编写一个章节的长篇小说。

<a href="https://idx.google.com/import?url=https://github.com/google-gemini/go-dreaming-of-adventure-sample">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://cdn.idx.dev/btn/open_dark_32@2x.png">
  <source media="(prefers-color-scheme: light)" srcset="https://cdn.idx.dev/btn/open_light_32@2x.png">
  <img height="32" alt="Open in IDX" src="https://cdn.idx.dev/btn/open_purple_32@2x.png">
</picture>
</a>

## 环境设置

此示例应用程序可以在Project IDX中打开，也可以在本地开发环境中运行。

## 项目IDX

1. 在 Project IDX 中打开此 repo：
    - [Open in Project IDX](https://idx.google.com/import?url=https://github.com/google-gemini/go-dreaming-of-adventure-sample)
    - 等待导入过程完成
    -  打开 IDX 面板并单击使用 Gemini API 集成的“身份验证”。
    - 一旦通过身份验证，单击即可获取密钥，该密钥将被复制到您的键盘。
    - 将密钥添加到 `.idx/dev.nix` 中的环境变量部分。

2. 打开一个新的终端窗口：
    - 打开命令面板（CTRL/CMD-SHIFT-P）
    - 开始输入 **terminal**
    - 选择**终端：创建新终端**
    - 运行`go run`。

## 本地开发环境

1. 克隆此存储库: `git clone https://github.com/google-gemini/go-dreaming-of-adventure-sample`

2. 验证是否安装了 Go 1.22 或更高版本：
    - 使用以下方法验证版本`go version`
    - 根据需要安装Go，参见: https://go.dev/doc/install


## 运行示例

1. 获取 Gemini API 密钥
    - 启动 Google AI Studio: https://aistudio.google.com/
    - **点击**获取 **API 密钥**

2. `API_KEY`在环境变量中设置 **API Key**
    - `export API_KEY=<your_api_key>`

3. 编译并运行程序:
    - `go run .`

4. 当被问到“你想梦见什么？”时，用一些有趣的事情来回答.
    - 例如输入: `I want to dream about unicode`
