package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	Blockchain []schema.BlockInterface
)

func GetBlockChainList(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response *schema.Response
	addHeaders(w, r)
	response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Blockchain list processed succesfully", BlockChain: Blockchain}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug("GetBlockChainList response : %s", string(b))
	fmt.Fprintf(w, "%s", string(b))
}

func GetBlockChain(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response *schema.Response
	var hold []schema.BlockInterface
	vars := mux.Vars(r)
	addHeaders(w, r)
	index, _ := strconv.Atoi(vars["index"])
	hold = append(hold, Blockchain[index])
	response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Blockchain list processed succesfully", BlockChain: hold}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug("GetBlockChain response : %s", string(b))
	fmt.Fprintf(w, "%s", string(b))
}

func WriteBlockChain(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response *schema.Response
	var in schema.InputData
	addHeaders(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		conn.Error("WriteBlockChain %v", err)
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not read body data %v\n", err), BlockChain: Blockchain}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.MarshalIndent(response, "", "	")
		fmt.Fprintf(w, "%s", string(b))
		return
	}
	e := json.Unmarshal(body, &in)
	if e != nil {
		conn.Error("WriteBlockChain %v", e)
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not unmarshal json input data %v\n", err), BlockChain: Blockchain}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.MarshalIndent(response, "", "	")
		fmt.Fprintf(w, "%s", string(b))
		return
	}
	newBlock, e := generateBlock(Blockchain[len(Blockchain)-1], in)
	if e != nil {
		conn.Error("WriteBlockChain %v", e)
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "KO", Message: fmt.Sprintf("Could not write to blockchain %v\n", err), BlockChain: Blockchain}
		w.WriteHeader(http.StatusInternalServerError)
		b, _ := json.MarshalIndent(response, "", "	")
		fmt.Fprintf(w, "%s", string(b))
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
	}
	response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data processed succesfully", BlockChain: Blockchain}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug("WriteBlockChain response : %s", string(b))
	fmt.Fprintf(w, "%s", string(b))
}

func Init(w http.ResponseWriter, r *http.Request) {
	if len(Blockchain) == 0 {
		t := time.Now()
		genesisBlock := schema.BlockInterface{0, "Genesis", "A", "B", "C", t.String(), "", ""}
		Blockchain = append(Blockchain, genesisBlock)
		fmt.Fprintf(w, "%s", "{ \"created genesis block\" }")
		return
	}
	fmt.Fprintf(w, "%s", "{ \"NOP genesis block created\" }")
}

func generateBlock(oldBlock schema.BlockInterface, in schema.InputData) (schema.BlockInterface, error) {
	var newBlock schema.BlockInterface
	t := time.Now()
	if in.ObjectA == "" || in.ObjectB == "" || in.ObjectC == "" {
		return newBlock, errors.New("object blocks can't be empty")
	}
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.MetaInfo = in.MetaInfo
	newBlock.ObjectA = in.ObjectA
	newBlock.ObjectB = in.ObjectB
	newBlock.ObjectC = in.ObjectC
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock schema.BlockInterface) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func replaceChain(newBlocks []schema.BlockInterface) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func calculateHash(block schema.BlockInterface) string {
	record := fmt.Sprintf("%d", block.Index) + block.Timestamp + block.ObjectA + block.ObjectB + block.ObjectC + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
