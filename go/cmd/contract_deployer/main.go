package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andantan/vault-lab/go/internal/config"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		configPath   string
		contractPath string
		alias        string
	)
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.StringVar(&contractPath, "contract", "", "path to the .sol file to deploy (relative to project root)")
	flag.StringVar(&contractPath, "c", "", "path to the .sol file to deploy (relative to project root)")
	flag.StringVar(&alias, "deployer", "", "key alias to use as deployer")
	flag.StringVar(&alias, "d", "", "key alias to use as deployer")
	flag.Parse()

	if contractPath == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: -contract/-c is required")
		os.Exit(1)
	}
	if alias == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: -deployer/-d is required")
		os.Exit(1)
	}

	root, err := findProjectRoot()
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(root, configPath))
	if err != nil {
		return err
	}

	key, err := cfg.KeyByAlias(alias)
	if err != nil {
		return err
	}

	ctx := context.Background()

	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return fmt.Errorf("connect rpc: %w", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain id: %w", err)
	}

	privateKey, from, err := keyToAddress(key.PrivateKey)
	if err != nil {
		return err
	}
	fmt.Println("Deployer:", from.Hex())

	bytecode, err := loadBytecode(root, contractPath)
	if err != nil {
		return err
	}
	fmt.Println("Bytecode size:", len(bytecode))

	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return fmt.Errorf("get nonce: %w", err)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    nil,
		Value: big.NewInt(0),
		Data:  bytecode,
	})
	if err != nil {
		return fmt.Errorf("estimate gas: %w", err)
	}
	gasLimit = gasLimit + gasLimit/5

	tipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("suggest gas tip cap: %w", err)
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("get latest header: %w", err)
	}

	feeCap := new(big.Int).Add(
		new(big.Int).Mul(header.BaseFee, big.NewInt(2)),
		tipCap,
	)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        nil,
		Value:     big.NewInt(0),
		Data:      bytecode,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}

	if err = client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("send tx: %w", err)
	}

	fmt.Println("RPC:", cfg.RPCURL)
	fmt.Println("Chain ID:", chainID.String())
	fmt.Println("Deploy tx:", signedTx.Hash().Hex())

	receipt, err := waitReceipt(ctx, client, signedTx.Hash(), 30*time.Second)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("deployment failed: status=%d", receipt.Status)
	}

	code, err := client.CodeAt(ctx, receipt.ContractAddress, nil)
	if err != nil {
		return fmt.Errorf("get deployed code: %w", err)
	}

	fmt.Println("Contract address:", receipt.ContractAddress.Hex())
	fmt.Println("Block number:", receipt.BlockNumber.String())
	fmt.Println("Gas used:", receipt.GasUsed)
	fmt.Println("Runtime code size:", len(code), "bytes")

	return nil
}

// solcjs output filename: path separators and dots replaced with underscores, suffixed with contract name.
// e.g. contracts/vault/MultiAccountVault.sol -> contracts_vault_MultiAccountVault_sol_MultiAccountVault.bin
func binFilename(solPath string) string {
	normalized := strings.ReplaceAll(solPath, "/", "_")
	normalized = strings.ReplaceAll(normalized, ".", "_")
	contractName := strings.TrimSuffix(filepath.Base(solPath), ".sol")
	return normalized + "_" + contractName + ".bin"
}

func loadBytecode(root, contractPath string) ([]byte, error) {
	binPath := filepath.Join(root, "build", binFilename(contractPath))
	binBytes, err := os.ReadFile(binPath)
	if err != nil {
		return nil, fmt.Errorf("read bytecode: %w", err)
	}

	binHex := strings.TrimPrefix(strings.TrimSpace(string(binBytes)), "0x")
	bytecode, err := hex.DecodeString(binHex)
	if err != nil {
		return nil, fmt.Errorf("decode bytecode: %w", err)
	}

	return bytecode, nil
}

func keyToAddress(privateKeyHex string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyHex = strings.TrimPrefix(strings.TrimSpace(privateKeyHex), "0x")

	privateKey, err := gethcrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("parse private key: %w", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("invalid public key type")
	}

	return privateKey, gethcrypto.PubkeyToAddress(*publicKey), nil
}

func waitReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	deadline := time.Now().Add(timeout)

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for receipt: %s", txHash.Hex())
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if exists(filepath.Join(dir, "config.yaml")) && exists(filepath.Join(dir, "contracts")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found from %s", wd)
		}

		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
