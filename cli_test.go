package main

import (
	"testing"
	"fmt"
	"net/http"
	"io/ioutil"
	"math/big"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
	"net/smtp"
	"time"
)

func TestGetBalance(t *testing.T)  {
	cli := &CLI{}
	cli.getBalance(string("Ivan"))
}

func BenchmarkSelfDefineKeyPair(b *testing.B) {

	timeS := time.Now()
	hex := "8888888888888888888888888888888800000000000000000000000000000000"//1HT7xU2Ngenf7D4yocz2SAcnNLW7rK8d4E
	//fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
	fmt.Println("hex key is :", hex, "len:", len(hex))
	private, public := SelfDefineKeyPair(hex)
	w := Wallet{private, public}
	fmt.Printf("public key:%s\n", w.GetAddress())
	fmt.Println("spend:", time.Since(timeS))
}

func TestApi(t *testing.T)  {
	resp, err := http.Get("https://chain.api.btc.com/v3/address/15PFuJ2gKoB9QfmG1rBGVC3L8bFwhxC6oQ")
	//resp, err := http.Get("https://bitinfocharts.com/top-100-richest-bitcoin-addresses-5.html")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("resp:%s\n", b)

}

func TestBigInt(t *testing.T)  {
	hex := "1000000000000000000000000000000000000000000000000000000000000000"
	num, b := big.NewInt(0).SetString(hex, 16)
	if !b {
		panic("error number")
	}
	num = num.Div(num, big.NewInt(2))
	fmt.Println("big number is: ", num.String())
}

func TestParseHtml(t *testing.T)  {

	address := make(map[string]struct {}, 2100)

	for i := 1; i <=31 ;i++  {
		fileName := fmt.Sprintf("C:/gopath/src/blockchain_go/html/Addresses%d.html", i)
		fmt.Println("Parse file:", fileName)
		fd, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(fd)
		if err != nil {
			panic(err)
		}

		// Find the review items
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			band,b := s.Attr("href")
			if !b {
				return
			}
			if strings.Contains(band, "https://bitinfocharts.com/bitcoin/address/") {
				addr := s.Text()
				if len(addr) < 0 {
					return
				}
				if  addr[0] != '1'{
					return
				}
				address[addr] = struct{}{}
			}
		})
		fd.Close()
		fmt.Println("address coumt is ", len(address))
	}


	fmt.Println("address count is" ,len(address))
	targetfd, err := os.OpenFile("C:/gopath/src/blockchain_go/html/addr.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		 panic(err)
	}
	for key, _ := range address{
		targetfd.Write([]byte(key))
		targetfd.Write([]byte("\n"))
	}
	targetfd.Close()
}

func TestReadFile(t *testing.T)  {
	data, err := ioutil.ReadFile("C:/gopath/src/blockchain_go/html/addr.txt")
	if err != nil {
		panic(err)
	}
	str := string(data)
	strs := strings.Split(str, "\n")
	for idx, val := range strs {
		if len(val) <= 0 {
			continue
		}
		fmt.Println(idx, " address:", val)
	}
}

func TestSendMail(t *testing.T)  {
	// 邮箱地址
	UserEmail := "zm53373581@163.com"
	// 端口号，:25也行
	Mail_Smtp_Port := ":25"
	//邮箱的授权码，去邮箱自己获取
	Mail_Password := "fuckhack2015"
	// 此处填写SMTP服务器
	Mail_Smtp_Host := "smtp.163.com"
	auth := smtp.PlainAuth("", UserEmail, Mail_Password, Mail_Smtp_Host)
	to := []string{UserEmail}
	nickname := "昵称"
	user := UserEmail

	subject := "testz主题"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "邮件内容."
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(Mail_Smtp_Host+Mail_Smtp_Port, auth, user, to, msg)
	if err != nil {
		fmt.Println("send mail error:", err)
	}
}

func TestCreateWallet(t *testing.T)  {
	S := time.Now()
	wallet := NewWallet()
	fmt.Printf("big number is: %s\n", wallet.GetAddress())
	fmt.Println("spend:", time.Since(S))
}
func BenchmarkNewWallet(b *testing.B) {
	S := time.Now()
	wallet := NewWallet()
	fmt.Printf("big number is: %s pri:%s\n", wallet.GetAddress(), wallet.PrivateKey.D.String())
	fmt.Println("spend:", time.Since(S))
}

func TestBigN(t *testing.T)  {
	start := NewWallet().PrivateKey.D
	start = big.NewInt(0).Add(start, big.NewInt(1))
	hex := fmt.Sprintf("%x", start.Bytes())
	fmt.Printf("hex key is :%s\n", hex)
}