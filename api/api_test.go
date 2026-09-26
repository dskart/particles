package api

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/require"
	golangSsh "golang.org/x/crypto/ssh"
)

func TestParseSSHHostKeyOpenSSH(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	block, err := golangSsh.MarshalPrivateKey(priv, "")
	require.NoError(t, err)

	signer, err := parseSSHHostKey(string(pem.EncodeToMemory(block)))
	require.NoError(t, err)
	require.Equal(t, golangSsh.KeyAlgoED25519, signer.PublicKey().Type())

	_, err = parseSSHHostKey("sm:not-a-key")
	require.Error(t, err)
}
