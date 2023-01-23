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

## Failed conversions

If processing fails or an interrupt signal (e.g. Ctrl+C on the CLI) is received, files will be renamed to allow easy identification of failed attempts and easy retries.

For example, if an error occurs during video conversion (or it's interrupted), these files:

```
video.mkv
video.mkv.log
```

Will automatically be renamed to:

```
video.failed.mkv
video.failed.mkv.log
```

Since Test Clips checks for the existence of successful files prior to working, simply re-run matrices of video parameters (like the `for` loops suggested above) to re-attempt failures while skipping previous successes.

If you retry processing and it fails again, the original failed files will be overwritten.
