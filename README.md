# test-tinygo-kint

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
$(ghq root)/github.com/tinygo-org/tinygo/build/tinygo flash --target teensy41 --size short --stack-size 8kb ./hid_keyboard
```
