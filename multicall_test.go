package multicall_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"multicall"
	"strings"

	//"github.com/alethio/web3-multicall-go/multicall"
	"testing"
	"time"
)

func TestExampleViwCall(t *testing.T) {
	eth, err := getETH("https://mainnet.infura.io/v3/17ed7fe26d014e5b9be7dfff5368c69d")
	vc := multicall.NewViewCall(
		"key.1",
		"0x47ac5182603f29a05bac8a56280b22a793503c97", //baobab USDC
		"balanceOf(address)(uint256)",
		[]interface{}{"0xff0e11e079abc08ac3cae86bd7937c6a45d1e917"},
		//[]interface{}{}, // empty arguments
	)
	//vc.Validate()
	vcs := multicall.ViewCalls{vc}
	mc, _ := multicall.New(eth)
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		panic(err)
	}

	resJson, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(resJson))
	fmt.Println(res)
	fmt.Println(err)

}

func getETH(url string) (ethrpc.ETHInterface, error) {
	provider, err := httprpc.New(url)
	if err != nil {
		return nil, err
	}
	provider.SetHTTPTimeout(5 * time.Second)
	return ethrpc.New(provider)
}

type callArgs struct {
	Target   [20]byte
	CallData []byte
}

func TestAbiEncode(t *testing.T) {
	addressTy, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
	}

	bytes, _ := arguments.Pack(
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		[32]byte{'I', 'D', '1'},
		big.NewInt(42),
	)

	//var buf []byte
	hash := crypto.Keccak256Hash(bytes)
	//hash := .NewKeccak256()
	//hash.Write(bytes)
	//buf = hash.Sum(buf)
	//log.Println(hexutil.Encode(buf))
	log.Println(hash.Hex())

	balanceOfID := crypto.Keccak256Hash([]byte("balanceOf(address)"))
	methodCallData := balanceOfID.Hex()[:10]
	log.Println("balanceOf", methodCallData) //0x70a08231

	address := "0xff0e11e079abc08ac3cae86bd7937c6a45d1e917" // user account
	argsSuffix, err := toByteArray(address)
	if err != nil {
		t.Error(err)
	}

	target := "0x47ac5182603f29a05bac8a56280b22a793503c97"
	targetBytes, err := toByteArray(target)
	if err != nil {
		t.Error(err)
	}

	payload := make([]byte, 0)
	payload = append(payload, targetBytes[:]...)
	payload = append(payload, methodCallData...)
	payload = append(payload, argsSuffix[:]...)

	fmt.Println(hex.EncodeToString(payload))

	//payloadArgs = callArgs{targetBytes, payload}

	//packed, err := abi.Pack("", big.NewInt(1), big.NewInt(2))
	//if err != nil {
	//	t.Error(err)
	//}

	// pack???
	//vc := multicall.NewViewCall(
	//	"key.1",
	//	"0xff0e11e079abc08ac3cae86bd7937c6a45d1e917", //user account
	//	"balanceOf(address)",
	//	"",
	//)
}

func toByteArray(address string) ([20]byte, error) {
	var addressBytes [20]byte
	address = strings.Replace(address, "0x", "", -1)
	addressBytesSlice, err := hex.DecodeString(address)
	if err != nil {
		return addressBytes, err
	}

	copy(addressBytes[:], addressBytesSlice[:])
	return addressBytes, nil
}
