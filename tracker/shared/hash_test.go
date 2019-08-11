package shared

import (
	"math/rand"
	"testing"
)

func dbWithHashes(count int) *PeerDatabase {
	var db PeerDatabase
	db.make()
	InitExpvar(&db)

	peerid := PeerID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	peer := Peer{
		Complete: true,
		IP:       PeerIP{1, 2, 3, 4},
		Port:     4321,
		LastSeen: 1234567890,
	}

	var h Hash
	for i := 0; i < count; i++ {
		hash := make([]byte, 20)
		rand.Read(hash)
		copy(h[:], hash)

		db.Save(&peer, &h, &peerid)
	}

	return &db
}

func dbWithPeers(count int) (*PeerDatabase, Hash) {
	var db PeerDatabase
	db.make()
	InitExpvar(&db)

	bytes := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	hash := Hash(bytes)
	peer := Peer{
		Complete: true,
		IP:       PeerIP{1, 2, 3, 4},
		Port:     4321,
		LastSeen: 1234567890,
	}

	var p PeerID
	for i := 0; i < count; i++ {
		peerid := make([]byte, 20)
		rand.Read(peerid)
		copy(p[:], peerid)

		db.Save(&peer, &hash, &p)
	}

	return &db, hash
}

func benchmarkHashes(b *testing.B, count int) {
	db := dbWithHashes(count)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.Hashes()
	}
}

func BenchmarkHashes0(b *testing.B)      { benchmarkHashes(b, 0) }
func BenchmarkHashes5000(b *testing.B)   { benchmarkHashes(b, 5000) }
func BenchmarkHashes50000(b *testing.B)  { benchmarkHashes(b, 50000) }
func BenchmarkHashes500000(b *testing.B) { benchmarkHashes(b, 500000) }

// more/less peers doesn't change performance
const numPeers = 1000

func benchmarkPeerList(b *testing.B, cap int) {
	db, hash := dbWithPeers(numPeers)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.PeerList(&hash, cap, false)
	}
}

func BenchmarkPeerList50(b *testing.B)  { benchmarkPeerList(b, 50) }
func BenchmarkPeerList100(b *testing.B) { benchmarkPeerList(b, 100) }
func BenchmarkPeerList200(b *testing.B) { benchmarkPeerList(b, 200) }
func BenchmarkPeerList400(b *testing.B) { benchmarkPeerList(b, 400) }

func benchmarkPeerListNopeerid(b *testing.B, cap int) {
	db, hash := dbWithPeers(numPeers)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.PeerList(&hash, cap, true)
	}
}

func BenchmarkPeerListNopeerid50(b *testing.B)  { benchmarkPeerListNopeerid(b, 50) }
func BenchmarkPeerListNopeerid100(b *testing.B) { benchmarkPeerListNopeerid(b, 100) }
func BenchmarkPeerListNopeerid200(b *testing.B) { benchmarkPeerListNopeerid(b, 200) }
func BenchmarkPeerListNopeerid400(b *testing.B) { benchmarkPeerListNopeerid(b, 400) }

func benchmarkPeerListBytes(b *testing.B, cap int) {
	db, hash := dbWithPeers(numPeers)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.PeerListBytes(&hash, cap)
	}
}

func BenchmarkPeerListBytes50(b *testing.B)  { benchmarkPeerListBytes(b, 50) }
func BenchmarkPeerListBytes100(b *testing.B) { benchmarkPeerListBytes(b, 100) }
func BenchmarkPeerListBytes200(b *testing.B) { benchmarkPeerListBytes(b, 200) }
func BenchmarkPeerListBytes400(b *testing.B) { benchmarkPeerListBytes(b, 400) }

func benchmarkHashStats(b *testing.B, peers int) {
	db, hash := dbWithPeers(peers)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.HashStats(&hash)
	}
}

func BenchmarkHashStats100(b *testing.B) { benchmarkHashStats(b, 100) }
func BenchmarkHashStats1000(b *testing.B) { benchmarkHashStats(b, 1000) }
