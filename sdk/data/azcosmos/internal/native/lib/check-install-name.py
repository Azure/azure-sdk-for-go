import struct, sys

# Reads LC_ID_DYLIB from a thin 64-bit Mach-O. dyld resolves an absolute install name directly and
# never consults the loading binary's rpaths, so an absolute one here is both a dependency on a
# directory that may not exist and somewhere a local user could plant a library.
LC_ID_DYLIB = 0xd
path = sys.argv[1]
f = open(path, 'rb').read()
magic, = struct.unpack_from('<I', f, 0)
if magic not in (0xfeedfacf,):
    sys.exit("%s: not a thin little-endian 64-bit Mach-O (magic %#x)" % (path, magic))
ncmds, = struct.unpack_from('<I', f, 16)
off = 32
for _ in range(ncmds):
    cmd, cmdsize = struct.unpack_from('<II', f, off)
    if cmd == LC_ID_DYLIB:
        nameoff, = struct.unpack_from('<I', f, off + 8)
        start = off + nameoff
        name = f[start:f.index(b'\x00', start)].decode()
        print("install name: %s" % name)
        if not name.startswith('@rpath/'):
            sys.exit("%s: install name must start with @rpath/, got %r" % (path, name))
        sys.exit(0)
    off += cmdsize
sys.exit("%s: no LC_ID_DYLIB" % path)
