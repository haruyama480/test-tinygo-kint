# test-tinygo-kint

kinT (Teensy 4.1) 上で TinyGo USB キーボードを検証するリポジトリ。USB デバイスは TinyGo PR [#5691](https://github.com/tinygo-org/tinygo/pull/5691)。ビルドは PR [#5704](https://github.com/tinygo-org/tinygo/pull/5704)（#5691 の USB コミットを含む）を checkout する。

| ディレクトリ | 内容 |
| --- | --- |
| [`hid-keyboard-g1/`](hid-keyboard-g1/) | Goal 1: USB FS + HID Report Protocol。設計は [`hid-keyboard-g1/design.md`](hid-keyboard-g1/design.md) |
| [`hid-keyboard-g2/`](hid-keyboard-g2/) | Goal 2: HID Boot Protocol。設計は [`hid-keyboard-g2/design.md`](hid-keyboard-g2/design.md) |
| `tinygo-keyboard/` | `sago35/tinygo-keyboard` + Vial。別バイナリ |
| `usb-midi/` | 同じマトリクスの MIDI。別バイナリ |

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

Goal 1:

```
$(ghq root)/github.com/tinygo-org/tinygo/build/tinygo flash --target teensy41 --size short --stack-size 8kb ./hid-keyboard-g1
```

Goal 2:

```
$(ghq root)/github.com/tinygo-org/tinygo/build/tinygo flash --target teensy41 --size short --stack-size 8kb ./hid-keyboard-g2
```

VID:PID は `16C0:0483`。`tinygo flash` の 1200 bps リセットが使える。失敗したら Teensy プログラムボタンで HalfKay。

マトリクス座標を CDC に出すときは対象ディレクトリの `main.go` で `const debug = true` にして焼き直す。

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

入力は `hid-keyboard-g1/keymap.go`（QMK `default_pretty` / my-customize）。親指中央は左右 Command。

`SET_PROTOCOL` / `GET_PROTOCOL` は ACK する（Goal 1 では送信フォーマットは 8 バイトのまま）。`machine/usb/hid` は import しない。
