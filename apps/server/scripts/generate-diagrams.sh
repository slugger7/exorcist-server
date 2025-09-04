#!/usr/bin/env bash
echo "Generating diagrams"

SOURCE_DIR="./diagrams/src"
DEST_DIR="./diagrams/out"

if [ ! -d "$SOURCE_DIR" ]; then
    echo "Source directory does not exist."
    exit 1
fi

mkdir -p "$DEST_DIR"

process_file() {
    local source_file="$1"
    local dest_file="$2"
    
    d2 "$source_file" "$dest_file"
}

find "$SOURCE_DIR" -type f -name "*.d2" | while read -r source_file; do
    relative_path="${source_file#$SOURCE_DIR/}"
    dest_file="${DEST_DIR}/${relative_path%.$SOURCE_EXT}.svg"

    mkdir -p "$(dirname "$dest_file")"

    process_file "$source_file" "$dest_file"
done

echo "Processing completed."
