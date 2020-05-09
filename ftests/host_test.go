package ftests

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jingweno/upterm/host/api/swagger/models"
)

func testHostSessionCreatedCallback(t *testing.T, hostURL, nodeAddr string) {
	h := &Host{
		Command:      []string{"bash", "--norc"},
		ForceCommand: []string{"vim"},
		PrivateKeys:  []string{HostPrivateKey},
		SessionCreatedCallback: func(session *models.APIGetSessionResponse) error {
			if want, got := []string{"bash", "--norc"}, session.Command; !cmp.Equal(want, got) {
				t.Fatalf("want=%s got=%s:\n%s", want, got, cmp.Diff(want, got))
			}
			if want, got := []string{"vim"}, session.ForceCommand; !cmp.Equal(want, got) {
				t.Fatalf("want=%s got=%s:\n%s", want, got, cmp.Diff(want, got))
			}

			checkSessionPayload(t, session, hostURL, nodeAddr)
			return nil
		},
	}

	if err := h.Share(hostURL); err != nil {
		t.Fatal(err)
	}
	defer h.Close()
}

func testHostFailToShareWithoutPrivateKey(t *testing.T, hostURL, nodeAddr string) {
	h := &Host{
		Command: []string{"bash"},
	}
	err := h.Share(hostURL)
	if err == nil {
		t.Fatal("expect error")
	}

	if !strings.Contains(err.Error(), "Permission denied (publickey)") {
		t.Fatalf("expect permission denied error: %s", err)
	}
}
