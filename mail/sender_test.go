package mail

import (
	"testing"

	"github.com/mnakhaev/simplebank/config"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	cfg, err := config.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	subject := "Test email"
	content := `
	<h1>Hello</h1>
	<p>This is a test email</p>
`
	to := []string{"sometestemail@corp.hahaha.com"}
	attachFiles := []string{"../readme.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
