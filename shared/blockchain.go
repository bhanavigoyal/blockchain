package pkg

type Blockchain struct {
	Head     *Block
	Height   int
	Balances map[string]int //add decimal

	// Branchlen int
	// Status    BlockchainStatus
	// Blocks    map[string]*Block
	// Forks     []*Block       //heads of forkchains

}

/*
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
*/

func (chain *Blockchain) CreateNewBlock() *Block {
	if chain.Height == 0 {
		block := chain.genesisBlock()
		return block
	}
	prevBlock := chain.Head.Header.CurrBlockHash
	newBlockHeader := NewBlockHeader(prevBlock)
	newBlock := NewBlockTemplate(newBlockHeader)
	return newBlock
}

func (chain *Blockchain) genesisBlock() *Block {
	//implement genesis block
	prevBlock := []byte("0")
	newBlockHeader := NewBlockHeader(prevBlock)
	newBlock := NewBlockTemplate(newBlockHeader)
	chain.AddMinedBlock(newBlock)

	//make a genesis wallet
	//add some initial coins to it

	return newBlock

}

func (chain *Blockchain) AddMinedBlock(block *Block) {
	chain.Head = block
	chain.Height++
}
