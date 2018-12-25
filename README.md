
# PngClean
Removes metadata from a PNG image

# Build
```
go build -o png *.go
chmod +x clean
```

# Run
```
# Display infos
./png info picture.png
# Removes metadata
./png clean picture.png
```

# Contributing
You might want to consider taking a look at the PNG specifications before diving into the code.
Especially chunk format, signature and mandatory chunks (IHDR, PLTE, IDAT, IEND)
