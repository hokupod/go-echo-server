package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスの Content-Type を JSON に設定
	w.Header().Set("Content-Type", "application/json")

	// リクエストボディを読み取る
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// レスポンス用の構造体を定義
	response := struct {
		Method string      `json:"method"`
		Path   string      `json:"path"`
		Body   string      `json:"body,omitempty"`
		Params interface{} `json:"params,omitempty"`
	}{
		Method: r.Method,
		Path:   r.URL.Path,
	}

	// リクエストメソッドに応じて処理を分岐
	if r.Method == "POST" {
		response.Body = string(body)
	} else if r.Method == "GET" {
		response.Params = r.URL.Query()
	}

	// 構造体を JSON にエンコードしてレスポンスとして送信
	json.NewEncoder(w).Encode(response)
}

func main() {
	// フラグ（コマンドライン引数）の定義
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	// ルートにハンドラーを登録
	http.HandleFunc("/", echoHandler)

	// サーバーを指定されたポートで起動
	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting echo-server at %s\n", address)
	http.ListenAndServe(address, nil)
}
