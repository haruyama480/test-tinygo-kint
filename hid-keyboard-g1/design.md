# kinT TinyGo HID キーボードファームウェア設計

| 項目 | 内容 |
| --- | --- |
| 文書タイトル | kinT + Teensy 4.1 TinyGo HID キーボード（Goal 1） |
| 著者 | TBD |
| 日付 | 2026-09-20 |
| ステータス | Implemented |
| 対象ディレクトリ | `hid-keyboard-g1/` |
| 対象ハードウェア | [kinT (kint41)](https://github.com/kinx-project/kint) = Kinesis Advantage コントローラ置換 + Teensy 4.1 (MIMXRT1062) |

本設計のスコープは **Goal 1: USB Full Speed + HID Report Protocol** である。実装は `hid-keyboard-g1/`。USB セットアップは Goal 2（HID Boot Protocol）を `hid-keyboard-g2/` で足しても組み直さない。USB High Speed は TinyGo `machine/usb` の 64 バイトエンドポイント前提に阻まれるため Goal 3 とし、ここでは扱わない。

---

## Overview

`hid-keyboard-g1/` は kinT + Teensy 4.1 向け TinyGo HID キーボードである。`machine.ConfigureUSBEndpoint` で CDC + HID Boot Keyboard を自前登録し、`keyboard.Port()` / `machine/usb/hid` は使わない。HID インタフェースは Boot subclass + Keyboard protocol を名乗り、入力は Report ID なしの 8 バイト。Goal 1 は常に Report Protocol として送り、`SET_PROTOCOL` は値を保持して ACK する。

キーマップは QMK `kinesis/keymaps/default_pretty`（my-customize）を `keymap.go` に持つ。マトリクスは行ドライブ Low、列プルアップ。CDC を残し、`tinygo flash` の 1200 bps リセットを維持する。

---

## Background & Motivation

### 現状

| パス | 役割 |
| --- | --- |
| `hid-keyboard-g1/` | Goal 1 ファーム（本設計） |
| `hid-keyboard-g2/` | Goal 2（Boot Protocol） |
| `tinygo-keyboard/` | `sago35/tinygo-keyboard` + Vial。内部は `k.Port()`。製品経路ではない |
| `usb-midi/` | 同じマトリクスの MIDI。ピン参照 |
| QMK `kinesis/kint41` | 実用ファーム。マトリクスとキーマップの参照 |

TinyGo は PR [#5704](https://github.com/tinygo-org/tinygo/pull/5704) を checkout してビルドする（[#5691](https://github.com/tinygo-org/tinygo/pull/5691) USB を含む）。Goal 1 は `machine.Flash` を呼ばない。

### #5691 が決めるハードウェア制約

`machine_mimxrt1062_usb.go` より:

- USB1 デバイスコントローラを TinyGo `machine/usb` に接続している。
- DMA 用 dQH / dTD / EP バッファは DTCM ではなく **4 KiB 非キャッシュ OCRAM**（リンカの `_usb_dma_start`）。USB DMA は DTCM を触れない。
- `PORTSC1.PFSC` で **Full Speed (12 Mbps) に固定**。`machine/usb` が 64 バイト EP 前提のため。
- ハードウェア EP バッファも 64 バイト。`SendUSBInPacket` は 64 バイト超を拒否する。
- VID:PID デフォルト `16C0:0483`（`board_teensy41.go`）。シリアルはチップユニーク ID。
- `NumberOfUSBEndpoints = 8`。
- バスリセット時に IN 転送中だった EP は `usbTxCancelled` 経由で `initEndpoint` 時に TxHandler を呼ぶ。TxHandler を登録しない実装では単に再送をメインループに任せる。

### TinyGo 標準 HID キーボードが Goal 2 と衝突する理由

`src/machine/usb/descriptor/hid.go` の `descriptor.CDCHID`:

- Interface 2: Class HID, **Subclass 0, Protocol 0**（Boot 非対応）。
- Report descriptor がキーボード（Report ID 2）+ マウス（ID 1）+ Consumer（ID 3）の複合。
- キーボード入力は **9 バイト**（Report ID + 8 バイト）。Boot Protocol は Report ID を禁止し、固定 8 バイトを要求する。
- `src/machine/usb/hid/hid.go` の `setupHandler` は `SET_IDLE` のみ。`SET_PROTOCOL` / `GET_PROTOCOL` / `GET_REPORT` / `SET_REPORT` は未処理（未処理は EP0 stall）。
- `keyboard.Port()` はパッケージ `init()` で `hid.SetHandler` → 上記ディスクリプタを登録する。import した時点で USB 構成が奪われる。

したがって Goal 1 から **`machine/usb/hid/keyboard` も `machine/usb/hid` も import しない**。

### 制約

- `tinygo-keyboard` は Vial・`k.Port()` 前提で、Boot Protocol を足すと USB 層を作り直す。
- QMK kint41 は HS + NKRO。TinyGo FS 64 バイト EP とは前提が違うので、マトリクスとキーマップだけ移植する。

---

## Goals & Non-Goals

### Goals（Goal 1）

1. kinT 実機で、QMK `default_pretty`（my-customize）相当の **単層 QWERTY** が macOS / Linux の HID キーボードとして入力できる。
2. USB **Full Speed**、HID **Report Protocol**。入力レポートは 8 バイト boot-compatible（6KRO + modifier）。
3. USB 登録は `machine.ConfigureUSBEndpoint`。`keyboard.Port()` 禁止。
4. `SET_PROTOCOL` / `GET_PROTOCOL` / `SET_IDLE` / `GET_IDLE` / `SET_REPORT` / `GET_REPORT` の **フックが Goal 1 で存在する**。プロトコル切替のレポート分岐は Goal 2。ただし 8 バイト固定なので Goal 1 の送信経路は Goal 2 でもそのまま使える。
5. TinyGo デフォルトの **USB CDC を共存**させる（デバッグと 1200 bps リセット）。
6. Caps/Num/Scroll/Compose LED を `SET_REPORT` および Interrupt OUT で点灯できる。

### Non-Goals

| 項目 | 扱い |
| --- | --- |
| HID Boot Protocol の BIOS/UEFI 実地検証 | Goal 2 |
| NKRO / ビットマップレポート | Goal 2 以降。内部状態は全キー分持つが、送信は 6KRO |
| USB High Speed (480 Mbps) | Goal 3。TinyGo `machine/usb` 側の変更が必要 |
| Vial / VIA / QMK 互換動的キーマップ | 対象外。`tinygo-keyboard/` に残す |
| マウス・Consumer・MIDI・マクロ・コンボ・タップダンス | 対象外 |
| キーマップの Flash 保存 | 対象外（`machine.Flash` は呼ばない。ビルド自体は README の #5704 を使う） |
| 複数レイヤ / `MO()` / ホールド | Goal 1 は単層。配列サイズだけ 2 レイヤ分確保 |
| JIS 専用キー（無変換 0x8B 等）のキーキャップ完全再現 | ホスト側レイアウトに任せる。ファームは QMK US HID usage |
| `-serial uart` | 非対応。`serial.usb` が無いと `initUSB()` 自体が走らない |
| Windows ホスト | Goal 1 非目標。`16C0:0483` の CDC+HID 複合は Windows の `Usbccgp` / `hidclass` 選択を混乱させうる。macOS / Linux のみ |

---

## Key Decisions

1. **`keyboard.Port()` / `machine/usb/hid` を使わない。**  
   `init()` 副作用で `descriptor.CDCHID` が入り、Boot 非対応かつ `SET_PROTOCOL` が stall する。Goal 2 で USB を書き直すことになる。

2. **`machine.ConfigureUSBEndpoint` で CDC + HID Boot Keyboard 複合ディスクリプタを自前登録する。**  
   パターンは `tinygo-keyboard/via.go` と同じ。グローバル `descriptor.CDCHID` は破壊しない（Vial 実装は in-place で `Configuration[2..4]` / `HID[3]` を書いており、他パッケージと衝突する）。HID class descriptor の `wDescriptorLength` は `descriptor.ClassHIDType{data: ...}` では書けない（`data` が非公開）。`descriptor.Append` した conf に対して `FindClassHIDType` でパッチする。

3. **HID インタフェースは Goal 1 から Boot subclass=1, protocol=1 を名乗る。**  
   Goal 2 で subclass を変えると OS のデバイス認識が変わる。Report Protocol がデフォルト（HID 1.11 Appendix F）なので、Goal 1 の「Report Protocol のみ話す」と矛盾しない。

4. **入力レポートは Report ID なし 8 バイト（modifier, reserved, 6 keys）。**  
   Boot Protocol と同じレイアウト。Goal 1 も Goal 2 も同じパケットを送れる。NKRO は将来 Report Protocol 側にだけ足す。

5. **6KRO。内部状態は押下ビットマップ。**  
   pack は usage `0x04`–`0xDF` を昇順に最大 6 個入れる（`0xE0`–`0xE7` は modifier ビット。`0x00`–`0x03` はキーマップに置かない）。7 キー目以降は HID ErrorRollOver (`0x01`) を 6 スロット全てに入れる（HID spec Appendix C）。modifier は 6KRO に数えない。

6. **CDC を残す。VID:PID は Teensy デフォルト `16C0:0483` を維持。**  
   `targets/teensy41.json` の `serial-port: ["16c0:0483"]` と 1200 bps リセットが切れると、アプリから HalfKay に入れずボタン操作が必要になる。製品用 PID（QMK の `1209:345C` や Teensyduino keyboard `16C0:04D2`）は Goal 2 以降の選択肢。

7. **キーマップは QMK `default_pretty`（my-customize、US HID usage）、単層。**  
   定数は `keycode.go`。文字の見え方はホスト OS のキーボードレイアウトで決まる。親指中央は左右 Command（`KC_LGUI` / `KC_RGUI`）。

8. **マトリクスは QMK COL2ROW と同じ極性: 行ドライブ Low、列プルアップ、押下は列 Low。**  
   QMK `quantum/matrix.c` の `DIODE_DIRECTION == COL2ROW` は `select_row()`（行を出力 Low）→ 列ピンを読む。`keyboard.json` の `"diode_direction": "COL2ROW"` は `hid-keyboard-g1/matrix.go` / `usb-midi/main.go` と一致する。残差は unselect だけ: QMK 既定は行を Hi-Z+pull-up（`MATRIX_UNSELECT_DRIVE_HIGH` は kint41 では未定義）、本ファームは非選択行を High のまま出力する。tinygo-keyboard の `InvertDiode(true)`（行 High / 列 pulldown）は使わない。  
   遅延は **unselect 後に DWT/CYCCNT busy-wait 20 µs**（`kint41.c` の `matrix_output_unselect_delay`。列が HIGH に戻るのを待つ）。select 後は短い settle（1–5 µs、同じく DWT）。スキャン内側で `time.Sleep` は使わない。

9. **デバウンスは 5 ms の per-key defer（連続一致）。スキャン周期 1 ms。**  
   QMK は `debounce: 50` + `sym_eager_pk` と非常に長い。Goal 1 は 5 ms を初期値にし、定数で変更可能にする。チャタリングが出たら上げる。

10. **USB クラスハンドラは Goal 1 で実装する（スタブではない ACK）。**  
    `SET_PROTOCOL` は値を保存して ZLP。送信フォーマットは 8 バイトのまま。`SET_REPORT` は LED に反映する（kinT にインジケータがあるため、Goal 1 でやる価値がある）。`GET_REPORT` は HID 必須リクエストなので last-sent 入力レポート／現在の LED を返す。

11. **ファームアクションは `0xF000` 以降。**  
    `KC_BOOTLOADER` は `EnterBootloader()`。`KC_MEH` / `KC_LCAG` / `KC_HYPR` は複数 modifier。`KC_DM_*` は Goal 1 では no-op。my-customize キーマップでは `[13,5]` は `KC_DM_RSTP`。`packBoot` は HID usage 以外をレポートに入れない。

12. **登録は `package main` の `init()`。`machine/usb/hid/keyboard` を import しない。**  
    TinyGo は `serial.usb` 時に `InitSerial` → `initUSB` → `EnableUSBCDC` → `USBDev.Configure`（Attach）を main より前に行う。ディスクリプタ置換が `main()` だと、ホストが CDC-only で列挙するレースがある。`init()` で `ConfigureUSBEndpoint` する（tinygo-keyboard `via.go` と同じ）。

---

## Proposed Design

### 全体像

```mermaid
flowchart TB
  subgraph hw [kinT hardware]
    Matrix["15 x 7 switch matrix"]
    LEDs["Caps/Num/Scroll/Compose LEDs"]
    PowerLED["D13 Teensy power LED"]
  end

  subgraph fw [hid-keyboard-g1 firmware]
    Scan["matrix scan + debounce"]
    KM["keymap layer0"]
    State["pressed bitmap + modifiers"]
    Pack["pack 8-byte boot report"]
    USBKB["usb.go: descriptors + setup/rx"]
    CDC["TinyGo USB CDC EP1/EP2"]
  end

  subgraph host [USB host FS 12 Mbps]
    OS["HID class driver Report Protocol"]
    Serial["CDC ACM debug / 1200bps reset"]
  end

  Matrix --> Scan --> KM --> State --> Pack
  Pack -->|"EP3 IN interrupt 8 B"| OS
  OS -->|"EP3 OUT or SET_REPORT LED"| USBKB --> LEDs
  USBKB -->|"GET/SET_PROTOCOL IDLE REPORT"| Pack
  CDC --> Serial
  PowerLED -.-> fw
```

### パッケージ / ファイル配置

実装ディレクトリは `hid-keyboard-g1/`。`tinygo flash --target teensy41 ./hid-keyboard-g1` が通るよう、**全て `package main`**。

```
hid-keyboard-g1/
  design.md      本設計
  main.go        メインループ、デバッグ印刷
  usb.go         複合ディスクリプタ、ConfigureUSBEndpoint、setup/rx
  hid.go         protocol/idle 状態、8 バイト pack、送信、GET_REPORT 用キャッシュ
  matrix.go      ピン、走査、デバウンス
  keycode.go     HID usage 定数とファームアクション
  keymap.go      QMK default_pretty (my-customize) の layer0
  led.go         インジケータ GPIO（active low）と power LED
```

`go.mod` の `github.com/sago35/tinygo-keyboard` は `tinygo-keyboard/` 用であり、`hid-keyboard-g1/` は依存しない。`tinygo-keyboard/` と `usb-midi/` は別 main。

### USB スタック配線（`keyboard.Port()` なし）

#### 初期化順

```mermaid
sequenceDiagram
  participant RT as TinyGo runtime
  participant M as machine.InitSerial
  participant CDC as machine/usb/cdc
  participant U as hid-keyboard-g1 init
  participant H as USB host

  RT->>M: serial.usb
  M->>CDC: EnableUSBCDC
  CDC->>CDC: ConfigureUSBEndpoint(descriptor.CDC, EP1/EP2)
  M->>M: USBDev.Configure / Attach
  Note over H: この時点では CDC-only ディスクリプタ<br/>列挙レースは短い
  U->>U: ConfigureUSBEndpoint(cdcHIDBoot, EP3, HID setup)
  Note over U: usbDescriptor を複合に置換<br/>endPoints に EP3 を追加<br/>CDC EP 登録はそのまま
  H->>M: GET_DESCRIPTOR device/config
  H->>M: SET_CONFIGURATION
  M->>M: initEndpoint(EP1,2,3)
  H->>U: HID class setup (SET_IDLE / GET_PROTOCOL ...)
```

`ConfigureUSBEndpoint`（`src/machine/usb.go`）は:

- `usbDescriptor = desc` で **置換**
- `endPoints` / handler 配列へ **追記**（CDC 分は消えない）
- `usbSetupHandler[Index]` は上書き

そのため HID 側は **EP3 だけ** 登録する。CDC の EP1/EP2 と `CDC_ACM_INTERFACE` ハンドラは `EnableUSBCDC` のものを使う。複合ディスクリプタの CDC 部分は、既に登録済みの EP 番号・インタフェース番号と一致させる。

#### インタフェース / エンドポイント割り当て

TinyGo 定数（`src/machine/usb/usb.go`）に合わせる。

| 番号 | 用途 | クラス | 備考 |
| --- | --- | --- | --- |
| IF 0 | CDC ACM | 0x02/0x02/0x01 | 既存 |
| IF 1 | CDC Data | 0x0A | 既存 |
| IF 2 | HID Keyboard | 0x03 / **0x01 Boot** / **0x01 Keyboard** | 本設計。標準 `InterfaceHID` は 0x00/0x00 なので **コピーして書き換える** |

| EP | 方向 | 転送 | wMaxPacket | bInterval | 用途 |
| --- | --- | --- | --- | --- | --- |
| 0 | IN/OUT | Control | 64 | — | 標準 + クラス setup |
| 1 | IN | Interrupt | 16 | 16 | CDC ACM 通知（既存） |
| 2 | OUT | Bulk | 64 | 0 | CDC OUT（既存） |
| 2 | IN | Bulk | 64 | 0 | CDC IN（既存） |
| 3 | IN | Interrupt | **8** | **1** | HID 入力レポート |
| 3 | OUT | Interrupt | **8** | **1** | HID LED 出力レポート |

Boot Keyboard の Interrupt IN は HID 1.11 Appendix B で **wMaxPacketSize = 8**。Goal 2 の BIOS がここを見る。ハードウェア dQH は `initEndpoint` が常に 64 バイトで切るが、送るのは 8 バイトだけなのでディスクリプタ上は 8 でよい。`bInterval = 1` は FS で 1 ms。QMK kint41 の `USB_POLLING_INTERVAL_MS 1` に合わせる（あちらは HS マイクロフレームの話なので、FS では 1 ms が最短）。

#### ディスクリプタ組み立て

`descriptor.CDCHID` を mutate しない。CDC 部分は TinyGo の部品を `descriptor.Append` でつなぎ、HID だけ差し替える。

`ClassHIDType.data` は非公開なので `descriptor.ClassHIDType{data: class}` はコンパイルできない。`ClassLength` は組み立て後の conf を `FindClassHIDType` で探す。グローバル `classHID` の既定 `wDescriptorLength` は `0x91`（マウス+キーボード+Consumer）であり、本 Boot レポートの長さではない。

`InterfaceHID.Bytes()` はプロセスグローバル配列への alias なので、subclass/protocol を書く前にコピーする。`Append` はコピーを返す。`EndpointIN`/`OUT` は EP 番号ごとのグローバルを書くが、本ディスクリプタは EP1 IN / EP2 IN / EP2 OUT / EP3 IN / EP3 OUT の 5 スロットで重複しない。同じ EP を 1 つの `[][]byte` リテラルで二度 `EndpointIN` すると alias するので、そうしない。

```go
// usb.go
func buildDescriptor() descriptor.Descriptor {
    iface := append([]byte(nil), descriptor.InterfaceHID.Bytes()...)
    iface[6] = 0x01 // bInterfaceSubClass: Boot
    iface[7] = 0x01 // bInterfaceProtocol: Keyboard

    report := bootKeyboardReportDescriptor() // Report ID なし

    conf := descriptor.Append([][]byte{
        descriptor.ConfigurationCDCHID.Bytes(), // bNumInterfaces=3 済み
        descriptor.InterfaceAssociationCDC.Bytes(),
        descriptor.InterfaceCDCControl.Bytes(),
        descriptor.ClassSpecificCDCHeader.Bytes(),
        descriptor.ClassSpecificCDCACM.Bytes(),
        descriptor.ClassSpecificCDCUnion.Bytes(),
        descriptor.ClassSpecificCDCCallManagement.Bytes(),
        descriptor.EndpointIN(descriptor.EndpointEP1, descriptor.TransferTypeInterrupt, 0x10, 0x10).Bytes(),
        descriptor.InterfaceCDCData.Bytes(),
        descriptor.EndpointOUT(descriptor.EndpointEP2, descriptor.TransferTypeBulk, 0x40, 0x00).Bytes(),
        descriptor.EndpointIN(descriptor.EndpointEP2, descriptor.TransferTypeBulk, 0x40, 0x00).Bytes(),
        iface,
        descriptor.ClassHID.Bytes(), // 長さは次の FindClassHIDType で上書き
        descriptor.EndpointIN(descriptor.EndpointEP3, descriptor.TransferTypeInterrupt, 8, 1).Bytes(),
        descriptor.EndpointOUT(descriptor.EndpointEP3, descriptor.TransferTypeInterrupt, 8, 1).Bytes(),
    })
    h, err := descriptor.FindClassHIDType(conf, descriptor.ClassHID.Bytes())
    if err != nil {
        // init 時点。ディスクリプタ組み立て失敗は起動不能
        panic(err)
    }
    h.ClassLength(uint16(len(report)))

    return descriptor.Descriptor{
        Device:        descriptor.DeviceCDC.Bytes(), // 0xEF/0x02/0x01 IAD。Configure() がこのグローバルを mutate する（TinyGo 既存動作）
        Configuration: conf,
        HID:           map[uint16][]byte{usb.HID_INTERFACE: report},
    }
}
```

`GET_DESCRIPTOR` Device 時に TinyGo が `usbDescriptor.Configure(vid, pid)` を呼び、VID/PID と `wTotalLength` を書き込む。`FindClassHIDType` は ClassHID の末尾 2 バイト（ClassLength）を除いて検索するので、パッチ後も再検索できる。

HID report descriptor のキー（`usbDescriptor.HID[wIndex]`）は **インタフェース番号 2**。`sendDescriptor` の `TypeHIDReport` 分岐が `setup.WIndex` で引く。

文字列:

```go
usb.Manufacturer = "kinT"
usb.Product      = "kinT TinyGo"
// Serial は空のまま → mimxrtSerialNumber()（チップ ID）
// VendorID / ProductID は 0 のまま → board_teensy41.go の 16C0:0483
```

#### HID レポートディスクリプタ（Boot Keyboard、Report ID なし）

HID 1.11 Appendix B.1 に沿う。LED は 5 ビット（Num/Caps/Scroll/Compose/Kana）。kinT の 4 LED + 未使用 Kana に一致する。TinyGo 標準の「LED 3 ビット + Report ID 2」は使わない。

```go
func bootKeyboardReportDescriptor() []byte {
    return descriptor.Append([][]byte{
        descriptor.HIDUsagePageGenericDesktop,
        descriptor.HIDUsageDesktopKeyboard,
        descriptor.HIDCollectionApplication,

        // Input: modifier 8 bits
        descriptor.HIDUsagePageKeyboard,
        descriptor.HIDUsageMinimum(224),
        descriptor.HIDUsageMaximum(231),
        descriptor.HIDLogicalMinimum(0),
        descriptor.HIDLogicalMaximum(1),
        descriptor.HIDReportSize(1),
        descriptor.HIDReportCount(8),
        descriptor.HIDInputDataVarAbs,

        // Input: reserved 8 bits
        descriptor.HIDReportCount(1),
        descriptor.HIDReportSize(8),
        descriptor.HIDInputConstVarAbs,

        // Output: 5 LED bits + 3 padding
        descriptor.HIDUsagePageLED,
        descriptor.HIDUsageMinimum(1),
        descriptor.HIDUsageMaximum(5),
        descriptor.HIDReportCount(5),
        descriptor.HIDReportSize(1),
        descriptor.HIDOutputDataVarAbs,
        descriptor.HIDReportCount(1),
        descriptor.HIDReportSize(3),
        descriptor.HIDOutputConstVarAbs,

        // Input: 6-key array
        descriptor.HIDUsagePageKeyboard,
        descriptor.HIDUsageMinimum(0),
        descriptor.HIDUsageMaximum(255),
        descriptor.HIDLogicalMinimum(0),
        descriptor.HIDLogicalMaximum(255),
        descriptor.HIDReportSize(8),
        descriptor.HIDReportCount(6),
        descriptor.HIDInputDataAryAbs,

        descriptor.HIDCollectionEnd,
    })
}
```

`FindClassHIDType` で書いた `ClassLength` は `len(report)` と一致していなければならない。

#### `ConfigureUSBEndpoint` 呼び出し

```go
func init() {
    usb.Manufacturer = "kinT"
    usb.Product = "kinT TinyGo"

    machine.ConfigureUSBEndpoint(buildDescriptor(),
        []usb.EndpointConfig{
            {
                Index:     usb.HID_ENDPOINT_OUT, // 3
                IsIn:      false,
                Type:      usb.ENDPOINT_TYPE_INTERRUPT,
                RxHandler: hidRxLEDs,
            },
            {
                Index: usb.HID_ENDPOINT_IN, // 3
                IsIn:  true,
                Type:  usb.ENDPOINT_TYPE_INTERRUPT,
                // TxHandler は置かない。メインループが最新スナップショットを送る
            },
        },
        []usb.SetupConfig{
            {Index: usb.HID_INTERFACE, Handler: hidSetup},
        },
    )
}
```

`main()` で `machine.USBDev.Configure` を再度呼ぶ必要はない（`initcomplete` なら no-op）。

### HID レポート形式

```
byte 0  modifier  bit0 LCtrl, 1 LShift, 2 LAlt, 3 LGUI,
                  bit4 RCtrl, 5 RShift, 6 RAlt, 7 RGUI
byte 1  reserved  0x00
byte 2..7  HID usage page 7 keycodes, 0x00 = empty
```

Goal 2 の Boot Protocol もこの 8 バイトである。Goal 1 の Report Protocol も同じディスクリプタなので同じ 8 バイトを送る。

将来 NKRO を足すとき:

- Boot (`protocol == 0`): 今の 8 バイトのまま
- Report (`protocol == 1`): Report ID 付きビットマップ等に変更 → **ディスクリプタ追加が必要**
- そのため `packReport` は `protocol` を見て分岐できるシグネチャにする。Goal 1 は常に 8 バイトの静的配列を返す（ヒープ `[]byte` を毎スキャン確保しない）。

```go
func packReport(pressed *[32]byte, mods uint8, protocol uint8) (out [8]byte) {
    // Goal 1: protocol を無視して boot 8 バイト
    return packBoot(pressed, mods)
}

func packBoot(pressed *[32]byte, mods uint8) (out [8]byte) {
    out[0] = mods
    n := 0
    // usage 0x04–0xDF を昇順。0xE0–0xE7 は mods 側。0x00–0x03 はキーマップに置かない。
    for usage := uint8(0x04); usage <= 0xDF; usage++ {
        if pressed[usage/8]&(1<<(usage%8)) == 0 {
            continue
        }
        n++
        if n > 6 {
            for i := 2; i < 8; i++ {
                out[i] = 0x01 // ErrorRollOver
            }
            return out
        }
        out[1+n] = usage // slots out[2]..out[7]
    }
    return out
}
```

6KRO 超過時は 6 スロットを全て `0x01` にする。部分的に 6 キーだけ送ると、ホストが「7 キー目が離された」と誤解する。

### コントロール転送（Goal 1 で配線、Goal 2 で意味が増える）

クラスリクエストは IF 2 宛。`handleEP0Setup` は standard 以外を `usbSetupHandler[WIndex]` に渡す。`BmRequestType` は TinyGo 定数で:

- Host→Device class interface = `0x21`（`usb.SET_REPORT_TYPE` / `REQUEST_HOSTTODEVICE_CLASS_INTERFACE`）
- Device→Host class interface = `0xA1`（`REQUEST_DEVICETOHOST_CLASS_INTERFACE`）

`bRequest`（`src/machine/usb/usb.go`）と `wValue` / `wLength`:

| bRequest | 値 | wValueH | wValueL | wLength | Goal 1 |
| --- | --- | --- | --- | --- | --- |
| GET_REPORT | 1 | 1 Input / 2 Output / 3 Feature | Report ID（Boot は 0） | ホスト要求長 | Input: `hidLastIN` の先頭 `min(8, WLength)` バイト。Output: LED 1 バイト。Feature または未知 type / 非 0 の Report ID は `false`（EP0 stall）。 |
| GET_IDLE | 2 | — | Report ID（無視して 1 つの idle） | 1 | `hidIdle` 1 バイト |
| GET_PROTOCOL | 3 | 0 | 0 | 1 | `hidProtocol`（0 boot / 1 report） |
| SET_REPORT | 9 | 2 Output（Feature は未使用） | Report ID 0 | 0 または 1+ | `WLength == 0`: ZLP して true（`ReceiveUSBControlPacket` は呼ばない。IRQ 内で `usbSpinLimit` スピンして stall する）。Output: データステージ先頭バイトを `applyLEDs`。それ以外の type は stall。 |
| SET_IDLE | 10 | duration（4 ms 単位） | Report ID | 0 | `hidIdle = WValueH`。0 = 変化時のみ。ZLP |
| SET_PROTOCOL | 11 | 0 | 0 boot / 1 report | 0 | `hidProtocol = WValueL`。ZLP。**Goal 1 は送信フォーマットを変えない** |

```go
const (
    protocolBoot   = 0
    protocolReport = 1
)

var (
    hidProtocol uint8 = protocolReport // 電源投入後は Report（HID Appendix F）
    hidIdle     uint8 = 0              // 変化時のみ。Goal 2 で Boot 時 default 125 を検討
    hidLEDs     uint8
    hidLastIN   [8]byte // GET_REPORT / 変化検出用。IRQ と main で共有
)

func hidSetup(setup usb.Setup) bool {
    switch setup.BmRequestType {
    case usb.REQUEST_HOSTTODEVICE_CLASS_INTERFACE:
        switch setup.BRequest {
        case usb.SET_IDLE:
            hidIdle = setup.WValueH
            machine.SendZlp()
            return true
        case usb.SET_PROTOCOL:
            hidProtocol = setup.WValueL
            machine.SendZlp()
            return true
        case usb.SET_REPORT:
            if setup.WLength == 0 {
                machine.SendZlp()
                return true
            }
            if setup.WValueH != 2 { // Output
                return false
            }
            b, err := machine.ReceiveUSBControlPacket()
            if err != nil {
                return false
            }
            applyLEDs(b[0])
            machine.SendZlp()
            return true
        }
    case usb.REQUEST_DEVICETOHOST_CLASS_INTERFACE:
        switch setup.BRequest {
        case usb.GET_IDLE:
            return machine.SendUSBInPacket(0, []byte{hidIdle})
        case usb.GET_PROTOCOL:
            return machine.SendUSBInPacket(0, []byte{hidProtocol})
        case usb.GET_REPORT:
            if setup.WValueL != 0 { // boot は report ID 0
                return false
            }
            switch setup.WValueH {
            case 1: // Input
                var local [8]byte
                state := interrupt.Disable()
                local = hidLastIN
                interrupt.Restore(state)
                n := int(setup.WLength)
                if n > 8 {
                    n = 8
                }
                return machine.SendUSBInPacket(0, local[:n])
            case 2: // Output
                return machine.SendUSBInPacket(0, []byte{hidLEDs})
            default:
                return false // Feature 等は stall
            }
        }
    }
    return false
}

// Interrupt OUT。mimxrt は dQH を 64 バイトで prime する。descriptor の wMaxPacket=8 でも
// ハンドラは 1..64 を許容する。len==0 は無視。アロケーション禁止。
func hidRxLEDs(b []byte) {
    if len(b) == 0 {
        return
    }
    applyLEDs(b[0])
}
```

注意:

- `ReceiveUSBControlPacket` の戻りは `[cdcLineInfoSize]byte`（**7 バイト**）。LED 1 バイトには足りる。`WLength == 0` で呼んではいけない（USB IRQ 内で最大 5e6 スピン → stall）。
- TinyGo `sendDescriptor` は HID **class** descriptor（type `0x21`）の単独 GET_DESCRIPTOR を扱わず ZLP を返す。通常ホストは configuration 内の HID descriptor を使う。BIOS が type 0x21 を要求したら Goal 2 で TinyGo 側修正が必要（本ファームの class setupHandler には standard GET_DESCRIPTOR は来ない）。
- Interrupt OUT と SET_REPORT は同じ `applyLEDs` に集約する。
- **USB リセット後の protocol/idle:** HID 7.2.6 はリセット後 Report Protocol を要求する。TinyGo `handleUSBBusReset` はアドレスと `usbConfiguration`（非公開）だけを 0 にし、アプリ状態は残る。`InitEndpointComplete` もリセットで false に戻らない。アプリから `usbConfiguration` は読めない。Goal 1 はこれを **既知の TinyGo ギャップ** とし、protocol/idle をバスリセットで戻さない。Goal 2 の BIOS が SET_PROTOCOL(0) のあとリセットして GET_PROTOCOL が Boot のままなら、TinyGo にリセットフックを足す。

### マトリクススキャン

#### ピン（QMK `kint41/keyboard.json` = 現行 TinyGo コード）

```go
var rows = []machine.Pin{ // 15
    machine.D8, machine.D9, machine.D10, machine.D11, machine.D7,
    machine.D16, machine.D5, machine.D3, machine.D4, machine.D1,
    machine.D0, machine.D2, machine.D17, machine.D23, machine.D21,
}
var cols = []machine.Pin{ // 7
    machine.D18, machine.D14, machine.D15,
    machine.D20, machine.D22, machine.D19, machine.D6,
}
```

QMK の `LINE_PINx` は Teensy デジタルピン番号 `Dx` に一致する。

#### ダイオードとドライブ方向

QMK JSON は `"diode_direction": "COL2ROW"`。QMK `quantum/matrix.c` でこれは **行を出力 Low（`select_row`）、列を入力プルアップで読む**。`ROW2COL` が列ドライブである。本ファームと同じ極性。

| 実装 | 選択 | 非選択 | 入力 | 押下 |
| --- | --- | --- | --- | --- |
| QMK COL2ROW（`matrix.c`） | 行 Low | 行 Hi-Z+pull-up（kint41 は `MATRIX_UNSELECT_DRIVE_HIGH` 未定義） | 列 Pullup | 列 Low |
| `hid-keyboard-g1/matrix.go`, `usb-midi/main.go` | 行 Low | 行 High（出力のまま） | 列 Pullup | `!col.Get()` |
| `tinygo-keyboard` + `InvertDiode(true)` | 行 High | 行を pulldown 入力に戻す | 列 Pulldown | `col.Get()` |

Goal 1 は **行 Low 選択、非選択行は High 出力、列 Pullup、押下は Low**。QMK との差は unselect が Hi-Z ではなく drive-High なことだけ。drive-High の方が列の立ち上がりは速いが、kint41 が測った 20 µs は Hi-Z 復帰向けなので、同じ 20 µs を unselect 後に置く（過剰でも害は少ない）。tinygo-keyboard の反転スキームは使わない。

ダイオードはスイッチと直列なので、同時押しでもゴーストはハードウェアが防ぐ。ソフトウェア側のアンチゴーストは実装しない。未実装交差点は常にオープン。キーマップは `KC_NO`。

#### セトル / unselect 遅延（DWT busy-wait。`time.Sleep` 禁止）

`kint41.c` がオーバーライドしているのは **`matrix_output_unselect_delay`**（行を離したあと「wait for all Col signals to go HIGH」）である。QMK COL2ROW の順は:

1. `select_row`（行 Low）
2. `matrix_output_select_delay()` — kint41 は `GPIO_INPUT_PIN_DELAY 0` なのでほぼ 0
3. 列を読む
4. `unselect_row`
5. **`matrix_output_unselect_delay` = 20 µs**（5 µs では不足、一部 KB600LF+stapelberg+teensy41 は 20 µs）

20 µs を `row.Low()` の直後に置くのは kint41 の計測対象ではない。Goal 1 は **unselect（`row.High()`）のあとに 20 µs**、select 後に短い settle（**5 µs**。高速 GPIO なら 1 µs でも足りうるが定数は 5）。

`time.Sleep(20 * time.Microsecond)` は使わない。TinyGo `time.Sleep` は協調スケジューラ（`addSleepTask` + `task.Pause`）。mimxrt の `sleepTicks` は PIT 24 MHz で `(last-curr)/pitCyclesPerMicro` するため、サブ〜24 µs は 0 カウントになり、`hasScheduler` なら `timerSleep` がすぐ戻る。QMK が DWT `CYCCNT` を使った理由と同じ。TinyGo は `initSysTick` で DWT を既に有効化している（`runtime_mimxrt1062_time.go`）。

TinyGo `Pin.Get` は GPIO6–9（AHB、`getGPIO()` + `initPins` GPR26–29）を使う。遅延の理由はレジスタ速度ではなくマトリクスのアナログ RC である。

```go
// matrix.go。Teensy 4.1 TinyGo CORE_FREQ = 600 MHz。
const cyclesPerUs = 600

var dwtCYCCNT = (*volatile.Register32)(unsafe.Pointer(uintptr(0xe0001004)))

func delayUs(us uint32) {
    start := dwtCYCCNT.Get()
    cycles := us * cyclesPerUs
    for dwtCYCCNT.Get()-start < cycles { // uint32 ラップ回り込み
    }
}
```

スキャン内側でスケジューラに yield しない。1 ms 周期の残り待ちだけ `time.Sleep` してよい。15 行 × (5+20) µs = 375 µs で 1 ms 周期に収まる。

#### デバウンス

```
scan interval     = 1 ms
debounce ticks    = 5     // 5 ms
algorithm         = per-key defer
```

各交差点にカウンタを持つ。生値が `debounced` と違う間だけカウントし、5 連続で一致したら遷移。QMK の eager（押下即登録、その後 50 ms 無視）より遅延は増えるが誤入力は減る。定数 `debounceTicks` で変更する。

ゴースト: ダイオード付きなので「3 点で 4 点目が立つ」現象は起きない。空セルを読んでもプルアップのまま。

#### スキャン疑似コード

```go
func (m *Matrix) scanOnce() {
    for r, row := range rows {
        row.Low()
        delayUs(5) // select settle
        for c, col := range cols {
            raw := !col.Get()
            m.debounce(r, c, raw)
        }
        row.High() // unselect: drive High（QMK 既定の Hi-Z ではない）
        delayUs(20) // unselect delay: 列が HIGH に戻るのを待つ
    }
}
```

起動時: 行は Output + High、列は InputPullup。`kint41.c` に倣い power LED `D13` を Output High。

### キーマップ

#### 表現

HID usage とファームアクションは `keycode.go` の `KC_*` 定数。layer0 の 105 セルは `keymap.go` が正。`machine/usb/hid/keyboard` は import しない。modifier `KC_LCTL`–`KC_RGUI` は pack 時にビットへ折り畳む。`KC_MEH` / `KC_LCAG` / `KC_HYPR` は複数 modifier ビット。`packBoot` は HID usage 以外をレポートに出さない。

Goal 1 は layer 0 のみ読む。layer 1 は全 `KC_NO`。

#### QMK LAYOUT → (row, col)

1. `keyboards/kinesis/kint41/keyboard.json` の `layouts.LAYOUT.layout` を配列順に辿る。`matrix: [row, col]` がセル。
2. 同じ順で `keyboards/kinesis/keymaps/default_pretty/keymap.c` の `LAYOUT(...)` を割り当てる。ASCII アートは使わない。
3. インデックスは `row*7+col`。正本は `keymap.go`。

物理キー 86。親指クラスタ（my-customize）:

| (r,c) | キー |
| --- | --- |
| 5,6 | LCtl |
| 5,5 | LCAG（Ctrl+Alt+GUI） |
| 9,6 | PgDn |
| 8,5 | PgUp |
| 3,5 | Esc |
| 8,6 | Bspc |
| 3,6 | Space |
| 4,6 | LGUI（左 Command） |
| 2,5 | LAlt |
| 6,6 | MEH（Ctrl+Shift+Alt） |
| 7,5 | RGUI（右 Command） |
| 6,5 | Enter |

ファンクション行の `KC_DM_*` は Goal 1 では no-op。`KC_BOOTLOADER` 定数はあるが、このキーマップでは割り当てていない。

#### US vs JP

ファームウェアが送るのは HID usage であり、ASCII ではない。QMK default の `KC_EQL` / `KC_QUOT` / `KC_GRV` は usage `0x2E` / `0x34` / `0x35` で、`tinygo-keyboard` の `KeyHat` / `KeyColon` / `KeyHankaku` と同じ値である。ホストを JIS 配列にすれば同じ物理キーが `^` / `:` / 半角になる。Goal 1 は QMK テーブルをそのまま送り、JIS 専用 usage（`0x87`–`0x8B` 等）は置かない。

### LED

QMK `indicators`（`on_state: 0` = active low）:

| 機能 | ピン | HID LED bit |
| --- | --- | --- |
| Caps Lock | D12 | bit 1 |
| Num Lock | D26 | bit 0 |
| Scroll Lock | D25 | bit 2 |
| Compose | D24 | bit 3 |
| Kana | （無し） | bit 4 は無視 |
| Power | D13 | USB とは無関係。起動時 High |

`applyLEDs(v uint8)` がビットをピンに書く。USB IRQ から呼ばれるので GPIO 書き込みだけ（`Pin.Set` は mimxrt で `DR_SET`/`DR_CLEAR`）。アロケーション禁止。`hidLEDs = v` を保持して GET_REPORT Output に使う。

### メインループ

```mermaid
flowchart LR
  A["every 1 ms"] --> B["scan: select 5us, unselect 20us DWT"]
  B --> C["debounce"]
  C --> D["keymap: usage or KC_BOOTLOADER"]
  D --> E["update bitmap + mods"]
  E --> F{"pack 8-byte"}
  F --> G{"changed or idle elapsed?"}
  G -->|yes| H["SendUSBInPacket EP3"]
  H -->|busy| I["retry next tick with latest"]
  H -->|ok| J["hidLastIN = report"]
  G -->|no| A
```

送信方針:

- **変化時に送る**（Goal 1 の電源投入 default は `hidIdle == 0`）。HID 7.2.4 はキーボードに 500 ms（duration 125）を推奨するが、macOS/Linux のキーリピートは OS 側なので 0 が正しい。Goal 2 の BIOS は idle 再送に頼ることがある。SET_IDLE は実装するので BIOS が設定すれば効く。Goal 2 で Boot 時 default を 125 に変えてよい（送信経路はそのまま）。
- `hidIdle > 0` なら最後の成功送信から `hidIdle * 4 ms` 経過で再送。
- `SendUSBInPacket` が false なら次の 1 ms で **最新** スナップショットを再試行。キューしない。
- 全キー離したときは 8 ゼロを必ず 1 回送る。
- pack 結果は静的 `[8]byte`（`hidCurrent`）。成功送信後、`interrupt.Disable` 下で `hidLastIN = hidCurrent`。GET_REPORT は同様にローカル `[8]byte` へコピーしてから EP0 送信。スキャン／送信パスで `make` / ヒープ `[]byte` を出さない。

`KC_BOOTLOADER` はエッジ（押した瞬間）で `machine.EnterBootloader()`（HalfKay、`bkpt #251`）。その前に空レポートを送れると望ましいが、ブートローダ突入後は USB が切れる。

デバッグ:

```go
const debug = false
```

- `debug == true` のとき押下/離鍵で CDC に `r,c idx usage` を `println`。キーログをシリアルに残さないため既定は `false`。

周期: スキャン ~0.4 ms + 処理。1 ms 周期の残りだけ `time.Sleep`。スキャン内側では DWT。USB bInterval=1 ms に対し 1–2 ms に 1 レポートなら十分。

期待負荷: キー入力バーストでも EP3 は 8 バイト / 1 ms = 8 KB/s。FS 12 Mbps に対して無視できる。CPU は 600 MHz 相当の Cortex-M7 でマトリクス走査は問題にならない。

### Goal 2 / Goal 3

Goal 2（HID Boot Protocol）は `hid-keyboard-g2/design.md`。Goal 1 が渡すもの: Boot subclass、8 バイト Report ID なし、SET/GET_PROTOCOL の ACK、`packReport(protocol)` シグネチャ。

Goal 3（USB HS）は TinyGo `EndpointPacketSize` / `PFSC` の変更が先。送信は `machine.SendUSBInPacket` のみ。

---

## API / Interface Changes

外部 Go API は無い（`package main` ファームウェア）。ホストから見た USB インタフェースは次のとおり。

ホストから見た USB:

- IF0/1 CDC
- IF2 HID **Boot Keyboard**
- キーボード入力 **8 バイト**、マウス/consumer なし
- クラスリクエスト: GET/SET_REPORT, GET/SET_IDLE, GET/SET_PROTOCOL
- キー入力: `keymap.go` の HID usage

`tinygo-keyboard/` とは USB 構成が異なる。同じデバイスに両方を同時に焼かない。

---

## Data Model Changes

永続化なし。RAM のみ。

| データ | サイズ | 備考 |
| --- | --- | --- |
| 生マトリクス | 15×7 bool | スキャン結果 |
| デバウンスカウンタ | 105 uint8 | 0..debounceTicks |
| 安定押下 | 105 bool | |
| キーマップ | 2×105 `Keycode` (`uint16`) | Flash に書かない |
| 押下ビットマップ | 32 バイト（256 usage） | pack 用。静的 |
| `hidCurrent` / `hidLastIN` | 各 `[8]byte` | 静的。8 バイトコピーは `interrupt.Disable` 下。uint8 の protocol/idle/leds は Cortex-M7 でアトミック |

マイグレーション不要。QMK EEPROM / Vial 定義とは非互換。

---

## Alternatives Considered

### A. `keyboard.Port()` を使う

- 利点: TinyGo 標準の文字入力ヘルパが使える。
- 欠点: Report ID 付き 9 バイト、subclass 0、`SET_PROTOCOL` stall。Goal 2 で USB 層を破棄することになる。**不採用**。

### B. `sago35/tinygo-keyboard` を製品にする

- 利点: マトリクス・レイヤ・Vial が既にある。本リポジトリに JP 寄りの `tinygo-keyboard/` がある。
- 欠点: 内部が `k.Port()` + マウス。Vial 用 vendor HID を EP6/7 に足すだけで、Boot Protocol は扱わない。Goal 1 の「手でディスクリプタを持つ」方針と逆。参照実装としては残す。**製品経路としては不採用**。

### C. 本設計（自前 ConfigureUSBEndpoint + 8 バイト boot-compatible）

- 利点: Goal 2 が setupHandler の分岐追加で済む。CDC 共存。QMK キーマップを素直に写せる。
- 欠点: TinyGo 標準 HID の文字入力ヘルパ（UTF-8 → キーコード）が使えない。ファームは usage を直接持つので問題にならない。
- **採用**。

### D. Goal 1 から NKRO ビットマップ

- 利点: 親指クラスタ + 修飾の同時押しで 6KRO を超えても死なない。QMK kint41 は `nkro: true`。
- 欠点: Report ID または非 boot ディスクリプタが必要で、Boot Protocol と二重管理になる。Goal 1 の同時押し実使用は modifier + 数キーで 6KRO に収まる。内部ビットマップだけ先に持つ。**Goal 1 では不採用、Goal 2 オプション**。

### E. HID-only（CDC を落とす）

- 利点: インタフェースがキーボードだけ。一部 BIOS が複合デバイスを嫌う可能性。
- 欠点: `tinygo flash` のシリアルポート検出と 1200 bps リセットが死ぬ。デバッグ印刷が UART ピンに逃げる。`-serial uart` では `initUSB` が走らず HID も死ぬ。**不採用**。Goal 2 で BIOS が複合を拒否したら、そのとき HID-only ディスクリプタを切替可能にしておく（EP 登録は HID だけ自前、CDC は EnableCDC 済みでもディスクリプタから消せばホストは使わない）。Goal 1 では切替コードを書かない。

### F. VID/PID を QMK と同じ `1209:345C`、または Teensyduino keyboard `16C0:04D2` にする

- 利点: OS がキーボードとして分類しやすい。QMK と入れ替えても同じ ID。
- 欠点: TinyGo の `serial-port 16c0:0483` と不一致。複合 CDC+HID で PID をキーボード用にすると、逆にシリアルツールが迷う。**Goal 1 は 16C0:0483 維持**。変えるなら flash 手順（ボタン必須化 or ターゲット JSON 更新）とセット。

---

## Security & Privacy Considerations

| 脅威 | 深刻度 | 緩和 |
| --- | --- | --- |
| 本デバイスは HID キーボードであり、任意キーをホストに注入できる | 仕様 | 物理アクセス前提。無線なし |
| CDC デバッグがキーログになる | 中 | 既定 `debug = false` |
| `KC_BOOTLOADER` で HalfKay に落ち、任意ファームを焼ける | 中 | 物理キーが必要。QMK も `QK_BOOT` がある |
| USB 複合デバイスとしての不正 descriptor | 低 | 固定 ROM。動的書き換えなし |
| Flash キーマップ改ざん | — | Goal 1 は永続化しない |
| リモート HID 攻撃 | — | バスパワー有線のみ |

認証・ペアリングはない。ホストの HID スタックを信頼する。

`ReceiveUSBControlPacket` は 7 バイト固定で、EP0 の過大 SET_REPORT は切り捨て。Goal 1 の LED 1 バイトには十分。

---

## Observability

メトリクスサーバは無い。観測手段:

1. **USB CDC**（`println`）  
   - 起動時に VID/PID、ディスクリプタ種別を 1 行。  
   - `debug` 時のみマトリクス座標。  
   - `SendUSBInPacket` 連続失敗回数（USB 未 configure / EP busy）。
2. **インジケータ LED**  
   - Caps 等がホストと同期していることが HID 往復の可視確認。  
   - Power D13 は「ファームが main に入った」。
3. **ホスト側**  
   - macOS: `system_profiler SPUSBDataType`、`ioreg -p IOUSB -l -w 0`  
   - Linux: `lsusb -v -d 16c0:0483`、`usbhid-dump`、`dmesg`  
   - レポート中身: `hidapitester --vidpid 16C0:0483 --open --read-input`

アラートは無し。入力不能時の切り分けは「CDC が残っているか（USB スタック生存）」「HID インタフェースがあるか（ディスクリプタ）」「マトリクス debug が反応するか（GPIO）」。

QMK kint41 は `-DCORTEX_ENABLE_WFI_IDLE=FALSE`（`rules.mk`、kinx-project/kint#77）。TinyGo mimxrt は `runtime_mimxrt1062_time.go` で `wfi` がコメントアウト（`TODO: causes hardfault`）し、`runtime_mimxrt1062_clock.go` は run mode にして WFI が SoC をロックしないようにしている。idle の `time.Sleep` は PIT/`ticks()` をビジーポールする。WFI を再有効化しない。キー欠けが出たら `SendUSBInPacket` busy、unselect 遅延、`println` によるスケジューラ飢餓を見る。

---

## Rollout Plan

サービス段階リリースではない。フラッシュ単位。

1. Teensy プログラムボタンで HalfKay に入れる（失敗時の逃げ）。`KC_BOOTLOADER` が動けばボタンなし。
2. README どおり TinyGo #5704（#5691 USB を含む）で `./hid-keyboard-g1` を flash。
3. ホストで CDC が出ること、続けて HID Boot Keyboard が出ることを確認。
4. 既知フレーズ（`asdf`、親指 Backspace、修飾+キー）を入力。
5. Caps Lock で D12 が反転すること。
6. 問題があれば QMK kint41 を焼き戻す。

フラグ:

- `debug` const（既定 `false`）。
- `debounceTicks`。
- VID/PID は Goal 1 では触らない。

ロールバック: HalfKay + 直前の `.hex`（QMK または本ファーム）。PID を変えていなければ `tinygo flash` の自動リセットも生きる。

---

## Testing / 検証

対象: kinT + Teensy 4.1、ホスト **macOS と Linux**。Windows（`Usbccgp` / `hidclass`）は Goal 1 非目標。

### USB 列挙

Linux 例:

```
lsusb -d 16c0:0483
lsusb -v -d 16c0:0483 | sed -n '/bInterfaceClass/,/bInterval/p'
```

確認項目:

- IAD 付き CDC（IF0/1）
- IF2: `bInterfaceClass 3`, `bInterfaceSubClass 1`, `bInterfaceProtocol 1`
- EP `83` Interrupt IN wMaxPacketSize=8, bInterval=1
- EP `03` Interrupt OUT 8 バイト
- iProduct `kinT TinyGo`

macOS:

```
system_profiler SPUSBDataType | grep -A20 'kinT'
```

### HID レポート

```
hidapitester --vidpid 16C0:0483 --list
hidapitester --vidpid 16C0:0483 --open --read-input
```

キー 1 つで 8 バイト（例: `A` 押下 `00 00 04 00 00 00 00 00`、離鍵は全 0）。Report ID 先頭が付いていないこと（9 バイトになっていないこと）が `keyboard.Port()` を使っていない証拠。

### レイアウト

ホストを **US 配列** にして `keymap.go` どおりタイプできること:

- 左上 `=`、数字行、QWERTY
- 親指 Space / Enter、左右 Command（`[4,6]` / `[7,5]`）

ホストを JIS にすると同じ usage が JIS 文字になる。それは成功（ファームは usage しか送っていない）。

### Boot vs Report（Goal 1 では記録のみ）

Linux:

```
# GET_PROTOCOL は usbhid 経由では見えにくい。usbmon / Wireshark USBPcap で
# bmRequestType=0x21 bRequest=0x0B (SET_PROTOCOL) を観察
```

Goal 1 合格条件: SET_PROTOCOL が stall しない（ZLP で ACK）。値によるレポート形状変化は Goal 2。

### LED

Caps Lock トグルで D12 が反転。Num/Scroll はホストがビットを送れば D26/D25。Compose はホスト依存。

### 回帰

- `tinygo flash --target teensy41 --size short --stack-size 8kb ./hid-keyboard-g1` が 1200 bps リセットで成功する（CDC 生存）。
- `tinygo-keyboard/` は本設計の変更で壊さない（別 main）。

---

## Risks

| ID | リスク | 深刻度 | 緩和 |
| --- | --- | --- | --- |
| R1 | `keyboard.Port()` を誤 import し `init()` が CDCHID を再登録 | 高 | `hid-keyboard-g1` から `machine/usb/hid` を import しない |
| R2 | `ConfigureUSBEndpoint` がディスクリプタを置換するため、CDC 部分を落とすとシリアル消失 | 高 | 複合ディスクリプタを CDC 部品から組む。flash 後に tty が見えること |
| R3 | Attach が machine init で先行し、一瞬 CDC-only 列挙 | 低 | アプリ `init()` で即座に置換。実害が出たら Detach → 置換 → Attach |
| R4 | TinyGo が HID class GET_DESCRIPTOR (0x21) を ZLP | 低（OS）/ 中（BIOS） | Goal 1 の OS は config 内 HID desc で足りる。Goal 2 で必要なら TinyGo パッチ |
| R5 | `ReceiveUSBControlPacket` が 7 バイト固定 | 低 | LED 1 バイト。NKRO feature report を足すときは TinyGo 側拡張 |
| R6 | FS 強制 + 64 バイト EP。HS 化できない | 情報 | Goal 3。QMK 比でポーリングは 1 ms が下限 |
| R7 | 6KRO 超過 | 低 | ErrorRollOver。内部ビットマップがあるので Goal 2 で NKRO 可 |
| R8 | unselect が QMK Hi-Z ではなく drive-High | 低 | 同じ COL2ROW 極性。ゴーストや隣接誤検出が出たら unselect 20 µs を延長 |
| R9 | QMK debounce 50 ms を無視して 5 ms にしたことによるチャタ | 低 | 定数化。症状が出たら上げる |
| R10 | キー欠け | 中 | TinyGo mimxrt は WFI を使っていない（コメントアウト + run mode）。WFI を再有効化しない。欠けたら `SendUSBInPacket` busy / unselect 遅延 / `println` 飢餓を疑う |
| R11 | `SendUSBInPacket` が busy で最新以外を失う | 低 | スナップショット再送。押しっぱなしは次の idle または変化で回復 |
| R12 | グローバル HID/Interface/Device 配列 | 低 | `InterfaceHID` はコピーしてから subclass を書く。`ClassLength` は `FindClassHIDType` で組み立て後の conf をパッチ。`DeviceCDC.Bytes()` は GET_DESCRIPTOR Device 時に TinyGo `Configure()` が mutate する（既存動作）。EndpointIN は本ディスクリプタでは 5 スロットが衝突しない |
| R13 | 親指キーマップを取り違える | 中 | 正本は `keymap.go`（my-customize `default_pretty`） |
| R14 | USB 無し TinyGo でビルド | 高 | README どおり #5704 を checkout（#5691 USB を含む）。`machine_mimxrt1062_usb.go` の存在を確認。Goal 1 は `machine.Flash` を呼ばない |
| R15 | `SET_REPORT` で `WLength==0` のとき `ReceiveUSBControlPacket` が USB IRQ をブロック | 高 | `WLength==0` では呼ばず ZLP |
| R16 | `hidLastIN` の 8 バイトが IRQ と main で撕裂 | 低 | `interrupt.Disable` 下で `[8]byte` コピー。pack は静的配列 |

---

## Open Questions

なし。`usb.Product` は `kinT TinyGo`。キーマップは `keymap.go`。

---

## References

- TinyGo PR [#5691](https://github.com/tinygo-org/tinygo/pull/5691) USB device (CDC) for Teensy 4.x
- TinyGo PR [#5704](https://github.com/tinygo-org/tinygo/pull/5704) flash（#5691 にブロック、USB コミットを含む。README のビルド経路。Goal 1 は `machine.Flash` 未使用）
- ローカル TinyGo: `src/machine/machine_mimxrt1062_usb.go`, `src/machine/usb.go`, `src/machine/usb/descriptor/hid.go`（`FindClassHIDType`）, `src/machine/usb/hid/keyboard/keyboard.go`, `src/machine/board_teensy41.go`, `src/runtime/runtime_mimxrt1062_time.go`
- `github.com/sago35/tinygo-keyboard` `via.go`（`ConfigureUSBEndpoint` の先例。Vial 用で Boot 非対応。CDCHID を in-place 改変）
- QMK `keyboards/kinesis/kint41/{keyboard.json,kint41.c,config.h,rules.mk}`
- QMK キーマップ正本: `keyboards/kinesis/keymaps/default_pretty/keymap.c`（my-customize）
- QMK `quantum/matrix.c` COL2ROW = `select_row`
- 本リポジトリ `hid-keyboard-g1/`, `tinygo-keyboard/`, `usb-midi/`
- [USB HID 1.11](https://www.usb.org/sites/default/files/documents/hid1_11.pdf) §7.2 クラスリクエスト, Appendix B Boot Interface, Appendix C keyboard
- kinT ハードウェア: https://github.com/kinx-project/kint

---

## Implementation status

Goal 1 は `hid-keyboard-g1/` に実装済み。Goal 2 は `hid-keyboard-g2/`。
