# しぐれうい CLI
しぐれういボタンをスクレイピングしてCLIで鳴らすだけのアプリです。

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

ビルドには以下のソフトウェアが必要です。

- `bash`
- `go`
- 各種ユーティリティ(`tr`/`grep`/`sed`/`awk`/`mktemp`/`make`)
- `curl`
- `jq`

また、 Go の以下のモジュールを使っています。

- [gopxl/Beep](https://github.com/gopxl/beep)
- [ebitengine/oto](https://github.com/ebitengine/oto)
- [urfave/cli/v2](https://github.com/urfave/cli)
