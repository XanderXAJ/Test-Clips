# Test Clips

Video clips for use in testing video and media encoders.

Currently this project only works with the SVT AV1 encoder.

## Usage

To see the available options:

```shell
go run . -h
```

To perform one conversion, with some (example) options specified:

```shell
go run . -i HaruhiSchool.mkv -crf 10 -p 12 -fg 0
```

Pro-tip: If you want to perform multiple conversions with one of the options changing, use your shell's `for` loop (`bash` shown):

```bash
for p in {4..12}; do go run . -i FFXVForwardVistas.m2ts -crf 30 -p $p -fg 0; done
```

Or in a matrix (`bash` shown) -- this example tests presets 4-12, every fourth CRF (i.e. 10,14, 18, etc.) and a hand-picked selection of film grain values:

```bash
for crf in {10..63..4}; do for p in {4..12}; do for fg in 0 1 8 15; do go run . -o output -i HaruhiSchool.mkv -crf $crf -p $p -fg $fg; done; done; done
```
