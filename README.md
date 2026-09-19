# test-tinygo-kint

kinT (Teensy 4.1) 上で TinyGo USB キーボードを検証するリポジトリ。USB デバイスは TinyGo PR [#5691](https://github.com/tinygo-org/tinygo/pull/5691)。ビルドは PR [#5704](https://github.com/tinygo-org/tinygo/pull/5704)（#5691 の USB コミットを含む）を checkout する。

Goal 1 ファームは `hid-keyboard/`（CDC + HID Boot Keyboard、Report Protocol、QMK kint41 default キーマップ）。設計は [`hid-keyboard/design.md`](hid-keyboard/design.md)。`tinygo-keyboard/` と `usb-midi/` は別バイナリ。

## build patched tinygo

```
git clone github.com/tinygo-org/tinygo
cd tinygo
gh pr checkout 5704

make llvm-source
make llvm-build
make
cp $(tinygo env TINYGOROOT)/src/device/nxp/mimxrt1062* src/device/nxp/ 2>/dev/null
./build/tinygo version
```

## build and write firmware

```
$(ghq root)/github.com/tinygo-org/tinygo/build/tinygo flash --target teensy41 --size short --stack-size 8kb ./hid-keyboard
```

VID:PID は `16C0:0483`。`tinygo flash` の 1200 bps リセットが使える。失敗したら Teensy プログラムボタンで HalfKay。キー `[13,5]`（QMK `QK_BOOT`）でも HalfKay に入る。

マトリクス座標を CDC に出すときは `hid-keyboard/main.go` の `const debug = true` にして焼き直す。

## Goal 1 確認

列挙（Linux）:

```
lsusb -d 16c0:0483
lsusb -v -d 16c0:0483
```

IF2 が `bInterfaceClass 3`, `bInterfaceSubClass 1`, `bInterfaceProtocol 1`。EP `83` / `03` Interrupt、wMaxPacketSize 8、bInterval 1。iProduct は `kinT TinyGo`。

macOS:

```
system_profiler SPUSBDataType | grep -A20 'kinT'
```

入力（US レイアウトホスト）: `asdf`、`,` `.` `/`、左上 `=`、親指 Space / Enter / Backspace / Delete、修飾+文字。Caps Lock で D12 が反転すること。

`SET_PROTOCOL` / `GET_PROTOCOL` は ACK する（送信フォーマットは 8 バイトのまま）。`machine/usb/hid` は import しない。
