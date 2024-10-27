package pkg

type Blockchain struct {
	Head      *Block
	Height    int
	Branchlen int
	Status    BlockchainStatus
	Blocks    map[string]*Block
	Forks     []*Block //heads of forkchains
}

type BlockchainStatus int

const (
	Active      BlockchainStatus = iota //current active chain
	ValidHeader                         //node observers a potential reorganization but curr chain becomes the longest
	ValidFork                           //node performed the reorganization
)

func (s BlockchainStatus) String() string {
	switch s {
	case Active:
		return "Active"
	case ValidHeader:
		return "ValidHeader"
	case ValidFork:
		return "ValidFork"
	default:
		return "Unknown"
	}
}

// called after new block is mined and added to the blockchain
func (chain *Blockchain) CreateNewBlock() *Block {
	prevBlock := chain.Head.Header.CurrBlockHash
	newBlockHeader := NewBlockHeader(prevBlock)
	newBlock := NewBlockTemplate(newBlockHeader)
	return newBlock

}

// func ValidBlock(){}
// func addminedblock(){}
