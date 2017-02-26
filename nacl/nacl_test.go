package nacl

import (
	"bytes"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryptDecrypt(t *testing.T) {
	Convey("We can test encrypting and decrypting bytes using secretbox (NaCL)", t, func() {
		pad := make([]byte, 32)
		_, err := rand.Read(pad)
		So(err, ShouldBeNil)

		b := []byte("This is a message we'd like to encrypt")
		k := []byte("super weak key")

		out, err := Encrypt(pad, k, b)
		So(err, ShouldBeNil)
		So(bytes.Equal(b, out), ShouldBeFalse)

		msg, err := Decrypt(pad, k, out)
		So(err, ShouldBeNil)

		So(bytes.Equal(b, msg), ShouldBeTrue)

		pad, out, err = RandomPadEncrypt(k, b)
		So(bytes.Equal(b, out), ShouldBeFalse)
		msg, err = Decrypt(pad, k, out)
		So(bytes.Equal(b, msg), ShouldBeTrue)
	})
}