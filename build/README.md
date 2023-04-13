# Get the emsdk repo
git clone https://github.com/emscripten-core/emsdk.git

# Enter that directory
cd emsdk

# Fetch the latest version of the emsdk (not needed the first time you clone)
git pull

# Download and install the latest SDK tools.
./emsdk install latest

# Make the "latest" SDK "active" for the current user. (writes .emscripten file)
./emsdk activate latest

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

# Make the cmake target for Emscripten
cmake .. -DCMAKE_TOOLCHAIN_FILE=$EMSDK/upstream/emscripten/cmake/Modules/Platform/Emscripten.cmake -DCMAKE_CXX_FLAGS="-std=c++14 -g" -DCMAKE_EXE_LINKER_FLAGS="-static -sERROR_ON_UNDEFINED_SYMBOLS=0 -s WASM=1 -s ALLOW_MEMORY_GROWTH=1 -s STANDALONE_WASM=1"

# Make xpdf
make

# Copy the built wasm binary to the root of the project
cp xpdf/*.wasm ../../../


