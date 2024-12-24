package mail

import (
	"testing"

	"github.com/nhan-ngo-usf/NBank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Won million dollars"
	content := `
	<h1>You've just won 1 million dollars. Please click the file below to redeem your rewards.</h1>
	`
	to := []string{"anhvan@usf.edu"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}