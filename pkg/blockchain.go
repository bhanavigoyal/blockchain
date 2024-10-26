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

// make a function which combines txns for data and after a height calls addnewblock
// the function can be like ki if len <=10 then add txn and after that call newblock
func (chain *Blockchain) CreateNewBlock(data *Transaction) {
	prevBlock := chain.Head.Header.CurrBlockHash
	newBlockHeader := NewBlockHeader(prevBlock)
	newBlock := NewBlock(newBlockHeader, data)
	newBlock.Header.MineBlock()
}

// func ValidBlock(){}
// func addminedblock(){}
