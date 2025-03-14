package user

import (
	"fmt"
	"sync"
)

// VoteManager holds the user mapping and votes per index (restaurant).
type VoteManager struct {
	mux             sync.RWMutex
	userToIndexMap  map[string]int   // maps userID to bit index
	restaurantVotes map[int][]uint64 // maps restaurantIndex to dynamic bitset (slice of bits that represent the likes for each restaurant)
}

// NewVoteManager creates a VoteManager.
func NewVoteManager() *VoteManager {
	return &VoteManager{
		userToIndexMap:  make(map[string]int),
		restaurantVotes: make(map[int][]uint64),
	}
}

// getUserIndex returns the index for a user, adding the user if needed.
func (vm *VoteManager) getUserIndex(userID string) int {
	if idx, ok := vm.userToIndexMap[userID]; ok {
		return idx
	}
	idx := len(vm.userToIndexMap)
	vm.userToIndexMap[userID] = idx
	return idx
}

// ensureBitsetSize ensures the bitset slice is long enough for the given user index.
func ensureBitsetSize(bitset []uint64, userIndex int) []uint64 {
	// normalize the userIndex to the range 0-63
	wordIdx := userIndex / 64
	if wordIdx < len(bitset) {
		return bitset
	}
	// Expand slice to hold at least wordIdx+1 words.
	newBits := make([]uint64, wordIdx+1)
	copy(newBits, bitset)
	return newBits
}

func (vm *VoteManager) AddUser(userID string) error {
	vm.mux.Lock()
	defer vm.mux.Unlock()
	if _, ok := vm.userToIndexMap[userID]; ok {
		return fmt.Errorf("user already exists")
	}
	vm.userToIndexMap[userID] = len(vm.userToIndexMap)
	return nil
}

/*
SetVote sets a user's vote for a restaurant and returns whether all users have liked it.

Returns true if the restaurant is now fully liked by all users, false otherwise.
*/
func (vm *VoteManager) SetVote(restaurantIndex int, userID string, vote bool) bool {
	vm.mux.Lock()

	// Get or create user index.
	userIndex := vm.getUserIndex(userID)
	// Get the index of the chunk of votes in the bitset the user's vote belongs to.
	chunckIdx := userIndex / 64
	// Get the position of the vote within the chunk.
	bitPos := uint(userIndex % 64)

	// Get or create the bitset for this restaurant.
	restaurantBitset, exists := vm.restaurantVotes[restaurantIndex]
	if !exists {
		restaurantBitset = make([]uint64, chunckIdx+1)
	}

	// Ensure bitset is large enough.
	restaurantBitset = ensureBitsetSize(restaurantBitset, userIndex)

	mask := uint64(1) << bitPos
	if vote {
		// Set the bit.
		restaurantBitset[chunckIdx] |= mask
	} else {
		// Clear the bit.
		restaurantBitset[chunckIdx] &^= mask
	}

	vm.restaurantVotes[restaurantIndex] = restaurantBitset
	vm.mux.Unlock()

	return vm.allLiked(restaurantIndex)
}

// AllLiked checks if every registered user has liked a specific restaurant.
// (Assumes that a missing vote is equivalent to not liking.)
func (vm *VoteManager) allLiked(restaurantIndex int) bool {
	vm.mux.RLock()
	defer vm.mux.RUnlock()

	bitset, exists := vm.restaurantVotes[restaurantIndex]
	if !exists {
		return false
	}
	totalUsers := len(vm.userToIndexMap)
	wordsNeeded := (totalUsers + 63) / 64

	// Ensure our bitset is up to date.
	if len(bitset) < wordsNeeded {
		return false
	}

	// For every user, the corresponding bit must be set.
	// We'll compare word by word.
	for w := 0; w < wordsNeeded; w++ {
		var mask uint64
		if w < wordsNeeded-1 {
			mask = ^uint64(0) // all 64 bits
		} else {
			// Last word: only the lower bits corresponding to remaining users.
			remaining := totalUsers - w*64
			mask = (uint64(1) << remaining) - 1
		}
		if bitset[w]&mask != mask {
			return false
		}
	}
	return true
}
