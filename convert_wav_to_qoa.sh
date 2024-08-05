#!/bin/bash

# This script ensures all the sound files used in this project are the same
# bit depth and sample rate, and then converts them to QOA format.

# Ensure sox and goqoa are installed and available in the PATH
if ! command -v sox &> /dev/null; then
    echo "sox could not be found, please install it."
    exit 1
fi

if ! command -v goqoa &> /dev/null; then
    echo "goqoa could not be found, please install it."
    exit 1
fi

for wav in $1*.wav; do
    # Check the bit depth of the WAV file using soxi
    bit_depth=$(soxi -b "$wav")
    sample_rate=$(soxi -r "$wav")

    if [ "$bit_depth" -eq 8 ]; then
        echo "Converting 8-bit WAV to 16-bit: $wav"
        # Convert the 8-bit WAV to 16-bit WAV
        temp_wav="${wav%.wav}_16bit.wav"
        sox "$wav" -b 16 "$temp_wav"

        # Convert the 16-bit WAV to QOA
        goqoa convert -v "$temp_wav" "$(basename "$wav" .wav).qoa"

        # Remove the temporary 16-bit WAV file
        rm "$temp_wav"
    elif [ "$sample_rate" -ne 44100 ]; then
        echo "Converting $sample_rate sample rate WAV to 44100: $wav"

        temp_wav="${wav%.wav}_44100.wav"
        sox "$wav" -r 44100 "$temp_wav" vol 0.9

        # Convert the 16-bit WAV to QOA
        goqoa convert -v "$temp_wav" "$(basename "$wav" .wav).qoa"

        # Remove the temporary 16-bit WAV file
        rm "$temp_wav"

    else
        # Directly convert the WAV to QOA if it's already 16-bit
        goqoa convert -v "$wav" "$(basename "$wav" .wav).qoa"
    fi
done
