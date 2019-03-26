package asset

import (
    "fmt"
    "github.com/sprinkle-it/donut/buffer"
    "hash/crc32"
)

// Creates the release manifest for a cache. The release manifest contains the checksum and version of each index
// manifest contained within the cache. This is used so that the client can verify if the cache it has locally is the
// same version as the one being distributed.
func CreateReleaseManifest(cache *Cache) ([]byte, error) {
    manifest := buffer.NewByteBuffer(cache.IndexCount() * 8)
    for i := 0; i < cache.IndexCount(); i++ {
        archive, err := cache.Get(Manifests, uint16(i))
        if err != nil {
            return nil, err
        }

        _ = manifest.PutUint32(crc32.Checksum(archive, crc32.IEEETable))

        b, err := DecompressArchive(archive)
        if err != nil {
            return nil, err
        }

        indexManifest := buffer.ByteBuffer{Bytes:b}

        // Oldschool only supports manifest format 5 and 6.
        format, _ := indexManifest.GetUint8()
        if format < 5 || format > 6 {
            return nil, fmt.Errorf("asset: unsupported manifest format %d", format)
        }

        var version uint32
        if format > 5 {
            version, _ = indexManifest.GetUint32()
        }

        _ = manifest.PutUint32(version)
    }
    return manifest.Bytes, nil
}

// Computes the checksum of each index manifest contained in the cache.
func GetManifestChecksums(cache *Cache) ([]uint32, error) {
    checksums := make([]uint32, len(cache.indexes))
    for i := 0; i < cache.IndexCount(); i++ {
        archive, err := cache.Get(Manifests, uint16(i))
        if err != nil {
            return nil, err
        }
        checksums[i] = crc32.ChecksumIEEE(archive)
    }
    return checksums, nil
}
