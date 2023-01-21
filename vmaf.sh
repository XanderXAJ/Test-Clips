#!/usr/bin/env bash
reference=$1
distorted=$2

ffmpeg -i "$reference" -i "$distorted" \
	-lavfi "[0:v]setpts=PTS-STARTPTS[reference];
					[1:v]setpts=PTS-STARTPTS[distorted];
					[distorted][reference]libvmaf=log_fmt=json:log_path='${distorted}.vmaf.json':n_threads=16:feature='name=psnr|name=float_ssim|name=float_ms_ssim|name=float_ansnr'" \
	-f null -
