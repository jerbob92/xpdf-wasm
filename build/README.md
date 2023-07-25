Get the emsdk repo

git clone https://github.com/emscripten-core/emsdk.git

# Enter that directory
cd emsdk

# Fetch the latest version of the emsdk (not needed the first time you clone)
git pull

# Checkout the correct version
git checkout 3.1.44

# Download and install the SDK tools.
./emsdk install 3.1.44

# Make the SDK version active for the current user. (writes .emscripten file)
./emsdk activate 3.1.44

# Activate PATH and other environment variables in the current terminal
source ./emsdk_env.sh

# Go to the emscripten directory
cd upstream/emscripten

# Remove the cache dir
rm -Rf cache

# Apply our emscripten WASI patch
patch -p1 < ../../../emscripten.patch

# Go back to the root of the build dir
cd ../../../

# Download xpdf
wget https://dl.xpdfreader.com/xpdf-latest.tar.gz

# Extract xpdf to the directory named xpfd
mkdir -p xpdf && tar -xvf xpdf-latest.tar.gz --strip 1 -C xpdf

# Go into the xpdf directory
cd xpdf

# Create a build directory
mkdir -p build

# Go into the build directory
cd build

# Install emscripten ports
embuilder build zlib
embuilder build libpng
embuilder build libjpeg
embuilder build freetype

# Make the cmake target for Emscripten
# Use -g instead of -O2 to get debug symbols.
emcmake cmake .. -DCMAKE_TOOLCHAIN_FILE=$EMSDK/upstream/emscripten/cmake/Modules/Platform/Emscripten.cmake -DCMAKE_CXX_FLAGS="-std=c++14 -O2 -DLOAD_FONTS_FROM_MEM=1" -DCMAKE_EXE_LINKER_FLAGS="-static -sERROR_ON_UNDEFINED_SYMBOLS=0 -s WASM=1 -s ALLOW_MEMORY_GROWTH=1 -s STANDALONE_WASM=1 -s USE_FREETYPE=1 -s USE_ZLIB=1 -s USE_LIBPNG=1 -s USE_LIBJPEG=1"

# Make xpdf
emmake make

# Copy the built wasm binary to the root of the project
cp xpdf/*.wasm ../../../wasm


