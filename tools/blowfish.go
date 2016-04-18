package main

import (
  "fmt"
  "encoding/base64"
  "crypto/cipher"
  "golang.org/x/crypto/blowfish"
)

func blowfishChecksize(pt []byte) []byte {
  // Blowfish用のブロックサイズが `blowfish.BlockSize` で扱えるようにする
  modules := len(pt) % blowfish.BlockSize
  if modules != 0 {
    padlen := blowfish.BlockSize - modules
    for i := 0;i < padlen; i++ {
      pt = append(pt, 0)
    }
  }
  return pt
}

func blowfishEncrypt(ppt, key []byte) []byte {
  // Cipherを初期化してBlowfishで暗号化を実行する
  ecipher, err := blowfish.NewCipher(key)
  if err != nil {
    panic(err)
  }
  ciphertext := make([]byte, blowfish.BlockSize+len(ppt))
  eiv := ciphertext[:blowfish.BlockSize]
  ecbc := cipher.NewCBCEncrypter(ecipher, eiv)
  ecbc.CryptBlocks(ciphertext[blowfish.BlockSize:],ppt)
  return ciphertext
}

func blowfishDecrypt(et, key []byte) []byte {
  // 暗号化されたデータを復号化する
  dcipher, err := blowfish.NewCipher(key)
  if err != nil {
    panic(err)
  }
  div := et[:blowfish.BlockSize]
  decrypted := et[blowfish.BlockSize:]
  if len(decrypted) % blowfish.BlockSize != 0{
    panic("decrypted is not a multiple of blowfish.BlockSize")
  }
  dcbc := cipher.NewCBCDecrypter(dcipher, div)
  dcbc.CryptBlocks(decrypted, decrypted)
  return decrypted
}

func main() {
  var secret string
  fmt.Print("Enter SecretKey : ")
  fmt.Scanln(&secret)
  var salt string
  fmt.Print("Enter Salt : ")
  fmt.Scanln(&salt)
  var pass string
  fmt.Print("Enter PassKey : ")
  fmt.Scanln(&pass)
  text := []byte(fmt.Sprintf("%s:%s", pass, salt))
  fmt.Printf("PassKey : Salt is %s\n", text)
  key := []byte(secret)
  fmt.Printf("SecretKey is %s\n", key)
  pad := blowfishChecksize(text)
  enced := blowfishEncrypt(pad, key)
  fmt.Printf("encrypted is %x\n", enced)
  bases := base64.StdEncoding.EncodeToString(enced)
  fmt.Printf("base64encoded is %s\n", bases)
  dec_base, _ := base64.StdEncoding.DecodeString(bases)
  fmt.Printf("base64decoded is %x\n", dec_base)
  deced := blowfishDecrypt(dec_base, key)
  fmt.Printf("decrypted is %s\n", deced)
}
