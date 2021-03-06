# goenc
[![Go Report Card](https://goreportcard.com/badge/github.com/alistanis/goenc)](https://goreportcard.com/report/github.com/alistanis/goenc)
[![GoDoc](https://godoc.org/github.com/alistanis/goenc?status.svg)](https://godoc.org/github.com/alistanis/goenc)
[![codecov](https://codecov.io/gh/alistanis/goenc/branch/master/graph/badge.svg)](https://codecov.io/gh/alistanis/goenc)

Encryption and Decryption functions for Go made easy. Encryption should be as simple as calling Encrypt(key, data) and Decrypt(key, data).

###Note: I am in the process of trying to get this reviewed - use at your own risk

# API

The API is built around the `BlockCipher` interface and the `Session` struct.

`BlockCipher` can be used to encrypt simple messages or small files. 

```go
// BlockCipher represents a cipher that encodes and decodes chunks of data at a time
type BlockCipher interface {
	Encrypt(key, plaintext []byte) ([]byte, error)
   	Decrypt(key, ciphertext []byte) ([]byte, error)
   	KeySize() int
}
```

`Session` can be used to perform key exchanges and send secure messages over a "channel" (`io.ReadWriter`)
It also natively performs key derivation, can handle key exchanges, and can prevent replay attaacks. // that is a joke

###Note: Session has been temporarily removed in order to author a more secure version.

```go
// Session represents a session that can be used to pass messages over a secure channel
type Session struct {
   	Cipher   *Cipher
   	Channel
   	lastSent uint32
   	lastRecv uint32
   	sendKey  *[32]byte
   	recvKey  *[32]byte
}
```

All internal packages implement the `BlockCipher` interface with a `Cipher` struct, allowing for flexibility when working with the `BlockCipher` interface.

Each package can also be used directly, as the `Cipher` struct in each simply wraps public functions.

This package supports the following types of encryption, all using the Go stdlib with exception for NaCL, which uses [Secretbox from the /x/crypto package](https://godoc.org/golang.org/x/crypto/nacl/secretbox):

---
* AES-CBC
* AES-CFB
* AES-CTR 
* AES-GCM
* NaCL - with a user provided pad


* ####All types support encryption + authentication

---

#####Example Usage of package functions
    
```go
key := []byte("some 32 byte key") // obviously this would fail without being 32 bytes
ciphertext, err := gcm.Encrypt(key, []byte("super secret message"))
if err != nil {
    return err
}
plaintext, err := gcm.Decrypt(key, ciphertext)
if err != nil {
    return err  
}
fmt.Println(plaintext) // super secret message
```       
#####Example Usage of the main Cipher struct and a BlockCipher interface - this will perform key derivation

```go
c, err := goenc.NewCipher(goenc.CBC, goenc.InteractiveComplexity)
if err != nil {
    return err
}
ciphertext, err := c.Encrypt(key, []byte("super secret message"))
if err != nil {
    return err       
}
plaintext, err := c.Decrypt(key, ciphertext)
if err != nil {
    return err
}    
fmt.Println(plaintext) // super secret message
```
    
#####Example Usages of Session

Note: Retries and connection breaking are not shown here

######As a server   

```go
// wait until a client connects and performs a key exchange
s, err := goenc.Listen(readWriter, cipher)
// if exchange is bad or none was given, we return
if err != nil {
    return err
}

// s is now a session on the given readWriter (underlying conn) and can wait to receive messages
for {
    msg, err := s.Receive()
    is err != nil {
        // check for closed connection here and break if it is (not shown)
        someErrChan <- err
        continue
    }
    
    msg, err := someMsgParsingFunc(msg)
    if err != nil {
        // garbled message
        someErrChan <- err
        continue
    }
    
    switch msg.Type {
        case SomeCoolThing:
            err = s.Send(someConstMessage)
            if err != nil {
                someErrChan <- err
            }
        default:
            // successfully parsed but we don't know what to do, probably retry parsing
    }
}
```

######As a client

```go    
// initial connection to underlying conn of readWriter
s, err := goenc.Dial(readWriter, cipher)
if err != nil {
    return err
}

// send an initial message
err := s.Send(someMessage)
        is err != nil {
            return err
        }
for {
       
        // wait for response
        msg, err := s.Receive()
        if err != nil {
            someErrChan <- err
            continue
        }
        
        msg, err = someMsgParsingFunc(msg)
        if err != nil {
            // garbled message
            someErrChan <- err
            continue
        }
        
        switch msg.Type {
            case SomeCoolThing:
                err = s.Send(someConstMessage)
                if err != nil {
                    someErrChan <- err
                }
            default:
                // successfully parsed but we don't know what to do, probably retry parsing
        }
    }
```

#SSH Package

The ssh package contains convenience functions for generating and parsing ssh keys. They are a wrapper around the /x/crypto package's ssh package.

TODO
---
```
1. [ ] Get project reviewed (if you are a security expert interested in reviewing this, please contact me and let me know if you find anything)
2. [ ] More complete documentation with examples
    *  [ ] Document full examples of package functions and the small differences they have
    *  [ ] Document SenderID functions in the GCM package and give a real world example
3. [ ] Implement SenderID functions in packages other than GCM
4. [ ] Give user level control over when/how key derivation takes place
    *  The way it works now on a session is that the key will be derived for every message - this is slow, but potentially more secure
       * If one algo has a flaw in which a prior key is discovered, only that message could be read
       * That should still be left up to the user
    *  [ ] Allow user given salt   
```        
        
#Special Thanks

A very special thanks to [Kyle Isom](https://github.com/kisom), whose book provided a very good jumping off point for starting this library.

You can find his book here: [Practical Cryptography with Go](https://leanpub.com/gocrypto/)