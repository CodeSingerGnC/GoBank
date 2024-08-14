package mail

import (
	"testing"

	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithSina(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewSinaSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1> Hello world </h1>
	<p> This is a test message from <a href="https://github.com/CodeSingerGnC"> CodeSingerGnC </a></p>
	`

	to := []string{"test@example.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}