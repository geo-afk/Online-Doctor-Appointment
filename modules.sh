find . -name '*.go' | while read file; do
    pkg=$(dirname "$file" | sed 's|^\./||')         # Remove leading ./ for clarity
    if [ "$pkg" = "." ]; then
        echo "$file => $(go list -m)/"
    else
        echo "$file => $(go list -m)/$pkg"
    fi
done
