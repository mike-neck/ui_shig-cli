# しぐれうい CLI
しぐれういボタンをスクレイピングしてCLIで鳴らすだけのアプリです。

入手方法
---

### ダウンロード

- [GitHub レポジトリー](https://github.com/mike-neck/ui_shig-cli) の画面右にある [Release](https://github.com/mike-neck/ui_shig-cli/releases) の [Tags](https://github.com/mike-neck/ui_shig-cli/tags) から最新版をダウンロードします 
- Windows で Intel CPU の場合は `windows-amd64` という名前の zip ファイルをダウンロードしてください
- Mac で M1 系の CPU の場合は `darwin-arm64` という名前の zip ファイルをダウンロードしてください

### ビルドする

- 下の方の[「ビルド方法」](#ビルド方法)に記載されています

使い方
---

次のコマンドを入力すると、「イラストレーターのしぐれういと申しま～す」と再生されます。

```shell
ui_shig say illust
```

コマンドには以下のようなものがあります。

|  コマンド   | 動作                   |
|:-------:|:---------------------|
|  `say`  | しぐれういの声が再生されます       |
| `list`  | どのような音声があるかリストを表示します |
| `cache` | しぐれういの音声をダウンロードします   |
| `help`  | コマンドのヘルプを表示します       |

ビルド方法
---

以下でビルドできます。ビルドすると、 `bin` ディレクトリーに `ui_shig` コマンドができます。

```shell
make init   # ビルドに必要なデータを集めます
make build  # ツールをビルドします
```

### 依存ツール・ライブラリー

ビルドには以下のソフトウェアが必要です。

- `bash`
- `go`
- 各種ユーティリティ(`tr`/`grep`/`sed`/`awk`/`mktemp`/`make`)
- `curl`
- `jq`

Go の以下のモジュールを使っています。

- [gopxl/Beep](https://github.com/gopxl/beep)
- [ebitengine/oto](https://github.com/ebitengine/oto)
- [urfave/cli/v2](https://github.com/urfave/cli)

また、 Linux でビルドする場合 [ebitengine/oto](https://github.com/ebitengine/oto) が以下のライブラリーに依存しているので、事前にインストールしてください。

- `libasound2-dev`
