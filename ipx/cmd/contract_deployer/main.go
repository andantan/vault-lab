package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andantan/evmlab/core"
	"github.com/andantan/evmlab/core/types"
	"github.com/andantan/evmlab/internal/config"
	"github.com/andantan/evmlab/internal/rpc"
	"github.com/andantan/evmlab/internal/util"
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
	flag.StringVar(&alias, "deployer", "", "deployer address")
	flag.StringVar(&alias, "d", "", "deployer address")
	flag.Parse()

	if contractPath == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: -contract/-c is required")
		os.Exit(1)
	}
	if alias == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: -deployer/-d address is required")
		os.Exit(1)
	}

	root, err := util.FindProjectRoot()
	if err != nil {
		return err
	}

	cfg, err := config.Load(filepath.Join(root, configPath))
	if err != nil {
		return err
	}

	key, err := cfg.KeyByAddress(alias)
	if err != nil {
		return err
	}

	evmKey, err := core.DeriveKeyFromHex(key.PrivateKey, key.PublicKey, key.Address)
	if err != nil {
		return fmt.Errorf("derive key: %w", err)
	}
	deployerAddr := evmKey.Address.String()
	fmt.Println("Deployer:", deployerAddr)

	ctx := context.Background()
	client := rpc.NewClient(cfg.RPCURL)

	chainIDHex, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain id: %w", err)
	}
	chainID, err := util.HexToBigInt(chainIDHex)
	if err != nil {
		return fmt.Errorf("parse chain id: %w", err)
	}

	bytecode, err := loadBytecode(root, contractPath)
	if err != nil {
		return err
	}
	fmt.Println("Bytecode size:", len(bytecode))

	nonceHex, err := client.GetTransactionCount(ctx, deployerAddr, "pending")
	if err != nil {
		return fmt.Errorf("get nonce: %w", err)
	}
	nonce, err := util.HexToUint64(nonceHex)
	if err != nil {
		return fmt.Errorf("parse nonce: %w", err)
	}

	gasLimitHex, err := client.EstimateGas(ctx, map[string]string{
		"from": deployerAddr,
		"data": "0x" + hex.EncodeToString(bytecode),
	}, "latest")
	if err != nil {
		return fmt.Errorf("estimate gas: %w", err)
	}
	gasLimit, err := util.HexToUint64(gasLimitHex)
	if err != nil {
		return fmt.Errorf("parse gas limit: %w", err)
	}
	gasLimit = gasLimit + gasLimit/5

	tipCapHex, err := client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		return fmt.Errorf("get tip cap: %w", err)
	}
	tipCap, err := util.HexToBigInt(tipCapHex)
	if err != nil {
		return fmt.Errorf("parse tip cap: %w", err)
	}

	block, err := client.BlockByNumber(ctx, "latest")
	if err != nil {
		return fmt.Errorf("get latest block: %w", err)
	}
	baseFee, err := util.HexToBigInt(block["baseFeePerGas"].(string))
	if err != nil {
		return fmt.Errorf("parse base fee: %w", err)
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)

	tx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		GasLimit:  gasLimit,
		To:        nil,
		Value:     big.NewInt(0),
		Data:      bytecode,
	}

	unsigned, err := core.Codec.EncodeDynamicFeeUnsigned(tx)
	if err != nil {
		return fmt.Errorf("encode unsigned tx: %w", err)
	}

	txHash := core.Hasher.Hash(unsigned)
	sig, err := core.Signer.Sign(txHash, *evmKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}

	rawTxBytes, err := core.Codec.EncodeDynamicFeeSigned(tx, sig)
	if err != nil {
		return fmt.Errorf("encode signed tx: %w", err)
	}

	txHashHex, err := client.SendRawTransaction(ctx, "0x"+hex.EncodeToString(rawTxBytes))
	if err != nil {
		return fmt.Errorf("send tx: %w", err)
	}

	fmt.Println("RPC:", cfg.RPCURL)
	fmt.Println("Chain ID:", chainID.String())
	fmt.Println("Deploy tx:", txHashHex)

	receipt, err := client.WaitForReceipt(ctx, txHashHex, 30*time.Second)
	if err != nil {
		return err
	}

	if receipt["status"] != "0x1" {
		return fmt.Errorf("deployment failed: status=%s", receipt["status"])
	}

	contractAddr := receipt["contractAddress"].(string)
	blockNumber := receipt["blockNumber"].(string)
	gasUsed := receipt["gasUsed"].(string)

	codeHex, err := client.GetCode(ctx, contractAddr, "latest")
	if err != nil {
		return fmt.Errorf("get deployed code: %w", err)
	}
	codeSize := (len(strings.TrimPrefix(codeHex, "0x"))) / 2

	fmt.Println("Contract address:", contractAddr)
	fmt.Println("Block number:", blockNumber)
	fmt.Println("Gas used:", gasUsed)
	fmt.Println("Runtime code size:", codeSize, "bytes")

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

// contracts/vault/Foo.sol -> vault
func contractSubdir(solPath string) string {
	return strings.TrimPrefix(filepath.Dir(solPath), "contracts/")
}

func loadBytecode(root, contractPath string) ([]byte, error) {
	binPath := filepath.Join(root, "build", contractSubdir(contractPath), binFilename(contractPath))
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
