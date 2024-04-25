#!/bin/bash

BASE_URL="http://localhost:17000"
figure_width_norm=$(echo "scale=2; 300/800" | bc)
figure_height_norm=$(echo "scale=2; 280/800" | bc)
offset_x=$(echo "scale=2; $figure_width_norm / 2" | bc)
offset_y=$(echo "scale=2; $figure_height_norm / 2" | bc)

Xstart=$offset_x
Ystart=$offset_y

step=$(echo "scale=2; 20/800" | bc)

while true; do
    if [ "$direction" == "down_right" ]; then
        newX=$(echo "scale=2; $Xstart + $step" | bc)
        newY=$(echo "scale=2; $Ystart + $step" | bc)
        if (( $(echo "$newX > 1 - $offset_x" | bc) )) || (( $(echo "$newY > 1 - $offset_y" | bc) )); then
            direction="up_left"
            newX=$(echo "1 - $offset_x" | bc)
            newY=$(echo "1 - $offset_y" | bc)
        fi
    else
        newX=$(echo "scale=2; $Xstart - $step" | bc)
        newY=$(echo "scale=2; $Ystart - $step" | bc)
        # If the figure reaches the left or top border, change direction
        if (( $(echo "$newX < $offset_x" | bc) )) || (( $(echo "$newY < $offset_y" | bc) )); then
            direction="down_right"
            newX=$offset_x
            newY=$offset_y
        fi
    fi

    Xstart=$newX
    Ystart=$newY

    curl -X POST "$BASE_URL" -d "reset"
    curl -X POST "$BASE_URL" -d "green"
    curl -X POST "$BASE_URL" -d "figure $Xstart $Ystart"
    curl -X POST "$BASE_URL" -d "update"

    sleep 0.01
done
