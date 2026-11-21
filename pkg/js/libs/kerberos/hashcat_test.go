package kerberos

import (
	"testing"

	"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/stretchr/testify/require"
)

func TestHashcatRejectsShortCipher(t *testing.T) {
	short := make([]byte, 8)

	var asrep messages.ASRep
	asrep.EncPart.Cipher = short
	_, err := ASRepToHashcat(asrep)
	require.Error(t, err)

	var ticket messages.Ticket
	ticket.EncPart.Cipher = short
	_, err = TGStoHashcat(ticket, "user")
	require.Error(t, err)

	asrep.EncPart.Cipher = make([]byte, 20)
	got, err := ASRepToHashcat(asrep)
	require.NoError(t, err)
	require.Contains(t, got, "$krb5asrep$")
}
